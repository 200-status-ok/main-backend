package UseCase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/Model"
	"github.com/200-status-ok/main-backend/src/MainService/Repository"
	"github.com/200-status-ok/main-backend/src/MainService/Repository/ElasticSearch"
	"github.com/200-status-ok/main-backend/src/MainService/Token"
	Utils2 "github.com/200-status-ok/main-backend/src/MainService/Utils"
	"github.com/200-status-ok/main-backend/src/MainService/View"
	"github.com/200-status-ok/main-backend/src/pkg/elasticsearch"
	"github.com/200-status-ok/main-backend/src/pkg/pgsql"
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func generateSecretKeyForNewUser(user string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "main-backend",
		AccountName: user,
	})

	return key.Secret(), err
}

type SendOTPRequest struct {
	Username string `json:"username" binding:"required,min=11,max=30"`
}

func SendOTPResponse(c *gin.Context) {
	var user SendOTPRequest
	var redisClient *Utils2.RedisClient
	appEnv := os.Getenv("APP_ENV2")

	fmt.Println(appEnv)
	if appEnv == "development" {
		redisClient = Utils2.NewRedisClient(utils.ReadFromEnvFile(".env", "LOCAL_REDIS_HOST"),
			utils.ReadFromEnvFile(".env", "LOCAL_REDIS_PORT"),
			"", 0)
	} else if appEnv == "production" {
		redisClient = Utils2.NewRedisClient(utils.ReadFromEnvFile(".env", "PRODUCTION_REDIS_HOST"),
			utils.ReadFromEnvFile(".env", "PRODUCTION_REDIS_PORT"),
			utils.ReadFromEnvFile(".env", "PRODUCTION_REDIS_PASSWORD"), 0)
	} else if appEnv == "testing" {
		redisClient = Utils2.NewRedisClient("redis", "6379", "", 0)
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if Utils2.UsernameValidation(user.Username) == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}

	secretKey, _ := generateSecretKeyForNewUser(user.Username)
	err := redisClient.Set(user.Username, secretKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	OTP, _ := totp.GenerateCodeCustom(secretKey, time.Now(), totp.ValidateOpts{
		Period:    120,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	messageBroker := Utils2.MessageClient{}
	if appEnv == "development" {
		err = messageBroker.ConnectBroker(utils.ReadFromEnvFile(".env", "RABBITMQ_LOCAL_CONNECTION"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else if appEnv == "production" {
		err = messageBroker.ConnectBroker(utils.ReadFromEnvFile(".env", "RABBITMQ_PROD_CONNECTION"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else if appEnv == "testing" {
		err = messageBroker.ConnectBroker("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if Utils2.UsernameValidation(user.Username) == 0 {
		msg := "email/login/" + OTP + "/" + user.Username
		err = messageBroker.PublishOnQueue([]byte(msg), "email_notification")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		msg := "sms/login/" + OTP + "/" + user.Username
		err = messageBroker.PublishOnQueue([]byte(msg), "sms_notification")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	messageBroker.Close()
	if appEnv == "testing" {
		c.JSON(http.StatusOK, gin.H{"OTP": OTP})
		return
	}
	View.LoginMessageView("OTP sent to registered email/phone", c)
}

type VerifyOTPRequest struct {
	Username string `json:"username" binding:"required,min=11,max=30"`
	OTP      string `json:"otp" binding:"required,len=6"`
}

func VerifyOtpResponse(c *gin.Context) {
	var verifyReq VerifyOTPRequest
	var redisClient *Utils2.RedisClient
	appEnv := os.Getenv("APP_ENV2")

	if appEnv == "development" {
		redisClient = Utils2.NewRedisClient(utils.ReadFromEnvFile(".env", "LOCAL_REDIS_HOST"),
			utils.ReadFromEnvFile(".env", "LOCAL_REDIS_PORT"),
			"", 0)
	} else if appEnv == "production" {
		redisClient = Utils2.NewRedisClient(utils.ReadFromEnvFile(".env", "PRODUCTION_REDIS_HOST"),
			utils.ReadFromEnvFile(".env", "PRODUCTION_REDIS_PORT"),
			utils.ReadFromEnvFile(".env", "PRODUCTION_REDIS_PASSWORD"), 0)
	} else if appEnv == "testing" {
		redisClient = Utils2.NewRedisClient("redis", "6379", "", 0)
	}

	userRepository := Repository.NewUserRepository(pgsql.GetDB())
	if err := c.ShouldBindJSON(&verifyReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if Utils2.UsernameValidation(verifyReq.Username) == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}
	var userId uint
	userExist, err := userRepository.FindByUsername(verifyReq.Username)
	if err != nil && err.Error() != "user not found" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	secretKey, err := redisClient.Get(verifyReq.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid, _ := totp.ValidateCustom(
		verifyReq.OTP,
		secretKey,
		time.Now().UTC(),
		totp.ValidateOpts{
			Period:    120,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		},
	)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid OTP"})
		return
	}

	if userExist == nil {
		newUser := &Model.User{
			Username:            verifyReq.Username,
			Posters:             nil,
			MarkedPosters:       nil,
			OwnConversations:    nil,
			MemberConversations: nil,
			Wallet:              0.0,
			Payments:            nil,
		}
		user, err := userRepository.UserCreate(newUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = userRepository.CommitChanges()
		if err != nil {
			return
		}
		userId = user.ID
	} else {
		userId = userExist.ID
	}

	var jwtSecret string
	if appEnv == "development" || appEnv == "production" {
		jwtSecret = utils.ReadFromEnvFile(".env", "JWT_SECRET")
	} else if appEnv == "testing" {
		jwtSecret = "qwertyuiopasdfghjklzxcvbnm123456"
	}

	jwtMaker, err := Token.NewJWTMaker(jwtSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, _, err := jwtMaker.MakeToken(verifyReq.Username, uint64(userId), "User", time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	View.LoginUserView(userId, token, c)
}

type GoogleLoginAndroidRequest struct {
	Email string `form:"email" binding:"required"`
}

func GoogleLoginAndroidResponse(c *gin.Context) {
	var googleLoginReq GoogleLoginAndroidRequest
	if err := c.ShouldBindQuery(&googleLoginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userRepository := Repository.NewUserRepository(pgsql.GetDB())
	userExist, err := userRepository.FindByUsername(googleLoginReq.Email)
	if err != nil && err.Error() != "user not found" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userExist == nil {
		user := &Model.User{
			Username:            googleLoginReq.Email,
			Posters:             nil,
			MarkedPosters:       nil,
			OwnConversations:    nil,
			MemberConversations: nil,
			Wallet:              0.0,
			Payments:            nil,
		}
		_, err := userRepository.UserCreate(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = userRepository.CommitChanges()
		if err != nil {
			fmt.Println(err)
			return
		}
		userExist, err = userRepository.FindByUsername(googleLoginReq.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	var jwtSecret string
	appEnv := os.Getenv("APP_ENV2")
	if appEnv == "development" || appEnv == "production" {
		jwtSecret = utils.ReadFromEnvFile(".env", "JWT_SECRET")
	} else if appEnv == "testing" {
		jwtSecret = "qwertyuiopasdfghjklzxcvbnm123456"
	}

	jwtMaker, err := Token.NewJWTMaker(jwtSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, _, err := jwtMaker.MakeToken(googleLoginReq.Email, uint64(userExist.ID), "User", time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	View.LoginUserView(userExist.ID, token, c)
}

type GoogleLoginRequest struct {
	RedirectURI string `form:"redirect_uri" binding:"required"`
}

var RedirectURI string

func OAuth2LoginResponse(c *gin.Context) {
	var googleLoginReq GoogleLoginRequest
	if err := c.ShouldBindQuery(&googleLoginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	RedirectURI = googleLoginReq.RedirectURI
	url := GetGoogleAuthURL("random-state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

type GoogleCallbackRes struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func GoogleCallbackResponse(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	fmt.Println(RedirectURI)
	response, err := GetGoogleUserInfo(code, state)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var googleRes GoogleCallbackRes
	err = json.Unmarshal(response, &googleRes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !googleRes.VerifiedEmail {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email not verified"})
		return
	}

	userRepository := Repository.NewUserRepository(pgsql.GetDB())
	var userID uint64
	userExist, err := userRepository.FindByUsername(googleRes.Email)
	if err != nil && err.Error() != "user not found" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userExist == nil {
		newUser := &Model.User{
			Username:            googleRes.Email,
			Posters:             nil,
			MarkedPosters:       nil,
			OwnConversations:    nil,
			MemberConversations: nil,
			Wallet:              0.0,
			Payments:            nil,
		}
		createdUser, err := userRepository.UserCreate(newUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID = uint64(createdUser.ID)
	} else {
		userID = uint64(userExist.ID)
	}

	jwtMaker, err := Token.NewJWTMaker(utils.ReadFromEnvFile(".env", "JWT_SECRET"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, _, err := jwtMaker.MakeToken(googleRes.Email, userID, "User", time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, RedirectURI+"?token="+token)
}

func GetUserByIdResponse(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*Token.Payload)
	userRepository := Repository.NewUserRepository(pgsql.GetDB())
	user, err := userRepository.FindById(uint(payload.UserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	View.GetUserByIdView(*user, c)
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required,min=11,max=50"`
}

func UpdateUserByIdResponse(c *gin.Context) {
	var updateUserReq UpdateUserRequest
	payload := c.MustGet("authorization_payload").(*Token.Payload)
	userRepository := Repository.NewUserRepository(pgsql.GetDB())

	if err := c.ShouldBindJSON(&updateUserReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if Utils2.UsernameValidation(updateUserReq.Username) == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}
	updateUser, err := userRepository.UserUpdate(&Model.User{
		Username: updateUserReq.Username,
	}, uint(payload.UserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = userRepository.CommitChanges()
	if err != nil {
		fmt.Println(err)
		return
	}
	View.GetUserByIdView(*updateUser, c)
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=11,max=30"`
}

func CreateUserResponse(c *gin.Context) {
	var createUserReq CreateUserRequest
	userRepository := Repository.NewUserRepository(pgsql.GetDB())
	if err := c.ShouldBindJSON(&createUserReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if Utils2.UsernameValidation(createUserReq.Username) == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}
	user, err := userRepository.UserCreate(&Model.User{
		Username: createUserReq.Username,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	View.GetUserByIdView(*user, c)
}

func DeleteUserByIdResponse(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*Token.Payload)
	userRepository := Repository.NewUserRepository(pgsql.GetDB())
	esDeletePostersByUserId := ElasticSearch.NewPosterES(elasticsearch.GetElastic())
	err := userRepository.DeleteUser(uint(payload.UserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if os.Getenv("APP_ENV2") != "testing" {
		err = esDeletePostersByUserId.DeletePosterByUserID(int(payload.UserID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = userRepository.CommitChanges()
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		err := userRepository.RoleBackChanges()
		if err != nil {
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"response": "user deleted"})
}

type PaymentRequest struct {
	Amount float64 `form:"amount" binding:"required,min=1"`
	Url    string  `form:"url" binding:"required"`
}
type data struct {
	Merchant    string  `json:"merchant"`
	CallbackURL string  `json:"callbackUrl"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

func PaymentResponse(c *gin.Context) {
	var paymentReq PaymentRequest
	payload := c.MustGet("authorization_payload").(*Token.Payload)
	if err := c.ShouldBindQuery(&paymentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var merchant = "zibal"
	d := data{
		Merchant:    merchant,
		CallbackURL: paymentReq.Url,
		Description: "payment",
		Amount:      paymentReq.Amount,
	}
	jsonData, err := json.Marshal(d)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	var url = "https://gateway.zibal.ir/v1/request"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var structResult map[string]interface{}
	br := bytes.NewReader(body)
	decodedJson := json.NewDecoder(br)
	decodedJson.UseNumber()
	err = decodedJson.Decode(&structResult)
	if err != nil {
		fmt.Println(err)
		return
	}
	var resultNumber = structResult["result"]
	var trackId = structResult["trackId"]
	trackIdStringValue := fmt.Sprint(trackId)
	resultStringValue := fmt.Sprint(resultNumber)
	if resultStringValue == "100" {
		PaymentRepository := Repository.NewPaymentRepository(pgsql.GetDB())
		_, err := PaymentRepository.CreatePayment(Model.Payment{
			Amount:  paymentReq.Amount,
			UserID:  uint(payload.UserID),
			TrackID: trackIdStringValue,
			Status:  "pending",
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"trackID": trackIdStringValue,
			"redirect": "https://gateway.zibal.ir/start/" + trackIdStringValue})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
		return
	}
}

type PaymentVerifyRequest struct {
	TrackId string `form:"track_id" binding:"required"`
}
type dataVerify struct {
	Merchant string `json:"merchant"`
	TrackId  string `json:"trackId"`
}

func PaymentVerifyResponse(c *gin.Context) {
	var paymentVerifyReq PaymentVerifyRequest
	if err := c.ShouldBindQuery(&paymentVerifyReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	paymentRepository := Repository.NewPaymentRepository(pgsql.GetDB())
	var merchant = "zibal"
	payment, err := paymentRepository.GetPaymentByTrackID(paymentVerifyReq.TrackId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d := dataVerify{
		Merchant: merchant,
		TrackId:  paymentVerifyReq.TrackId,
	}
	var url = "https://gateway.zibal.ir/v1/verify"
	jsonData, err := json.Marshal(d)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var structResult map[string]interface{}
	br := bytes.NewReader(body)
	decodedJson := json.NewDecoder(br)
	decodedJson.UseNumber()
	err = decodedJson.Decode(&structResult)
	if err != nil {
		fmt.Println(err)
		return
	}
	var resultNumber = structResult["result"]
	var status = structResult["status"]
	resultStringValue := fmt.Sprint(resultNumber)
	statusStringValue := fmt.Sprint(status)
	if resultStringValue == "100" {
		if statusStringValue == "1" {
			payment.Status = "paid"
			userRep := Repository.NewUserRepository(pgsql.GetDB())
			_, err := userRep.UpdateWallet(payment.UserID, payment.Amount)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			updatePayment, err := paymentRepository.UpdatePayment(payment.ID, payment)
			if err != nil {
				return
			}
			fmt.Println(updatePayment)
			c.JSON(http.StatusOK, gin.H{"message": "payment is successful"})
			return
		} else {
			payment.Status = "failed"
			_, err := paymentRepository.UpdatePayment(payment.ID, payment)
			if err != nil {
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "payment is not successful"})
			return
		}
	} else {
		payment.Status = "failed"
		_, err := paymentRepository.UpdatePayment(payment.ID, payment)
		if err != nil {
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment is not successful"})
		return
	}
}

func GetTransactionsResponse(c *gin.Context) {
	payload := c.MustGet("authorization_payload").(*Token.Payload)
	paymentRepository := Repository.NewPaymentRepository(pgsql.GetDB())
	payments, err := paymentRepository.GetPaymentsOfUser(uint(payload.UserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	View.GetUserPaymentsView(payments, c)
	return
}

var googleClientId = utils.ReadFromEnvFile(".env", "GOOGLE_CLIENT_ID")
var googleClientSecret = utils.ReadFromEnvFile(".env", "GOOGLE_CLIENT_SECRET")

func GetGoogleOauthConfig() *oauth2.Config {
	redirectGoogleUri := ""
	if os.Getenv("APP_ENV2") == "production" {
		redirectGoogleUri = utils.ReadFromEnvFile(".env", "PRODUCTION_REDIRECT_GOOGLE_URL")
	} else if os.Getenv("APP_ENV2") == "development" {
		redirectGoogleUri = utils.ReadFromEnvFile(".env", "LOCAL_REDIRECT_GOOGLE_URL")
	}
	var googleOauthConfig = &oauth2.Config{
		RedirectURL:  redirectGoogleUri,
		ClientID:     googleClientId,
		ClientSecret: googleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	return googleOauthConfig
}

func GetGoogleAuthURL(state string) string {
	return GetGoogleOauthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func GetGoogleUserInfo(code string, state string) ([]byte, error) {
	if state != "random-state" {
		return nil, errors.New("invalid oauth state")
	}
	token, err := GetGoogleOauthConfig().Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New("code exchange wrong: " + err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, errors.New("failed getting user info: " + err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("failed reading response body: " + err.Error())
	}
	return contents, nil
}
