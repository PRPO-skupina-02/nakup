package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PRPO-skupina-02/common/database"
	"github.com/PRPO-skupina-02/common/xtesting"
	"github.com/PRPO-skupina-02/nakup/db"
	"github.com/PRPO-skupina-02/nakup/models"
	"github.com/PRPO-skupina-02/nakup/services"
	"github.com/stretchr/testify/assert"
)

func TestPurchasesList(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	service := services.NewMockTimeSlotService()
	r := TestingRouter(t, db, service)

	tests := []struct {
		name          string
		status        int
		params        string
		reservationID string
	}{
		{
			name:          "ok",
			status:        http.StatusOK,
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "ok-paginated",
			status:        http.StatusOK,
			params:        "?limit=1&offset=1",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "ok-sort",
			status:        http.StatusOK,
			params:        "?sort=-updated_at",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "ok-paginated-sort",
			status:        http.StatusOK,
			params:        "?limit=2&offset=1&sort=updated_at",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "ok-no-purchases",
			status:        http.StatusOK,
			reservationID: "ea0b7f96-ddc9-11f0-9635-23efd36396bd",
		},
		{
			name:          "invalid-reservation-id",
			status:        http.StatusNotFound,
			reservationID: "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name:          "nil-reservation-id",
			status:        http.StatusBadRequest,
			reservationID: "00000000-0000-0000-0000-000000000000",
		},
		{
			name:          "malformed-reservation-id",
			status:        http.StatusBadRequest,
			reservationID: "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/nakup/reservations/%s/purchases%s", testCase.reservationID, testCase.params)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodGet, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
		})
	}
}

func TestPurchasesCreate(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	service := services.NewMockTimeSlotService()
	r := TestingRouter(t, db, service)

	tests := []struct {
		name          string
		body          PurchaseRequest
		status        int
		reservationID string
	}{
		{
			name: "ok",
			body: PurchaseRequest{
				Type:              string(models.Food),
				Name:              "Candy Bar",
				Count:             3,
				PricePerItemCents: 250,
			},
			status:        http.StatusCreated,
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "validation-errors",
			body: PurchaseRequest{
				Type:              "INVALID",
				Name:              "AB",
				Count:             0,
				PricePerItemCents: -100,
			},
			status:        http.StatusBadRequest,
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "no-body",
			status:        http.StatusBadRequest,
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "invalid-reservation-id",
			body: PurchaseRequest{
				Type:              string(models.Drink),
				Name:              "Water",
				Count:             1,
				PricePerItemCents: 200,
			},
			status:        http.StatusNotFound,
			reservationID: "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name: "nil-reservation-id",
			body: PurchaseRequest{
				Type:              string(models.Drink),
				Name:              "Water",
				Count:             1,
				PricePerItemCents: 200,
			},
			status:        http.StatusBadRequest,
			reservationID: "00000000-0000-0000-0000-000000000000",
		},
		{
			name: "malformed-reservation-id",
			body: PurchaseRequest{
				Type:              string(models.Drink),
				Name:              "Water",
				Count:             1,
				PricePerItemCents: 200,
			},
			status:        http.StatusBadRequest,
			reservationID: "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/nakup/reservations/%s/purchases", testCase.reservationID)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodPost, testCase.body)
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			ignoreResp := xtesting.ValuesCheckers{
				"id":         xtesting.ValueUUID(),
				"created_at": xtesting.ValueTimeInPastDuration(time.Second),
				"updated_at": xtesting.ValueTimeInPastDuration(time.Second),
			}

			ignorePurchases := xtesting.GenerateValueCheckersForArrays(map[string]xtesting.ValueChecker{"ID": xtesting.ValueUUID(), "CreatedAt": xtesting.ValueTime(), "UpdatedAt": xtesting.ValueTime()}, 10)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w, ignoreResp)
			xtesting.AssertGoldenDatabaseTable(t, db.Order("name"), []models.Purchase{}, ignorePurchases)
		})
	}
}

func TestPurchasesShow(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	service := services.NewMockTimeSlotService()
	r := TestingRouter(t, db, service)

	tests := []struct {
		name          string
		status        int
		purchaseID    string
		reservationID string
	}{
		{
			name:          "ok",
			status:        http.StatusOK,
			purchaseID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "purchase-from-different-reservation",
			status:        http.StatusNotFound,
			purchaseID:    "dddddddd-dddd-dddd-dddd-dddddddddddd",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "invalid-purchase-id",
			status:        http.StatusNotFound,
			purchaseID:    "01234567-0123-0123-0123-0123456789ab",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "nil-purchase-id",
			status:        http.StatusBadRequest,
			purchaseID:    "00000000-0000-0000-0000-000000000000",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "malformed-purchase-id",
			status:        http.StatusBadRequest,
			purchaseID:    "000",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "invalid-reservation-id",
			status:        http.StatusNotFound,
			purchaseID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			reservationID: "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name:          "nil-reservation-id",
			status:        http.StatusBadRequest,
			purchaseID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			reservationID: "00000000-0000-0000-0000-000000000000",
		},
		{
			name:          "malformed-reservation-id",
			status:        http.StatusBadRequest,
			purchaseID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			reservationID: "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/nakup/reservations/%s/purchases/%s", testCase.reservationID, testCase.purchaseID)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodGet, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
		})
	}
}

