package handler

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"net/http"
)

const (
	S3_REGION       = "eu-central-1"
	S3_BUCKET       = "TODO"
	S3_HOST         = "https://TODO.s3.eu-central-1.amazonaws.com"
	DESTINATION_DIR = "TODO/"
)

/*
	expects credentials in env:
	AWS_SECRET_ACESS_KEY
	AWS_ACCESS_KEY_ID

*/

// returns full url to image
func storeS3(buffer []byte, fname string) (string, error) {

	sess, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION)})

	if err != nil {
		return "", err
	}

	var size int64 = int64(len(buffer))

	var destDir string = DESTINATION_DIR + fname

	res, err := s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(S3_BUCKET),
		Key:                  aws.String(destDir),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("inline"),
		ServerSideEncryption: aws.String("AES256"),
	})

	log.Println("S3-Upload: ", res)

	if err != nil {
		log.Panic("S3-Upload error: ", err)
		return "", err
	}

	url := S3_HOST + "/" + destDir

	return url, nil
}
