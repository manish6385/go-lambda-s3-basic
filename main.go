package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3session *s3.S3

const (
	REGION      = "us-east-1"
	BUCKET_NAME = "g33kzone-lambda-test"
)

func init() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	})))
}

func main() {
	lambda.Start(Handler)
}

func Handler(event InputEvent) (string, error) {

	image := GetImage(event.Link)
	_, err := s3session.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(image),
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(event.Key),
	})
	if err != nil {
		return "", errors.New("something went wrong")
	}

	return "S3 Upload successful...", nil
}

type InputEvent struct {
	Link string `json:"link"`
	Key  string `json:"key"`
}

func GetImage(url string) (bytes []byte) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return bytes
}
