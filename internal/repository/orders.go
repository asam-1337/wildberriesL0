package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
	"github.com/asam-1337/wildberriesL0/internal/localErrors"
	"github.com/jackc/pgx/v5"
)

const (
	ordersTable = "orders"
)

type Repository interface {
	Insert(ctx context.Context, order entity.Order) error
	SelectById(ctx context.Context, id string) (entity.Order, error)
	SelectAll(ctx context.Context) ([]entity.Order, error)
}

type orderDto struct {
	Id        string `db:"id"`
	OrderData []byte `db:"order_data"`
}

func dtoFromOrder(order entity.Order) (orderDto, error) {
	data, err := json.Marshal(&order)
	if err != nil {
		return orderDto{}, fmt.Errorf("cant marshal order to json: %s", err.Error())
	}
	dto := orderDto{
		Id:        order.OrderUID,
		OrderData: data,
	}

	return dto, nil
}

func orderFromDto(dto orderDto) (entity.Order, error) {
	order := entity.Order{}
	err := json.Unmarshal(dto.OrderData, &order)
	if err != nil {
		return entity.Order{}, fmt.Errorf("cant unmarshal json to order: %s", err.Error())
	}

	return order, nil
}

type OrdersRepository struct {
	pool Runner
}

func NewOrdersRepository(pool Runner) *OrdersRepository {
	return &OrdersRepository{
		pool: pool,
	}
}

func (r *OrdersRepository) Insert(ctx context.Context, order entity.Order) error {
	dto, err := dtoFromOrder(order)
	if err != nil {
		return err
	}

	sql, values, err := squirrel.Insert(ordersTable).
		Columns("id", "order_data").
		Values(dto.Id, dto.OrderData).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	ct, err := r.pool.Exec(ctx, sql, values...)
	if err != nil {
		if ct.RowsAffected() == 0 {
			return err
		}
		return err
	}

	return err
}

func (r *OrdersRepository) SelectById(ctx context.Context, id string) (entity.Order, error) {
	sql, values, err := squirrel.Select("order_data").
		From(ordersTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return entity.Order{}, err
	}

	row := r.pool.QueryRow(ctx, sql, values...)
	dto := orderDto{}
	err = row.Scan(&dto.OrderData)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{}, localErrors.ErrNotFound
		}
		return entity.Order{}, err
	}

	order, err := orderFromDto(dto)
	if err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (r *OrdersRepository) SelectAll(ctx context.Context) ([]entity.Order, error) {
	sql, values, err := squirrel.Select("order_data").
		From(ordersTable).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sql, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]entity.Order, 0)
	for rows.Next() {
		dto := orderDto{}
		err = rows.Scan(&dto.OrderData)
		if err != nil {
			return nil, err
		}

		order, err := orderFromDto(dto)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	return orders, nil
}
