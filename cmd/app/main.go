package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/summaryProcessor"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func handler(cctx context.Context, s3Event events.S3Event) error {
	for _, record := range s3Event.Records {
		s3 := record.S3
		bucket := s3.Bucket.Name
		key := s3.Object.Key

		// Download the file from S3
		fileBytes, err := downloadFileFromS3(bucket, key)
		if err != nil {
			fmt.Printf("Error downloading file %s/%s from S3: %s\n", bucket, key, err.Error())
			return err
		}
		// Save the file locally
		filePath := "/tmp/" + key
		err = os.MkdirAll("/tmp/archivosCSV", 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %s\n", err.Error())
			return err
		}
		err = ioutil.WriteFile(filePath, fileBytes, 0600)
		if err != nil {
			fmt.Printf("Error saving file to local file system: %s\n", err.Error())
			return err
		}

		// Call your function with the local file path
		summaryProcessor.SummaryReportGenerator(filePath)
	}

	return nil
}

func downloadFileFromS3(bucket string, key string) ([]byte, error) {
	// Create a new S3 downloader
	downloader := s3manager.NewDownloader(session.Must(session.NewSession()))

	// Create a buffer to write the S3 object to
	buffer := aws.NewWriteAtBuffer([]byte{})

	// Download the file from S3
	_, err := downloader.Download(buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" && os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		//  solo si no está en un entorno de Lambda
		summaryProcessor.SummaryReportGenerator("/app/client1.csv")
		summaryProcessor.SummaryReportGenerator("/app/client2.csv")
		return
	}

	lambda.Start(handler)

}
