package repository

import (
	"context"
	"fmt"
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
	"github.com/asam-1337/wildberriesL0/internal/localErrors"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var orders = []entity.Order{
	{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: entity.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "publish@gmail.com",
		},
		Payment: entity.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []entity.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "publish",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		OofShard:          "1",
	},
	{
		OrderUID:    "b563feb7b2b84b6sdfs",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: entity.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "publish@gmail.com",
		},
		Payment: entity.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []entity.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "publish",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		OofShard:          "1",
	},
}

func TestOrdersRepository_Insert(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewOrdersRepository(mock)
	dto, err := dtoFromOrder(orders[0])
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec("INSERT INTO orders").
		WithArgs(dto.Id, dto.OrderData).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	err = repo.Insert(context.Background(), orders[0])
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec("INSERT INTO orders").
		WithArgs(dto.Id, dto.OrderData).
		WillReturnError(fmt.Errorf("already exist"))
	err = repo.Insert(context.Background(), orders[0])
	assert.Error(t, err)
}

func TestOrdersRepository_SelectById(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Error(err)
	}
	defer mock.Close()

	repo := NewOrdersRepository(mock)
	rows := mock.NewRows([]string{"order_data"})
	for _, val := range orders {
		dto, err := dtoFromOrder(val)
		if err != nil {
			t.Error(err)
		}
		rows = rows.AddRow(dto.OrderData)
	}

	for _, val := range orders {
		mock.ExpectQuery("SELECT order_data FROM orders WHERE").
			WithArgs(val.OrderUID).
			WillReturnRows(rows)

		order, err := repo.SelectById(context.Background(), val.OrderUID)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(order, val) {
			t.Errorf("results not match,\nwant: %v,\nhave:%v", val, order)
		}
	}

	//query error
	mock.ExpectQuery("SELECT order_data FROM orders WHERE").
		WithArgs("order_uid").
		WillReturnError(fmt.Errorf("bad query"))
	_, err = repo.SelectById(context.Background(), "order_uid")
	assert.Error(t, err)

	//ErrNoRows
	mock.ExpectQuery("SELECT order_data FROM orders WHERE").
		WithArgs("order_uid").
		WillReturnError(pgx.ErrNoRows)
	_, err = repo.SelectById(context.Background(), "order_uid")
	assert.ErrorIs(t, err, localErrors.ErrNotFound)

	//scan error
	rows = pgxmock.NewRows([]string{"id", "order_data"})
	dto, err := dtoFromOrder(orders[0])
	rows.AddRow(dto.Id, dto.OrderData)
	mock.ExpectQuery("SELECT order_data FROM orders WHERE").
		WithArgs("order_uid").
		WillReturnRows(rows)
	_, err = repo.SelectById(context.Background(), "order_uid")
	assert.Error(t, err)
}

func TestOrdersRepository_SelectAll(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Error(err)
	}
	defer mock.Close()

	repo := NewOrdersRepository(mock)
	rows := pgxmock.NewRows([]string{"order_data"})
	for _, val := range orders {
		dto, err := dtoFromOrder(val)
		if err != nil {
			t.Error(err)
		}
		rows = rows.AddRow(dto.OrderData)
	}

	mock.ExpectQuery("SELECT order_data FROM orders").
		WillReturnRows(rows)
	ordrs, err := repo.SelectAll(context.Background())
	if !reflect.DeepEqual(ordrs, orders) {
		t.Errorf("results not match,\nwant: %v,\nhave:%v", orders, ordrs)
	}

	//query error
	mock.ExpectQuery("SELECT order_data FROM orders").
		WithArgs("order_uid").
		WillReturnError(fmt.Errorf("bad query"))
	_, err = repo.SelectAll(context.Background())
	assert.Error(t, err)

	//scan error
	rows = pgxmock.NewRows([]string{"id", "order_data"})
	dto, err := dtoFromOrder(orders[0])
	rows.AddRow(dto.Id, dto.OrderData)
	mock.ExpectQuery("SELECT order_data FROM orders").
		WithArgs("order_uid").
		WillReturnRows(rows)
	_, err = repo.SelectAll(context.Background())
	assert.Error(t, err)
}
