package customers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

// ErrNotFound возвращается, когда покупатель не найден.
var ErrNotFound = errors.New("item not found")

// ErrInternal возвращается, когда произошла внутренняя ошибка.
var ErrInternal = errors.New("internal error")

// Service описывает сервис работы с покупателями.
type Service struct {
	db *sql.DB
}

// NewService создает сервис.
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// Customer представляет информацию о покупателе.
type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

// ByID возвращает покупателя по идентификатору.
func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.db.QueryRowContext(ctx, `
	  Select id, name, phone, active, created FROM customers WHERE id = $1
	  `, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}
