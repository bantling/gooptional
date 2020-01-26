package gooptional

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
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

// EqualValue returns true if this OptionalString is present and contains the value passed
func (o OptionalString) EqualValue(val string) bool {
	if !o.present {
		return false
	}

	return o.value == val
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

// Empty returns true if this OptionalString is not present
func (o OptionalString) IsEmpty() bool {
	return !o.present
}

// Present returns true if this OptionalString is present
func (o OptionalString) IsPresent() bool {
	return o.present
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
