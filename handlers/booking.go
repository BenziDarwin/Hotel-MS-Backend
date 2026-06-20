package handlers

import (
	"time"

	"hotelmanagementsystem.com/v2/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BookingHandler struct {
	db *gorm.DB
}

type createBookingRequest struct {
	GuestID      string  `json:"guestId"`
	RoomID       string  `json:"roomId"`
	CheckInDate  string  `json:"checkInDate"`
	CheckOutDate string  `json:"checkOutDate"`
	GuestsCount  int     `json:"guestsCount"`
	TotalAmount  float64 `json:"totalAmount"`
	Source       string  `json:"source"`
	SpecialNote  string  `json:"specialNote"`
}

type updateBookingStatusRequest struct {
	Status string `json:"status"`
}

func NewBookingHandler(db *gorm.DB) *BookingHandler {
	return &BookingHandler{db: db}
}

func (h *BookingHandler) ListBookings(c *fiber.Ctx) error {
	var bookings []models.Booking
	if err := h.db.Preload("Guest").Preload("Room").Order("check_in_date asc").Find(&bookings).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to load bookings")
	}

	return c.JSON(fiber.Map{
		"data": bookings,
	})
}

func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	var payload createBookingRequest
	if err := c.BodyParser(&payload); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request payload")
	}

	if payload.GuestID == "" || payload.RoomID == "" || payload.CheckInDate == "" || payload.CheckOutDate == "" || payload.GuestsCount == 0 || payload.TotalAmount == 0 {
		return respondError(c, fiber.StatusBadRequest, "guestId, roomId, dates, guestsCount, and totalAmount are required")
	}

	checkInDate, err := time.Parse(time.RFC3339, payload.CheckInDate)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "checkInDate must be ISO-8601")
	}

	checkOutDate, err := time.Parse(time.RFC3339, payload.CheckOutDate)
	if err != nil {
		return respondError(c, fiber.StatusBadRequest, "checkOutDate must be ISO-8601")
	}

	booking := models.Booking{
		GuestID:      payload.GuestID,
		RoomID:       payload.RoomID,
		CheckInDate:  checkInDate,
		CheckOutDate: checkOutDate,
		Status:       models.BookingStatusReserved,
		GuestsCount:  payload.GuestsCount,
		TotalAmount:  payload.TotalAmount,
		Source:       payload.Source,
		SpecialNote:  payload.SpecialNote,
	}

	if err := h.db.Create(&booking).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create booking")
	}

	if err := h.db.Preload("Guest").Preload("Room").First(&booking, "id = ?", booking.ID).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "booking created but failed to reload details")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "booking created successfully",
		"data":    booking,
	})
}

func (h *BookingHandler) UpdateBookingStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	var payload updateBookingStatusRequest
	if err := c.BodyParser(&payload); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request payload")
	}

	if payload.Status == "" {
		return respondError(c, fiber.StatusBadRequest, "status is required")
	}

	var booking models.Booking
	if err := h.db.Preload("Room").First(&booking, "id = ?", id).Error; err != nil {
		return respondError(c, fiber.StatusNotFound, "booking not found")
	}

	booking.Status = payload.Status
	if err := h.db.Save(&booking).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to update booking")
	}

	if payload.Status == models.BookingStatusCheckedIn {
		if err := h.db.Model(&models.Room{}).Where("id = ?", booking.RoomID).Update("status", models.RoomStatusOccupied).Error; err != nil {
			return respondError(c, fiber.StatusInternalServerError, "booking updated but room sync failed")
		}
	}

	if payload.Status == models.BookingStatusCheckedOut || payload.Status == models.BookingStatusCancelled {
		if err := h.db.Model(&models.Room{}).Where("id = ?", booking.RoomID).Update("status", models.RoomStatusAvailable).Error; err != nil {
			return respondError(c, fiber.StatusInternalServerError, "booking updated but room sync failed")
		}
	}

	if err := h.db.Preload("Guest").Preload("Room").First(&booking, "id = ?", id).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to reload booking")
	}

	return c.JSON(fiber.Map{
		"message": "booking status updated",
		"data":    booking,
	})
}
