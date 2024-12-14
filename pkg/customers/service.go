package customers

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

// ErrNotFound возвращается, когда покупатель не найден.
var ErrNotFound = errors.New("item not found")

// ErrInternal возвращается, когда произошла внутренняя ошибка.
var ErrInternal = errors.New("internal error")

// ErrNotDeleted возвращает, когда удаление не произошло.
var ErrNotDeleted = errors.New("no rows was deleted")

// Service описывает сервис работы с покупателями.
type Service struct {
	pool *pgxpool.Pool
}

// NewService создает сервис.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
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

	err := s.pool.QueryRow(ctx, `
	  Select id, name, phone, active, created FROM customers WHERE id = $1
	  `, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}

// GetAll возвращает список всех.
func (s *Service) GetAll(ctx context.Context) ([]*Customer, error) {
	var customers []*Customer

	rows, err := s.pool.Query(ctx, `SELECT * FROM customers;`)
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

	rows, err := s.pool.Query(ctx, `SELECT * FROM customers WHERE active = true;`)
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

// Save создает или обновляет.
func (s *Service) Save(ctx context.Context, item *Customer) (*Customer, error) {
	if item.ID == 0 {
		// Создание нового клиента
		_, err := s.pool.Exec(ctx, `INSERT INTO customers (name, phone) VALUES ($1, $2);`, item.Name, item.Phone)
		if err != nil {
			return nil, err
		}
	} else {
		// Обновление существующего клиента
		_, err := s.pool.Exec(ctx, `UPDATE customers SET name = $1, phone = $2 WHERE id = $3;`, item.Name, item.Phone, item.ID)
		if err != nil {
			return nil, err
		}
	}

	// Возвращаем обновленный объект Customer
	return item, nil
}

// RemoveById удаляет пользователя по ID.
func (s *Service) RemoveByID(ctx context.Context, id int) error {
	result, err := s.pool.Exec(ctx, `DELETE FROM customers WHERE id = $1;`, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotDeleted
	}

	return nil
}

// Выставляет статус false по ID .
func (s *Service) BlockByID(ctx context.Context, id int64) error {
	item, err := s.ByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.pool.QueryRow(ctx, `
	  Select id, name, phone, active, created FROM customers WHERE id = $1
	  `, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if err != nil {
		log.Print(err)
		return ErrInternal
	}

	if !item.Active {
		return nil
	}
	item.Active = false

	_, err = s.pool.Exec(ctx, `
	  UPDATE customers SET active = $1 WHERE id = $2
	  `, item.Active, item.ID)

	if err != nil {
		return err
	}

	return nil

}

// Выставляет статус active по ID .
func (s *Service) UnBlockByID(ctx context.Context, id int64) error {
	item, err := s.ByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.pool.QueryRow(ctx, `
	  Select id, name, phone, active, created FROM customers WHERE id = $1
	  `, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if err != nil {
		log.Print(err)
		return ErrInternal
	}

	if item.Active {
		return nil
	}
	item.Active = true

	_, err = s.pool.Exec(ctx, `
	  UPDATE customers SET active = $1 WHERE id = $2
	  `, item.Active, item.ID)

	if err != nil {
		return err
	}

	return nil

}
