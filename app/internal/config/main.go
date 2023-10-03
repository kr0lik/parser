package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"parser/internal/infrastructure/mongodb"
	"parser/internal/infrastructure/selenium"
)

var main *Main

type Main struct {
	MongoDbHost     string
	MongoDbUsername string
	MongoDbPassword string
	SeleniumHost    string
}

func ReadMain() {
	dotEnvPath := "../.env"

	if err := godotenv.Load(dotEnvPath); err != nil {
		log.Print("Error loading .env file")
	}

	mongoDbHost, exists := os.LookupEnv("MONGODB_HOST")
	if !exists {
		panic(fmt.Errorf("MONGODB_HOST not exists"))
	}

	mongoDbUsername, exists := os.LookupEnv("MONGODB_USERNAME")
	if !exists {
		panic(fmt.Errorf("MONGODB_USERNAME not exists"))
	}

	mongoDbPassword, exists := os.LookupEnv("MONGODB_PASSWORD")
	if !exists {
		panic(fmt.Errorf("MONGODB_PASSWORD not exists"))
	}

	seleniumHost, exists := os.LookupEnv("SELENIUM_HOST")
	if !exists {
		panic(fmt.Errorf("SELENIUM_HOST not exists"))
	}

	main = &Main{
		MongoDbHost:     mongoDbHost,
		MongoDbUsername: mongoDbUsername,
		MongoDbPassword: mongoDbPassword,
		SeleniumHost:    seleniumHost,
	}
}

func ProvideMongodbClientOptions() *mongodb.ClientOptions {
	return &mongodb.ClientOptions{
		Host: main.MongoDbHost,
		User: main.MongoDbUsername,
		Pass: main.MongoDbPassword,
	}
}

func ProvideSeleniumServerOptions() *selenium.ServerOptions {
	return &selenium.ServerOptions{
		Host: main.SeleniumHost,
	}
}
