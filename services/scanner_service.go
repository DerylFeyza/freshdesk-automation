package services

import (
	"fmt"
	"os"

	"resty.dev/v3"
)

func AttachmentToText(attachment string) ([]byte, error) {
	api := os.Getenv("SCANNER_BASE_URL")

	client := resty.New()
	response, err := client.R().
		SetFormData(map[string]string{"url": attachment}).
		Post(fmt.Sprintf("%s/process", api))

	if err != nil {
		return nil, err
	}

	return response.Bytes(), nil
}
