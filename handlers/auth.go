package handlers

import (
	"time"

	"hotelmanagementsystem.com/v2/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db            *gorm.DB
	jwtSecret     string
	tokenTTLHours int
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(db *gorm.DB, jwtSecret string, tokenTTLHours int) *AuthHandler {
	return &AuthHandler{
		db:            db,
		jwtSecret:     jwtSecret,
		tokenTTLHours: tokenTTLHours,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var payload registerRequest
	if err := c.BodyParser(&payload); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request payload")
	}

	if payload.Name == "" || payload.Email == "" || payload.Password == "" {
		return respondError(c, fiber.StatusBadRequest, "name, email, and password are required")
	}

	var existing models.User
	if err := h.db.Where("email = ?", payload.Email).First(&existing).Error; err == nil {
		return respondError(c, fiber.StatusConflict, "user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to secure password")
	}

	role := payload.Role
	if role == "" {
		role = models.RoleManager
	}

	user := models.User{
		Name:         payload.Name,
		Email:        payload.Email,
		PasswordHash: string(hash),
		Role:         role,
	}

	if err := h.db.Create(&user).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create user")
	}

	token, err := h.generateToken(user)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create token")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user registered successfully",
		"token":   token,
		"user":    user,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var payload loginRequest
	if err := c.BodyParser(&payload); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request payload")
	}

	var user models.User
	if err := h.db.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		return respondError(c, fiber.StatusUnauthorized, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		return respondError(c, fiber.StatusUnauthorized, "invalid credentials")
	}

	token, err := h.generateToken(user)
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create token")
	}

	return c.JSON(fiber.Map{
		"message": "login successful",
		"token":   token,
		"user":    user,
	})
}

func (h *AuthHandler) generateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Duration(h.tokenTTLHours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}
