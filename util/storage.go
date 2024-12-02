package util

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/BATUCHKA/real-estate-back/database/models"
)

func S3TempFileUpload(bucketModel models.Bucket, filePath string, contentType string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(bucketModel.S3Region),
		Credentials: credentials.NewStaticCredentials(bucketModel.S3AccessKey, bucketModel.S3SecretKey, ""),
		Endpoint:    aws.String(bucketModel.S3EndPoint),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return err
	}

	// fileBytes, err := io.ReadAll(file)
	// if err != nil {
	// 	err = errors.New("s3 file reading error" + err.Error())
	// 	return err
	// }
	// fileType := http.DetectContentType(fileBytes)

	fileName := filepath.Base(file.Name())
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:       aws.String(bucketModel.S3Bucket),
		Key:          aws.String(fileName),
		Body:         file,
		StorageClass: aws.String("STANDARD"), // Set the appropriate storage class here
		Metadata: map[string]*string{
			"Content-Type": aws.String(contentType),
		},
	})
	if err != nil {
		log.Printf("uploader error: %v", err)
		return err
	}

	return nil
}

func S3FileRemove(bucketModel models.Bucket, fileName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(bucketModel.S3Region),
		Credentials: credentials.NewStaticCredentials(bucketModel.S3AccessKey, bucketModel.S3SecretKey, ""),
		Endpoint:    aws.String(bucketModel.S3EndPoint),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	svc := s3.New(sess)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketModel.S3Bucket), Key: aws.String(fileName)})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucketModel.S3Bucket),
		Key:    aws.String(fileName),
	})
	return nil
}

func S3FileDownload(bucketModel models.Bucket, fileName string) (io.ReadCloser, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(bucketModel.S3Region),
		Credentials: credentials.NewStaticCredentials(bucketModel.S3AccessKey, bucketModel.S3SecretKey, ""),
		Endpoint:    aws.String(bucketModel.S3EndPoint),
	})

	if err != nil {
		return nil, err
	}

	s3Client := s3.New(sess)

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketModel.S3Bucket),
		Key:    aws.String(fileName),
	}

	resp, err := s3Client.GetObject(input)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
