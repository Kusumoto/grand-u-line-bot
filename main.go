package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kataras/iris"
	"github.com/kusumoto/grand-u-line-notify/config"
	"github.com/kusumoto/grand-u-line-notify/utils"
	"github.com/spf13/viper"
)

// BaseResultRegisterMail model for base response from api register mail checker
type BaseResultRegisterMail struct {
	NotReceived []ResultRegisterMail
	Received    []ResultRegisterMail
	SendBack    []ResultRegisterMail
}

// ResultRegisterMail model for response from api resgiter mail checker
type ResultRegisterMail struct {
	Mobile       string
	UnitNumber   string
	Address      string
	TrackNo      string
	ParcelNumber string
	Sender       string
	SentTo       string
	Recipient    string
	Dispenser    string
	Title        string
	Status       string
	CreateDate   time.Time
	ReceivedDate time.Time
	ProjectCode  string
	ProjectID    string
}

func main() {

	app := iris.Default()
	app.Post("/webhook", func(ctx iris.Context) {
		fmt.Println(ctx.Request().Body)
		applicationRunner()
	})
	app.Run(iris.Addr(":8080"))
}

func applicationRunner() {
	var cachedRegisterMailAPI *BaseResultRegisterMail
	config := readAppConfig()
	for range time.Tick(30 * time.Second) {
		registerMailResult := getDataFromCheckRegisterMailAPI(config)
		fmt.Println(len(registerMailResult.Received))
		if cachedRegisterMailAPI == nil {
			cachedRegisterMailAPI = &registerMailResult
		} else {
			filteredMailResult := findNewRegisterMailService(*cachedRegisterMailAPI, registerMailResult)
			sendMessageToLine(filteredMailResult)
			cachedRegisterMailAPI = &registerMailResult
		}
	}
}

func getDataFromCheckRegisterMailAPI(config config.Config) BaseResultRegisterMail {
	var resultObject = new(BaseResultRegisterMail)
	err := utils.GetJSON(config.CheckRegisterMailAPIUrl, &resultObject)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err.Error())
	}
	return *resultObject
}

func findNewRegisterMailService(cachedResultRegisterMail BaseResultRegisterMail, currentResultRegisterMail BaseResultRegisterMail) []ResultRegisterMail {
	var filteredMailResult = []ResultRegisterMail{}
	for _, registerMail := range currentResultRegisterMail.Received {
		for _, cacheRegisterMail := range cachedResultRegisterMail.Received {
			if registerMail.ParcelNumber == cacheRegisterMail.ParcelNumber {
				fmt.Println(registerMail.ParcelNumber)
				continue
			}
		}
		filteredMailResult = append(filteredMailResult, registerMail)
	}
	return filteredMailResult
}

func sendMessageToLine(resultRegisterMail []ResultRegisterMail) {

}

func readAppConfig() config.Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	viper.SetDefault("http_timeout", 10)
	viper.SetDefault("delay", 5)

	var config config.Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return config
}
