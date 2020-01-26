package gooptional

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
)

var (
	emptyIntString = "OptionalInt"
)

// OptionalInt is a mostly immutable wrapper for an int value with a present flag.
// The only mutable operation is the implementation of the sql.Scanner interface.
type OptionalInt struct {
	value   int
	present bool
}

// OfInt returns an OptionalInt.
// If no value is provided, an empty OptionalInt is returned.
// Otherwise a new OptionalInt that wraps the value is returned.
func OfInt(value ...int) OptionalInt {
	opt := OptionalInt{}
	if len(value) == 0 {
		return opt
	}

	opt.value = value[0]
	opt.present = true
	return opt
}

// Equal returns true if:
// 1. This OptionalInt is empty and the OptionalInt passed is empty.
// 2. This OptionalInt is present and the OptionalInt passed is present and contains the same value.
func (o OptionalInt) Equal(opt OptionalInt) bool {
	if !o.present {
		return !opt.present
	}

	if !opt.present {
		return false
	}

	return o.value == opt.value
}

// EqualValue returns true if this OptionalInt is present and contains the value passed
func (o OptionalInt) EqualValue(val int) bool {
	if !o.present {
		return false
	}

	return o.value == val
}

// Filter applies the predicate to the value of this OptionalInt.
// Returns this OptionalInt only if this OptionalInt is present and the filter returns true for the value.
// Otherwise an empty OptionalInt is returned.
func (o OptionalInt) Filter(predicate func(int) bool) OptionalInt {
	if o.present && predicate(o.value) {
		return o
	}

	return OptionalInt{}
}

// Get returns the wrapped value and whether or not it is present.
// The value is only valid if the boolean is true.
func (o OptionalInt) Get() (int, bool) {
	return o.value, o.present
}

// IfPresent executes the consumer function with the wrapped value only if the value is present.
func (o OptionalInt) IfPresent(consumer func(int)) {
	if o.present {
		consumer(o.value)
	}
}

// Empty returns true if this OptionalInt is not present
func (o OptionalInt) IsEmpty() bool {
	return !o.present
}

// Present returns true if this OptionalInt is present
func (o OptionalInt) IsPresent() bool {
	return o.present
}

// Map the wrapped value with the given mapping function, which must return the same type.
// If this optional is not present, the function is not invoked and an empty OptionalInt is returned.
// Otherwise, a new OptionalInt wrapping the mapped value is returned.
func (o OptionalInt) Map(f func(int) int) OptionalInt {
	if o.present {
		return OfInt(f(o.value))
	}

	return OptionalInt{}
}

// MapToFloat maps the wrapped value to a float64 with the given mapping function.
// If this optional is not present, the function is not invoked and an empty OptionalFloat is returned.
// Otherwise, an OptionalFloat wrapping the mapped value is returned.
func (o OptionalInt) MapToFloat(f func(int) float64) OptionalFloat {
	if o.present {
		return OfFloat(f(o.value))
	}

	return OptionalFloat{}
}

// MapToOptional maps the wrapped value with the given mapping function, which may return a different type.
// If this optional is not present, the function is not invoked and an empty Optional is returned.
// If this optional is present and the map function returns a zero value, an empty Optional is returned.
// Otherwise, an Optional wrapping the mapped value is returned.
// The mapping function result is determined to be zero by reflect.Value.IsZero().
func (o OptionalInt) MapToOptional(f func(int) interface{}) Optional {
	if o.present {
		v := f(o.value)
		if !reflect.ValueOf(v).IsZero() {
			return Of(v)
		}
	}

	return Optional{}
}

// MapToString the wrapped value to a string with the given mapping function.
// If this optional is not present, the function is not invoked and an empty OptionalString is returned.
// Otherwise, an OptionalString wrapping the mapped value is returned.
func (o OptionalInt) MapToString(f func(int) string) OptionalString {
	if o.present {
		return OfString(f(o.value))
	}

	return OptionalString{}
}

// MustGet returns the unwrapped value and panics if it is not present
func (o OptionalInt) MustGet() int {
	if !o.present {
		panic(notPresentError)
	}

	return o.value
}

// OrElse returns the wrapped value if it is present, else it returns the given value
func (o OptionalInt) OrElse(value int) int {
	if o.present {
		return o.value
	}

	return value
}

// OrElseGet returns the wrapped value if it is present, else it returns the result of the given function
func (o OptionalInt) OrElseGet(supplier func() int) int {
	if o.present {
		return o.value
	}

	return supplier()
}

// OrElsePanic returns the wrapped value if it is present, else it panics with the result of the given function
func (o OptionalInt) OrElsePanic(f func() error) int {
	if o.present {
		return o.value
	}

	panic(f())
}

// Scan is database/sql Scanner interface, allowing users to read null query columns into an OptionalInt.
// This is the only method that modifies an OptionalInt.
// The result will be same whether or not the OptionalInt was initially empty.
// If the value is not compatible with sql.NullInt64, an error will be thrown.
func (o *OptionalInt) Scan(src interface{}) error {
	var val sql.NullInt64
	if err := val.Scan(src); err != nil {
		return err
	}

	o.value = int(val.Int64)
	o.present = true
	return nil
}

// String returns fmt.Sprintf("OptionalInt (%v)", wrapped value) if it is present, else "OptionalInt" if it is empty.
func (o OptionalInt) String() string {
	if o.present {
		return fmt.Sprintf("OptionalInt (%v)", o.value)
	}

	return emptyIntString
}

// Value is the database/sql/driver/Valuer interface, allowing users to write an OptionalInt into a column.
func (o OptionalInt) Value() (driver.Value, error) {
	if !o.present {
		return nil, nil
	}

	return o.value, nil
}
