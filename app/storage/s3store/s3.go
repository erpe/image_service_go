package s3store

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/erpe/image_service_go/app/config"
	"image"
	_ "image/gif"
	_ "image/png"
	//_ "image/tif"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	S3_REGION string
	S3_BUCKET string
	S3_HOST   string
)

func init() {
	s3cfg := config.GetConfig().S3
	S3_REGION = s3cfg.Region
	S3_BUCKET = s3cfg.Bucket
	S3_HOST = s3cfg.Host
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

	var key string = fname

	res, err := s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(S3_BUCKET),
		Key:                  aws.String(key),
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

	url := S3_HOST + "/" + key

	return url, nil
}

func UnlinkImage(fname string) error {

	sess, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION)})

	if err != nil {
		return err
	}

	key := fname

	svc := s3.New(sess)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(S3_BUCKET),
		Key:    aws.String(key),
	}

	result, err := svc.DeleteObject(input)

	if err != nil {
		log.Println("ERROR: " + err.Error())
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println("AWS-ERR: " + aerr.Error())
				return err
			}
		} else {
			log.Println("ERR: " + err.Error())
			return err
		}
	}
	log.Println("S3-Delete: ", result)
	return nil
}

func ReadImage(fname string) (image.Image, string, error) {

	var img image.Image

	data, err := ReadImageBytes(fname)

	if err != nil {
		return img, "", err
	}

	img, format, err := image.Decode(bytes.NewReader(data))

	if err != nil {
		log.Println("ERROR - s3.ReadImage ", err)
		return img, format, err
	} else {
		return img, format, nil
	}
}

func ReadImageBytes(fname string) ([]byte, error) {

	sess, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION)})

	if err != nil {
		return []byte(""), err
	}

	key := fname

	svc := s3.New(sess)

	input := &s3.GetObjectInput{
		Bucket: aws.String(S3_BUCKET),
		Key:    aws.String(key),
	}

	result, err := svc.GetObject(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				log.Println("S3-ERROR NO SUCH KEY: ", s3.ErrCodeNoSuchKey, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println("S3-ERROR: ", err.Error())
		}
		return []byte(""), err
	}

	body, err := ioutil.ReadAll(result.Body)

	if err != nil {
		log.Println("Error reading result: ", err.Error())
		return []byte(""), err
	}

	return body, nil
}
