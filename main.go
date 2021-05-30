package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const DEFAULT_FEED_LENGTH = 7

// xml parse
type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PublishDate string `xml:"pubDate"`
}

func pleaseBeNoError(err error) {
	if err != nil {
		println("Oh no we got an error, ", err.Error())
		os.Exit(1)
	}
}

// json serialize
type NewsArchive struct {
	LastUpdated string     `json:"lastUpdated"`
	NewsFeeds   []NewsFeed `json:"news"`
}

type NewsFeed struct {
	Title         string         `json:"title"`
	NewsFeedItems []NewsFeedItem `json:"items"`
}

type NewsFeedItem struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	PublishDate string `json:"pubDate"`
}

func main() {
	feeds := []string{
		"https://www.deutschlandfunk.de/die-nachrichten.353.de.rss",
		"https://www.deutschlandfunk.de/politik.1499.de.rss",
		"https://www.swr.de/export:xml:rss/swraktuell/baden-wuerttemberg/swraktuell-bw-100.html",
		"http://apod.nasa.gov/apod.rss",
		"https://xkcd.com/rss.xml",
		"http://blog.acolyer.org/feed/",
		"http://planet.debian.org/rss20.xml",
	}

	newsArchive := NewsArchive{
		LastUpdated: time.Now().String(),
	}

	fmt.Println("<head> <style> h1,h2 { line-height: 1.3; } body {font-size: x-large; line-height: 1.5; font-family: \"Helvetica Neue\", Helvetica, Arial, sans-serif; } </style> </head>")
	fmt.Println(`<link rel="manifest" href="/manifest.json">`)
	fmt.Printf("<h1>RSS News generated at %s</h1>", time.Now().Local().Format(time.Kitchen))
	for _, feedUrl := range feeds {
		resp, err := http.Get(feedUrl)
		pleaseBeNoError(err)
		defer resp.Body.Close()
		content, err := io.ReadAll(resp.Body)
		pleaseBeNoError(err)
		var newsFeed Feed
		err = xml.Unmarshal(content, &newsFeed)
		pleaseBeNoError(err)
		fmt.Printf("\n\n<h2>%s</h2>\n", newsFeed.Channel.Title)

		currentFeed := NewsFeed{
			Title: newsFeed.Channel.Title,
		}

		// Limit the number of items for readability, thus not using a range loop
		for i := 0; i < upperIndexBound(len(newsFeed.Channel.Items)); i++ {
			newsFeedItem := newsFeed.Channel.Items[i]
			if newsFeedItem.Link == "" || newsFeedItem.Title == "" {
				s := fmt.Sprintf("Can't add entry with link '%s' and title '%s'.", newsFeedItem.Link, newsFeedItem.Title)
				println(s)
			}
			if newsFeedItem.PublishDate == "" {
				fmt.Printf("<a href=\"%s\">%s</a><br/>\n", newsFeedItem.Link, newsFeedItem.Title)
			} else {
				fmt.Printf("<p><a href=\"%s\">%s <i>%s</i></a><p/>\n", newsFeedItem.Link, newsFeedItem.Title, formatDate(newsFeedItem.PublishDate))
			}
			currentFeed.NewsFeedItems = append(currentFeed.NewsFeedItems, NewsFeedItem(newsFeedItem))
		}
		newsArchive.NewsFeeds = append(newsArchive.NewsFeeds, currentFeed)
	}

	newsFeedJSON, err := json.MarshalIndent(newsArchive, "", "  ")
	pleaseBeNoError(err)
	err = ioutil.WriteFile("news-archive.json", newsFeedJSON, 0644)
	pleaseBeNoError(err)
}

func formatDate(inputDateString string) string {
	if inputDateString == "" {
		return ""
	}

	date, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", inputDateString)
	if err != nil {
		panic("fixme " + err.Error())
	}

	return date.Local().Format("2006-01-02 15:04:05")
}

func upperIndexBound(feedLen int) int {
	if feedLen < DEFAULT_FEED_LENGTH {
		return feedLen
	}

	return DEFAULT_FEED_LENGTH
}
