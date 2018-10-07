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

// LineMessageWebhook base line response with webhook
type LineMessageWebhook struct {
	ReplyToken string `json:"replyToken"`
	Type       string `json:"type"`
	Timestamp  int    `json:"timestamp"`
}

// LineWebHookEvent base line webhook event
type LineWebHookEvent struct {
	Events []LineMessageWebhook `json:"events"`
}

// ResultTypeReturn is base cover api
type ResultTypeReturn struct {
	BaseResultRegisterMail []BaseResultRegisterMail `json:"result"`
}

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
	CreateDate   string
	ReceivedDate string
	ProjectCode  string
	ProjectID    string
}

func main() {
	app := iris.Default()
	app.Post("/webhook", func(ctx iris.Context) {
		webhook := &LineWebHookEvent{}
		if err := ctx.ReadJSON(webhook); err != nil {
			panic(err.Error())
		} else {
			fmt.Println(webhook.Events[0].ReplyToken)
			applicationRunner()
		}
	})
	app.Run(iris.Addr(":8080"))
}

func applicationRunner(replyToken string, userId string) {
	fmt.Println(replyToken)
	var cachedRegisterMailAPI *BaseResultRegisterMail
	config := readAppConfig()
	for range time.Tick(30 * time.Second) {
		registerMailResult := getDataFromCheckRegisterMailAPI(config)
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
	var resultObject = new(ResultTypeReturn)
	err := utils.GetJSON(config.CheckRegisterMailAPIUrl, &resultObject)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err.Error())
	}
	return resultObject.BaseResultRegisterMail[0]
}

func findNewRegisterMailService(cachedResultRegisterMail BaseResultRegisterMail, currentResultRegisterMail BaseResultRegisterMail) []ResultRegisterMail {
	var filteredMailResult = []ResultRegisterMail{}
	for _, registerMail := range currentResultRegisterMail.Received {
		for _, cacheRegisterMail := range cachedResultRegisterMail.Received {
			if registerMail.ParcelNumber == cacheRegisterMail.ParcelNumber {
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
