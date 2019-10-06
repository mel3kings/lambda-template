package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/mel3kings/lambda-template/lambdahandler"
	"github.com/mel3kings/lambda-template/router"
	"log"
	"os"
)

func main() {
	environment := loadEnv()
	if environment == "develop" {
		router.NewRouter()
		select {}
	} else {
		lambda.Start(lambdahandler.HandleRequest)
	}
}

func loadEnv() string{
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "develop"
		_ = os.Setenv("ENV", "develop")
	}

	err := godotenv.Load("config/" + environment + ".env")
	if err != nil {
		log.Fatalf("Error loading .env file; Err: %v", err)
	}
	return environment
}
