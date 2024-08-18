package controllers

import (
	"fmt"
	"go-jwt/db"
	"go-jwt/helpers"
	"go-jwt/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := db.Init()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		var count int64
		db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This email already exists"})
			return
		}

		db.Model(&models.User{}).Where("phone = ?", user.Phone).Count(&count)
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This phone number already exists"})
			return
		}

		password := helpers.HashPassword(*user.Password)
		user.Password = &password

		user.Created_at = time.Now()
		user.Updated_at = time.Now()

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, fmt.Sprintf("%d", user.ID))
		user.Token = &token
		user.Refresh_token = &refreshToken

		db.Save(&user)

		c.JSON(http.StatusOK, gin.H{"user_id": user.User_id})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := db.Init()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Where("email = ?", user.Email).First(&foundUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or password is incorrect"})
			return
		}

		passwordIsValid, msg := helpers.VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, fmt.Sprintf("%d", foundUser.ID))
		helpers.UpdateAllTokens(token, refreshToken, fmt.Sprintf("%d", foundUser.ID))

		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		db := db.Init()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}
		startIndex := (page - 1) * recordPerPage

		var users []models.User
		var totalUsers int64

		db.Model(&models.User{}).Count(&totalUsers)
		db.Offset(startIndex).Limit(recordPerPage).Find(&users)

		c.JSON(http.StatusOK, gin.H{"total_count": totalUsers, "user_items": users})
	}
}

func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := db.Init()

		userId := c.Param("user_id")

		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		if err := db.First(&user, userId).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": user})
	}
}
