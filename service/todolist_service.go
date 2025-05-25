package service

import (
	"log"
	"net/http"
	"strings"

	postgresdb "github.com/M1123Ananda/tododo/db"
	"github.com/M1123Ananda/tododo/model"
	"github.com/M1123Ananda/tododo/utils"
	"github.com/gin-gonic/gin"
)

func CreateToDoItem(ctx *gin.Context) {
	var req model.CreateToDoItemRequest

	ctx.BindJSON(&req)

	DB := postgresdb.DB
	if DB == nil {
		log.Panic("DB is not initialized")
		ctx.IndentedJSON(http.StatusInternalServerError, model.AuthError{Error: "Internal Error"})
		return
	}

	bearer := ctx.Request.Header.Get("Authorization")
	if bearer == "" {
		ctx.IndentedJSON(http.StatusUnauthorized, model.AuthError{Error: "Unauthorized"})
		return
	} else {

		parts := strings.Split(bearer, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.IndentedJSON(http.StatusUnauthorized, model.AuthError{Error: "Invalid Authorization Header"})
			return
		}
		token := parts[1]

		claims, err := utils.VerifyToken(token)
		if err != nil {
			log.Default().Printf("Failed to verify token %v", err)
			ctx.IndentedJSON(http.StatusUnauthorized, model.AuthError{Error: "Unauthorized"})
			return
		}

		item := &model.ToDo{
			UserEmail:   claims.Email,
			Title:       req.Title,
			Description: req.Description,
		}

		tx := DB.Create(&item)
		if tx.Error != nil {
			log.Panic("Todo Item could not be inserted into DB")
			ctx.IndentedJSON(http.StatusInternalServerError, model.AuthError{Error: "Internal Error"})
			return
		}

		ctx.IndentedJSON(http.StatusOK, model.
			CreateToDoItemResponse{ID: int(item.ID),
			Title:       item.Title,
			Description: item.Description})
	}
}
