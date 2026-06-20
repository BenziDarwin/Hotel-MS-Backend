package database

import (
	"fmt"
	"strings"
	"time"

	"hotelmanagementsystem.com/v2/models"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(driver string, dsn string, sqlitePath string) (*gorm.DB, error) {
	switch strings.ToLower(strings.TrimSpace(driver)) {
	case "", "sqlite":
		return gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	case "postgres", "postgresql":
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}
}

func Seed(db *gorm.DB) error {
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		admin := models.User{
			Name:         "Front Office Manager",
			Email:        "admin@aurorapalms.com",
			PasswordHash: string(hash),
			Role:         models.RoleManager,
		}

		if err := db.Create(&admin).Error; err != nil {
			return err
		}
	}

	var roomCount int64
	if err := db.Model(&models.Room{}).Count(&roomCount).Error; err != nil {
		return err
	}

	if roomCount == 0 {
		rooms := []models.Room{
			{Number: "101", Type: "Deluxe King", Floor: 1, Rate: 145, Capacity: 2, Status: models.RoomStatusOccupied},
			{Number: "102", Type: "Twin Suite", Floor: 1, Rate: 165, Capacity: 3, Status: models.RoomStatusAvailable},
			{Number: "201", Type: "Garden Villa", Floor: 2, Rate: 220, Capacity: 4, Status: models.RoomStatusCleaning},
			{Number: "202", Type: "Executive Loft", Floor: 2, Rate: 245, Capacity: 2, Status: models.RoomStatusReserved},
			{Number: "301", Type: "Family Residence", Floor: 3, Rate: 280, Capacity: 5, Status: models.RoomStatusAvailable},
			{Number: "302", Type: "Panorama Suite", Floor: 3, Rate: 310, Capacity: 2, Status: models.RoomStatusMaintenance},
		}

		if err := db.Create(&rooms).Error; err != nil {
			return err
		}
	}

	var guestCount int64
	if err := db.Model(&models.Guest{}).Count(&guestCount).Error; err != nil {
		return err
	}

	if guestCount == 0 {
		guests := []models.Guest{
			{FirstName: "Amina", LastName: "Nabwire", Email: "amina@example.com", Phone: "+256700000101", Nationality: "Ugandan"},
			{FirstName: "Daniel", LastName: "Kimani", Email: "daniel@example.com", Phone: "+254711111202", Nationality: "Kenyan"},
			{FirstName: "Grace", LastName: "Mensah", Email: "grace@example.com", Phone: "+233200000303", Nationality: "Ghanaian"},
		}

		if err := db.Create(&guests).Error; err != nil {
			return err
		}
	}

	var bookingCount int64
	if err := db.Model(&models.Booking{}).Count(&bookingCount).Error; err != nil {
		return err
	}

	if bookingCount == 0 {
		var guests []models.Guest
		var rooms []models.Room
		if err := db.Find(&guests).Error; err != nil {
			return err
		}
		if err := db.Find(&rooms).Error; err != nil {
			return err
		}

		now := time.Now()
		bookings := []models.Booking{
			{
				GuestID:      guests[0].ID,
				RoomID:       rooms[0].ID,
				CheckInDate:  now.AddDate(0, 0, -1),
				CheckOutDate: now.AddDate(0, 0, 2),
				Status:       models.BookingStatusCheckedIn,
				GuestsCount:  2,
				TotalAmount:  435,
				Source:       "Website",
				SpecialNote:  "Late arrival",
			},
			{
				GuestID:      guests[1].ID,
				RoomID:       rooms[3].ID,
				CheckInDate:  now.AddDate(0, 0, 1),
				CheckOutDate: now.AddDate(0, 0, 4),
				Status:       models.BookingStatusReserved,
				GuestsCount:  1,
				TotalAmount:  735,
				Source:       "Walk-in",
				SpecialNote:  "Airport pickup",
			},
			{
				GuestID:      guests[2].ID,
				RoomID:       rooms[4].ID,
				CheckInDate:  now.AddDate(0, 0, 3),
				CheckOutDate: now.AddDate(0, 0, 6),
				Status:       models.BookingStatusReserved,
				GuestsCount:  4,
				TotalAmount:  840,
				Source:       "Travel Agent",
				SpecialNote:  "Family cot",
			},
		}

		if err := db.Create(&bookings).Error; err != nil {
			return err
		}
	}

	var settingsCount int64
	if err := db.Model(&models.HotelSettings{}).Count(&settingsCount).Error; err != nil {
		return err
	}

	if settingsCount == 0 {
		settings := models.HotelSettings{
			HotelName: "Aurora Palms Hotel",
		}

		if err := db.Create(&settings).Error; err != nil {
			return err
		}
	}

	var categoryCount int64
	if err := db.Model(&models.IncomeCategory{}).Count(&categoryCount).Error; err != nil {
		return err
	}

	if categoryCount == 0 {
		categories := []models.IncomeCategory{
			{Name: "Room stay", Type: models.IncomeCategoryTypeHotel, Description: "Guest room accommodation"},
			{Name: "Walk-in stay", Type: models.IncomeCategoryTypeHotel, Description: "Direct check-in revenue"},
			{Name: "Conference", Type: models.IncomeCategoryTypeOther, Description: "Conference hall income"},
			{Name: "Restaurant", Type: models.IncomeCategoryTypeOther, Description: "Food and beverage revenue"},
			{Name: "Transport", Type: models.IncomeCategoryTypeOther, Description: "Transfer and shuttle charges"},
		}

		if err := db.Create(&categories).Error; err != nil {
			return err
		}
	}

	var incomeCount int64
	if err := db.Model(&models.IncomeRecord{}).Count(&incomeCount).Error; err != nil {
		return err
	}

	if incomeCount == 0 {
		var categories []models.IncomeCategory
		var users []models.User
		var bookings []models.Booking
		if err := db.Find(&categories).Error; err != nil {
			return err
		}
		if err := db.Find(&users).Error; err != nil {
			return err
		}
		if err := db.Preload("Guest").Find(&bookings).Error; err != nil {
			return err
		}

		if len(categories) > 0 && len(users) > 0 {
			records := []models.IncomeRecord{
				{
					Type:            models.IncomeCategoryTypeOther,
					Title:           "Conference hall rental",
					CategoryID:      categories[2].ID,
					GuestName:       "Lakeside Events",
					Amount:          620,
					Notes:           "Half-day booking",
					RecordedAt:      time.Now().AddDate(0, 0, -3),
					ReceiptNumber:   "RCPT-OTH-0001",
					CreatedByUserID: users[0].ID,
				},
			}

			if len(bookings) > 0 {
				guestName := bookings[0].Guest.FirstName + " " + bookings[0].Guest.LastName
				records = append(records, models.IncomeRecord{
					Type:            models.IncomeCategoryTypeHotel,
					Title:           "Stay receipt",
					CategoryID:      categories[0].ID,
					BookingID:       &bookings[0].ID,
					GuestName:       guestName,
					Amount:          bookings[0].TotalAmount,
					Notes:           "Seeded from booking",
					RecordedAt:      bookings[0].CheckInDate,
					ReceiptNumber:   "RCPT-HTL-0001",
					CreatedByUserID: users[0].ID,
				})
			}

			if err := db.Create(&records).Error; err != nil {
				return err
			}
		}
	}

	var expenditureCount int64
	if err := db.Model(&models.ExpenditureRecord{}).Count(&expenditureCount).Error; err != nil {
		return err
	}

	if expenditureCount == 0 {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			return err
		}

		if len(users) > 0 {
			records := []models.ExpenditureRecord{
				{
					Title:           "Kitchen restock",
					Category:        "supplies",
					Vendor:          "Fresh Valley Foods",
					Amount:          290,
					Notes:           "Breakfast produce",
					RecordedAt:      time.Now().AddDate(0, 0, -3),
					ReceiptNumber:   "RCPT-EXP-0001",
					CreatedByUserID: users[0].ID,
				},
				{
					Title:           "Generator servicing",
					Category:        "maintenance",
					Vendor:          "PowerCore Systems",
					Amount:          410,
					Notes:           "Quarterly service",
					RecordedAt:      time.Now().AddDate(0, 0, -5),
					ReceiptNumber:   "RCPT-EXP-0002",
					CreatedByUserID: users[0].ID,
				},
			}

			if err := db.Create(&records).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
