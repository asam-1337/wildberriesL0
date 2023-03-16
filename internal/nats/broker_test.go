package nats

import (
	"context"
	"fmt"
	"github.com/asam-1337/wildberriesL0/config"
	"github.com/asam-1337/wildberriesL0/internal/cache"
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
	repository "github.com/asam-1337/wildberriesL0/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

type validTestCase struct {
	descr  string
	data   []byte
	expect func(t assert.TestingT, value bool, msgAndArgs ...interface{}) bool
}

func TestBroker_validateMessage(t *testing.T) {
	cases := []validTestCase{
		{
			descr:  "valid 1",
			data:   []byte("{\n    \"order_uid\": \"b563feb7b2b84b6test\",\n    \"track_number\": \"WBILMTESTTRACK\",\n    \"entry\": \"WBIL\",\n    \"delivery\": {\n        \"name\": \"Test Testov\",\n        \"phone\": \"+9720000000\",\n        \"zip\": \"2639809\",\n        \"city\": \"Kiryat Mozkin\",\n        \"address\": \"Ploshad Mira 15\",\n        \"region\": \"Kraiot\",\n        \"email\": \"publish@gmail.com\"\n    },\n    \"payment\": {\n        \"transaction\": \"b563feb7b2b84b6\",\n        \"request_id\": \"\",\n        \"currency\": \"USD\",\n        \"provider\": \"wbpay\",\n        \"amount\": 1817,\n        \"payment_dt\": 1637907727,\n        \"bank\": \"alpha\",\n        \"delivery_cost\": 100,\n        \"goods_total\": 1,\n        \"custom_fee\": 0\n    },\n    \"items\": [\n        {\n            \"chrt_id\": 9934930,\n            \"track_number\": \"WBILMTESTTRACK\",\n            \"price\": 453,\n            \"rid\": \"ab4219087a764ae0btest\",\n            \"name\": \"Mascaras\",\n            \"sale\": 30,\n            \"size\": \"0\",\n            \"total_price\": 317,\n            \"nm_id\": 2389212,\n            \"brand\": \"Vivienne Sabo\",\n            \"status\": 202\n        }\n    ],\n    \"locale\": \"en\",\n    \"internal_signature\": \"\",\n    \"customer_id\": \"publish\",\n    \"delivery_service\": \"meest\",\n    \"shardkey\": \"9\",\n    \"sm_id\": 99,\n    \"date_created\": \"2021-11-26T06:22:19Z\",\n    \"oof_shard\": \"1\"\n}"),
			expect: assert.True,
		},
		{
			descr:  "valid 2",
			data:   []byte("{\n    \"order_uid\": \"b563feb7b2b84b651sd\",\n    \"track_number\": \"WBILMTESTTRACK\",\n    \"entry\": \"WBIL\",\n    \"delivery\": {\n        \"name\": \"Test Testov\",\n        \"phone\": \"+9720000000\",\n        \"zip\": \"2639809\",\n        \"city\": \"Kiryat Mozkin\",\n        \"address\": \"Ploshad Mira 15\",\n        \"region\": \"Kraiot\",\n        \"email\": \"publish@gmail.com\"\n    },\n    \"payment\": {\n        \"transaction\": \"b563feb7b2b84b6\",\n        \"request_id\": \"\",\n        \"currency\": \"USD\",\n        \"provider\": \"wbpay\",\n        \"amount\": 1817,\n        \"payment_dt\": 1637907727,\n        \"bank\": \"alpha\",\n        \"delivery_cost\": 100,\n        \"goods_total\": 1,\n        \"custom_fee\": 0\n    },\n    \"items\": [\n        {\n            \"chrt_id\": 9934930,\n            \"track_number\": \"WBILMTESTTRACK\",\n            \"price\": 453,\n            \"rid\": \"ab4219087a764ae0btest\",\n            \"name\": \"Mascaras\",\n            \"sale\": 30,\n            \"size\": \"0\",\n            \"total_price\": 317,\n            \"nm_id\": 2389212,\n            \"brand\": \"Vivienne Sabo\",\n            \"status\": 202\n        }\n    ],\n    \"locale\": \"en\",\n    \"internal_signature\": \"\",\n    \"customer_id\": \"publish\",\n    \"delivery_service\": \"meest\",\n    \"shardkey\": \"9\",\n    \"sm_id\": 99,\n    \"date_created\": \"2021-11-26T06:22:19Z\",\n    \"oof_shard\": \"1\"\n}"),
			expect: assert.True,
		},
		{
			descr:  "invalid json",
			data:   []byte("jfnjanjj';jnajhvakcv;a;avkskmkmd,fsdf"),
			expect: assert.False,
		},
		{
			descr:  "invalid uid",
			data:   []byte("{\n    \"order_uid\": \"b563feb7b2b84b6aaaaaa\",\n    \"track_number\": \"WBILMTESTTRACK\",\n    \"entry\": \"WBIL\",\n    \"delivery\": {\n        \"name\": \"Test Testov\",\n        \"phone\": \"+9720000000\",\n        \"zip\": \"2639809\",\n        \"city\": \"Kiryat Mozkin\",\n        \"address\": \"Ploshad Mira 15\",\n        \"region\": \"Kraiot\",\n        \"email\": \"publish@gmail.com\"\n    },\n    \"payment\": {\n        \"transaction\": \"b563feb7b2b84b6\",\n        \"request_id\": \"\",\n        \"currency\": \"USD\",\n        \"provider\": \"wbpay\",\n        \"amount\": 1817,\n        \"payment_dt\": 1637907727,\n        \"bank\": \"alpha\",\n        \"delivery_cost\": 100,\n        \"goods_total\": 1,\n        \"custom_fee\": 0\n    },\n    \"items\": [\n        {\n            \"chrt_id\": 9934930,\n            \"track_number\": \"WBILMTESTTRACK\",\n            \"price\": 453,\n            \"rid\": \"ab4219087a764ae0btest\",\n            \"name\": \"Mascaras\",\n            \"sale\": 30,\n            \"size\": \"0\",\n            \"total_price\": 317,\n            \"nm_id\": 2389212,\n            \"brand\": \"Vivienne Sabo\",\n            \"status\": 202\n        }\n    ],\n    \"locale\": \"en\",\n    \"internal_signature\": \"\",\n    \"customer_id\": \"publish\",\n    \"delivery_service\": \"meest\",\n    \"shardkey\": \"9\",\n    \"sm_id\": 99,\n    \"date_created\": \"2021-11-26T06:22:19Z\",\n    \"oof_shard\": \"1\"\n}"),
			expect: assert.False,
		},
		{
			descr:  "no order_uid",
			data:   []byte("{\n    \"track_number\": \"WBILMTESTTRACK\",\n    \"entry\": \"WBIL\",\n    \"delivery\": {\n        \"name\": \"Test Testov\",\n        \"phone\": \"+9720000000\",\n        \"zip\": \"2639809\",\n        \"city\": \"Kiryat Mozkin\",\n        \"address\": \"Ploshad Mira 15\",\n        \"region\": \"Kraiot\",\n        \"email\": \"publish@gmail.com\"\n    },\n    \"payment\": {\n        \"transaction\": \"b563feb7b2b84b6\",\n        \"request_id\": \"\",\n        \"currency\": \"USD\",\n        \"provider\": \"wbpay\",\n        \"amount\": 1817,\n        \"payment_dt\": 1637907727,\n        \"bank\": \"alpha\",\n        \"delivery_cost\": 100,\n        \"goods_total\": 1,\n        \"custom_fee\": 0\n    },\n    \"items\": [\n        {\n            \"chrt_id\": 9934930,\n            \"track_number\": \"WBILMTESTTRACK\",\n            \"price\": 453,\n            \"rid\": \"ab4219087a764ae0btest\",\n            \"name\": \"Mascaras\",\n            \"sale\": 30,\n            \"size\": \"0\",\n            \"total_price\": 317,\n            \"nm_id\": 2389212,\n            \"brand\": \"Vivienne Sabo\",\n            \"status\": 202\n        }\n    ],\n    \"locale\": \"en\",\n    \"internal_signature\": \"\",\n    \"customer_id\": \"publish\",\n    \"delivery_service\": \"meest\",\n    \"shardkey\": \"9\",\n    \"sm_id\": 99,\n    \"date_created\": \"2021-11-26T06:22:19Z\",\n    \"oof_shard\": \"1\"\n}"),
			expect: assert.False,
		},
		{
			descr:  "bad fields",
			data:   []byte("{\n    \"orderuid\": \"b563feb7b2b84b6test\",\n    \"tracking\": \"123\",\n    \"entry\": \"WBIL\",\n    \"delivery\": {\n        \"surname\": \"Test Testov\",\n        \"email\": \"+9720000000\",\n        \"code\": \"2639809\",\n        \"city\": \"Kiryat Mozkin\",\n        \"address\": \"Ploshad Mira 15\",\n        \"region\": \"Kraiot\",\n        \"email\": \"publish@gmail.com\"\n    },\n    \"payment\": {\n        \"transaction\": \"b563feb7b2b84b6\",\n        \"request_id\": \"\",\n        \"currency\": \"USD\",\n        \"provider\": \"wbpay\",\n        \"amount\": 1817,\n        \"payment_dt\": 1637907727,\n        \"bank\": \"alpha\",\n        \"delivery_cost\": 100,\n        \"goods_total\": 1,\n        \"custom_fee\": 0\n    },\n    \"items\": [\n        {\n            \"chrt_id\": 9934930,\n            \"track_number\": \"WBILMTESTTRACK\",\n            \"price\": 453,\n            \"rid\": \"ab4219087a764ae0btest\",\n            \"name\": \"Mascaras\",\n            \"sale\": 30,\n            \"size\": \"0\",\n            \"total_price\": 317,\n            \"nm_id\": 2389212,\n            \"brand\": \"Vivienne Sabo\",\n            \"status\": 202\n        }\n    ],\n    \"locale\": \"en\",\n    \"internal_signature\": \"\",\n    \"customer_id\": \"publish\",\n    \"delivery_service\": \"meest\",\n    \"shardkey\": \"9\",\n    \"sm_id\": 99,\n    \"date_created\": \"2021-11-26T06:22:19Z\",\n    \"oof_shard\": \"1\"\n}"),
			expect: assert.False,
		},
	}

	for _, c := range cases {
		fmt.Println("test_case:", c.descr)
		_, ok := validate(c.data)
		c.expect(t, ok)
		fmt.Println("result:", ok)
	}
}

func TestBroker_Subscribe(t *testing.T) {
	stan := config.StanConfig{
		ClusterID: "publish-cluster",
		ClientID:  "subscriber56",
		ChannelID: "publish-channel",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := repository.NewMockRepository(ctrl)

	c := cache.NewCache(time.Hour, time.Hour)
	broker, err := NewBroker(stan, c, repo)
	if err != nil {
		t.Fatal(err)
	}
	defer broker.Close()

	orders[0].DateCreated = time.Now()
	repo.EXPECT().Insert(context.Background(), orders[0]).Return(nil)

	err = broker.Subscribe()
	if err != nil {
		t.Fatal(err)
	}
	defer broker.Unsubscribe()

	err = broker.Publish(orders[0])
	if err != nil {
		t.Error(err)
	}
	time.Sleep(2 * time.Second)
}
