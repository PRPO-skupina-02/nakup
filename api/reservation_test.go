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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestReservationsList(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	validator := services.NewMockTimeSlotValidator()
	r := TestingRouter(t, db, validator)

	tests := []struct {
		name   string
		status int
		params string
	}{
		{
			name:   "ok",
			status: http.StatusOK,
		},
		{
			name:   "ok-paginated",
			status: http.StatusOK,
			params: "?limit=1&offset=1",
		},
		{
			name:   "ok-sort",
			status: http.StatusOK,
			params: "?sort=-updated_at",
		},
		{
			name:   "ok-paginated-sort",
			status: http.StatusOK,
			params: "?limit=2&offset=1&sort=updated_at",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/nakup/reservations%s", testCase.params)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodGet, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
		})
	}
}

func TestReservationsCreate(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	validator := services.NewMockTimeSlotValidator()
	r := TestingRouter(t, db, validator)

	theaterID := uuid.MustParse("bae209f6-d059-11f0-b2a4-cbf992c2eb6d")
	roomID := uuid.MustParse("925c2358-df46-11f0-a38e-abe580bde3d1")
	timeSlotID := uuid.MustParse("9d71d7fd-d88e-41a1-86dc-21b7f2550295")

	validator.AddValidTimeSlot(theaterID, roomID, timeSlotID)

	tests := []struct {
		name   string
		body   ReservationRequest
		status int
	}{
		{
			name: "ok",
			body: ReservationRequest{
				TimeSlotID: timeSlotID,
				TheaterID:  theaterID,
				RoomID:     roomID,
				Type:       models.Online,
				Row:        7,
				Col:        12,
			},
			status: http.StatusCreated,
		},
		{
			name: "validation-errors",
			body: ReservationRequest{
				TimeSlotID: timeSlotID,
				TheaterID:  theaterID,
				RoomID:     roomID,
				Type:       "INVALID",
				Row:        0,
				Col:        -5,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "invalid-timeslot",
			body: ReservationRequest{
				TimeSlotID: uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
				TheaterID:  theaterID,
				RoomID:     roomID,
				Type:       models.Online,
				Row:        5,
				Col:        10,
			},
			status: http.StatusNotFound,
		},
		{
			name:   "no-body",
			status: http.StatusBadRequest,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := "/api/v1/nakup/reservations"

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodPost, testCase.body)
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			ignoreResp := xtesting.ValuesCheckers{
				"id":         xtesting.ValueUUID(),
				"created_at": xtesting.ValueTimeInPastDuration(time.Second),
				"updated_at": xtesting.ValueTimeInPastDuration(time.Second),
			}

			ignoreReservations := xtesting.GenerateValueCheckersForArrays(map[string]xtesting.ValueChecker{"ID": xtesting.ValueUUID(), "CreatedAt": xtesting.ValueTime(), "UpdatedAt": xtesting.ValueTime()}, 10)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w, ignoreResp)
			xtesting.AssertGoldenDatabaseTable(t, db.Order("time_slot_id"), []models.Reservation{}, ignoreReservations)
		})
	}
}

func TestReservationsShow(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	validator := services.NewMockTimeSlotValidator()
	r := TestingRouter(t, db, validator)

	tests := []struct {
		name   string
		status int
		id     string
	}{
		{
			name:   "ok",
			status: http.StatusOK,
			id:     "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:   "invalid-id",
			status: http.StatusNotFound,
			id:     "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name:   "nil-id",
			status: http.StatusBadRequest,
			id:     "00000000-0000-0000-0000-000000000000",
		},
		{
			name:   "malformed-id",
			status: http.StatusBadRequest,
			id:     "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/nakup/reservations/%s", testCase.id)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodGet, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
		})
	}
}

