package service

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	postgresdb "github.com/M1123Ananda/tododo/db"
	"github.com/M1123Ananda/tododo/model"
	"github.com/M1123Ananda/tododo/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateToDoItem(ctx *gin.Context) {
	var req model.CreateToDoItemRequest

	ctx.BindJSON(&req)

	DB := postgresdb.DB
	if DB == nil {
		log.Default().Printf("DB is not initialized")
		ctx.IndentedJSON(http.StatusInternalServerError, model.RequestError{Error: "Internal Error"})
		return
	}

	bearer := ctx.Request.Header.Get("Authorization")
	if bearer == "" {
		ctx.IndentedJSON(http.StatusUnauthorized, model.RequestError{Error: "Unauthorized"})
		return
	} else {

		parts := strings.Split(bearer, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.IndentedJSON(http.StatusUnauthorized, model.RequestError{Error: "Invalid Authorization Header"})
			return
		}
		token := parts[1]

		claims, err := utils.VerifyToken(token)
		if err != nil {
			log.Default().Printf("Failed to verify token: %v", err)
			ctx.IndentedJSON(http.StatusUnauthorized, model.RequestError{Error: "Unauthorized"})
			return
		}

		item := &model.ToDo{
			UserEmail:   claims.Email,
			Title:       req.Title,
			Description: req.Description,
		}

		tx := DB.Create(&item)
		if tx.Error != nil {
			log.Default().Printf("Todo Item could not be inserted into DB: %v", err)
			ctx.IndentedJSON(http.StatusInternalServerError, model.RequestError{Error: "Internal Error"})
			return
		}

		ctx.IndentedJSON(http.StatusOK, model.
			CreateToDoItemResponse{ID: int(item.ID),
			Title:       item.Title,
			Description: item.Description})
	}
}

func GetToDoItem(id int) (*model.ToDo, error) {
	var item model.ToDo

	DB := postgresdb.DB
	if DB == nil {
		log.Default().Printf("DB is not initialized")
		return nil, errors.New("DB is not initialized")
	}

	err := DB.First(&item, id).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func UpdateToDoItem(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, model.RequestError{Error: "id must be an integer"})
		return
	}

	var req model.UpdateToDoItemRequest
	ctx.BindJSON(&req)

	DB := postgresdb.DB
	if DB == nil {
		log.Default().Printf("DB is not initialized")
		ctx.IndentedJSON(http.StatusInternalServerError, model.RequestError{Error: "Internal Error"})
		return
	}

	bearer := ctx.Request.Header.Get("Authorization")
	if bearer == "" {
		ctx.IndentedJSON(http.StatusUnauthorized, model.RequestError{Error: "Unauthorized"})
		return
	}

	token, err := utils.GetTokenFromBearer(bearer)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, model.RequestError{Error: "Invalid Authorization Header"})
		return
	}

	claims, err := utils.VerifyToken(token)
	if err != nil {
		log.Default().Printf("Failed to verify token %v", err)
		ctx.IndentedJSON(http.StatusUnauthorized, model.RequestError{Error: "Unauthorized"})
		return
	}

	item, err := GetToDoItem(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.IndentedJSON(http.StatusNotFound, model.RequestError{Error: "item to update not found"})
		return
	} else if err != nil {
		log.Default().Printf("cannot get item from DB: %v", err)
		ctx.IndentedJSON(http.StatusInternalServerError, model.RequestError{Error: "Internal Error"})
		return
	}

	if item.UserEmail != claims.Email {
		ctx.IndentedJSON(http.StatusUnauthorized, model.RequestError{Error: "Unauthorized"})
		return
	}

	newItem := model.ToDo{
		Title:       req.Title,
		Description: req.Description,
	}

	tx := DB.Model(item).Updates(newItem)
	if tx.Error != nil {
		log.Default().Printf("Todo Item could not be updated: %v", err)
		ctx.IndentedJSON(http.StatusInternalServerError, model.RequestError{Error: "Internal Error"})
		return
	}

	updatedItem, err := GetToDoItem(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.IndentedJSON(http.StatusNotFound, model.RequestError{Error: "newly updated item not found"})
		return
	} else if err != nil {
		log.Default().Printf("cannot get new item from DB: %v", err)
		ctx.IndentedJSON(http.StatusInternalServerError, model.RequestError{Error: "Internal Error"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, model.UpdateToDoItemResponse{
		ID:          int(updatedItem.ID),
		Title:       updatedItem.Title,
		Description: updatedItem.Description,
	})
}

func DeleteToDoItem(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, model.RequestError{Error: "id must be an integer"})
		return
	}

	var req model.UpdateToDoItemRequest
	ctx.BindJSON(&req)

	DB := postgresdb.DB
	if DB == nil {
		log.Default().Printf("DB is not initialized: %v", err)
		ctx.IndentedJSON(http.StatusInternalServerError, model.RequestError{Error: "Internal Error"})
		return
	}

	bearer := ctx.Request.Header.Get("Authorization")
	if bearer == "" {
		ctx.IndentedJSON(http.StatusUnauthorized, model.RequestError{Error: "Unauthorized"})
		return
	}

	token, err := utils.GetTokenFromBearer(bearer)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, model.RequestError{Error: "Invalid Authorization Header"})
		return
	}

	claims, err := utils.VerifyToken(token)
	if err != nil {
		log.Default().Printf("Failed to verify token: %v", err)
		ctx.IndentedJSON(http.StatusUnauthorized, model.RequestError{Error: "Unauthorized"})
		return
	}

	item, err := GetToDoItem(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.IndentedJSON(http.StatusNotFound, model.RequestError{Error: "item to update not found"})
		return
	} else if err != nil {
		log.Default().Printf("cannot get item from DB: %v", err)
		ctx.IndentedJSON(http.StatusInternalServerError, model.RequestError{Error: "Internal Error"})
		return
	}

	if item.UserEmail != claims.Email {
		ctx.IndentedJSON(http.StatusUnauthorized, model.RequestError{Error: "Unauthorized"})
		return
	}

	tx := DB.Delete(&item)
	if tx.Error != nil {
		log.Default().Printf("Todo Item could not be deleted: %v", err)
		ctx.IndentedJSON(http.StatusInternalServerError, model.RequestError{Error: "Internal Error"})
		return
	}

	ctx.IndentedJSON(http.StatusNoContent, model.DeleteToDoItemResponse{
		Success: true,
	})
}

