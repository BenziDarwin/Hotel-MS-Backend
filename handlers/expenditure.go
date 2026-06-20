package handlers

import (
	"fmt"
	"strings"
	"time"

	"hotelmanagementsystem.com/v2/middleware"
	"hotelmanagementsystem.com/v2/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ExpenditureHandler struct {
	db *gorm.DB
}

type createExpenditureRecordRequest struct {
	Title      string  `json:"title"`
	Category   string  `json:"category"`
	Vendor     string  `json:"vendor"`
	Amount     float64 `json:"amount"`
	Notes      string  `json:"notes"`
	RecordedAt string  `json:"recordedAt"`
}

func NewExpenditureHandler(db *gorm.DB) *ExpenditureHandler {
	return &ExpenditureHandler{db: db}
}

func (h *ExpenditureHandler) ListExpenditureRecords(c *fiber.Ctx) error {
	var records []models.ExpenditureRecord
	if err := h.db.Preload("CreatedByUser").Order("recorded_at desc").Find(&records).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to load expenditure records")
	}

	return c.JSON(fiber.Map{
		"data": records,
	})
}

func (h *ExpenditureHandler) CreateExpenditureRecord(c *fiber.Ctx) error {
	var payload createExpenditureRecordRequest
	if err := c.BodyParser(&payload); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request payload")
	}

	payload.Title = strings.TrimSpace(payload.Title)
	payload.Category = strings.TrimSpace(payload.Category)
	payload.Vendor = strings.TrimSpace(payload.Vendor)
	payload.Notes = strings.TrimSpace(payload.Notes)

	if payload.Title == "" || payload.Category == "" || payload.Amount <= 0 {
		return respondError(c, fiber.StatusBadRequest, "title, category, and amount are required")
	}

	recordedAt := time.Now()
	if payload.RecordedAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, payload.RecordedAt)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "recordedAt must be ISO-8601")
		}
		recordedAt = parsedTime
	}

	userID := middleware.UserIDFromToken(c)
	if userID == "" {
		return respondError(c, fiber.StatusUnauthorized, "unable to resolve current user")
	}

	record := models.ExpenditureRecord{
		Title:           payload.Title,
		Category:        payload.Category,
		Vendor:          payload.Vendor,
		Amount:          payload.Amount,
		Notes:           payload.Notes,
		RecordedAt:      recordedAt,
		ReceiptNumber:   fmt.Sprintf("RCPT-EXP-%d", time.Now().UnixNano()),
		CreatedByUserID: userID,
	}

	if err := h.db.Create(&record).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create expenditure record")
	}

	if err := h.db.Preload("CreatedByUser").First(&record, "id = ?", record.ID).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "expenditure record created but reload failed")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "expenditure recorded",
		"data":    record,
	})
}
