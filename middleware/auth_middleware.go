package middleware

import (
	"fmt"
	"net/http"
	"t-card/config/app_config"
	"t-card/database"
	"t-card/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(ctx *gin.Context) {
	// Get the cookie off req
	tokenString, err := ctx.Cookie("Authorization")

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	// Decode the cookie
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(app_config.SECRET_KEY), nil
	})

	if err != nil {
		fmt.Println(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		// Find the user with token sub
		var user models.User
		database.DB.First(&user, claims["sub"])

		if *user.ID == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		// Attach user to the req
		ctx.Set("user", user)

	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	// Continue
	ctx.Next()
}

func RequireEmployer(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	userData := user.(models.User)
	if *userData.Role == "employer" {
		ctx.Next()
	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "only employer can this work",
		})
		return
	}
}

func RequireFreelancer(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	userData := user.(models.User)
	if *userData.Role == "freelancer" {
		ctx.Next()
	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "only freelancer can this work",
		})
	}
}
