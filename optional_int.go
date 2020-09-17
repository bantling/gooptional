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

// NotEqual returns the opposite of Equal
func (o OptionalInt) NotEqual(opt OptionalInt) bool {
	if !o.present {
		return opt.present
	}

	if !opt.present {
		return true
	}

	return o.value != opt.value
}

// EqualValue returns true if this OptionalInt is present and contains the value passed
func (o OptionalInt) EqualValue(val int) bool {
	if !o.present {
		return false
	}

	return o.value == val
}

// NotEqualValue returns the opposite of EqualValue
func (o OptionalInt) NotEqualValue(val int) bool {
	if !o.present {
		return true
	}

	return o.value != val
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

// FilterNot applies the inverted predicate to the value of this OptionalInt.
// Returns this OptionalInt only if this OptionalInt is present and the filter returns false for the value.
// Otherwise an empty OptionalInt is returned.
func (o OptionalInt) FilterNot(predicate func(int) bool) OptionalInt {
	if o.present && (!predicate(o.value)) {
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

// IfEmpty executes the function only if the value is not present.
func (o OptionalInt) IfEmpty(f func()) {
	if !o.present {
		f()
	}
}

// IfPresentOrElse executes the consumer function with the wrapped value if the value is present, otherwise executes the function of no args.
func (o OptionalInt) IfPresentOrElse(consumer func(int), f func()) {
	if o.present {
		consumer(o.value)
	} else {
		f()
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

// FlatMap operates like Map, except that the mapping function already returns an OptionalInt, which is returned as is.
func (o OptionalInt) FlatMap(f func(int) OptionalInt) OptionalInt {
	if o.present {
		return f(o.value)
	}

	return OptionalInt{}
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

// FlatMapTo operates like MapTo, except that the mapping function already returns an OptionalInt, which is returned as is.
func (o OptionalInt) FlatMapTo(f func(int) Optional) Optional {
	if o.present {
		return f(o.value)
	}

	return Optional{}
}

// MapTo maps the wrapped value with the given mapping function, which may return a different type.
// If this optional is not present, the function is not invoked and an empty Optional is returned.
// If this optional is present and the map function returns a zero value, an empty Optional is returned.
// Otherwise, an Optional wrapping the mapped value is returned.
// The mapping function result is determined to be zero by reflect.Value.IsZero().
func (o OptionalInt) MapTo(f func(int) interface{}) Optional {
	if o.present {
		v := f(o.value)
		if !reflect.ValueOf(v).IsZero() {
			return Of(v)
		}
	}

	return Optional{}
}

// FlatMapToFloat operates like MapToFloat, except that the mapping function already returns an OptionalInt, which is returned as is.
func (o OptionalInt) FlatMapToFloat(f func(int) OptionalFloat) OptionalFloat {
	if o.present {
		return f(o.value)
	}

	return OptionalFloat{}
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

// FlatMapToString operates like MapToString, except that the mapping function already returns an OptionalString, which is returned as is.
func (o OptionalInt) FlatMapToString(f func(int) OptionalString) OptionalString {
	if o.present {
		return f(o.value)
	}

	return OptionalString{}
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
