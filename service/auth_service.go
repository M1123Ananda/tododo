package service

import (
	"errors"
	"log"
	"net/http"

	postgresdb "github.com/M1123Ananda/tododo/db"
	"github.com/M1123Ananda/tododo/model"
	"github.com/M1123Ananda/tododo/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUser(ctx *gin.Context) {
	var req model.RegisterRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Panic("Cannot Bind JSON to Request")
		ctx.IndentedJSON(http.StatusInternalServerError, model.AuthError{Error: "Internal Error"})
		return
	}

	DB := postgresdb.DB
	if DB == nil {
		log.Panic("DB is not initialized")
		ctx.IndentedJSON(http.StatusInternalServerError, model.AuthError{Error: "Internal Error"})
		return
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Panicf("Failed to Hash Password: %v", err)
		ctx.IndentedJSON(http.StatusInternalServerError, model.AuthError{Error: "Internal Error"})
	}

	tx := DB.Create(&model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: passwordHash,
	})
	if tx.Error != nil {
		log.Panic("User could not be inserted into DB")
		ctx.IndentedJSON(http.StatusInternalServerError, model.AuthError{Error: "Internal Error"})
		return
	}

	token, err := utils.GenerateToken(req.Email)
	if err != nil {
		log.Panicf("Failed to generate token: %v", err)
		ctx.IndentedJSON(http.StatusInternalServerError, model.AuthError{Error: "Internal Error"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, model.AuthenticateResponse{Token: token})
}

func LoginUser(ctx *gin.Context) {
	var req model.LoginRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Panic("Cannot Bind JSON to Request")
		ctx.IndentedJSON(http.StatusInternalServerError, model.AuthError{Error: "Internal Error"})
		return
	}

	DB := postgresdb.DB
	if DB == nil {
		log.Panic("DB is not initialized")
		ctx.IndentedJSON(http.StatusInternalServerError, model.AuthError{Error: "Internal Error"})
		return
	}

	var user model.User

	tx := DB.Where(model.User{Email: req.Email}).First(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			ctx.IndentedJSON(http.StatusUnauthorized, model.AuthError{Error: "Incorrect Credentials"})
			return
		} else {
			log.Panicf("DB Error: %v", tx.Error)
			ctx.IndentedJSON(http.StatusInternalServerError, model.AuthError{Error: "Internal Error"})
			return
		}
	}

	correctPassword := utils.VerifyPassword(req.Password, user.Password)

	if correctPassword {
		token, err := utils.GenerateToken(req.Email)
		if err != nil {
			log.Panicf("Failed to generate token: %v", err)
			ctx.IndentedJSON(http.StatusInternalServerError, model.AuthError{Error: "Internal Error"})
			return
		}
		ctx.IndentedJSON(http.StatusOK, model.AuthenticateResponse{Token: token})
		return
	} else {
		ctx.IndentedJSON(http.StatusUnauthorized, model.AuthError{Error: "Incorrect Credentials"})
		return
	}
}
