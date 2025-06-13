package config

import "github.com/joho/godotenv"
import "github.com/sirupsen/logrus"

func LoadConfig() {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("error loading .env file^ %v", err)
	}
}
