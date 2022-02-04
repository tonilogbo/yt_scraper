package video

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type YTVideoInfo struct {
	Title        string `dynamodbav:"title"`
	Owner        string `dynamodbav:"owner"`
	ThumbnailUrl string `dynamodbav:"thumbnailUrl"`
	VideoID      string `dynamodbav:"videoId"`
	Length       string `dynamodbav:"length"`
}

var (
	ErrExistingVideo = errors.New("video exists")
	ErrCouldNotMarshal = errors.New("could not marshal")
	ErrPutFail = errors.New("could not put")
	ErrExistingTable = errors.New("table Exists")
	ErrCreateTable = errors.New("problem creating table")
)
// Add a video to DynamoDB
// 1. Check if video in DB, search for videoID.
func FetchVideo(videoId, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*YTVideoInfo, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"videoId":{
				S: aws.String(videoId),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New("failed to fetch")
	}
	if result.Item == nil {
		return nil, errors.New("no video found")
	}

	item := new(YTVideoInfo)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New("failed to unmarshal")
	}
	return item, nil
}
// 2. If in DE, leave. If not, put Item with w/c tag
func AddVideo(video YTVideoInfo, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	existingVideo, _ := FetchVideo(video.VideoID, tableName, dynaClient)
	if existingVideo != nil {
		return ErrExistingVideo
	}

	av, err := dynamodbattribute.MarshalMap(video)

	if err != nil {
		return ErrCouldNotMarshal
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return ErrPutFail
	}
	return nil
}
func NewTable(tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {

	describeInput := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
	existingTable , err := dynaClient.DescribeTable(describeInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				return err
			default:
				fmt.Println(aerr.Error())
				return err
			}
		}
	}
	fmt.Println("output: ", existingTable)
	if existingTable != nil && existingTable.Table != nil && existingTable.Table.TableStatus != nil && *existingTable.Table.TableStatus == "ACTIVE" {
		return ErrExistingTable
	}
	fmt.Println("Creating input for table")
	 input := &dynamodb.CreateTableInput{
        AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
                AttributeName: aws.String("videoId"),
                AttributeType: aws.String("S"),
            },
        },
        KeySchema: []*dynamodb.KeySchemaElement{
            {
                AttributeName: aws.String("videoId"),
                KeyType:       aws.String("HASH"),
            },
        },
        ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
            ReadCapacityUnits:  aws.Int64(1),
            WriteCapacityUnits: aws.Int64(1),
        },
        TableName: aws.String(tableName),
    }

	result, err := dynaClient.CreateTable(input)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return ErrCreateTable
	}
	if result != nil && result.TableDescription != nil && result.TableDescription.TableStatus != nil && *result.TableDescription.TableStatus == "CREATING" {
		ready := false
		timeTaken := 0
		for !ready {
			existingTable, _ = dynaClient.DescribeTable(describeInput)
			if existingTable != nil && existingTable.Table != nil && existingTable.Table.TableStatus != nil && *existingTable.Table.TableStatus == "ACTIVE" {
				ready = true
				fmt.Println("Took about " + strconv.Itoa(timeTaken) + " to get online")
			} else {
				time.Sleep(2 * time.Second)
				timeTaken = timeTaken + 2
			}
		}
		
	}
	return err
}