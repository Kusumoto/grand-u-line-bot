package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kataras/iris"
	"github.com/kusumoto/grand-u-line-notify/config"
	"github.com/kusumoto/grand-u-line-notify/utils"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/spf13/viper"
)

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

var userID string

func main() {
	go startChecker()
	app := iris.Default()
	app.Post("/webhook", func(ctx iris.Context) {
		request := ctx.Request()
		initApplication(request)
		ctx.JSON(iris.StatusOK)
	})
	app.Run(iris.Addr(":8080"))
}

func initApplication(request *http.Request) {
	config := readAppConfig()
	bot, err := linebot.New(config.LineChannelSecret, config.LineAccessToken)
	if err != nil {
		fmt.Println(1)
		log.Fatal(err)
		fmt.Println(err.Error())
	}
	sendHelloToLine(bot, request)
}

func startChecker() {
	var cachedRegisterMailAPI BaseResultRegisterMail
	config := readAppConfig()
	bot, err := linebot.New(config.LineChannelSecret, config.LineAccessToken)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err.Error())
	}

	for range time.Tick(2 * time.Second) {
		if userID != "" {
			registerMailResult := getDataFromCheckRegisterMailAPI(config)
			filteredMailResult := findNewRegisterMailService(cachedRegisterMailAPI, registerMailResult)
			sendMessageToLine(filteredMailResult, bot)
			cachedRegisterMailAPI = registerMailResult
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

func buildFlexMessage(registerMailObject ResultRegisterMail) *linebot.FlexMessage {
	flexHeaderComponent := []linebot.FlexComponent{}
	rootFlexBodyComponent := []linebot.FlexComponent{}

	flexHeaderComponent = append(flexHeaderComponent, &linebot.TextComponent{Type: linebot.FlexComponentTypeText, Text: "มีพัสดุใหม่", Size: linebot.FlexTextSizeTypeXl, Weight: linebot.FlexTextWeightTypeBold, Color: "#1DB446"})
	flexHeaderComponent = append(flexHeaderComponent, &linebot.SeparatorComponent{Type: linebot.FlexComponentTypeSeparator, Margin: linebot.FlexComponentMarginTypeXl})

	rootFlexBodyComponent = append(rootFlexBodyComponent, buildChildBox("รายการพัสดุ", registerMailObject.Title))
	rootFlexBodyComponent = append(rootFlexBodyComponent, buildChildBox("ส่งถึง", registerMailObject.SentTo))
	rootFlexBodyComponent = append(rootFlexBodyComponent, buildChildBox("Tracking No.", registerMailObject.TrackNo))
	rootFlexBodyComponent = append(rootFlexBodyComponent, buildChildBox("ผู้ส่ง", registerMailObject.Sender))
	rootFlexBodyComponent = append(rootFlexBodyComponent, buildChildBox("ผู้รับพัสดุ", registerMailObject.Recipient))
	rootFlexBodyComponent = append(rootFlexBodyComponent, buildChildBox("วันที่รับพัสดุเข้า", registerMailObject.CreateDate))

	headerBox := &linebot.BoxComponent{
		Type:     linebot.FlexComponentTypeBox,
		Layout:   linebot.FlexBoxLayoutTypeVertical,
		Contents: flexHeaderComponent,
	}
	bodyBox := &linebot.BoxComponent{
		Type:     linebot.FlexComponentTypeBox,
		Layout:   linebot.FlexBoxLayoutTypeVertical,
		Contents: rootFlexBodyComponent,
	}
	flexContainerTemplate := &linebot.BubbleContainer{
		Type:   linebot.FlexContainerTypeBubble,
		Header: headerBox,
		Body:   bodyBox,
	}
	template := linebot.NewFlexMessage("การแจ้งเตือนพัสดุ", flexContainerTemplate)
	return template
}

func buildChildBox(title string, value string) linebot.FlexComponent {
	childFlexBodyContentComponent := []linebot.FlexComponent{}
	childFlexBodyContentComponent = append(childFlexBodyContentComponent, &linebot.TextComponent{Type: linebot.FlexComponentTypeText, Text: title, Align: linebot.FlexComponentAlignTypeStart})
	childFlexBodyContentComponent = append(childFlexBodyContentComponent, &linebot.TextComponent{Type: linebot.FlexComponentTypeText, Text: value, Align: linebot.FlexComponentAlignTypeEnd})

	childFlexBodyComponent := &linebot.BoxComponent{
		Type:     linebot.FlexComponentTypeBox,
		Layout:   linebot.FlexBoxLayoutTypeHorizontal,
		Spacing:  linebot.FlexComponentSpacingTypeMd,
		Contents: childFlexBodyContentComponent,
	}
	return childFlexBodyComponent
}

func findNewRegisterMailService(cachedResultRegisterMail BaseResultRegisterMail, currentResultRegisterMail BaseResultRegisterMail) []ResultRegisterMail {
	var filteredMailResult = []ResultRegisterMail{}
	isCached := false
	for _, registerMail := range currentResultRegisterMail.NotReceived {
		for _, cacheRegisterMail := range cachedResultRegisterMail.NotReceived {
			if registerMail.ParcelNumber == cacheRegisterMail.ParcelNumber {
				isCached = true
				break
			}
		}
		if !isCached {
			filteredMailResult = append(filteredMailResult, registerMail)
		}
		isCached = false
	}
	return filteredMailResult
}

func sendHelloToLine(botClient *linebot.Client, request *http.Request) {
	messages := linebot.NewTextMessage("ยินดีต้อนรับจ้า ระบบเริ่มทำงานแล้ว!")

	events, err := botClient.ParseRequest(request)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err.Error())
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			replyToken := event.ReplyToken
			userID = event.Source.UserID
			_, err := botClient.ReplyMessage(replyToken, messages).Do()
			if err != nil {
				log.Fatal(err)
				fmt.Println(err.Error())
			}
		}
	}
}

func sendMessageToLine(resultRegisterMail []ResultRegisterMail, botClient *linebot.Client) {
	for _, mailResult := range resultRegisterMail {
		_, err := botClient.PushMessage(userID, buildFlexMessage(mailResult)).Do()
		if err != nil {
			log.Fatal(err)
			fmt.Println(err.Error())
		}
	}
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
