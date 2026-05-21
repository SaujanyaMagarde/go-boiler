package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/SaujanyaMagarde/go-server/internal/database"
	"github.com/google/uuid"
)

func startScrapping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scrapping on %v goroutine every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)

		if err != nil {
			log.Printf("error fetching feeds %v", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)

			go scrapefeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapefeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedFetched(context.Background(), feed.ID)

	if err != nil {
		log.Printf("error fetching feeds %v", err)
		return
	}

	RSSFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("error fetching feeds %v", err)
		return
	}

	for _, item := range RSSFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{String: item.Description, Valid: true}
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("error parsing date %v", err)
			continue
		}

		 _, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			FeedID:      feed.ID,
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: pubAt,
		})

		if err != nil{
			if strings.Contains(err.Error(),"duplicate key"){
				continue;
			}
			log.Printf("error creating post %v", err)
		}
	}
	log.Printf("Fetched %v posts from %v", len(RSSFeed.Channel.Item), feed.Name)
}
