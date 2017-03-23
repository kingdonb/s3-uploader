package main

import (
	"fmt"
	"net/http"
	"os"

	minio "github.com/minio/minio-go"
)

func main() {
	fmt.Println("Starting s3-uploader...")

	s3Endpoint := os.Getenv("S3_ENDPOINT")
	sslEnvvar := os.Getenv("USE_SSL")
	useSSL := false
	if sslEnvvar != "" {
		useSSL = true
	}
	if s3Endpoint == "" {
		s3Endpoint = "s3.amazonaws.com"
		useSSL = true
	}

	minioClient, err := minio.New(s3Endpoint, os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), useSSL)
	if err != nil {
		fmt.Printf("error: minioClient failed to connect: %s", err)
		os.Exit(1)
	}

	// open test.jpg (or the file override) for upload to s3
	fullPath := "/upload/test.jpg"
	fileOverride := os.Getenv("FILE_OVERRIDE")
	if fileOverride != "" {
		fullPath = fileOverride
	}

	file, err := os.Open(fullPath)
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Failed to get file info: %s", err)
		os.Exit(1)
	}
	size := fileInfo.Size()

	buffer := make([]byte, size)

	// read in file and generate content length and type for PutObjectInput struct
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Printf("Failed to read file: %s", err)
		os.Exit(1)
	}
	fileType := http.DetectContentType(buffer)
	//fileType := "application/octet-stream"

	// read s3 bucket name from envvar
	bucketName := os.Getenv("S3_BUCKET_NAME")
	// exit 1 if S3_BUCKET_NAME envvar is empty
	if bucketName == "" {
		fmt.Printf("error: envvar S3_BUCKET_NAME empty\n")
		os.Exit(1)
	}
	fmt.Printf("Bucket name: %s\n", bucketName)

	ok, err := minioClient.BucketExists(bucketName)
	if err != nil {
		fmt.Printf("failed to determine if bucket exists!")
		os.Exit(1)
	}
	if !ok {
		fmt.Printf("bucket does not exist!")
		os.Exit(1)
	}

	fmt.Printf("Full path: %s\n", fullPath)
	fmt.Printf("File name: %s\n", fileInfo.Name())
	fmt.Printf("File type: %s\n", fileType)
	fmt.Printf("File size: %d\n", size)

	// upload to s3
	resp, err := minioClient.FPutObject(bucketName, fileInfo.Name(), fullPath, fileType)
	if err != nil {
		fmt.Printf("Bad response: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Upload successful. Response: \n%d\n", resp)
}
