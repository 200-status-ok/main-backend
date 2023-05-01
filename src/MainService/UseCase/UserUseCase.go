package UseCase

import (
	"encoding/json"
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/MainService/Model"
	"github.com/403-access-denied/main-backend/src/MainService/Repository"
	"github.com/403-access-denied/main-backend/src/MainService/Token"
	Utils2 "github.com/403-access-denied/main-backend/src/MainService/Utils"
	"github.com/403-access-denied/main-backend/src/MainService/View"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
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
			Username:      verifyReq.Username,
			Posters:       nil,
			MarkedPosters: nil,
			ChatRooms:     nil,
			Conversations: nil,
		}
		err = userRepository.UserCreate(newUser)
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

func OAuth2LoginResponse(c *gin.Context) {
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
	response, err := Utils2.GetGoogleUserInfo(code, state, c)

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
	c.JSON(http.StatusOK, gin.H{"response": googleRes})
}
