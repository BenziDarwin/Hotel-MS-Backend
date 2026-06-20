package handlers

import (
	"strings"

	"hotelmanagementsystem.com/v2/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type IncomeCategoryHandler struct {
	db *gorm.DB
}

type createIncomeCategoryRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func NewIncomeCategoryHandler(db *gorm.DB) *IncomeCategoryHandler {
	return &IncomeCategoryHandler{db: db}
}

func (h *IncomeCategoryHandler) ListCategories(c *fiber.Ctx) error {
	categoryType := normalizeIncomeType(c.Query("type"))

	query := h.db.Order("name asc")
	if categoryType != "" {
		query = query.Where("type = ?", categoryType)
	}

	var categories []models.IncomeCategory
	if err := query.Find(&categories).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to load income categories")
	}

	return c.JSON(fiber.Map{
		"data": categories,
	})
}

func (h *IncomeCategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var payload createIncomeCategoryRequest
	if err := c.BodyParser(&payload); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request payload")
	}

	payload.Name = strings.TrimSpace(payload.Name)
	payload.Description = strings.TrimSpace(payload.Description)
	payload.Type = normalizeIncomeType(payload.Type)

	if payload.Name == "" || payload.Type == "" {
		return respondError(c, fiber.StatusBadRequest, "name and type are required")
	}

	category := models.IncomeCategory{
		Name:        payload.Name,
		Type:        payload.Type,
		Description: payload.Description,
	}

	if err := h.db.Create(&category).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create income category")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "income category created",
		"data":    category,
	})
}

func normalizeIncomeType(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case models.IncomeCategoryTypeHotel:
		return models.IncomeCategoryTypeHotel
	case models.IncomeCategoryTypeOther:
		return models.IncomeCategoryTypeOther
	default:
		return ""
	}
}
