package helper

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

func UploadFilesToGCS(c echo.Context, fileHeader *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	bucketName := os.Getenv("BUCKET_NAME")
	key := os.Getenv("BUCKET_SA")

	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Fatal("Can't decode service account key")
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(keyBytes))
	if err != nil {
		log.Fatal("Can't connect to Google Cloud Storage")
	}

	currentTime := time.Now().UTC()

	year := currentTime.Year()
	month := int(currentTime.Month())
	day := currentTime.Day()
	hour := currentTime.Hour()
	minute := currentTime.Minute()
	second := currentTime.Second()

	formattedTime := fmt.Sprintf("%04d%02d%02d-%d%02d%02d", year, month, day, second, minute, hour)

	filePath := formattedTime + "-" + fileHeader.Filename

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	obj := client.Bucket(bucketName).Object(filePath)

	// create a writer for the object
	wc := obj.NewWriter(ctx)

	// upload
	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	// generate URL
	objectAttrs, err := obj.Attrs(ctx)
	if err != nil {
		return "", err
	}

	URL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectAttrs.Name)

	return URL, nil
}

func DeleteFilesFromGCS(filename string) error {
	ctx := context.Background()

	key, err := base64.StdEncoding.DecodeString(os.Getenv("BUCKET_SA"))
	if err != nil {
		return err
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(key))
	if err != nil {
		return err
	}
	defer client.Close()

	bucket := client.Bucket(os.Getenv("BUCKET_NAME"))
	object := bucket.Object(filename)

	if err := object.Delete(ctx); err != nil {
		return err
	}

	return nil
}
