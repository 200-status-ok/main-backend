package MessageCli

import (
	"errors"
	"fmt"
	"github.com/403-access-denied/main-backend/src/WorkerService/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/WorkerService/Repository"
	"github.com/403-access-denied/main-backend/src/WorkerService/Repository/ElasticSearch"
	"github.com/403-access-denied/main-backend/src/WorkerService/Utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MessageClient struct {
	Connection *amqp.Connection
}

func (client *MessageClient) ConnectBroker(connectionString string) error {
	if connectionString == "" {
		return errors.New("connectionString is empty")
	}
	var err error
	client.Connection, err = amqp.DialConfig(connectionString, amqp.Config{
		Heartbeat: 90 * time.Second,
	})
	if err != nil {
		return err
	}
	return nil
}

func (client *MessageClient) Subscribe(exchangeName string, exchangeType string, consumerName string) error {
	channel, err := client.Connection.Channel()
	if err != nil {
		return err
	}
	err = channel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = channel.QueueBind(
		queue.Name,
		exchangeName,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = channel.Qos(1, 0, false)
	if err != nil {
		return err
	}
	msgs, err := channel.Consume(
		queue.Name,
		consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	// TODO: refactor this
	var forever chan struct{}
	go func() {
		for d := range msgs {
			fmt.Println("Received a message: ", string(d.Body))
		}
	}()
	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

type CustomArray []string

func (client *MessageClient) SubscribeOnQueue(queueName string, consumerName string, db *gorm.DB) error {
	if client.Connection == nil {
		return errors.New("connection is nil")
	}
	channel, err := client.Connection.Channel()
	if err != nil {
		return err
	}
	queue, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = channel.Qos(1, 0, false)
	if err != nil {
		return err
	}
	message, err := channel.Consume(
		queue.Name,
		consumerName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	var forever = make(chan struct{})

	for d := range message {
		fmt.Println("Received a message: ", string(d.Body))
		if queue.Name == "email_notification" || queue.Name == "sms_notification" {
			arr := CustomArray{}
			arr = strings.Split(string(d.Body), "/")
			arr.SendingNotification()
		} else if queue.Name == "nsfw-validation" {
			posterID, _ := strconv.ParseUint(string(d.Body), 10, 64)
			PhotoTextValidation(posterID, db)
		}
	}
	<-forever
	return nil
}

func (client *MessageClient) Close() {
	if client.Connection != nil {
		client.Connection.Close()
	}
}

func (a CustomArray) SendingNotification() {
	if a[0] == "email" {
		emailService := Utils.NewEmail("mhmdrzsmip@gmail.com", a[2],
			"Sending OTP code", "کد تایید ورود به سامانه همینجا: "+a[1],
			Utils.ReadFromEnvFile(".env", "GOOGLE_SECRET"))
		err := emailService.SendEmailWithGoogle()
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if a[0] == "sms" {
		pattern := map[string]string{
			"code": a[1],
		}
		otpSms := Utils.NewSMS(Utils.ReadFromEnvFile(".env", "API_KEY"), pattern)
		err := otpSms.SendSMSWithPattern(a[2], Utils.ReadFromEnvFile(".env", "OTP_PATTERN_CODE"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func PhotoTextValidation(posterID uint64, db *gorm.DB) {
	posterRepository := Repository.NewPosterRepository(db)
	var photoUrls []string
	var posterTexts *Repository.PosterResult
	isBadPhoto := false
	isBadText := false

	var wg sync.WaitGroup
	wg.Add(2)
	photoUrls, posterTexts, err := posterRepository.GetImagesTextsPosterByID(uint(posterID))
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for _, imgUrl := range photoUrls {
			url := fmt.Sprintf("https://api.apilayer.com/nudity_detection/url?url=%s", imgUrl)
			client := &http.Client{}
			req, err := http.NewRequest("GET", url, nil)
			req.Header.Set("apikey", "z232GHVwPAec88LzsqdBjUhL5BZDgvGp")

			if err != nil {
				fmt.Println(err)
			}
			res, err := client.Do(req)
			if res.Body != nil {
				defer res.Body.Close()
			}
			body, err := ioutil.ReadAll(res.Body)
			splitStr := strings.Split(string(body), ",")
			splitStr2 := strings.Split(splitStr[0], ": ")
			adultResult, err := strconv.Atoi(splitStr2[1])
			if adultResult > 1 {
				isBadPhoto = true
				break
			} else {
				isBadPhoto = false
			}
		}
		wg.Done()
	}()

	go func() {
		file, err := ioutil.ReadFile("Utils/data.txt")
		if err != nil {
			fmt.Println(err)
			return
		}

		newString := strings.ReplaceAll(string(file), "\n", "")
		newString2 := strings.ReplaceAll(newString, "\"", "")
		newString3 := strings.ReplaceAll(newString2, "\r", "")
		splitStr := strings.Split(newString3, ",")

		splitTitle := strings.Split(posterTexts.Title, " ")
		splitDescription := strings.Split(posterTexts.Description, " ")

		concatenated := append(splitTitle, splitDescription...)

		for _, v := range concatenated {
			for j, _ := range splitStr {
				if v == splitStr[j] {
					isBadText = true
					break
				}
			}
			if isBadText {
				break
			}
		}
		wg.Done()
	}()

	wg.Wait()
	state := ""

	if !isBadPhoto && !isBadText {
		state = "accepted"
	} else {
		state = "rejected"
	}

	err = posterRepository.UpdatePosterState(uint(posterID), state)
	if err != nil {
		fmt.Println(err)
		return
	}

	esPosterCli := ElasticSearch.NewPosterES(DBConfiguration.GetElastic())
	err = esPosterCli.UpdatePosterState(state, int(posterID))
	if err != nil {
		fmt.Println(err)
		return
	}
}
