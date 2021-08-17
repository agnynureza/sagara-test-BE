package client

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go"
)

func InitCloudinary() (*cloudinary.Cloudinary, error) {
	cloudName := os.Getenv("CLOUD_NAME")
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	var cld, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to intialize cloudinary, %v", err)
	}
	log.Println("Cloudinary Connect 🚀")
	return cld, nil
}
