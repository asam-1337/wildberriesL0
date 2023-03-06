package templateRepo

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/asam-1337/wildberriesL0/internal/pgxbalancer"
	"reflect"
)

func getDtoInfo(dto interface{}) ([]string, []interface{}) {
	pointerVal := reflect.ValueOf(dto)
	val := reflect.Indirect(pointerVal)
	typ := val.Type()
	numFeilds := val.NumField()
	columns := make([]string, 0, numFeilds)
	values := make([]interface{}, 0, numFeilds)
	for i := 0; i < numFeilds; i++ {
		values = append(values, val.Field(i).Interface())
		columns = append(columns, typ.Field(i).Tag.Get("db"))
	}
	return columns, values
}

func getDtoColumns(dto interface{}) []string {
	pointerVal := reflect.ValueOf(dto)
	val := reflect.Indirect(pointerVal)
	typ := val.Type()
	numFeilds := val.NumField()
	columns := make([]string, 0, numFeilds)

	for i := 0; i < numFeilds; i++ {
		columns = append(columns, typ.Field(i).Tag.Get("db"))
	}
	return columns
}

type repo interface {
	GetRunner() pgxbalancer.Runner
	TableName() string
}

func Create[T any](ctx context.Context, r repo, value T) {
	col, val := getDtoInfo(value)
	sql, v, err := squirrel.Insert(r.TableName()).
		Columns(col...).Values(val...).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {

	}
}
