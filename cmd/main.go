package main

import (
	"errors"
	"fmt"

	"github.com/Username/Project-Name/pkg/scraper"
	"github.com/Username/Project-Name/pkg/video"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var(
	dynaClient dynamodbiface.DynamoDBAPI
)

func main() {
	region := "eu-west-2"
	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	dynaClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

const tableName = "Videos"
const stopLength = 10
// Steps
// Loop
//    Scrape a search
//    For each video scraped:
//      Check if id has been added or checked already
//      Try to add to DB
//        If in DB, add id to checked ids
//        If added, add id to added ids
//      If number of required ids is reached, return ids added
//    If out of scraped videos, scrape next search

func handler () {
	var addedVideoIds []string
	var checkedVideoIds []string
	fmt.Println("Creating table")
	err := video.NewTable(tableName, dynaClient)
	if err != nil && !errors.Is(err, video.ErrExistingTable) {
		fmt.Printf("error: %s", err.Error())
		fmt.Println("Nah, brev")
		return
	}
	fmt.Println("Table exists")
	fmt.Println("About to scrape")
	scraper.ScrapeYTSearch("midwest emo riff", "EgQIAxAB", func(videos scraper.Contents) {
		scraper.AddVideosToDB(videos, tableName, &addedVideoIds, &checkedVideoIds, stopLength, dynaClient)
	})
	fmt.Println(len(addedVideoIds))
}