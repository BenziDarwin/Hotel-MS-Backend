package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
)

const (
	RoomStatusAvailable   = "available"
	RoomStatusOccupied    = "occupied"
	RoomStatusReserved    = "reserved"
	RoomStatusCleaning    = "cleaning"
	RoomStatusMaintenance = "maintenance"
)

const (
	BookingStatusReserved   = "reserved"
	BookingStatusCheckedIn  = "checked_in"
	BookingStatusCheckedOut = "checked_out"
	BookingStatusCancelled  = "cancelled"
)

const (
	IncomeCategoryTypeHotel = "hotel"
	IncomeCategoryTypeOther = "other"
)

type User struct {
	ID           string    `gorm:"primaryKey;size:36" json:"id"`
	Name         string    `gorm:"size:120;not null" json:"name"`
	Email        string    `gorm:"uniqueIndex;size:120;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"size:30;not null;default:manager" json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Guest struct {
	ID          string    `gorm:"primaryKey;size:36" json:"id"`
	FirstName   string    `gorm:"size:80;not null" json:"firstName"`
	LastName    string    `gorm:"size:80;not null" json:"lastName"`
	Email       string    `gorm:"uniqueIndex;size:120;not null" json:"email"`
	Phone       string    `gorm:"size:40;not null" json:"phone"`
	Nationality string    `gorm:"size:80" json:"nationality"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Room struct {
	ID        string    `gorm:"primaryKey;size:36" json:"id"`
	Number    string    `gorm:"uniqueIndex;size:20;not null" json:"number"`
	Type      string    `gorm:"size:80;not null" json:"type"`
	Floor     int       `gorm:"not null" json:"floor"`
	Rate      float64   `gorm:"type:decimal(10,2);not null" json:"rate"`
	Capacity  int       `gorm:"not null" json:"capacity"`
	Status    string    `gorm:"size:30;not null;default:available" json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Booking struct {
	ID           string    `gorm:"primaryKey;size:36" json:"id"`
	GuestID      string    `gorm:"size:36;not null" json:"guestId"`
	RoomID       string    `gorm:"size:36;not null" json:"roomId"`
	Guest        Guest     `json:"guest"`
	Room         Room      `json:"room"`
	CheckInDate  time.Time `gorm:"not null" json:"checkInDate"`
	CheckOutDate time.Time `gorm:"not null" json:"checkOutDate"`
	Status       string    `gorm:"size:30;not null;default:reserved" json:"status"`
	GuestsCount  int       `gorm:"not null" json:"guestsCount"`
	TotalAmount  float64   `gorm:"type:decimal(10,2);not null" json:"totalAmount"`
	Source       string    `gorm:"size:80" json:"source"`
	SpecialNote  string    `gorm:"size:255" json:"specialNote"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type HotelSettings struct {
	ID         string    `gorm:"primaryKey;size:36" json:"id"`
	HotelName  string    `gorm:"size:160;not null" json:"hotelName"`
	HotelImage string    `gorm:"size:255" json:"hotelImage"`
	Currency   string    `gorm:"size:12;not null;default:USD" json:"currency"`
	Phone      string    `gorm:"size:60" json:"phone"`
	Email      string    `gorm:"size:120" json:"email"`
	Address    string    `gorm:"size:255" json:"address"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type IncomeCategory struct {
	ID          string    `gorm:"primaryKey;size:36" json:"id"`
	Name        string    `gorm:"size:120;not null" json:"name"`
	Type        string    `gorm:"size:20;not null;index" json:"type"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type IncomeRecord struct {
	ID              string         `gorm:"primaryKey;size:36" json:"id"`
	Type            string         `gorm:"size:20;not null;index" json:"type"`
	Title           string         `gorm:"size:160;not null" json:"title"`
	CategoryID      string         `gorm:"size:36;not null" json:"categoryId"`
	Category        IncomeCategory `json:"category"`
	BookingID       *string        `gorm:"size:36" json:"bookingId"`
	Booking         *Booking       `json:"booking,omitempty"`
	GuestName       string         `gorm:"size:160" json:"guestName"`
	Amount          float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Notes           string         `gorm:"size:255" json:"notes"`
	RecordedAt      time.Time      `gorm:"not null;index" json:"recordedAt"`
	ReceiptNumber   string         `gorm:"size:40;not null;uniqueIndex" json:"receiptNumber"`
	CreatedByUserID string         `gorm:"size:36;not null" json:"createdByUserId"`
	CreatedByUser   User           `gorm:"foreignKey:CreatedByUserID" json:"createdByUser"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

type ExpenditureRecord struct {
	ID              string    `gorm:"primaryKey;size:36" json:"id"`
	Title           string    `gorm:"size:160;not null" json:"title"`
	Category        string    `gorm:"size:60;not null;index" json:"category"`
	Vendor          string    `gorm:"size:160" json:"vendor"`
	Amount          float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Notes           string    `gorm:"size:255" json:"notes"`
	RecordedAt      time.Time `gorm:"not null;index" json:"recordedAt"`
	ReceiptNumber   string    `gorm:"size:40;not null;uniqueIndex" json:"receiptNumber"`
	CreatedByUserID string    `gorm:"size:36;not null" json:"createdByUserId"`
	CreatedByUser   User      `gorm:"foreignKey:CreatedByUserID" json:"createdByUser"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	u.ID = uuid.NewString()
	return nil
}

func (g *Guest) BeforeCreate(_ *gorm.DB) error {
	g.ID = uuid.NewString()
	return nil
}

func (r *Room) BeforeCreate(_ *gorm.DB) error {
	r.ID = uuid.NewString()
	return nil
}

func (b *Booking) BeforeCreate(_ *gorm.DB) error {
	b.ID = uuid.NewString()
	return nil
}

func (s *HotelSettings) BeforeCreate(_ *gorm.DB) error {
	s.ID = uuid.NewString()
	return nil
}

func (c *IncomeCategory) BeforeCreate(_ *gorm.DB) error {
	c.ID = uuid.NewString()
	return nil
}

func (r *IncomeRecord) BeforeCreate(_ *gorm.DB) error {
	r.ID = uuid.NewString()
	return nil
}

func (e *ExpenditureRecord) BeforeCreate(_ *gorm.DB) error {
	e.ID = uuid.NewString()
	return nil
}