func GetToDoItems(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		log.Default().Printf("Error Parsing Page Parameter: %v", err)
		ctx.IndentedJSON(http.StatusBadRequest, model.RequestError{Error: "page must be an integer and > 0"})
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		log.Default().Printf("Error Parsing Limit Parameter: %v", err)
		ctx.IndentedJSON(http.StatusBadRequest, model.RequestError{Error: "limit must be an integer and > 0"})
		return
	}

	DB := postgresdb.DB
	if DB == nil {
		log.Default().Printf("DB is not initialized: %v", err)
		ctx.IndentedJSON(http.StatusInternalServerError, model.RequestError{Error: "Internal Error"})
		return
	}

	bearer := ctx.Request.Header.Get("Authorization")
	if bearer == "" {
		ctx.IndentedJSON(http.StatusUnauthorized, model.RequestError{Error: "Unauthorized"})
		return
	}

	token, err := utils.GetTokenFromBearer(bearer)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, model.RequestError{Error: "Invalid Authorization Header"})
		return
	}

	claims, err := utils.VerifyToken(token)
	if err != nil {
		log.Default().Printf("Failed to verify token: %v", err)
		ctx.IndentedJSON(http.StatusUnauthorized, model.RequestError{Error: "Unauthorized"})
		return
	}

	var items []model.GetToDoItemsData

	tx := DB.Model(&model.ToDo{}).Select("id", "title", "description").Where("user_email = ?", claims.Email).Offset((page * limit) - limit).Limit(limit).Scan(&items)	
	if tx.Error != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.IndentedJSON(http.StatusNotFound, model.RequestError{Error: "item to update not found"})
				return
			} else {
				log.Default().Printf("cannot get items from DB: %v", err)
				ctx.IndentedJSON(http.StatusInternalServerError, model.RequestError{Error: "Internal Error"})
				return
			}
		}

	ctx.IndentedJSON(http.StatusOK, model.GetToDoItemsResponse{
		Data:  items,
		Page:  page,
		Limit: limit,
		Total: int(tx.RowsAffected),
	})
}
