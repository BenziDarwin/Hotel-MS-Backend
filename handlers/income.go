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

type IncomeHandler struct {
	db *gorm.DB
}

type createIncomeRecordRequest struct {
	Type       string  `json:"type"`
	Title      string  `json:"title"`
	CategoryID string  `json:"categoryId"`
	BookingID  string  `json:"bookingId"`
	GuestName  string  `json:"guestName"`
	Amount     float64 `json:"amount"`
	Notes      string  `json:"notes"`
	RecordedAt string  `json:"recordedAt"`
}

func NewIncomeHandler(db *gorm.DB) *IncomeHandler {
	return &IncomeHandler{db: db}
}

func (h *IncomeHandler) ListIncomeRecords(c *fiber.Ctx) error {
	incomeType := normalizeIncomeType(c.Query("type"))

	query := h.db.Preload("Category").Preload("Booking.Guest").Preload("Booking.Room").Preload("CreatedByUser").Order("recorded_at desc")
	if incomeType != "" {
		query = query.Where("type = ?", incomeType)
	}

	var records []models.IncomeRecord
	if err := query.Find(&records).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to load income records")
	}

	return c.JSON(fiber.Map{
		"data": records,
	})
}

func (h *IncomeHandler) CreateIncomeRecord(c *fiber.Ctx) error {
	var payload createIncomeRecordRequest
	if err := c.BodyParser(&payload); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request payload")
	}

	payload.Type = normalizeIncomeType(payload.Type)
	payload.Title = strings.TrimSpace(payload.Title)
	payload.CategoryID = strings.TrimSpace(payload.CategoryID)
	payload.BookingID = strings.TrimSpace(payload.BookingID)
	payload.GuestName = strings.TrimSpace(payload.GuestName)
	payload.Notes = strings.TrimSpace(payload.Notes)

	if payload.Type == "" || payload.Title == "" || payload.CategoryID == "" || payload.Amount <= 0 {
		return respondError(c, fiber.StatusBadRequest, "type, title, categoryId, and amount are required")
	}

	recordedAt := time.Now()
	if payload.RecordedAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, payload.RecordedAt)
		if err != nil {
			return respondError(c, fiber.StatusBadRequest, "recordedAt must be ISO-8601")
		}
		recordedAt = parsedTime
	}

	var category models.IncomeCategory
	if err := h.db.First(&category, "id = ?", payload.CategoryID).Error; err != nil {
		return respondError(c, fiber.StatusBadRequest, "income category not found")
	}

	if category.Type != payload.Type {
		return respondError(c, fiber.StatusBadRequest, "income category type mismatch")
	}

	userID := middleware.UserIDFromToken(c)
	if userID == "" {
		return respondError(c, fiber.StatusUnauthorized, "unable to resolve current user")
	}

	record := models.IncomeRecord{
		Type:            payload.Type,
		Title:           payload.Title,
		CategoryID:      payload.CategoryID,
		GuestName:       payload.GuestName,
		Amount:          payload.Amount,
		Notes:           payload.Notes,
		RecordedAt:      recordedAt,
		ReceiptNumber:   h.nextReceiptNumber(payload.Type),
		CreatedByUserID: userID,
	}

	if payload.BookingID != "" {
		var booking models.Booking
		if err := h.db.Preload("Guest").First(&booking, "id = ?", payload.BookingID).Error; err != nil {
			return respondError(c, fiber.StatusBadRequest, "booking not found")
		}

		record.BookingID = &booking.ID
		if record.GuestName == "" {
			record.GuestName = strings.TrimSpace(booking.Guest.FirstName + " " + booking.Guest.LastName)
		}
	}

	if record.GuestName == "" {
		return respondError(c, fiber.StatusBadRequest, "guestName is required")
	}

	if err := h.db.Create(&record).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create income record")
	}

	if err := h.db.Preload("Category").Preload("Booking.Guest").Preload("Booking.Room").Preload("CreatedByUser").First(&record, "id = ?", record.ID).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "income record created but reload failed")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "income recorded",
		"data":    record,
	})
}

func (h *IncomeHandler) nextReceiptNumber(recordType string) string {
	prefix := "OTH"
	if recordType == models.IncomeCategoryTypeHotel {
		prefix = "HTL"
	}

	return fmt.Sprintf("RCPT-%s-%d", prefix, time.Now().UnixNano())
}
