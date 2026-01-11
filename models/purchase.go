package models

import (
	"time"

	"github.com/PRPO-skupina-02/common/request"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseType string

const (
	Food  PurchaseType = "FOOD"
	Drink PurchaseType = "DRINK"
	Snack PurchaseType = "SNACK"
)

type Purchase struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ReservationID uuid.UUID

	Type              PurchaseType
	Name              string
	Count             int
	PricePerItemCents int

	Reservation Reservation `gorm:"foreignKey:ReservationID" json:"-"`
}

func (p *Purchase) Create(tx *gorm.DB) error {
	if err := tx.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (p *Purchase) Save(tx *gorm.DB) error {
	if err := tx.Save(p).Error; err != nil {
		return err
	}
	return nil
}

func GetReservationPurchases(tx *gorm.DB, reservationID uuid.UUID, pagination *request.PaginationOptions, sort *request.SortOptions) ([]Purchase, int, error) {
	var purchases []Purchase

	query := tx.Model(&Purchase{}).Where("reservation_id = ?", reservationID).Session(&gorm.Session{})

	if err := query.Scopes(request.PaginateScope(pagination), request.SortScope(sort)).Find(&purchases).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return purchases, int(total), nil
}

func GetPurchase(tx *gorm.DB, reservationID, purchaseID uuid.UUID) (Purchase, error) {
	reservation := Purchase{
		ID:            purchaseID,
		ReservationID: reservationID,
	}

	if err := tx.Where(&reservation).First(&reservation).Error; err != nil {
		return reservation, err
	}

	return reservation, nil
}

func DeletePurchase(tx *gorm.DB, reservationID, purchaseID uuid.UUID) error {
	purchase := Purchase{
		ID:            purchaseID,
		ReservationID: reservationID,
	}

	if err := tx.Where(&purchase).First(&purchase).Error; err != nil {
		return err
	}

	if err := tx.Delete(&purchase).Error; err != nil {
		return err
	}
	return nil
}
