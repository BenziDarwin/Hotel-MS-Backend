package handlers

import (
	"hotelmanagementsystem.com/v2/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GuestHandler struct {
	db *gorm.DB
}

type createGuestRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Nationality string `json:"nationality"`
}

func NewGuestHandler(db *gorm.DB) *GuestHandler {
	return &GuestHandler{db: db}
}

func (h *GuestHandler) ListGuests(c *fiber.Ctx) error {
	var guests []models.Guest
	if err := h.db.Order("created_at desc").Find(&guests).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to load guests")
	}

	return c.JSON(fiber.Map{
		"data": guests,
	})
}

func (h *GuestHandler) CreateGuest(c *fiber.Ctx) error {
	var payload createGuestRequest
	if err := c.BodyParser(&payload); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request payload")
	}

	if payload.FirstName == "" || payload.LastName == "" || payload.Email == "" || payload.Phone == "" {
		return respondError(c, fiber.StatusBadRequest, "firstName, lastName, email, and phone are required")
	}

	guest := models.Guest{
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Email:       payload.Email,
		Phone:       payload.Phone,
		Nationality: payload.Nationality,
	}

	if err := h.db.Create(&guest).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create guest")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "guest created successfully",
		"data":    guest,
	})
}
