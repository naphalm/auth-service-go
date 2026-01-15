package handlers

import (
	"auth-service/src/db"
	"auth-service/src/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(creds.Password), 12)

	_, err := db.DB.Exec(
		"INSERT INTO users (email, password) VALUES ($1,$2)",
		creds.Email, string(hash),
	)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user exists"})
		return
	}

	c.Status(http.StatusCreated)
}

func Login(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	var user models.User
	err := db.DB.QueryRow(
		"SELECT id,password FROM users WHERE email=$1",
		creds.Email,
	).Scan(&user.ID, &user.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	signed, _ := token.SignedString([]byte(secret))
	c.JSON(http.StatusOK, gin.H{"token": signed})
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "valid"})
}
