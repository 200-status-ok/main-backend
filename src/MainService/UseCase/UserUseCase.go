package UseCase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/MainService/Model"
	"github.com/403-access-denied/main-backend/src/MainService/Repository"
	"github.com/403-access-denied/main-backend/src/MainService/Token"
	Utils2 "github.com/403-access-denied/main-backend/src/MainService/Utils"
	"github.com/403-access-denied/main-backend/src/MainService/View"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
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

	if appEnv == "development" {
		redisClient = Utils2.NewRedisClient(Utils2.ReadFromEnvFile(".env", "LOCAL_REDIS_HOST"),
			Utils2.ReadFromEnvFile(".env", "LOCAL_REDIS_PORT"),
			"", 0)
	} else if appEnv == "production" {
		redisClient = Utils2.NewRedisClient(Utils2.ReadFromEnvFile(".env", "PRODUCTION_REDIS_HOST"),
			Utils2.ReadFromEnvFile(".env", "PRODUCTION_REDIS_PORT"),
			Utils2.ReadFromEnvFile(".env", "PRODUCTION_REDIS_PASSWORD"), 0)
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
		err = messageBroker.ConnectBroker(Utils2.ReadFromEnvFile(".env", "RABBITMQ_LOCAL_CONNECTION"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else if appEnv == "production" {
		err = messageBroker.ConnectBroker(Utils2.ReadFromEnvFile(".env", "RABBITMQ_PROD_CONNECTION"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if Utils2.UsernameValidation(user.Username) == 0 {
		msg := "email/" + OTP + "/" + user.Username
		err = messageBroker.PublishOnQueue([]byte(msg), "email_notification")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		msg := "sms/" + OTP + "/" + user.Username
		err = messageBroker.PublishOnQueue([]byte(msg), "sms_notification")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	messageBroker.Close()
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
		redisClient = Utils2.NewRedisClient(Utils2.ReadFromEnvFile(".env", "LOCAL_REDIS_HOST"),
			Utils2.ReadFromEnvFile(".env", "LOCAL_REDIS_PORT"),
			"", 0)
	} else if appEnv == "production" {
		redisClient = Utils2.NewRedisClient(Utils2.ReadFromEnvFile(".env", "PRODUCTION_REDIS_HOST"),
			Utils2.ReadFromEnvFile(".env", "PRODUCTION_REDIS_PORT"),
			Utils2.ReadFromEnvFile(".env", "PRODUCTION_REDIS_PASSWORD"), 0)
	}

	userRepository := Repository.NewUserRepository(DBConfiguration.GetDB())
	if err := c.ShouldBindJSON(&verifyReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if Utils2.UsernameValidation(verifyReq.Username) == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}

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
			Wallet:              0,
			Payments:            nil,
		}
		_, err = userRepository.UserCreate(newUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	jwtMaker, err := Token.NewJWTMaker(Utils2.ReadFromEnvFile(".env", "JWT_SECRET"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, _, err := jwtMaker.MakeToken(verifyReq.Username, time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	View.LoginUserView(token, c)
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
	url := Utils2.GetGoogleAuthURL("random-state")
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
	response, err := Utils2.GetGoogleUserInfo(code, state)

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

	userRepository := Repository.NewUserRepository(DBConfiguration.GetDB())
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
			Wallet:              0,
			Payments:            nil,
		}
		_, err = userRepository.UserCreate(newUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	jwtMaker, err := Token.NewJWTMaker(Utils2.ReadFromEnvFile(".env", "JWT_SECRET"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, _, err := jwtMaker.MakeToken(googleRes.Email, time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, RedirectURI+"?token="+token)
}

type GetUserByIdRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

func GetUserByIdResponse(c *gin.Context) {
	var getUserReq GetUserByIdRequest
	userRepository := Repository.NewUserRepository(DBConfiguration.GetDB())
	if err := c.ShouldBindUri(&getUserReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := userRepository.FindById(getUserReq.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	View.GetUserByIdView(*user, c)
}

func GetUsersResponse(c *gin.Context) {
	userRepository := Repository.NewUserRepository(DBConfiguration.GetDB())
	users, err := userRepository.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	View.GetUsersView(*users, c)
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required,min=11,max=30"`
}
type UpdateUserByIdRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

func UpdateUserByIdResponse(c *gin.Context) {
	var updateUserReq UpdateUserRequest
	var updateUserByIdReq UpdateUserByIdRequest
	userRepository := Repository.NewUserRepository(DBConfiguration.GetDB())
	if err := c.ShouldBindUri(&updateUserByIdReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&updateUserReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if Utils2.UsernameValidation(updateUserReq.Username) == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}
	user, err := userRepository.UserUpdate(&Model.User{
		Username: updateUserReq.Username,
	}, updateUserByIdReq.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	View.GetUserByIdView(*user, c)
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=11,max=30"`
}

func CreateUserResponse(c *gin.Context) {
	var createUserReq CreateUserRequest
	userRepository := Repository.NewUserRepository(DBConfiguration.GetDB())
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

type DeleteUserByIdRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

func DeleteUserByIdResponse(c *gin.Context) {
	var deleteUserByIdReq DeleteUserByIdRequest
	userRepository := Repository.NewUserRepository(DBConfiguration.GetDB())
	if err := c.ShouldBindUri(&deleteUserByIdReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := userRepository.DeleteUser(deleteUserByIdReq.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": "user deleted"})
}

type PaymentRequest struct {
	Amount float64 `form:"amount" binding:"required,min=1"`
	Id     int     `form:"id" binding:"required,min=1"`
}
type data struct {
	Merchant    string  `json:"merchant"`
	CallbackURL string  `json:"callbackUrl"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	OrderID     string  `json:"orderId"`
}

func PaymentResponse(c *gin.Context) {
	var paymentReq PaymentRequest
	if err := c.ShouldBindQuery(&paymentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var merchant = "zibal"
	d := data{
		Merchant:    merchant,
		CallbackURL: "https://example.com/callback",
		Description: "golang package",
		Amount:      paymentReq.Amount,
	}
	fmt.Println(d.OrderID)
	jsonData, err := json.Marshal(d)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	var url = "https://gateway.zibal.ir/v1/request"
	// post request to zibal with gin

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Map result to a struct to easily access parameters
	var structResult map[string]interface{}
	br := bytes.NewReader([]byte(string(body)))
	decodedJson := json.NewDecoder(br)
	decodedJson.UseNumber()
	err = decodedJson.Decode(&structResult)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Access response parameters
	var resultNumber = structResult["result"]
	var trackId = structResult["trackId"]
	// Change response parameters types to string
	trackIdStringValue := fmt.Sprint(trackId)
	resultStringValue := fmt.Sprint(resultNumber)
	fmt.Println(trackIdStringValue)
	fmt.Println(resultStringValue)
	// check if result is 100
	if resultStringValue == "100" {
		// redirect to zibal
		PaymentRepository := Repository.NewPaymentRepository(DBConfiguration.GetDB())
		_, err := PaymentRepository.CreatePayment(Model.Payment{
			Amount:  paymentReq.Amount,
			UserID:  uint(paymentReq.Id),
			TrackID: trackIdStringValue,
			Status:  "pending",
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusFound, "https://gateway.zibal.ir/start/"+trackIdStringValue)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
		return
	}
}

type VerifyRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}
type dataVerify struct {
	Merchant string `json:"merchant"`
	TrackId  string `json:"trackId"`
}

func PaymentVerifyResponse(c *gin.Context) {
	var verifyReq VerifyRequest
	paymentRepository := Repository.NewPaymentRepository(DBConfiguration.GetDB())
	if err := c.ShouldBindUri(&verifyReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var merchant = "zibal"
	payment, err := paymentRepository.GetPaymentById(int(verifyReq.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d := dataVerify{
		Merchant: merchant,
		TrackId:  payment.TrackID,
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

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Map result to a struct to easily access parameters
	var structResult map[string]interface{}
	br := bytes.NewReader([]byte(string(body)))
	decodedJson := json.NewDecoder(br)
	decodedJson.UseNumber()
	err = decodedJson.Decode(&structResult)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Access response parameters
	var resultNumber = structResult["result"]
	var paidAt = structResult["paidAt"]
	var status = structResult["status"]
	// Change response parameters types to string
	resultStringValue := fmt.Sprint(resultNumber)
	paidAtStringValue := fmt.Sprint(paidAt)
	statusStringValue := fmt.Sprint(status)
	fmt.Println(resultStringValue)
	fmt.Println(paidAtStringValue)
	fmt.Println(statusStringValue)
	if resultStringValue == "100" {
		if statusStringValue == "1" {
			payment.Status = "paid"
			userRep := Repository.NewUserRepository(DBConfiguration.GetDB())
			user, err := userRep.UpdateWallet(payment.UserID, payment.Amount)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			fmt.Println(user)
			updatePayment, err := paymentRepository.UpdatePayment(payment.ID, payment)
			if err != nil {
				return
			}
			fmt.Println(updatePayment)
			c.JSON(http.StatusOK, gin.H{"message": "payment is successful"})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "payment is not successful"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment is not successful"})
		return
	}
}
