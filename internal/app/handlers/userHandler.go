package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"resapi/internal/domain/models"
	"resapi/internal/domain/repo"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo repo.UserRepo
}

func NewUserHandler(userRepo repo.UserRepo) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (uh *UserHandler) CreateUserHandler(ctx *gin.Context) {

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid req"})
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name and email are required"})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("Error hashing password: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashPassword),
		CreatedAt: time.Now(),
	}

	if err := uh.userRepo.CreateUser(ctx.Request.Context(), &user); err != nil {
		log.Printf("Error create user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user create succesfully",
		"user":    user,
	})

}

func (uh *UserHandler) GetUser(ctx *gin.Context) {

	idStr := ctx.Param("id")

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := uh.userRepo.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	user.Password = ""

	ctx.JSON(http.StatusOK, user)
}

func (uh *UserHandler) UpdateUserHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	existingUser, err := uh.userRepo.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Name != "" {
		existingUser.Name = req.Name
	}
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.Password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
			return
		}
		existingUser.Password = string(hashPassword)
	}

	if err := uh.userRepo.UpdateUser(ctx.Request.Context(), existingUser); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		log.Printf("Error updating user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
		"user": gin.H{
			"id":    existingUser.ID,
			"name":  existingUser.Name,
			"email": existingUser.Email,
		},
	})
}

func (uh *UserHandler) DeleteUserHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	_, err = uh.userRepo.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	if err := uh.userRepo.DeleteUser(ctx.Request.Context(), userID); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		log.Printf("Error deleting user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user deleted successfully",
	})
}
