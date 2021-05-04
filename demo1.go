package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/satori/go.uuid"
	"log"
	"os"
	"time"
)

type WebFile struct {
	ID 			string
	FileName   	string
	Status 		string
	Created		time.Time
	Updated		time.Time
}

func main() {

	dir := "c:\\"
	fileName := "webdata.txt"

	localFile,err := os.Create(dir + fileName)

	if err !=nil {
		fmt.Println(err.Error())
	} else {
		_,_ = localFile.Write([]byte("web content data"))
	}
	_ = localFile.Close()
	bucket := "spike-bucket"


	file, err := os.Open(dir + fileName)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}

	defer file.Close()

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("cn-northwest-1")},
	)

	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	// for more information on configuring part size, and concurrency.
	//
	// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
	uploader := s3manager.NewUploader(sess)

	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),

		// Can also use the `filepath` standard library package to modify the
		// filename as need for an S3 object key. Such as turning absolute path
		// to a relative path.
		Key: aws.String(fileName),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: file,
	})
	//fmt.Println(resp.ETag, resp.Location, resp.UploadID, resp.VersionID)
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", fileName, bucket, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", fileName, bucket)

	item := WebFile{
		ID: GetUUID(),
		FileName:   fileName,
		Status: "created",
		Created: time.Now(),
		Updated: time.Now(),

	}
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	var svc *dynamodb.DynamoDB
	svc = dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Printf("Got error marshalling new webdata file: %s\n", err)
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


	fmt.Println("Successfully added '" + item.FileName + "' to table " + tableName)

}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func GetUUID() string {
	u2:= uuid.NewV4()
	return u2.String()
}
