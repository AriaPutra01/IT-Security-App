package initializers

import (
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var svc *s3.S3

func init() {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	svc = s3.New(sess)
}

func UploadToS3(file *multipart.FileHeader) error {
	fileContent, _ := file.Open()
	uploader := s3manager.NewUploader(session.Must(session.NewSession()))
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("your-bucket-name"),
		Key:    aws.String(file.Filename),
		Body:   fileContent,
	})
	return err
}

func DeleteFromS3(filename string) error {
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String("your-bucket-name"),
		Key:    aws.String(filename),
	})
	return err
}

func GeneratePresignedURL(filename string) (string, error) {
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("your-bucket-name"),
		Key:    aws.String(filename),
	})
	url, err := req.Presign(15 * time.Minute)
	return url, err
}

func ListFiles() ([]string, error) {
	var fileList []string
	result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("your-bucket-name"),
	})
	if err != nil {
		return nil, err
	}
	for _, item := range result.Contents {
		fileList = append(fileList, *item.Key)
	}
	return fileList, nil
}
