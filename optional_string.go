package gooptional

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
)

var (
	emptyStringString = "OptionalString"
)

// OptionalString is a mostly immutable wrapper for a string value with a present flag.
// The only mutable operation is the implementation of the sql.Scanner stringerface.
type OptionalString struct {
	value   string
	present bool
}

// OfString returns an OptionalString.
// If no value is provided, an empty OptionalString is returned.
// Otherwise a new OptionalString that wraps the value is returned.
func OfString(value ...string) OptionalString {
	opt := OptionalString{}
	if len(value) == 0 {
		return opt
	}

	opt.value = value[0]
	opt.present = true
	return opt
}

// Equal returns true if:
// 1. This OptionalString is empty and the OptionalString passed is empty.
// 2. This OptionalString is present and the OptionalString passed is present and contains the same value.
func (o OptionalString) Equal(opt OptionalString) bool {
	if !o.present {
		return !opt.present
	}

	if !opt.present {
		return false
	}

	return o.value == opt.value
}

// NotEqual returns the opposite of Equal
func (o OptionalString) NotEqual(opt OptionalString) bool {
	if !o.present {
		return opt.present
	}

	if !opt.present {
		return true
	}

	return o.value != opt.value
}

// EqualValue returns true if this OptionalString is present and contains the value passed
func (o OptionalString) EqualValue(val string) bool {
	if !o.present {
		return false
	}

	return o.value == val
}

// NotEqualValue returns the opposite of EqualValue
func (o OptionalString) NotEqualValue(val string) bool {
	if !o.present {
		return true
	}

	return o.value != val
}

// Filter applies the predicate to the value of this OptionalString.
// Returns this OptionalString only if this OptionalString is present and the filter returns true for the value.
// Otherwise an empty OptionalString is returned.
func (o OptionalString) Filter(predicate func(string) bool) OptionalString {
	if o.present && predicate(o.value) {
		return o
	}

	return OptionalString{}
}

// FilterNot applies the inverse predicate to the value of this OptionalString.
// Returns this OptionalString only if this OptionalString is present and the filter returns false for the value.
// Otherwise an empty OptionalString is returned.
func (o OptionalString) FilterNot(predicate func(string) bool) OptionalString {
	if o.present && (!predicate(o.value)) {
		return o
	}

	return OptionalString{}
}

// Get returns the wrapped value and whether or not it is present.
// The value is only valid if the boolean is true.
func (o OptionalString) Get() (string, bool) {
	return o.value, o.present
}

// IfPresent executes the consumer function with the wrapped value only if the value is present.
func (o OptionalString) IfPresent(consumer func(string)) {
	if o.present {
		consumer(o.value)
	}
}

// IfEmpty executes the function only if the value is not present.
func (o OptionalString) IfEmpty(f func()) {
	if !o.present {
		f()
	}
}

// IfPresentOrElse executes the consumer function with the wrapped value if the value is present, otherwise executes the function of no args.
func (o OptionalString) IfPresentOrElse(consumer func(string), f func()) {
	if o.present {
		consumer(o.value)
	} else {
		f()
	}
}

// Empty returns true if this OptionalString is not present
func (o OptionalString) IsEmpty() bool {
	return !o.present
}

// Present returns true if this OptionalString is present
func (o OptionalString) IsPresent() bool {
	return o.present
}

// FlatMap operates like Map, except that the mapping function already returns an OptionalString, which is returned as is.
func (o OptionalString) FlatMap(f func(string) OptionalString) OptionalString {
	if o.present {
		return f(o.value)
	}

	return OptionalString{}
}

// Map the wrapped value with the given mapping function, which must return the same type.
// If this optional is not present, the function is not invoked and an empty OptionalString is returned.
// Otherwise, a new OptionalString wrapping the mapped value is returned.
func (o OptionalString) Map(f func(string) string) OptionalString {
	if o.present {
		return OfString(f(o.value))
	}

	return OptionalString{}
}

// FlatMapTo operates like MapTo, except that the mapping function already returns an Optional, which is returned as is.
func (o OptionalString) FlatMapTo(f func(string) Optional) Optional {
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
func (o OptionalString) MapTo(f func(string) interface{}) Optional {
	if o.present {
		v := f(o.value)
		if !reflect.ValueOf(v).IsZero() {
			return Of(v)
		}
	}

	return Optional{}
}

// FlatMapToFloat operates like MapToFloat, except that the mapping function already returns an OptionalFloat, which is returned as is.
func (o OptionalString) FlatMapToFloat(f func(string) OptionalFloat) OptionalFloat {
	if o.present {
		return f(o.value)
	}

	return OptionalFloat{}
}

// MapToFloat maps the wrapped value to a float64 with the given mapping function.
// If this optional is not present, the function is not invoked and an empty OptionalFloat is returned.
// Otherwise, an OptionalFloat wrapping the mapped value is returned.
func (o OptionalString) MapToFloat(f func(string) float64) OptionalFloat {
	if o.present {
		return OfFloat(f(o.value))
	}

	return OptionalFloat{}
}

// FlatMapToInt operates like MapToInt, except that the mapping function already returns an OptionalInt, which is returned as is.
func (o OptionalString) FlatMapToInt(f func(string) OptionalInt) OptionalInt {
	if o.present {
		return f(o.value)
	}

	return OptionalInt{}
}

// MapToInt the wrapped value to an int with the given mapping function.
// If this optional is not present, the function is not invoked and an empty OptionalInt is returned.
// Otherwise, an OptionalInt wrapping the mapped value is returned.
func (o OptionalString) MapToInt(f func(string) int) OptionalInt {
	if o.present {
		return OfInt(f(o.value))
	}

	return OptionalInt{}
}

// MustGet returns the unwrapped value and panics if it is not present
func (o OptionalString) MustGet() string {
	if !o.present {
		panic(notPresentError)
	}

	return o.value
}

// OrElse returns the wrapped value if it is present, else it returns the given value
func (o OptionalString) OrElse(value string) string {
	if o.present {
		return o.value
	}

	return value
}

// OrElseGet returns the wrapped value if it is present, else it returns the result of the given function
func (o OptionalString) OrElseGet(supplier func() string) string {
	if o.present {
		return o.value
	}

	return supplier()
}

// OrElsePanic returns the wrapped value if it is present, else it panics with the result of the given function
func (o OptionalString) OrElsePanic(f func() error) string {
	if o.present {
		return o.value
	}

	panic(f())
}

// Scan is database/sql Scanner string, allowing users to read null query columns into an OptionalString.
// This is the only method that modifies an OptionalString.
// The result will be same whether or not the OptionalString was initially empty.
// If the value is not compatible with sql.NullString, an error will be thrown.
func (o *OptionalString) Scan(src interface{}) error {
	var val sql.NullString
	if err := val.Scan(src); err != nil {
		return err
	}

	o.value = val.String
	o.present = true
	return nil
}

// String returns fmt.Sprintf("OptionalString (%v)", wrapped value) if it is present, else "OptionalString" if it is empty.
func (o OptionalString) String() string {
	if o.present {
		return fmt.Sprintf("OptionalString (%v)", o.value)
	}

	return emptyStringString
}

// Value is the database/sql/driver/Valuer stringerface, allowing users to write an OptionalString stringo a column.
func (o OptionalString) Value() (driver.Value, error) {
	if !o.present {
		return nil, nil
	}

	return o.value, nil
}
