package mock

import (
	"context"
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/sschwartz96/syncapod/internal/database"
)

type DB struct {
	collectionMap map[string]([]interface{})
	l             *log.Logger
}

func (d *DB) Open(ctx context.Context) error {
	d.collectionMap = make(map[string]([]interface{}))
	return nil
}

func (d *DB) Close(ctx context.Context) error {
	return nil
}

func (d *DB) Insert(collection string, object interface{}) error {
	if object == nil {
		return errors.New("object is nil")
	}

	if d.collectionMap[collection] == nil {
		d.collectionMap[collection] = make([]interface{}, 1)
		d.collectionMap[collection][0] = object
	} else {
		d.collectionMap[collection] = append(d.collectionMap[collection], object)
	}

	return nil
}

func (d *DB) FindOne(collection string, object interface{}, filter *database.Filter, opts *database.Options) error {
	if d.collectionMap[collection] == nil {
		return errors.New("collection doesn not exist")
	}

	dataMap := d.collectionMap[collection]
	for _, data := range dataMap {
		if compareInterfaceToFilter(d.l, data, filter) {
			return setValue(object, data)
		}
	}

	return errors.New("no object found with filter")
}

func (d *DB) FindAll(collection string, slice interface{}, filter *database.Filter, opts *database.Options) error {
	if d.collectionMap[collection] == nil {
		return errors.New("collection does not not exist")
	}

	dataMap := d.collectionMap[collection]
	for _, data := range dataMap {
		if compareInterfaceToFilter(d.l, data, filter) {
			setValue(object, data)
		}
	}

	return errors.New("no object found with filter")
}

func (d *DB) Update(collection string, object interface{}, filter *database.Filter) error {
	panic("not implemented") // TODO: Implement
}

func (d *DB) Upsert(collection string, object interface{}, filter *database.Filter) error {
	panic("not implemented") // TODO: Implement
}

func (d *DB) Delete(collection string, filter *database.Filter) error {
	panic("not implemented") // TODO: Implement
}

func (d *DB) Search(collection string, search string, fields []string, object interface{}) error {
	panic("not implemented") // TODO: Implement
}

func compareInterfaceToFilter(l *log.Logger, a interface{}, filter *database.Filter) bool {
	aVal := reflect.ValueOf(a)

	if !aVal.IsValid() {
		l.Printf("not valid: %v\n", a)
		return false
	}

	for filterKey, filterVal := range *filter {
		for i := 0; i < aVal.NumField(); i++ {
			fieldVal := aVal.Field(i)
			fieldName := aVal.Type().Field(i).Name
			l.Printf("comparing fields: %s & %s. values are %v & %v\n", filterKey, fieldName, filterVal, fieldVal)
			if isLowerEqual(filterKey, fieldName) {
				l.Println("names are equal")
				if isEqual(l, reflect.ValueOf(filterVal), fieldVal) {
					l.Println("values are equal")
					break
				} else {
					return false
				}
			}
		}
	}
	return true
}

func isEqual(l *log.Logger, a, b reflect.Value) bool {
	if reflect.TypeOf(a) == reflect.TypeOf(b) {
		l.Printf("values: %v & %v\n", a.Interface(), b.Interface())
		return a.Interface() == b.Interface()
	}
	return false
}

func isLowerEqual(a, b string) bool {
	return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == 0
}

func setValue(object, data interface{}) error {
	if reflect.TypeOf(object).Kind() != reflect.Ptr {
		return errors.New("input object is not type pointer")
	}
	reflect.ValueOf(object).Elem().Set(reflect.ValueOf(data))
	return nil
}
