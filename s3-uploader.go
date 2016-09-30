package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	/*
	   grab AWS credentials from envvars
	       AWS_ACCESS_KEY_ID
	       AWS_SECRET_ACCESS_KEY
	*/
	creds := credentials.NewEnvCredentials()

	// test that creds have been read in correctly
	_, err := creds.Get()
	if err != nil {
		fmt.Printf("Error getting credentails: %s\n", err)
		os.Exit(1)
	}

	/*
	   create aws client config for use with s3
	   hardcoded region as it is irrelevant for s3
	*/
	config := aws.NewConfig().WithRegion("us-east-1").WithCredentials(creds)
	//create s3 client session with config
	s3Client := s3.New(session.New(), config)

	// open test.jpg for upload to s3
	file, err := os.Open("test.jpg")
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()

	buffer := make([]byte, size)

	// read in file and generate content length and type for PutObjectInput struct
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	// read s3 bucket name from envvar
	bucketName := os.Getenv("S3_BUCKET_NAME")
	// exit 1 if S3_BUCKET_NAME envvar is empty
	if bucketName == "" {
		fmt.Printf("error: envvar S3_BUCKET_NAME empty\n")
		os.Exit(1)
	}
	fmt.Printf("Bucket name: %s\n", bucketName)
	// upload path will be /upload/test.jpg
	path := "/upload/" + file.Name()

	// build params for s3 upload
	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
	// upload to s3
	resp, err := s3Client.PutObject(params)
	if err != nil {
		fmt.Printf("bad response: %s", err)
	}

	fmt.Printf("upload successful. response: \n%s\n", awsutil.StringValue(resp))
}
