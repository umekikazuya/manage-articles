package src

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	QiitaAPIEndpoint = "https://qiita.com/api/v2"
)

// QiitaClient はQiita APIと通信するためのクライアント
type QiitaClient struct {
	AccessToken string
	HTTPClient  *http.Client
}

// NewQiitaClient は新しいQiitaClientを作成します
func NewQiitaClient(accessToken string) *QiitaClient {
	return &QiitaClient{
		AccessToken: accessToken,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// Item はQiitaの投稿データを表します
type Item struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	URL        string    `json:"url"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Tags       []Tag     `json:"tags"`
	LikesCount int       `json:"likes_count"`
}

// Tag はQiitaの記事のタグを表します
type Tag struct {
	Name string `json:"name"`
}

// GetAuthenticatedUserItems は認証済みユーザーの投稿一覧を取得します
func (c *QiitaClient) GetAuthenticatedUserItems(page, perPage int) ([]Item, error) {
	url := fmt.Sprintf("%s/authenticated_user/items?page=%d&per_page=%d", QiitaAPIEndpoint, page, perPage)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("APIリクエストエラー: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("APIエラー: ステータスコード %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	var items []Item
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("JSONデコードエラー: %w", err)
	}

	return items, nil
}

// GetAllUserItems はすべてのユーザー投稿を取得します（ページネーション対応）
func (c *QiitaClient) GetAllUserItems() ([]Item, error) {
	var allItems []Item
	page := 1
	perPage := 100

	for {
		items, err := c.GetAuthenticatedUserItems(page, perPage)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, items...)

		if len(items) < perPage {
			break
		}

		page++
	}

	return allItems, nil
}
