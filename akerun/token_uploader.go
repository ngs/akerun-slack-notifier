package akerun

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang.org/x/oauth2"
)

// UploadToken uploads token to S3 bucket
func UploadToken(token *oauth2.Token) error {
	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)
	json, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}
	bucket := os.Getenv("S3_BUCKET")
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("token.json"),
		Body:   bytes.NewReader(json),
	})
	return err
}
