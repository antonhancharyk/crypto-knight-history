package kline

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	apperr "github.com/antonhancharyk/crypto-knight-history/internal/errors"
	"github.com/antonhancharyk/crypto-knight-history/internal/entity"
	"github.com/jmoiron/sqlx"
)

func TestGetKlines_InvalidFrom_ReturnsBadRequest(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := New(sqlxDB)

	_, err = repo.GetKlines(entity.GetKlinesQueryParams{
		From: "not-a-date", To: "2024-01-02 00:00:00", Interval: "1h",
	})
	if err == nil {
		t.Fatal("expected error for invalid from")
	}
	if !errors.Is(err, apperr.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest, got %v", err)
	}
}

func TestGetLastKlineByInterval_NoRows_ReturnsEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT \* FROM klines`).
		WithArgs("1h").
		WillReturnError(sql.ErrNoRows)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := New(sqlxDB)

	k, err := repo.GetLastKlineByInterval("1h")
	if err != nil {
		t.Fatalf("GetLastKlineByInterval: %v", err)
	}
	if k.OpenTime != 0 {
		t.Errorf("expected zero kline, got OpenTime=%d", k.OpenTime)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations: %v", err)
	}
}
