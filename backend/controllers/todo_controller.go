package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Zhaobo-Wang/go-projects/database"
	"github.com/Zhaobo-Wang/go-projects/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func respondJSON(c *gin.Context, status int, payload interface{}) {
	c.JSON(status, gin.H{"data": payload})
}

func respondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// --- Request structs with validation tags ---
type CreateTodoInput struct {
	Title       string  `json:"title" binding:"required,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
	Completed   *bool   `json:"completed" binding:"omitempty"`
}

type UpdateTodoInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

const defaultUserID uint = 1

// --- GetTodos: 支持分页、过滤和排序 ---
func GetTodos(c *gin.Context) {
	// 使用请求的 context，设置超时以避免长时间阻塞
	reqCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(reqCtx, 5*time.Second)
	defer cancel()

	userID := defaultUserID

	// 分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "5")
	sort := c.DefaultQuery("sort", "created_at desc") // e.g. "created_at desc" or "title asc"
	completedFilter := c.Query("completed")           // "true"/"false" 可选

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		respondError(c, http.StatusBadRequest, "invalid page parameter")
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		respondError(c, http.StatusBadRequest, "invalid page_size parameter (1-100)")
		return
	}

	var todos []models.Todo
	var total int64

	db := database.DB.WithContext(ctx).Model(&models.Todo{}).Where("user_id = ?", userID)

	if completedFilter != "" {
		if completedFilter == "true" {
			db = db.Where("completed = ?", true)
		} else if completedFilter == "false" {
			db = db.Where("completed = ?", false)
		} else {
			respondError(c, http.StatusBadRequest, "invalid completed filter; must be true or false")
			return
		}
	}

	if err := db.Count(&total).Error; err != nil {
		log.Printf("GetTodos count error: %v", err)
		respondError(c, http.StatusInternalServerError, "database error")
		return
	}

	offset := (page - 1) * pageSize
	if err := db.Order(sort).Offset(offset).Limit(pageSize).Find(&todos).Error; err != nil {
		log.Printf("GetTodos find error: %v", err)
		respondError(c, http.StatusInternalServerError, "database error")
		return
	}

	respondJSON(c, http.StatusOK, gin.H{
		"items": todos,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"total_pages": func() int {
				if total == 0 {
					return 0
				}
				p := int((total + int64(pageSize) - 1) / int64(pageSize))
				return p
			}(),
		},
	})
}

// --- CreateTodo: 使用验证、事务、返回 Location header ---
func CreateTodo(c *gin.Context) {
	reqCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(reqCtx, 5*time.Second)
	defer cancel()

	userID := defaultUserID

	var input CreateTodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	tx := database.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		log.Printf("CreateTodo begin tx error: %v", tx.Error)
		respondError(c, http.StatusInternalServerError, "database error")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	todo := models.Todo{
		Title:  input.Title,
		UserID: userID,
	}
	if input.Description != nil {
		todo.Description = *input.Description
	}
	if input.Completed != nil {
		todo.Completed = *input.Completed
	}

	if err := tx.Create(&todo).Error; err != nil {
		tx.Rollback()
		log.Printf("CreateTodo create error: %v", err)
		respondError(c, http.StatusInternalServerError, "database error")
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("CreateTodo commit error: %v", err)
		respondError(c, http.StatusInternalServerError, "database error")
		return
	}

	c.Header("Location", fmt.Sprintf("/todos/%d", todo.ID))
	respondJSON(c, http.StatusCreated, todo)
}

// --- GetTodo: 单个资源，带权限检查 & gorm ErrRecordNotFound 判断 ---
func GetTodo(c *gin.Context) {
	reqCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(reqCtx, 5*time.Second)
	defer cancel()

	userID := defaultUserID
	id := c.Param("id")

	var todo models.Todo
	if err := database.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(c, http.StatusNotFound, "todo not found or not authorized")
			return
		}
		log.Printf("GetTodo db error: %v", err)
		respondError(c, http.StatusInternalServerError, "database error")
		return
	}

	respondJSON(c, http.StatusOK, todo)
}

// --- UpdateTodo: 你已有实现的增强版（保持指针字段以便 PATCH） ---
func UpdateTodo(c *gin.Context) {
	reqCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(reqCtx, 5*time.Second)
	defer cancel()

	userID := defaultUserID
	id := c.Param("id")

	log.Printf("UpdateTodo called - id=%s, userID=%d", id, userID)

	// 1) 查找目标 todo
	var todo models.Todo
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(c, http.StatusNotFound, "todo not found or not authorized")
			return
		}
		log.Printf("UpdateTodo fetch error: %v", err)
		respondError(c, http.StatusInternalServerError, "database error")
		return
	}

	// 2) 绑定输入（指针字段）
	var input UpdateTodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Completed != nil {
		updates["completed"] = *input.Completed
	}

	if len(updates) == 0 {
		respondError(c, http.StatusBadRequest, "no fields to update")
		return
	}

	if err := database.DB.WithContext(ctx).Model(&todo).Updates(updates).Error; err != nil {
		log.Printf("UpdateTodo update error: %v", err)
		respondError(c, http.StatusInternalServerError, "database error")
		return
	}

	if err := database.DB.WithContext(ctx).First(&todo, todo.ID).Error; err != nil {
		log.Printf("UpdateTodo reload error: %v", err)
		respondError(c, http.StatusInternalServerError, "database error")
		return
	}

	respondJSON(c, http.StatusOK, todo)
}

func DeleteTodo(c *gin.Context) {
	reqCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(reqCtx, 5*time.Second)
	defer cancel()

	userID := defaultUserID
	id := c.Param("id")

	var todo models.Todo
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(c, http.StatusNotFound, "todo not found or not authorized")
			return
		}
		log.Printf("DeleteTodo fetch error: %v", err)
		respondError(c, http.StatusInternalServerError, "database error")
		return
	}

	if err := database.DB.WithContext(ctx).Delete(&todo).Error; err != nil {
		log.Printf("DeleteTodo delete error: %v", err)
		respondError(c, http.StatusInternalServerError, "database error")
		return
	}

	respondJSON(c, http.StatusOK, gin.H{
		"message": "Todo successfully deleted",
		"id": id,
	})
}
