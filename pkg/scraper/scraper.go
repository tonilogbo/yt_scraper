package scraper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Username/Project-Name/pkg/video"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/gocolly/colly/v2"
)

func ScrapeYTSearch(search, filter string, callback func(videos Contents)) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.youtube.com"),
	)
		c.OnHTML("html", func(e *colly.HTMLElement) {
		var dataString string
		e.ForEachWithBreak("script", func(i int, h *colly.HTMLElement) bool {
			if strings.Contains(h.Text, "ytInitialData") {
				dataString = strings.TrimPrefix(strings.TrimSuffix(h.Text, ";"), "var ytInitialData = ")
				return false
			}
			return true
		})
		if json.Valid([]byte(dataString)) {
			var y YtInitialData
			json.Unmarshal([]byte(dataString), &y)
			aa := y.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents[0].ItemSectionRenderer.Contents
			callback(aa)
		}
	})
	c.Visit(GetSearchURL(search, filter))
}

func AddVideosToDB (videos Contents, tableName string, addedVideoIds, checkedVideoIds *[]string, stopLength int, dynaClient dynamodbiface.DynamoDBAPI) {
	for _, s := range videos {
		fmt.Println("Extracting video info")
		vid, err := ExtractVideoInfo(s.VideoRenderer)
		if err != nil {
			fmt.Println("Failed extracting video info")
			continue
		}
		// TODO: Add validator logic to filter videos
		// Let's see if we've already added this video 
		addedVideoIndex := checkSlice(*addedVideoIds, vid.VideoID)
		if addedVideoIndex != -1 {
			fmt.Println("Id already added")
			continue
		}
		checkedVideoIndex := checkSlice(*checkedVideoIds, vid.VideoID)
		if checkedVideoIndex != -1 {
			fmt.Println("Id already checked")
			continue
		}
		fmt.Println("Attempting to add video " + vid.VideoID)
		err = video.AddVideo(vid, tableName, dynaClient)
		if err != nil {
			fmt.Printf("error: %s", err.Error())
			*checkedVideoIds = append(*checkedVideoIds, vid.VideoID)
		} else {
			fmt.Println(vid.VideoID + " added")
			*addedVideoIds = append(*addedVideoIds, vid.VideoID)
		}
		if stopLength > 0 && len(*addedVideoIds) == stopLength {
			fmt.Println("Added enough videos")
			break
		}
	}
	fmt.Println("All done", len(*addedVideoIds))
}

func checkSlice (slice []string, value string) int {
for i, v := range slice {
	if value == v {
		return i
		}
	}
	return -1
}

func ExtractVideoInfo(v VideoRenderer) (video.YTVideoInfo, error) {

	return video.YTVideoInfo{
		Title: v.Title.Runs[0].Text,
		Owner: v.OwnerText.Runs[0].Text,
		ThumbnailUrl: func() string {
			if len(v.Thumbnail.Thumbnails) > 1 {
				return v.Thumbnail.Thumbnails[1].Url
			} else {
				return v.Thumbnail.Thumbnails[0].Url
			}
		}(),
		VideoID: v.VideoId,
		Length:  v.LengthText.SimpleText,
	}, nil
}

func GetSearchURL(search, filter string) string {
	if filter == "" {
		return "https://www.youtube.com/results?search_query=" + strings.ReplaceAll(search, " ", "+")
	}
	return "https://www.youtube.com/results?search_query=" + strings.ReplaceAll(search, " ", "+") + "&sp=" + filter
}