package s3store

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/erpe/image_service_go/app/config"
	"log"
	"net/http"
)

var (
	S3_REGION string
	S3_BUCKET string
	S3_HOST   string
	S3_FOLDER string
)

func init() {
	s3cfg := config.GetConfig().S3
	S3_REGION = s3cfg.Region
	S3_BUCKET = s3cfg.Bucket
	S3_HOST = s3cfg.Host
	S3_FOLDER = s3cfg.Folder
}

/*
	expects credentials in env:
	AWS_SECRET_ACESS_KEY
	AWS_ACCESS_KEY_ID

*/

// returns full url to image
func SaveImage(buffer []byte, fname string) (string, error) {

	sess, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION)})

	if err != nil {
		return "", err
	}

	var size int64 = int64(len(buffer))

	var destDir string = S3_FOLDER + fname

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