func TestReservationsUpdate(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	validator := services.NewMockTimeSlotValidator()
	r := TestingRouter(t, db, validator)

	theaterID := uuid.MustParse("bae209f6-d059-11f0-b2a4-cbf992c2eb6d")
	roomID := uuid.MustParse("925c2358-df46-11f0-a38e-abe580bde3d1")
	timeSlotID1 := uuid.MustParse("eed99bc8-1fb4-443b-8287-a988a3bc4406")
	timeSlotID2 := uuid.MustParse("9d71d7fd-d88e-41a1-86dc-21b7f2550295")

	validator.AddValidTimeSlot(theaterID, roomID, timeSlotID1)
	validator.AddValidTimeSlot(theaterID, roomID, timeSlotID2)

	tests := []struct {
		name   string
		body   ReservationRequest
		status int
		id     string
	}{
		{
			name: "ok",
			body: ReservationRequest{
				TimeSlotID: timeSlotID1,
				TheaterID:  theaterID,
				RoomID:     roomID,
				Type:       models.Pos,
				Row:        10,
				Col:        15,
			},
			status: http.StatusOK,
			id:     "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "validation-errors",
			body: ReservationRequest{
				TimeSlotID: timeSlotID2,
				TheaterID:  theaterID,
				RoomID:     roomID,
				Type:       "INVALID",
				Row:        0,
				Col:        -1,
			},
			status: http.StatusBadRequest,
			id:     "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:   "no-body",
			status: http.StatusBadRequest,
			id:     "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name: "invalid-id",
			body: ReservationRequest{
				TimeSlotID: timeSlotID2,
				TheaterID:  theaterID,
				RoomID:     roomID,
				Type:       models.Online,
				Row:        5,
				Col:        8,
			},
			status: http.StatusNotFound,
			id:     "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name: "nil-id",
			body: ReservationRequest{
				TimeSlotID: timeSlotID2,
				TheaterID:  theaterID,
				RoomID:     roomID,
				Type:       models.Online,
				Row:        5,
				Col:        8,
			},
			status: http.StatusBadRequest,
			id:     "00000000-0000-0000-0000-000000000000",
		},
		{
			name: "malformed-id",
			body: ReservationRequest{
				TimeSlotID: timeSlotID2,
				TheaterID:  theaterID,
				RoomID:     roomID,
				Type:       models.Online,
				Row:        5,
				Col:        8,
			},
			status: http.StatusBadRequest,
			id:     "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/nakup/reservations/%s", testCase.id)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodPut, testCase.body)
			assert.NoError(t, err)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			ignoreResp := xtesting.ValuesCheckers{
				"updated_at": xtesting.ValueTimeInPastDuration(time.Second),
			}

			ignoreReservations := xtesting.GenerateValueCheckersForArrays(map[string]xtesting.ValueChecker{"UpdatedAt": xtesting.ValueTime()}, 10)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w, ignoreResp)
			xtesting.AssertGoldenDatabaseTable(t, db, []models.Reservation{}, ignoreReservations)
		})
	}
}

func TestReservationsDelete(t *testing.T) {
	db, fixtures := database.PrepareTestDatabase(t, db.FixtureFS, db.MigrationsFS)
	validator := services.NewMockTimeSlotValidator()
	r := TestingRouter(t, db, validator)

	tests := []struct {
		name   string
		status int
		id     string
	}{
		{
			name:   "ok",
			status: http.StatusNoContent,
			id:     "fb126c8c-d059-11f0-8fa4-b35f33be83b7",
		},
		{
			name:   "invalid-id",
			status: http.StatusNotFound,
			id:     "01234567-0123-0123-0123-0123456789ab",
		},
		{
			name:   "nil-id",
			status: http.StatusBadRequest,
			id:     "00000000-0000-0000-0000-000000000000",
		},
		{
			name:   "malformed-id",
			status: http.StatusBadRequest,
			id:     "000",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := fixtures.Load()
			assert.NoError(t, err)

			targetURL := fmt.Sprintf("/api/v1/nakup/reservations/%s", testCase.id)

			req := xtesting.NewTestingRequest(t, targetURL, http.MethodDelete, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.status, w.Code)
			xtesting.AssertGoldenJSON(t, w)
			xtesting.AssertGoldenDatabaseTable(t, db, []models.Reservation{}, nil)
		})
	}
}
