package handlers

import (
	"hotelmanagementsystem.com/v2/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RoomHandler struct {
	db *gorm.DB
}

type createRoomRequest struct {
	Number   string  `json:"number"`
	Type     string  `json:"type"`
	Floor    int     `json:"floor"`
	Rate     float64 `json:"rate"`
	Capacity int     `json:"capacity"`
	Status   string  `json:"status"`
}

type updateRoomStatusRequest struct {
	Status string `json:"status"`
}

func NewRoomHandler(db *gorm.DB) *RoomHandler {
	return &RoomHandler{db: db}
}

func (h *RoomHandler) ListRooms(c *fiber.Ctx) error {
	var rooms []models.Room
	if err := h.db.Order("number asc").Find(&rooms).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to load rooms")
	}

	return c.JSON(fiber.Map{
		"data": rooms,
	})
}

func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	var payload createRoomRequest
	if err := c.BodyParser(&payload); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request payload")
	}

	if payload.Number == "" || payload.Type == "" || payload.Floor == 0 || payload.Capacity == 0 || payload.Rate == 0 {
		return respondError(c, fiber.StatusBadRequest, "number, type, floor, capacity, and rate are required")
	}

	status := payload.Status
	if status == "" {
		status = models.RoomStatusAvailable
	}

	room := models.Room{
		Number:   payload.Number,
		Type:     payload.Type,
		Floor:    payload.Floor,
		Rate:     payload.Rate,
		Capacity: payload.Capacity,
		Status:   status,
	}

	if err := h.db.Create(&room).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to create room")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "room created successfully",
		"data":    room,
	})
}

func (h *RoomHandler) UpdateRoomStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	var payload updateRoomStatusRequest
	if err := c.BodyParser(&payload); err != nil {
		return respondError(c, fiber.StatusBadRequest, "invalid request payload")
	}

	if payload.Status == "" {
		return respondError(c, fiber.StatusBadRequest, "status is required")
	}

	var room models.Room
	if err := h.db.First(&room, "id = ?", id).Error; err != nil {
		return respondError(c, fiber.StatusNotFound, "room not found")
	}

	room.Status = payload.Status
	if err := h.db.Save(&room).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to update room status")
	}

	return c.JSON(fiber.Map{
		"message": "room status updated",
		"data":    room,
	})
}
