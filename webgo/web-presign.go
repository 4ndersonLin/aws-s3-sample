package main

import (
	//"bytes"
	"fmt"
	"html/template"
	"time"
	"log"
	//"io/ioutil"
	"net/http"
	//"encoding/json"

	
	"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	
)

var(
	myBucket string
	region string
	obj_name string
	
)

func indexHandle(w http.ResponseWriter, r *http.Request) {
	index := template.Must(template.ParseFiles("index-presign.html"))
	index.Execute(w, nil)
}

func presignHandle(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "POST" {
		r.ParseForm()
		obj_name := r.Form.Get("obj_name")
		myBucket := "aws-archi"
		region := "ap-northeast-1"
		sess, err := session.NewSession()
		if err != nil {
			fmt.Println("failed to create session,", err)
			return
		}
		
		// create s3 connection
		svc := s3.New(sess, &aws.Config{Region: aws.String(region)})

		// create a s3.getobjectrequest using local s3 sdk
		req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(myBucket),
			Key:    aws.String(obj_name),
		})

		// Sign url using local credential and set timeout = 15 mins
		urlStr, err := req.Presign(15 * time.Minute)
		if err != nil {
			log.Println("Failed to sign request", err)
		}
		
		fmt.Fprintln(w, obj_name)
		fmt.Fprint(w, "Please PUT object to ")
		fmt.Fprintln(w, urlStr)
		//fmt.Fprint(w, "POST done")

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	
}

func main() {
	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/presign", presignHandle)
	//http.HandleFunc("/upload", uploadHandle)
	http.ListenAndServe(":8080", nil)
}
