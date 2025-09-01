package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RssUrl struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	RssUrl      string    `json:"url" gorm:"column:rss_url"`
	IsSummarize bool      `json:"is_active" gorm:"column:is_summarize"`
	CreatedAt   time.Time `json:"created_at"`
}

func (RssUrl) TableName() string {
	return "rss_urls"
}

func connectDB() (*gorm.DB, error) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	return db, nil
}

func fetchRssURLs(db *gorm.DB) ([]RssUrl, error) {
	var feeds []RssUrl
	result := db.Find(&feeds)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch RSS feeds: %w", result.Error)
	}
	return feeds, nil
}

type Article struct {
	RssId       int       `json:"rss_id"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published"`
}

func (Article) TableName() string {
	return "articles"
}

func saveArticles(db *gorm.DB, articles []Article) error {
	for _, rssGroup := range groupArticlesByRssId(articles) {
		// 各RSS IDごとに最新の記事の公開日を取得
		var latestArticle Article
		result := db.Where("rss_id = ?", rssGroup[0].RssId).
			Order("published_at DESC").
			First(&latestArticle)

		var articlesToSave []Article
		for _, article := range rssGroup {
			// 最新の記事の公開日以降の記事のみ保存
			if result.Error == gorm.ErrRecordNotFound || article.PublishedAt.After(latestArticle.PublishedAt) {
				articlesToSave = append(articlesToSave, article)
			}
		}

		// 新しい記事を一括保存
		if len(articlesToSave) > 0 {
			if err := db.Create(&articlesToSave).Error; err != nil {
				return fmt.Errorf("failed to save articles for RSS ID %d: %w", rssGroup[0].RssId, err)
			}
		}
	}
	return nil
}

func groupArticlesByRssId(articles []Article) map[int][]Article {
	grouped := make(map[int][]Article)
	for _, article := range articles {
		grouped[article.RssId] = append(grouped[article.RssId], article)
	}
	return grouped
}

func fetchArticlesFromRSS(feed RssUrl) ([]Article, error) {
	fp := gofeed.NewParser()
	parsedFeed, err := fp.ParseURL(feed.RssUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed %s: %w", feed.RssUrl, err)
	}

	var articles []Article
	for _, item := range parsedFeed.Items {
		published := time.Now()
		if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		}

		article := Article{
			RssId:       feed.ID,
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			PublishedAt: published,
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func main() {
	envPath := filepath.Join("..", "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: Could not load .env file from %s: %v", envPath, err)
	}

	db, err := connectDB()
	if err != nil {
		log.Printf("データベース接続エラー: %v", err)
		return
	}

	feeds, err := fetchRssURLs(db)
	if err != nil {
		log.Printf("RSSフィード取得エラー: %v", err)
		return
	}

	if len(feeds) == 0 {
		log.Println("アクティブなRSSフィードが見つかりませんでした")
		return
	}

	var allArticles []Article
	for _, feed := range feeds {
		articles, err := fetchArticlesFromRSS(feed)
		if err != nil {
			log.Printf("Warning: %sからの記事取得に失敗: %v", feed.RssUrl, err)
			continue
		}
		allArticles = append(allArticles, articles...)
	}

	// 記事を保存
	if err := saveArticles(db, allArticles); err != nil {
		log.Printf("記事の保存に失敗: %v", err)
		return
	}

	log.Printf("合計 %d 件の新しい記事を正常に保存しました", len(allArticles))
}
