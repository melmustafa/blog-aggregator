package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/melmustafa/blog-aggregator/internal/database"
)

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PublishedAt string `xml:"pubDate,omitempty"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Item        []Item `xml:"item"`
}

func fetcher(
	wg *sync.WaitGroup,
	feedId uuid.UUID,
	feedUrl string,
	feedData map[uuid.UUID][]Channel,
	net *http.Client,
	mark func(context.Context, string) error,
) {
	defer wg.Done()
	req, _ := http.NewRequest("GET", feedUrl, nil)
	res, err := net.Do(req)
	if err != nil {
		log.Printf("something wrong with the request or the client: %v", err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		log.Printf("something wrong with the request or the client: %v", err)
		return
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("something wrong with the request or the client: %v", err)
		return
	}
	var result struct {
		Channel []Channel `xml:"channel"`
	}
	xml.Unmarshal(data, &result)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	mark(ctx, feedUrl)
	cancel()
	feedData[feedId] = result.Channel
}

func (cfg *apiConfig) FetchService(limit int32, duration time.Duration) {
	ticker := time.NewTicker(duration)
	fmt.Println("fetcher started")
	for {
		fmt.Println("ticker unblocked")
		nextFeeds, err := cfg.DB.GetNextFeedsToFetch(context.TODO(), limit)
		if err != nil {
			log.Printf("something wrong with the database: %v", err)
			continue
		}
		fmt.Println("feeds to retrieve fetched")
		feedData := make(map[uuid.UUID][]Channel)
		var wg sync.WaitGroup
		net := &http.Client{}
		for _, feed := range nextFeeds {
			wg.Add(1)
			go fetcher(&wg, feed.ID, feed.Url, feedData, net, cfg.DB.MarkFeedFetched)
		}
		wg.Wait()
		fmt.Println("feed fetched")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		for _, feed := range nextFeeds {
			for _, channel := range feedData[feed.ID] {
				for _, item := range channel.Item {
					publishedAt, _ := time.Parse(time.RFC1123, item.PublishedAt)
					if feed.LastFetchedAt.Valid && publishedAt.Before(feed.LastFetchedAt.Time) {
						continue
					}
					wg.Add(1)
					go func(item Item, feedID uuid.UUID, publishedAt time.Time) {
						fmt.Println("a new post is being created")
						defer wg.Done()
						description := sql.NullString{}
						if item.Description != "" {
							description.Valid = true
							description.String = item.Description
						}
						_, err := cfg.DB.CreatePost(
							ctx, database.CreatePostParams{
								ID:          uuid.New(),
								Title:       item.Title,
								Description: description,
								Url:         item.Link,
								PublishedAt: sql.NullTime{
									Valid: !publishedAt.Equal(time.Time{}),
									Time:  publishedAt,
								},
								FeedID: uuid.NullUUID{
									Valid: true,
									UUID:  feedID,
								},
							},
						)
						if err != nil {
							log.Printf("couldn't create a post record with the err %v\n", err)
							return
						}
						fmt.Println("new post created")
					}(item, feed.ID, publishedAt)
				}
			}
		}
		wg.Wait()
		cancel()
		fmt.Println("all posts saved")
		<-ticker.C
	}
}
