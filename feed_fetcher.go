package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}

type Channel struct {
	Item
	Items []Item `xml:"item"`
}

func fetcher(
	wg *sync.WaitGroup,
	feedUrl string,
	feedData map[string][]Channel,
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
	var feeds struct {
		Channels []Channel `xml:"channel"`
	}
	xml.Unmarshal(data, &feeds)
	mark(context.TODO(), feedUrl)
	feedData[feedUrl] = feeds.Channels
}

func (cfg *apiConfig) FetchService(limit int32, duration time.Duration) {
	ticker := time.NewTicker(duration)
	for {
		<-ticker.C
		nextFeeds, err := cfg.DB.GetNextFeedsToFetch(context.TODO(), limit)
		feedData := make(map[string][]Channel)
		if err != nil {
			log.Printf("something wrong with the database: %v", err)
			continue
		}
		var wg sync.WaitGroup
		net := &http.Client{}
		for _, feed := range nextFeeds {
			wg.Add(1)
			go fetcher(&wg, feed.Url, feedData, net, cfg.DB.MarkFeedFetched)
		}
		wg.Wait()
		for _, feed := range feedData {
			for _, channel := range feed {
				for _, item := range channel.Items {
					fmt.Printf("Blog Post Title: %s\n", item.Title)
					fmt.Printf("Blog Post Link: %s\n", item.Link)
					fmt.Printf("Blog Post Description: %s\n", item.Description)
				}
			}
		}
	}
}
