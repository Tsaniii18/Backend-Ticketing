package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    UserID                    string    `gorm:"primaryKey;size:191" json:"user_id"`
    Username                  string    `gorm:"unique;not null;size:100" json:"username"`
    Name                      string    `gorm:"not null" json:"name"`
    Email                     string    `gorm:"unique;not null;size:150" json:"email"`
    Password                  string    `gorm:"not null" json:"-"`
    Role                      string    `gorm:"not null;size:50" json:"role"`
    ProfilePic                string    `gorm:"type:text" json:"profile_pic"`
    Organization              *string   `gorm:"size:200" json:"organization,omitempty"`
    OrganizationType          *string   `gorm:"size:100" json:"organization_type,omitempty"`
    OrganizationDescription   *string   `gorm:"type:text" json:"organization_description,omitempty"`
    KTP                       *string   `gorm:"size:50" json:"ktp,omitempty"`
    RegisterStatus            string    `gorm:"default:pending;size:50" json:"register_status"`
    RefreshToken              *string   `gorm:"type:text" json:"-"`
    AccessToken               *string   `gorm:"type:text" json:"-"`
    CreatedAt                 time.Time `json:"created_at"`
    UpdatedAt                 time.Time `json:"updated_at"`
}

type Event struct {
    EventID          string    `gorm:"primaryKey;size:191" json:"event_id"`
    Name             string    `gorm:"not null;size:200" json:"name"`
    OwnerID          string    `gorm:"not null;size:191" json:"owner_id"`
    Status           string    `gorm:"default:pending;size:50" json:"status"`
    ApprovalComment  *string   `gorm:"type:text" json:"approval_comment"`
    DateStart        time.Time `gorm:"not null" json:"date_start"`
    DateEnd          time.Time `gorm:"not null" json:"date_end"`
    Location         string    `gorm:"not null" json:"location"`
    Description      string    `gorm:"type:text" json:"description"`
    Image            *string   `gorm:"type:text" json:"image"`
    Flyer            *string   `gorm:"type:text" json:"flyer"`
    Category         string    `gorm:"size:100" json:"category"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

type TicketCategory struct {
    TicketCategoryID string    `gorm:"primaryKey;size:191" json:"ticket_category_id"`
    EventID          string    `gorm:"not null;size:191" json:"event_id"`
    Price            float64   `gorm:"not null" json:"price"`
    Quota            int       `gorm:"not null" json:"quota"`
    Sold             int       `gorm:"default:0" json:"sold"`
    Description      string    `gorm:"type:text" json:"description"`
    DateStart        time.Time `gorm:"not null" json:"date_start"`
    DateEnd          time.Time `gorm:"not null" json:"date_end"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

type Report struct {
    ReportID        string    `gorm:"primaryKey;size:191" json:"report_id"`
    EventID         string    `gorm:"not null;size:191" json:"event_id"`
    OwnerID         string    `gorm:"not null;size:191" json:"owner_id"`
    TotalAttendant  int       `gorm:"default:0" json:"total_attendant"`
    TotalSales      float64   `gorm:"default:0" json:"total_sales"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

type TransactionHistory struct {
    TransactionID   string    `gorm:"primaryKey;size:191" json:"transaction_id"`
    OwnerID         string    `gorm:"not null;size:191" json:"owner_id"`
    EventID         string    `gorm:"not null;size:191" json:"event_id"`
    TransactionTime time.Time `json:"transaction_time"`
    TotalAmount     float64   `gorm:"not null" json:"total_amount"`
    Status          string    `gorm:"default:completed;size:50" json:"status"`
    CreatedAt       time.Time `json:"created_at"`
}

type Ticket struct {
    TicketID         string    `gorm:"primaryKey;size:191" json:"ticket_id"`
    EventID          string    `gorm:"not null;size:191" json:"event_id"`
    TicketCategoryID string    `gorm:"not null;size:191" json:"ticket_category_id"`
    OwnerID          string    `gorm:"not null;size:191" json:"owner_id"`
    Status           string    `gorm:"default:active;size:50" json:"status"`
    Code             string    `gorm:"unique;not null;size:255" json:"code"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

type Cart struct {
    CartID           string    `gorm:"primaryKey;size:191" json:"cart_id"`
    UserID           string    `gorm:"not null;size:191" json:"user_id"`
    TicketCategoryID string    `gorm:"not null;size:191" json:"ticket_category_id"`
    Quantity         int       `gorm:"not null" json:"quantity"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
    if user.UserID == "" {
        user.UserID = uuid.New().String()
    }
    return nil
}

func (event *Event) BeforeCreate(tx *gorm.DB) error {
    if event.EventID == "" {
        event.EventID = uuid.New().String()
    }
    return nil
}

func (tc *TicketCategory) BeforeCreate(tx *gorm.DB) error {
    if tc.TicketCategoryID == "" {
        tc.TicketCategoryID = uuid.New().String()
    }
    return nil
}

func (report *Report) BeforeCreate(tx *gorm.DB) error {
    if report.ReportID == "" {
        report.ReportID = uuid.New().String()
    }
    return nil
}

func (th *TransactionHistory) BeforeCreate(tx *gorm.DB) error {
    if th.TransactionID == "" {
        th.TransactionID = uuid.New().String()
    }
    return nil
}

func (ticket *Ticket) BeforeCreate(tx *gorm.DB) error {
    if ticket.TicketID == "" {
        ticket.TicketID = uuid.New().String()
    }
    return nil
}

func (cart *Cart) BeforeCreate(tx *gorm.DB) error {
    if cart.CartID == "" {
        cart.CartID = uuid.New().String()
    }
    return nil
}