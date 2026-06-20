package handlers

import (
	"time"

	"hotelmanagementsystem.com/v2/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
	var roomCount int64
	var occupiedCount int64
	var guestCount int64
	var arrivalCount int64
	var revenue float64
	var bookings []models.Booking
	var rooms []models.Room

	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	if err := h.db.Model(&models.Room{}).Count(&roomCount).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to count rooms")
	}
	if err := h.db.Model(&models.Room{}).Where("status = ?", models.RoomStatusOccupied).Count(&occupiedCount).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to count occupied rooms")
	}
	if err := h.db.Model(&models.Guest{}).Count(&guestCount).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to count guests")
	}
	if err := h.db.Model(&models.Booking{}).Where("check_in_date >= ? AND check_in_date < ?", startOfDay, endOfDay).Count(&arrivalCount).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to count arrivals")
	}
	if err := h.db.Model(&models.Booking{}).
		Where("status IN ?", []string{models.BookingStatusReserved, models.BookingStatusCheckedIn}).
		Select("coalesce(sum(total_amount), 0)").
		Scan(&revenue).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to compute revenue")
	}
	if err := h.db.Preload("Guest").Preload("Room").Order("created_at desc").Limit(5).Find(&bookings).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to load recent bookings")
	}
	if err := h.db.Order("number asc").Limit(6).Find(&rooms).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to load room snapshot")
	}

	occupancyRate := 0.0
	if roomCount > 0 {
		occupancyRate = float64(occupiedCount) / float64(roomCount) * 100
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"stats": fiber.Map{
				"totalRooms":       roomCount,
				"occupiedRooms":    occupiedCount,
				"totalGuests":      guestCount,
				"todayArrivals":    arrivalCount,
				"occupancyRate":    occupancyRate,
				"projectedRevenue": revenue,
			},
			"recentBookings": bookings,
			"roomSnapshot":   rooms,
		},
	})
}
