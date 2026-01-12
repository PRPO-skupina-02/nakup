package models

import (
	"time"

	"github.com/PRPO-skupina-02/common/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReservationType string

const (
	Online ReservationType = "ONLINE"
	Pos    ReservationType = "POS"
)

type Reservation struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	TimeSlotID uuid.UUID
	UserID     uuid.UUID
	Type       ReservationType

	Row int
	Col int

	Purchases []Purchase `gorm:"foreignKey:ReservationID" json:"-"`
}

func (r *Reservation) Create(tx *gorm.DB) error {
	if err := tx.Create(r).Error; err != nil {
		return err
	}
	return nil
}

func (r *Reservation) Save(tx *gorm.DB) error {
	if err := tx.Save(r).Error; err != nil {
		return err
	}
	return nil
}

func GetReservations(tx *gorm.DB, pagination *request.PaginationOptions, sort *request.SortOptions) ([]Reservation, int, error) {
	var reservations []Reservation

	query := tx.Model(&Reservation{}).Session(&gorm.Session{})

	if err := query.Scopes(request.PaginateScope(pagination), request.SortScope(sort)).Find(&reservations).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return reservations, int(total), nil
}

func GetUserReservations(tx *gorm.DB, userID uuid.UUID, pagination *request.PaginationOptions, sort *request.SortOptions) ([]Reservation, int, error) {
	var reservations []Reservation

	query := tx.Model(&Reservation{}).Where("user_id = ?", userID).Session(&gorm.Session{})

	if err := query.Scopes(request.PaginateScope(pagination), request.SortScope(sort)).Find(&reservations).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return reservations, int(total), nil
}

func GetReservation(tx *gorm.DB, id uuid.UUID) (Reservation, error) {
	reservation := Reservation{
		ID: id,
	}

	if err := tx.Where(&reservation).First(&reservation).Error; err != nil {
		return reservation, err
	}

	return reservation, nil
}

func CheckDuplicateReservation(tx *gorm.DB, timeSlotID uuid.UUID, row, col int, excludeID *uuid.UUID) (bool, error) {
	query := tx.Model(&Reservation{}).Where("time_slot_id = ? AND row = ? AND col = ?", timeSlotID, row, col)

	if excludeID != nil {
		query = query.Where("id != ?", excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func DeleteReservation(tx *gorm.DB, id uuid.UUID) error {
	reservation := Reservation{
		ID: id,
	}

	if err := tx.Where(&reservation).Preload("Purchases").First(&reservation).Error; err != nil {
		return err
	}

	for _, purchase := range reservation.Purchases {
		err := DeletePurchase(tx, id, purchase.ID)
		if err != nil {
			return err
		}
	}

	if err := tx.Delete(&reservation).Error; err != nil {
		return err
	}
	return nil
}
