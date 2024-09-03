package usecase

import (
	"SmartBook/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type ArticleUseCase struct {
	client *http.Client
}

func NewArticleUseCase() *ArticleUseCase {
	return &ArticleUseCase{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (u *ArticleUseCase) GetLatestArticles() ([]model.Article, error) {
	var articles []model.Article
	var errors []error

	hackerNewsArticles, err := u.getHackerNewsArticles()
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to get Hacker News articles: %w", err))
	} else {
		articles = append(articles, hackerNewsArticles...)
	}

	devToArticles, err := u.getDevToArticles()
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to get DEV.to articles: %w", err))
	} else {
		articles = append(articles, devToArticles...)
	}

	if len(articles) == 0 && len(errors) > 0 {
		return nil, fmt.Errorf("failed to fetch articles: %v", errors)
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].CreatedAt.After(articles[j].CreatedAt)
	})

	if len(articles) > 30 {
		return articles[:30], nil
	}
	return articles, nil
}

func (u *ArticleUseCase) getHackerNewsArticles() ([]model.Article, error) {
	resp, err := u.client.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}

	articles := make([]model.Article, 0, 30)
	for i := 0; i < 30 && i < len(ids); i++ {
		article, err := u.getHackerNewsArticleByID(ids[i])
		if err != nil {
			continue
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (u *ArticleUseCase) getHackerNewsArticleByID(id int) (model.Article, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	resp, err := u.client.Get(url)
	if err != nil {
		return model.Article{}, err
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
		return model.Article{}, err
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

func (u *ArticleUseCase) getDevToArticles() ([]model.Article, error) {
	url := "https://dev.to/api/articles?top=30"
	resp, err := u.client.Get(url)
	if err != nil {
		return nil, err
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
		return nil, err
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

func (u *ArticleUseCase) GetRecommendedArticles(userInterests []string) ([]model.Article, error) {
	// 現在は記事の題名とタグに興味を持っているかどうかでスコアを計算
	allArticles, err := u.GetLatestArticles()
	if err != nil {
		return nil, err
	}

	// 記事をスコア付け
	scoredArticles := make([]struct {
		Article model.Article
		Score   float64
	}, len(allArticles))

	for i, article := range allArticles {
		score := float64(article.Score)

		// ユーザーの興味に基づいてスコアを調整
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

	// スコアに基づいてソート
	sort.Slice(scoredArticles, func(i, j int) bool {
		return scoredArticles[i].Score > scoredArticles[j].Score
	})

	// 上位30記事を返す
	recommendedArticles := make([]model.Article, 0, 30)
	for i := 0; i < 30 && i < len(scoredArticles); i++ {
		recommendedArticles = append(recommendedArticles, scoredArticles[i].Article)
	}

	return recommendedArticles, nil
}

func (u *ArticleUseCase) GetArticleByID(id string) (*model.Article, error) {
	articles, err := u.GetLatestArticles()
	if err != nil {
		return nil, err
	}

	for _, article := range articles {
		if article.ID == id {
			return &article, nil
		}
	}

	return nil, errors.New("article not found")
}
