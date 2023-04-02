package UseCase

import (
	"encoding/json"
	"github.com/403-access-denied/main-backend/src/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/Model"
	"github.com/403-access-denied/main-backend/src/Repository"
	"github.com/403-access-denied/main-backend/src/Token"
	"github.com/403-access-denied/main-backend/src/Utils"
	"github.com/403-access-denied/main-backend/src/View"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"net/http"
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
	userRepository := Repository.NewUserRepository(DBConfiguration.GetDB())
	var user SendOTPRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if Utils.UsernameValidation(user.Username) == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}

	userExist, err := userRepository.FindByUsername(user.Username)
	if err != nil && err.Error() != "user not found" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userExist == nil {
		secretKey, _ := generateSecretKeyForNewUser(user.Username)
		newUser := &Model.User{
			Username:      user.Username,
			SecretKey:     secretKey,
			Posters:       nil,
			MarkedPosters: nil,
			Conversations: nil,
		}
		err := userRepository.UserCreate(newUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		otp, _ := totp.GenerateCode(secretKey, time.Now())
		if Utils.UsernameValidation(user.Username) == 0 {
			err = Utils.SendEmail(user.Username, otp)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		} else {
			err = Utils.SendOTP(user.Username, otp)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
		View.LoginMessageView("OTP sent to registered email/phone", c)
		return
	}
	otp, _ := totp.GenerateCode(userExist.SecretKey, time.Now())
	if Utils.UsernameValidation(user.Username) == 0 {
		err = Utils.SendEmail(user.Username, otp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		err = Utils.SendOTP(user.Username, otp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	View.LoginMessageView("OTP sent to registered email/phone", c)
	return
}

type VerifyOTPRequest struct {
	Username string `json:"username" binding:"required,min=11,max=30"`
	OTP      string `json:"otp" binding:"required,len=6"`
}

func VerifyOtpResponse(c *gin.Context) {
	var verifyReq VerifyOTPRequest
	userRepository := Repository.NewUserRepository(DBConfiguration.GetDB())
	if err := c.ShouldBindJSON(&verifyReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if Utils.UsernameValidation(verifyReq.Username) == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}

	user, err := userRepository.FindByUsername(verifyReq.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	secretKey := user.SecretKey
	valid := totp.Validate(verifyReq.OTP, secretKey)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid OTP"})
		return
	}
	jwtSecret, _ := Utils.ReadFromEnvFile(".env", "JWT_SECRET")
	jwtMaker, err := Token.NewJWTMaker(jwtSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, _, err := jwtMaker.MakeToken(user.Username, time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("token", token, 86400, "/", "localhost", false, true)
	View.LoginUserView(token, c)
}

func OAuth2LoginResponse(c *gin.Context) {
	url := Utils.GetGoogleAuthURL("random-state")
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
	response, err := Utils.GetGoogleUserInfo(code, state, c)

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
