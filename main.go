package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

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

func main() {
	feeds := []string{
		"https://www.deutschlandfunk.de/politik.1499.de.rss",
		"http://apod.nasa.gov/apod.rss",
		"https://xkcd.com/rss.xml",
		"http://blog.acolyer.org/feed/",
		"http://planet.debian.org/rss20.xml",
	}

	for _, feedUrl := range feeds {
		resp, err := http.Get(feedUrl)
		pleaseBeNoError(err)
		defer resp.Body.Close()
		content, err := io.ReadAll(resp.Body)
		pleaseBeNoError(err)
		var newsFeed Feed
		err = xml.Unmarshal(content, &newsFeed)
		pleaseBeNoError(err)
		fmt.Printf("\n\n<h1>%s</h1>\n", newsFeed.Channel.Title)
		for _, newsFeedItem := range newsFeed.Channel.Items {
			if newsFeedItem.PublishDate == "" {
				fmt.Printf("<a href=\"%s\">%s</a><br/>\n", newsFeedItem.Link, newsFeedItem.Title)

			} else {
				date, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", newsFeedItem.PublishDate)
				pleaseBeNoError(err)
				fmt.Printf("<a href=\"%s\">%s</a> <i>%s</i><br/>\n", newsFeedItem.Link, newsFeedItem.Title, date.Local().Format("2006-01-02 15:04:05"))
			}
		}
	}
}
