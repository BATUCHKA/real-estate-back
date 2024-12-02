package util

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"gitlab.com/steppelink/odin/odin-backend/database/models"
	"golang.org/x/image/webp"
)

func ProfileImageUpload(w http.ResponseWriter, r *http.Request) (hexName string, fileSize uint, mimeType string, err error) {
	isS3 := os.Getenv("MEDIA_S3")

	if err := r.ParseMultipartForm(32 << 20); err != nil {

		return "", 0, "", err
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error: ", err)
		return "", 0, "", err
	}
	defer file.Close()

	fileExtension := filepath.Ext(header.Filename)
	var img image.Image

	if fileExtension != ".svg" {
		fmt.Println("fileExtension: ", fileExtension)
		img, _, err = image.Decode(file)
		if err != nil {
			img, err = webp.Decode(file)
			if err != nil {
				err = errors.New("Unsupported image type & " + err.Error())
				return "", 0, "", err
			}
		}
	}

	var resizedImage *image.NRGBA
	var buf bytes.Buffer
	var fileByte int

	if fileExtension == ".svg" {
		fileSize = uint(header.Size)
	} else {
		resizedImage = imaging.Resize(img, 0, 400, imaging.Lanczos)
		if err = jpeg.Encode(&buf, resizedImage, nil); err != nil {
			return "", 0, "", err
		}
		fileByte = len(buf.Bytes())
		fileSize = uint(fileByte)
	}

	imageFile := make([]byte, 32)
	_, err = rand.Read(imageFile)
	if err != nil {
		return "", 0, "", err
	}

	hexName = hex.EncodeToString(imageFile)

	var folderRoute string
	if isS3 == "1" {
		folderRoute = "data/temp/" + hexName[0:2] + "/" + hexName[2:4]
	} else {
		folderRoute = "data/" + hexName[0:2] + "/" + hexName[2:4]
	}
	if err := os.MkdirAll(folderRoute, os.FileMode(0777)); err != nil {
		return "", 0, "", err
	}

	fileRoute, err := os.OpenFile(folderRoute+"/"+hexName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", 0, "", err
	}
	defer fileRoute.Close()

	if fileExtension == ".svg" {
		if _, err := io.Copy(fileRoute, file); err != nil {
			return "", 0, "", err
		}
	} else {
		if _, err := io.Copy(fileRoute, &buf); err != nil {
			return "", 0, "", err
		}
	}

	fileOpen, err := os.Open(fileRoute.Name())
	if err != nil {
		return "", 0, "", err
	}
	defer fileOpen.Close()

	fileBytes, err := os.ReadFile(fileRoute.Name())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fileType := http.DetectContentType(fileBytes)

	if isS3 == "1" {
		s3BucketModel := models.Bucket{
			S3Region:    os.Getenv("S3_REGION"),
			S3Bucket:    os.Getenv("S3_BUCKET_NAME"),
			S3EndPoint:  os.Getenv("S3_END_POINT"),
			S3AccessKey: os.Getenv("S3_ACCESS_KEY"),
			S3SecretKey: os.Getenv("S3_SECRET_KEY"),
		}
		if err := S3TempFileUpload(s3BucketModel, fileRoute.Name(), header.Header.Get("Content-Type")); err != nil {
			return "", 0, "", err
		}
		dirPath := filepath.Dir(filepath.Dir(fileRoute.Name()))
		if err := os.RemoveAll(dirPath); err != nil {
			return "", 0, "", err
		}
	}

	return hexName, fileSize, fileType, nil
}
