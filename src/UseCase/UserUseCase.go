package UseCase

import (
	"fmt"
	"github.com/403-access-denied/main-backend/src/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/Model"
	"github.com/403-access-denied/main-backend/src/Utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"time"
)

func generateSecretkeyfornewuser(user string) string {
	key, err := totp.Generate(totp.GenerateOpts{
		AccountName: user,
	})
	if err != nil {
		fmt.Println("Error generating Key:", err)
	}
	return key.Secret()
}

func CheckUserExists(username string) (bool, error) {
	var count int64
	err := DBConfiguration.GetDB().Model(&Model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func diagnoseString(str string) int {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	mobileRegex1 := `^\+98 9[0-9]{2} [0-9]{3} [0-9]{4}$`
	mobileRegex2 := `^\+989[0-9]{2}[0-9]{3}[0-9]{4}$`
	mobileRegex3 := `^\0 9[0-9]{2} [0-9]{3} [0-9]{4}$`
	mobileRegex4 := `^\09[0-9]{2}[0-9]{3}[0-9]{4}$`
	if match, _ := regexp.MatchString(emailRegex, str); match {
		return 0
	} else if match, _ := regexp.MatchString(mobileRegex1, str); match {
		return 1
	} else if match, _ := regexp.MatchString(mobileRegex2, str); match {
		return 1
	} else if match, _ := regexp.MatchString(mobileRegex3, str); match {
		return 1
	} else if match, _ := regexp.MatchString(mobileRegex4, str); match {
		return 1
	} else {
		return 2
	}
}

type UserRequest struct {
	Username string `json:"username" binding:"required,min=5,max=50"`
	//Username string `gorm:"type:varchar(50);not null;unique" json:"username"`
}

func LoginResponse(c *gin.Context) {
	var user UserRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if diagnoseString(user.Username) == 2 {
		c.JSON(400, gin.H{"error": "invalid phone or email"})
		return
	}
	//generate OTP
	//send OTP
	//sendOTP(username)
	//
	//
	//
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to registered email/phone"})
}
func VerifyOtpResponse(c *gin.Context) {
	var user UserRequest
	//check otp
	//checkOtp(username)
	//check if user exists
	exis, _ := CheckUserExists(user.Username)
	if exis {
		//login
		var foundUser Model.User
		DBConfiguration.GetDB().First(&foundUser, "username = ?", user.Username)
		if foundUser.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "error while finding user",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user logged in"})
	} else {
		//create user
		secretKey := generateSecretkeyfornewuser(user.Username)
		_ = secretKey
		user := Model.User{
			Model:         gorm.Model{},
			Username:      user.Username,
			Posters:       nil,
			Conversations: nil,
			MarkedPosters: nil,
		}
		res := DBConfiguration.GetDB().Create(&user)
		if res.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "error while creating user",
			})
			return
		}
		var foundUser Model.User
		DBConfiguration.GetDB().First(&foundUser, "username = ?", user.Username)
		if foundUser.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "error while creating user",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user created"})
	}
	//make jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	secret, _ := Utils.ReadFromEnvFile(".env", "JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while creating token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenString, 3600*24, "", "", false, true)
}

func ValidateResponse(c *gin.Context) {
	secret, _ := Utils.ReadFromEnvFile(".env", "JWT_SECRET")
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while validating token",
		})
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if time.Now().Unix() > int64(claims["exp"].(float64)) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "token expired",
			})
			return
		}
		var foundUser Model.User
		DBConfiguration.GetDB().First(&foundUser, "username = ?", claims["sub"])
		if foundUser.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "error while finding user",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user logged in"})
		c.Next()
	} else {
		fmt.Println(err)
	}
}

func LogedInResponse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "user logged in"})
}
