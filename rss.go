package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct{
	Channel struct{
		Title string `xml:"title"`
		Link string `xml:"link"`
		Description string `xml:"description"`
		Language string `xml:"language"`
		Item []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct{
	Title string `xml:"title"`
	Link string `xml:"link"`
	Description string `xml:"description"`
	PubDate string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSFeed , error){
	httpClient := http.Client{
		Timeout : 10*time.Second,
	} //made a http client with 10 sectimeout

	resp , err := httpClient.Get(url) //This is the actual network call
	if err != nil {
		return RSSFeed{} , err
	}
	defer resp.Body.Close() //"Right before this function finishes, make absolutely sure you close the network connection

	data , err := io.ReadAll(resp.Body) // read and store

	if err != nil{
		return RSSFeed{} , err
	}

	rssFeed := RSSFeed{}
	err = xml.Unmarshal(data , &rssFeed) //xml.Unmarshal takes all the raw XML text from data and perfectly maps it into your Go variables.
	if err != nil{
		return RSSFeed{} , err
	}

	return rssFeed , nil
}

