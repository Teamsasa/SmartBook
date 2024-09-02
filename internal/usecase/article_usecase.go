package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type Article struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Score       int    `json:"score"`
	By          string `json:"by"`
	Time        int64  `json:"time"`
	Description string `json:"description"`
}

type ArticleUseCase struct {
	client *http.Client
}

func NewArticleUseCase() *ArticleUseCase {
	return &ArticleUseCase{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (u *ArticleUseCase) GetArticles() ([]Article, error) {
	resp, err := u.client.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}

	articles := make([]Article, 0, 30)
	for i := 0; i < 30 && i < len(ids); i++ {
		article, err := u.getArticleByID(ids[i])
		if err != nil {
			continue
		}
		articles = append(articles, article)
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Score > articles[j].Score
	})

	return articles, nil
}

func (u *ArticleUseCase) getArticleByID(id int) (Article, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	resp, err := u.client.Get(url)
	if err != nil {
		return Article{}, err
	}
	defer resp.Body.Close()

	var article Article
	if err := json.NewDecoder(resp.Body).Decode(&article); err != nil {
		return Article{}, err
	}

	return article, nil
}

func (u *ArticleUseCase) GetArticleByID(id int) (*Article, error) {
	article, err := u.getArticleByID(id)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (u *ArticleUseCase) GetRecommendedArticles() ([]Article, error) {
	// この実装では、単純にトップ記事を "おすすめ" として返します
	// 実際のアプリケーションでは、ユーザーの興味に基づいてパーソナライズされた推薦を行うべきです
	return u.GetArticles()
}

func (u *ArticleUseCase) GetArticleContent(url string) (string, error) {
	resp, err := u.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// ここで、実際のコンテンツを取得するロジックを実装します
	// 例えば、HTMLをパースしてメインコンテンツを抽出するなど
	// この例では、簡単のために全体のHTMLを返しています
	var content []byte
	_, err = resp.Body.Read(content)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
