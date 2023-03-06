package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
	"github.com/asam-1337/wildberriesL0/internal/localErrors"
	"github.com/asam-1337/wildberriesL0/internal/pgxbalancer"
	"github.com/jackc/pgx/v5"
)

const (
	ordersTable = "orders"
)

type orderDto struct {
	Id        string `db:"id"`
	OrderData []byte `db:"order_data"`
}

func getOrdersDto(order entity.Order) (orderDto, error) {
	data, err := json.Marshal(&order)
	if err != nil {
		return orderDto{}, err
	}
	dto := orderDto{
		Id:        order.OrderUID,
		OrderData: data,
	}

	return dto, nil
}

type OrdersRepository struct {
	pgxbalancer.TransactionBalancer
}

func NewOrdersRepository(balancer pgxbalancer.TransactionBalancer) *OrdersRepository {
	return &OrdersRepository{
		balancer,
	}
}

func (r *OrdersRepository) GetRunner(ctx context.Context) pgxbalancer.Runner {
	return r.TransactionBalancer.GetRunner(ctx)
}

func (r *OrdersRepository) Insert(ctx context.Context, order entity.Order) error {
	dto, err := getOrdersDto(order)
	if err != nil {
		return err
	}

	runner := r.GetRunner(ctx)
	sql, values, err := squirrel.Insert(ordersTable).
		Columns("id", "order_data").
		Values(dto.Id, dto.OrderData).
		ToSql()
	ct, err := runner.Exec(ctx, sql, values)
	if ct.RowsAffected() == 0 {
		return localErrors.ErrAlreadyExists
	}

	return err
}

func (r *OrdersRepository) SelectById(ctx context.Context, id string) (entity.Order, error) {
	sql, values, err := squirrel.Select("id", "order_data").
		From(ordersTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return entity.Order{}, nil
	}

	runner := r.GetRunner(ctx)
	row := runner.QueryRow(ctx, sql, values...)

	dto := orderDto{}
	err = row.Scan(&dto.OrderData)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Order{}, localErrors.ErrNotFound
	}

	return dto, nil
}