package service

import (
	"log"
	"net/http"

	postgresdb "github.com/M1123Ananda/tododo/db"
	"github.com/M1123Ananda/tododo/model"
	"github.com/M1123Ananda/tododo/utils"
	"github.com/gin-gonic/gin"
)

func Registeruser(ctx *gin.Context) {
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

	tx := DB.Create(&model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
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

	ctx.IndentedJSON(http.StatusOK, model.RegisterResponse{Token: token})
}