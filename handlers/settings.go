package handlers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hotelmanagementsystem.com/v2/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SettingsHandler struct {
	db *gorm.DB
}

func NewSettingsHandler(db *gorm.DB) *SettingsHandler {
	return &SettingsHandler{db: db}
}

func (h *SettingsHandler) GetHotelSettings(c *fiber.Ctx) error {
	settings, err := h.ensureHotelSettings()
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to load hotel settings")
	}

	return c.JSON(fiber.Map{
		"data": settings,
	})
}

func (h *SettingsHandler) UpdateHotelSettings(c *fiber.Ctx) error {
	settings, err := h.ensureHotelSettings()
	if err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to load hotel settings")
	}

	hotelName := strings.TrimSpace(c.FormValue("hotelName"))
	if hotelName == "" {
		return respondError(c, fiber.StatusBadRequest, "hotelName is required")
	}

	currency := strings.ToUpper(strings.TrimSpace(c.FormValue("currency")))
	if currency == "" {
		return respondError(c, fiber.StatusBadRequest, "currency is required")
	}

	settings.HotelName = hotelName
	settings.Currency = currency
	settings.Phone = strings.TrimSpace(c.FormValue("phone"))
	settings.Email = strings.TrimSpace(c.FormValue("email"))
	settings.Address = strings.TrimSpace(c.FormValue("address"))

	fileHeader, err := c.FormFile("hotelImage")
	if err == nil && fileHeader != nil {
		imagePath, saveErr := saveHotelImage(c, fileHeader)
		if saveErr != nil {
			return respondError(c, fiber.StatusInternalServerError, "failed to store hotel image")
		}

		if settings.HotelImage != "" && settings.HotelImage != imagePath {
			_ = deleteStoredFile(settings.HotelImage)
		}

		settings.HotelImage = imagePath
	}

	if err := h.db.Save(&settings).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "failed to update hotel settings")
	}

	return c.JSON(fiber.Map{
		"message": "hotel settings updated",
		"data":    settings,
	})
}

func (h *SettingsHandler) ensureHotelSettings() (models.HotelSettings, error) {
	var settings models.HotelSettings
	err := h.db.First(&settings).Error
	if err == nil {
		return settings, nil
	}

	if err != gorm.ErrRecordNotFound {
		return settings, err
	}

	settings = models.HotelSettings{
		HotelName: "Hotel MS",
		Currency:  "USD",
		Phone:     "",
		Email:     "",
		Address:   "",
	}
	if createErr := h.db.Create(&settings).Error; createErr != nil {
		return settings, createErr
	}

	if settings.Currency == "" {
		settings.Currency = "USD"
		if saveErr := h.db.Save(&settings).Error; saveErr != nil {
			return settings, saveErr
		}
	}

	return settings, nil
}

func saveHotelImage(c *fiber.Ctx, fileHeader *multipart.FileHeader) (string, error) {
	uploadDir := filepath.Join("uploads", "hotel-branding")
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		ext = ".png"
	}

	fileName := fmt.Sprintf("hotel-%d%s", time.Now().UnixNano(), ext)
	targetPath := filepath.Join(uploadDir, fileName)

	return "/" + filepath.ToSlash(targetPath), c.SaveFile(fileHeader, targetPath)
}

func deleteStoredFile(publicPath string) error {
	trimmed := strings.TrimPrefix(publicPath, "/")
	if trimmed == "" {
		return nil
	}

	err := os.Remove(trimmed)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
