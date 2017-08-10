package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func indexHandle(w http.ResponseWriter, r *http.Request) {
	uploadTemplate := template.Must(template.ParseFiles("index.html"))
	uploadTemplate.Execute(w, nil)
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {
	file, head, _ := r.FormFile("file")
	defer file.Close()
	bts, _ := ioutil.ReadAll(file)
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}
	svc := s3.New(sess, &aws.Config{Region: aws.String("ap-northeast-1")})
	input := &s3.PutObjectInput{
		Body:                 bytes.NewReader(bts),
		Bucket:               aws.String("iloterry-yehuwes67yuk"),
		Key:                  aws.String(head.Filename),
		ServerSideEncryption: aws.String("AES256"),
	}

	result, err := svc.PutObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
	//ioutil.WriteFile(head.Filename, bytes, os.ModeAppend)
}

func main() {
	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/upload", uploadHandle)
	http.ListenAndServe(":8080", nil)
}
