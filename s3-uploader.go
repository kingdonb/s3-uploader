package main

import (
	"fmt"
	"net/http"
	"os"

	minio "github.com/minio/minio-go"
)

func main() {
	s3Endpoint := os.Getenv("S3_ENDPOINT")
	useSSL := false
	if s3Endpoint == "" {
		s3Endpoint = "https://s3.amazonaws.com"
		useSSL = true
	}

	minioClient, err := minio.New(s3Endpoint, os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), useSSL)
	if err != nil {
		fmt.Printf("error: minioClient failed to connect: %s", err)
		os.Exit(1)
	}

	// open test.jpg for upload to s3
	fileName := "test.jpg"
	file, err := os.Open(fmt.Sprintf("/upload/%s", fileName))
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()

	buffer := make([]byte, size)

	// read in file and generate content length and type for PutObjectInput struct
	file.Read(buffer)
	fileType := http.DetectContentType(buffer)

	// read s3 bucket name from envvar
	bucketName := os.Getenv("S3_BUCKET_NAME")
	// exit 1 if S3_BUCKET_NAME envvar is empty
	if bucketName == "" {
		fmt.Printf("error: envvar S3_BUCKET_NAME empty\n")
		os.Exit(1)
	}
	fmt.Printf("Bucket name: %s\n", bucketName)

	// upload to s3
	resp, err := minioClient.PutObject(bucketName, fileName, file, fileType)
	if err != nil {
		fmt.Printf("bad response: %s", err)
		os.Exit(1)
	}

	fmt.Printf("upload successful. response: \n%d\n", resp)
}
