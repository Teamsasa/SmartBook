package usecase

import (
	"SmartBook/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type ArticleFetcher interface {
	FetchArticles(ctx context.Context, limit int) ([]model.Article, error)
}

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, expiration time.Duration)
}

type ArticleUseCase struct {
	client            *http.Client
	hackerNewsFetcher ArticleFetcher
	devToFetcher      ArticleFetcher
	cache             Cache
}

func NewArticleUseCase(client *http.Client, cache Cache) *ArticleUseCase {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &ArticleUseCase{
		client:            client,
		hackerNewsFetcher: &HackerNewsFetcher{client: client},
		devToFetcher:      &DevToFetcher{client: client},
		cache:             cache,
	}
}

func (u *ArticleUseCase) GetAllArticles(ctx context.Context) ([]model.Article, error) {
	if cachedArticles, found := u.cache.Get("all_articles"); found {
		return cachedArticles.([]model.Article), nil
	}

	var articles []model.Article
	var mu sync.Mutex
	g, ctx := errgroup.WithContext(ctx)

	fetchAndAppend := func(fetcher ArticleFetcher, limit int) func() error {
		return func() error {
			fetchedArticles, err := fetcher.FetchArticles(ctx, limit)
			if err != nil {
				return fmt.Errorf("failed to fetch articles: %w", err)
			}
			mu.Lock()
			articles = append(articles, fetchedArticles...)
			mu.Unlock()
			return nil
		}
	}

	g.Go(fetchAndAppend(u.hackerNewsFetcher, 100))
	g.Go(fetchAndAppend(u.devToFetcher, 100))

	if err := g.Wait(); err != nil {
		return nil, err
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].CreatedAt.After(articles[j].CreatedAt)
	})

	u.cache.Set("all_articles", articles, 5*time.Minute)
	return articles, nil
}

func (u *ArticleUseCase) GetLatestArticles(ctx context.Context) ([]model.Article, error) {
	allArticles, err := u.GetAllArticles(ctx)
	if err != nil {
		return nil, err
	}

	if len(allArticles) > 30 {
		return allArticles[:30], nil
	}
	return allArticles, nil
}

func (u *ArticleUseCase) GetRecommendedArticles(ctx context.Context, userInterests []string) ([]model.Article, error) {
	allArticles, err := u.GetAllArticles(ctx)
	if err != nil {
		return nil, err
	}

	scoredArticles := make([]struct {
		Article model.Article
		Score   float64
	}, len(allArticles))

	for i, article := range allArticles {
		score := float64(article.Score)
		for _, interest := range userInterests {
			if strings.Contains(strings.ToLower(article.Title), strings.ToLower(interest)) {
				score += 100
			}
			for _, tag := range article.Tags {
				if strings.EqualFold(tag, interest) {
					score += 50
				}
			}
		}
		scoredArticles[i] = struct {
			Article model.Article
			Score   float64
		}{Article: article, Score: score}
	}

	sort.Slice(scoredArticles, func(i, j int) bool {
		return scoredArticles[i].Score > scoredArticles[j].Score
	})

	recommendedArticles := make([]model.Article, 0, 30)
	for i := 0; i < 30 && i < len(scoredArticles); i++ {
		recommendedArticles = append(recommendedArticles, scoredArticles[i].Article)
	}

	return recommendedArticles, nil
}

func (u *ArticleUseCase) GetArticleByID(ctx context.Context, id string) (*model.Article, error) {
	articles, err := u.GetAllArticles(ctx)
	if err != nil {
		return nil, err
	}

	for _, article := range articles {
		if article.ID == id {
			return &article, nil
		}
	}

	return nil, fmt.Errorf("article not found: %s", id)
}

type HackerNewsFetcher struct {
	client *http.Client
}

func (f *HackerNewsFetcher) FetchArticles(ctx context.Context, limit int) ([]model.Article, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://hacker-news.firebaseio.com/v0/topstories.json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch top stories: %w", err)
	}
	defer resp.Body.Close()

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	articles := make([]model.Article, 0, limit)
	for i := 0; i < limit && i < len(ids); i++ {
		article, err := f.getHackerNewsArticleByID(ctx, ids[i])
		if err != nil {
			continue
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (f *HackerNewsFetcher) getHackerNewsArticleByID(ctx context.Context, id int) (model.Article, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return model.Article{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return model.Article{}, fmt.Errorf("failed to fetch article: %w", err)
	}
	defer resp.Body.Close()

	var articleData struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		URL   string `json:"url"`
		Score int    `json:"score"`
		By    string `json:"by"`
		Time  int64  `json:"time"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&articleData); err != nil {
		return model.Article{}, fmt.Errorf("failed to decode article: %w", err)
	}

	return model.Article{
		ID:        fmt.Sprintf("hn_%d", articleData.ID),
		Title:     articleData.Title,
		URL:       articleData.URL,
		Score:     articleData.Score,
		Author:    articleData.By,
		CreatedAt: time.Unix(articleData.Time, 0),
		Source:    "Hacker News",
	}, nil
}

type DevToFetcher struct {
	client *http.Client
}

func (f *DevToFetcher) FetchArticles(ctx context.Context, limit int) ([]model.Article, error) {
	url := fmt.Sprintf("https://dev.to/api/articles?top=1&per_page=%d", limit)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch articles: %w", err)
	}
	defer resp.Body.Close()

	var devToArticles []struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		URL         string `json:"url"`
		PublishedAt string `json:"published_at"`
		User        struct {
			Name string `json:"name"`
		} `json:"user"`
		PositiveReactionsCount int      `json:"positive_reactions_count"`
		Tags                   []string `json:"tag_list"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&devToArticles); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	articles := make([]model.Article, 0, len(devToArticles))
	for _, a := range devToArticles {
		publishedTime, _ := time.Parse(time.RFC3339, a.PublishedAt)
		article := model.Article{
			ID:        fmt.Sprintf("dev_%d", a.ID),
			Title:     a.Title,
			URL:       a.URL,
			Score:     a.PositiveReactionsCount,
			Author:    a.User.Name,
			CreatedAt: publishedTime,
			Source:    "DEV.to",
			Tags:      a.Tags,
		}
		articles = append(articles, article)
	}

	return articles, nil
}