func TestPurchasesUpdate(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	service := services.NewMockTimeSlotService()
	r := TestingRouter(t, db, service)

	tests := []struct {
		name          string
		body          PurchaseRequest
		status        int
		purchaseID    string
		reservationID string
	}{
		{
			name: "ok",
			body: PurchaseRequest{
				Type:              string(models.Snack),
				Name:              "Updated Snack",
				Count:             5,
				PricePerItemCents: 300,
			},
			status:        http.StatusOK,
			purchaseID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "validation-errors",
			body: PurchaseRequest{
				Type:              "INVALID",
				Name:              "AB",
				Count:             0,
				PricePerItemCents: -100,
			},
			status:        http.StatusBadRequest,
			purchaseID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "no-body",
			status:        http.StatusBadRequest,
			purchaseID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "purchase-from-different-reservation",
			body: PurchaseRequest{
				Type:              string(models.Food),
				Name:              "Updated Food",
				Count:             2,
				PricePerItemCents: 400,
			},
			status:        http.StatusNotFound,
			purchaseID:    "dddddddd-dddd-dddd-dddd-dddddddddddd",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "invalid-purchase-id",
			body: PurchaseRequest{
				Type:              string(models.Food),
				Name:              "Updated Food",
				Count:             2,
				PricePerItemCents: 400,
			},
			status:        http.StatusNotFound,
			purchaseID:    "01234567-0123-0123-0123-0123456789ab",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "nil-purchase-id",
			body: PurchaseRequest{
				Type:              string(models.Food),
				Name:              "Updated Food",
				Count:             2,
				PricePerItemCents: 400,
			},
			status:        http.StatusBadRequest,
			purchaseID:    "00000000-0000-0000-0000-000000000000",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "malformed-purchase-id",
			body: PurchaseRequest{
				Type:              string(models.Food),
				Name:              "Updated Food",
				Count:             2,
				PricePerItemCents: 400,
			},
			status:        http.StatusBadRequest,
			purchaseID:    "000",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "invalid-reservation-id",
			body: PurchaseRequest{
				Type:              string(models.Food),
				Name:              "Updated Food",
				Count:             2,
				PricePerItemCents: 400,
			},
			status:        http.StatusNotFound,
			purchaseID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			reservationID: "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name: "nil-reservation-id",
			body: PurchaseRequest{
				Type:              string(models.Food),
				Name:              "Updated Food",
				Count:             2,
				PricePerItemCents: 400,
			},
			status:        http.StatusBadRequest,
			purchaseID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			reservationID: "00000000-0000-0000-0000-000000000000",
		},
		{
			name: "malformed-reservation-id",
			body: PurchaseRequest{
				Type:              string(models.Food),
				Name:              "Updated Food",
				Count:             2,
				PricePerItemCents: 400,
			},
			status:        http.StatusBadRequest,
			purchaseID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			reservationID: "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/nakup/reservations/%s/purchases/%s", testCase.reservationID, testCase.purchaseID)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodPut, testCase.body)
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			ignoreResp := xtesting.ValuesCheckers{
				"updated_at": xtesting.ValueTimeInPastDuration(time.Second),
			}

			ignorePurchases := xtesting.GenerateValueCheckersForArrays(map[string]xtesting.ValueChecker{"UpdatedAt": xtesting.ValueTime()}, 10)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w, ignoreResp)
			xtesting.AssertGoldenDatabaseTable(t, db, []models.Purchase{}, ignorePurchases)
		})
	}
}

func TestPurchasesDelete(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	service := services.NewMockTimeSlotService()
	r := TestingRouter(t, db, service)

	tests := []struct {
		name          string
		status        int
		purchaseID    string
		reservationID string
	}{
		{
			name:          "ok",
			status:        http.StatusNoContent,
			purchaseID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "purchase-from-different-reservation",
			status:        http.StatusNotFound,
			purchaseID:    "dddddddd-dddd-dddd-dddd-dddddddddddd",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "invalid-purchase-id",
			status:        http.StatusNotFound,
			purchaseID:    "01234567-0123-0123-0123-0123456789ab",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "nil-purchase-id",
			status:        http.StatusBadRequest,
			purchaseID:    "00000000-0000-0000-0000-000000000000",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "malformed-purchase-id",
			status:        http.StatusBadRequest,
			purchaseID:    "000",
			reservationID: "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:          "invalid-reservation-id",
			status:        http.StatusNotFound,
			purchaseID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			reservationID: "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name:          "nil-reservation-id",
			status:        http.StatusBadRequest,
			purchaseID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			reservationID: "00000000-0000-0000-0000-000000000000",
		},
		{
			name:          "malformed-reservation-id",
			status:        http.StatusBadRequest,
			purchaseID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			reservationID: "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/nakup/reservations/%s/purchases/%s", testCase.reservationID, testCase.purchaseID)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodDelete, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
			xtesting.AssertGoldenDatabaseTable(t, db, []models.Purchase{}, nil)
		})
	}
}
