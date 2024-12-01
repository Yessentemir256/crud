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

// GetAll возвращает список всех
func (s *Service) GetAll(ctx context.Context) ([]*Customer, error) {
	var customers []*Customer

	rows, err := s.db.QueryContext(ctx, `SELECT * FROM customers;`)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		customer := &Customer{}
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)
		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

// GetAllActive возвращает список всех активных
func (s *Service) GetAllActive(ctx context.Context) ([]*Customer, error) {
	var customers []*Customer

	rows, err := s.db.QueryContext(ctx, `SELECT * FROM customers WHERE active = true;`)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		customer := &Customer{}
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created); err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
		if !customer.Active {
			// Пропускаем неактивных клиентов
			continue
		}
		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return customers, nil
}
