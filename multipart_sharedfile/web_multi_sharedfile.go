package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func indexHandle(w http.ResponseWriter, r *http.Request) {
	uploadTemplate := template.Must(template.ParseFiles("index.html"))
	uploadTemplate.Execute(w, nil)
	// refer to index.html in the same folder
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {
	file, head, _ := r.FormFile("file")
	defer file.Close()
	bts, _ := ioutil.ReadAll(file)
	//read uploaded file
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-1"),
		Credentials: credentials.NewSharedCredentials("", "test-account"),
	}))
	// create session of s3 and loading credential from .aws folder
	// aws.Config
	// (Region:) to specify the region where the s3 bucket at
	// (Credentials:) to specify the credential "test-account" at the .aws/credentials
	svc := s3manager.NewUploader(sess)
	// create a s3 uploader service
	input := &s3manager.UploadInput{
		Bucket:               aws.String("bucket-name"),
		Key:                  aws.String(head.Filename),
		ServerSideEncryption: aws.String("AES256"),
		Body:                 bytes.NewReader(bts),
	}
	//upload file to s3 Bucket
	result, err := svc.Upload(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Println(result)
}

func main() {
	http.HandleFunc("/", indexHandle)
	// index handler
	http.HandleFunc("/upload", uploadHandle)
	// upload hadler
	http.ListenAndServe(":8080", nil)
	// http service Listen port
}
