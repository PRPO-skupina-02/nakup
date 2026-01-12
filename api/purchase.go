package api

import (
	"net/http"
	"time"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/PRPO-skupina-02/common/request"
	"github.com/PRPO-skupina-02/nakup/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PurchaseResponse struct {
	ID                uuid.UUID           `json:"id"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	Type              models.PurchaseType `json:"type"`
	Name              string              `json:"name"`
	Count             int                 `json:"count"`
	PricePerItemCents int                 `json:"price_per_item_cents"`
}

func newPurchaseResponse(purchase models.Purchase) PurchaseResponse {
	return PurchaseResponse{
		ID:                purchase.ID,
		CreatedAt:         purchase.CreatedAt,
		UpdatedAt:         purchase.UpdatedAt,
		Type:              purchase.Type,
		Name:              purchase.Name,
		Count:             purchase.Count,
		PricePerItemCents: purchase.PricePerItemCents,
	}
}

// PurchasesList
//
//	@Id				PurchasesList
//	@Summary		List purchases
//	@Description	List purchases
//	@Tags			purchases
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			reservationID	path		string	true	"Reservation ID"				Format(uuid)
//	@Param			limit			query		int		false	"Limit the number of responses"	Default(10)
//	@Param			offset			query		int		false	"Offset the first response"		Default(0)
//	@Param			sort			query		string	false	"Sort results"
//	@Success		200				{object}	request.PaginatedResponse{data=[]PurchaseResponse}
//	@Failure		400				{object}	middleware.HttpError
//	@Failure		404				{object}	middleware.HttpError
//	@Failure		500				{object}	middleware.HttpError
//	@Router			/reservations/{reservationID}/purchases [get]
func PurchasesList(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	reservation := GetContextReservation(c)
	pagination := request.GetNormalizedPaginationArgs(c)
	sort := request.GetSortOptions(c)

	purchases, total, err := models.GetReservationPurchases(tx, reservation.ID, pagination, sort)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := []PurchaseResponse{}

	for _, purchase := range purchases {
		response = append(response, newPurchaseResponse(purchase))
	}

	request.RenderPaginatedResponse(c, response, total)
}

type PurchaseRequest struct {
	Type              string `json:"type" binding:"required,oneof=FOOD DRINK SNACK" enums:"FOOD,DRINK,SNACK"`
	Name              string `json:"name" binding:"required,min=3"`
	Count             int    `json:"count" binding:"required,min=1"`
	PricePerItemCents int    `json:"price_per_item_cents" binding:"required,min=0"`
}

// PurchasesCreate
//
//	@Id				PurchasesCreate
//	@Summary		Create purchase
//	@Description	Create purchase
//	@Tags			purchases
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			reservationID	path		string			true	"Reservation ID"	Format(uuid)
//	@Param			request			body		PurchaseRequest	true	"request body"
//	@Success		200				{object}	PurchaseResponse
//	@Failure		400				{object}	middleware.HttpError
//	@Failure		404				{object}	middleware.HttpError
//	@Failure		500				{object}	middleware.HttpError
//	@Router			/reservations/{reservationID}/purchases [post]
func PurchasesCreate(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	reservation := GetContextReservation(c)

	var req PurchaseRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	purchase := models.Purchase{
		ID:                uuid.New(),
		ReservationID:     reservation.ID,
		Type:              models.PurchaseType(req.Type),
		Name:              req.Name,
		Count:             req.Count,
		PricePerItemCents: req.PricePerItemCents,
	}

	err = purchase.Create(tx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, newPurchaseResponse(purchase))
}

// PurchasesShow
//
//	@Id				PurchasesShow
//	@Summary		Show purchase
//	@Description	Show purchase
//	@Tags			purchases
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			reservationID	path		string	true	"Reservation ID"	Format(uuid)
//	@Param			purchaseID		path		string	true	"Purchase ID"		Format(uuid)
//	@Success		200				{object}	PurchaseResponse
//	@Failure		400				{object}	middleware.HttpError
//	@Failure		404				{object}	middleware.HttpError
//	@Failure		500				{object}	middleware.HttpError
//	@Router			/reservations/{reservationID}/purchases/{purchaseID} [get]
func PurchasesShow(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	reservation := GetContextReservation(c)
	id, err := request.GetUUIDParam(c, "purchaseID")
	if err != nil {
		_ = c.Error(err)
		return
	}

	purchase, err := models.GetPurchase(tx, reservation.ID, id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, newPurchaseResponse(purchase))
}

// PurchasesUpdate
//
//	@Id				PurchasesUpdate
//	@Summary		Update purchase
//	@Description	Update purchase
//	@Tags			purchases
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			reservationID	path		string			true	"Reservation ID"	Format(uuid)
//	@Param			purchaseID		path		string			true	"Purchase ID"		Format(uuid)
//	@Param			request			body		PurchaseRequest	true	"request body"
//	@Success		200				{object}	PurchaseResponse
//	@Failure		400				{object}	middleware.HttpError
//	@Failure		404				{object}	middleware.HttpError
//	@Failure		500				{object}	middleware.HttpError
//	@Router			/reservations/{reservationID}/purchases/{purchaseID} [put]
func PurchasesUpdate(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	reservation := GetContextReservation(c)
	id, err := request.GetUUIDParam(c, "purchaseID")
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req PurchaseRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	purchase, err := models.GetPurchase(tx, reservation.ID, id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	purchase.Type = models.PurchaseType(req.Type)
	purchase.Name = req.Name
	purchase.Count = req.Count
	purchase.PricePerItemCents = req.PricePerItemCents

	err = purchase.Save(tx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, newPurchaseResponse(purchase))
}

// PurchasesDelete
//
//	@Id				PurchasesDelete
//	@Summary		Delete purchase
//	@Description	Delete purchase
//	@Tags			purchases
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			reservationID	path	string	true	"Reservation ID"	Format(uuid)
//	@Param			purchaseID		path	string	true	"Purchase ID"		Format(uuid)
//	@Success		204
//	@Failure		400	{object}	middleware.HttpError
//	@Failure		404	{object}	middleware.HttpError
//	@Failure		500	{object}	middleware.HttpError
//	@Router			/reservations/{reservationID}/purchases/{purchaseID} [delete]
func PurchasesDelete(c *gin.Context) {
	tx := middleware.GetContextTransaction(c)
	reservation := GetContextReservation(c)
	id, err := request.GetUUIDParam(c, "purchaseID")
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = models.DeletePurchase(tx, reservation.ID, id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, "")
}
