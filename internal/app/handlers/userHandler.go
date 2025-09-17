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
