package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"time"
)

// Create struct to hold info about new item
type Item struct {
	ID string
	Year   int
	Title  string
	Plot   string
	Rating float64
}

type MyEvent struct {
	Version string `json:"version"`
	Id string `json:"id"`
	DetailType string `json:"detail-type"`
	Source string `json:"source"`
	Account string `json:"account"`
	Time time.Time `json:"time"`
	Region string `json:"region"`
	Resources []string `json:"resources"`
	Detail interface{} `json:"detail"`
}

var svc *dynamodb.DynamoDB

func init() {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc = dynamodb.New(sess)
}

func HandleRequest(ctx context.Context, event MyEvent) (string, error) {
	//// create the input configuration instance
	//input := &dynamodb.ListTablesInput{}
	//
	//fmt.Printf("Tables:\n")
	//
	//for {
	//	// Get the list of tables
	//	result, err := svc.ListTables(input)
	//	if err != nil {
	//		if aerr, ok := err.(awserr.Error); ok {
	//			switch aerr.Code() {
	//			case dynamodb.ErrCodeInternalServerError:
	//				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
	//			default:
	//				fmt.Println(aerr.Error())
	//			}
	//		} else {
	//			// Print the error, cast err to awserr.Error to get the Code and
	//			// Message from an error.
	//			fmt.Println(err.Error())
	//		}
	//		return "" , err
	//	}
	//
	//	for _, n := range result.TableNames {
	//		log.Println(*n)
	//	}
	//
	//	// assign the last read tablename as the start for our next call to the ListTables function
	//	// the maximum number of table names returned in a call is 100 (default), which requires us to make
	//	// multiple calls to the ListTables function to retrieve all table names
	//	input.ExclusiveStartTableName = result.LastEvaluatedTableName
	//
	//	if result.LastEvaluatedTableName == nil {
	//		break
	//	}
	//}

	item := Item{
		ID: "sssddd",
		Year:   2015,
		Title:  "The Big New Movie",
		Plot:   "Nothing happens at all.",
		Rating: 0.0,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}
	// Create item in table Movies
	tableName := "spike-table"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	year := strconv.Itoa(item.Year)

	log.Println("Successfully added '" + item.Title + "' (" + year + ") to table " + tableName)

	log.Println(event.Version)
	log.Println(event.Id)
	log.Println(event.DetailType)
	log.Println(event.Source)
	log.Println(event.Account)
	log.Println(event.Time)
	log.Println(event.Region)
	//log.Println(event.Resources)
	log.Println(event.Detail)


	return fmt.Sprintf("Hello %s!", event.Account ), nil
}


func main() {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create the cloudwatch events client
	svc := cloudwatchevents.New(sess)


	resultPutRule, _ := svc.PutRule(&cloudwatchevents.PutRuleInput{
		Name:               aws.String("schedule"),
		ScheduleExpression: aws.String("rate(1 minute)"),
		State: 				aws.String("ENABLED"),
	})
	fmt.Println("Rule ARN:", resultPutRule.GoString())
}