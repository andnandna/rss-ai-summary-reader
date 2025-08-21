package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func getRSSFeed(url string) (*RSS, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rss RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil, fmt.Errorf("failed to parse RSS XML: %w", err)
	}

	return &rss, nil
}

func main() {
	rssURL := "https://andna.dev/rss.xml"
	
	rss, err := getRSSFeed(rssURL)
	if err != nil {
		response := Response{
			Success: false,
			Error:   err.Error(),
		}
		jsonBytes, _ := json.Marshal(response)
		fmt.Println(string(jsonBytes))
		return
	}

	response := Response{
		Success: true,
		Data: map[string]interface{}{
			"title":       rss.Channel.Title,
			"link":        rss.Channel.Link,
			"description": rss.Channel.Description,
			"itemCount":   len(rss.Channel.Items),
			"items":       rss.Channel.Items,
		},
	}

	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		errorResponse := Response{
			Success: false,
			Error:   fmt.Sprintf("failed to marshal response: %v", err),
		}
		errorBytes, _ := json.Marshal(errorResponse)
		fmt.Println(string(errorBytes))
		return
	}

	fmt.Println(string(jsonBytes))
}
