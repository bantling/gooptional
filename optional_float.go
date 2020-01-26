package gooptional

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

var (
	emptyFloatString = "OptionalFloat"
)

// OptionalFloat is a mostly immutable wrapper for a float64 value with a present flag.
// The only mutable operation is the implementation of the sql.Scanner float64erface.
type OptionalFloat struct {
	value   float64
	present bool
}

// OfFloat returns an OptionalFloat.
// If no value is provided, an empty OptionalFloat is returned.
// Otherwise a new OptionalFloat that wraps the value is returned.
func OfFloat(value ...float64) OptionalFloat {
	opt := OptionalFloat{}
	if len(value) == 0 {
		return opt
	}

	opt.value = value[0]
	opt.present = true
	return opt
}

// Equal returns true if:
// 1. This OptionalFloat is empty and the OptionalFloat passed is empty.
// 2. This OptionalFloat is present and the OptionalFloat passed is present and contains the same value.
func (o OptionalFloat) Equal(opt OptionalFloat) bool {
	if !o.present {
		return !opt.present
	}

	if !opt.present {
		return false
	}

	return o.value == opt.value
}

// EqualValue returns true if this OptionalFloat is present and contains the value passed
func (o OptionalFloat) EqualValue(val float64) bool {
	if !o.present {
		return false
	}

	return o.value == val
}

// Filter applies the predicate to the value of this OptionalFloat.
// Returns this OptionalFloat only if this OptionalFloat is present and the filter returns true for the value.
// Otherwise an empty OptionalFloat is returned.
func (o OptionalFloat) Filter(predicate func(float64) bool) OptionalFloat {
	if o.present && predicate(o.value) {
		return o
	}

	return OptionalFloat{}
}

// Get returns the wrapped value and whether or not it is present.
// The value is only valid if the boolean is true.
func (o OptionalFloat) Get() (float64, bool) {
	return o.value, o.present
}

// IfPresent executes the consumer function with the wrapped value only if the value is present.
func (o OptionalFloat) IfPresent(consumer func(float64)) {
	if o.present {
		consumer(o.value)
	}
}

// Empty returns true if this OptionalFloat is not present
func (o OptionalFloat) IsEmpty() bool {
	return !o.present
}

// Present returns true if this OptionalFloat is present
func (o OptionalFloat) IsPresent() bool {
	return o.present
}

// Map the wrapped value with the given mapping function, which must return the same type.
// If this optional is not present, the function is not invoked and an empty OptionalFloat is returned.
// Otherwise, a new OptionalFloat wrapping the mapped value is returned.
func (o OptionalFloat) Map(f func(float64) float64) OptionalFloat {
	if o.present {
		return OfFloat(f(o.value))
	}

	return OptionalFloat{}
}

// MustGet returns the unwrapped value and panics if it is not present
func (o OptionalFloat) MustGet() float64 {
	if !o.present {
		panic(notPresentError)
	}

	return o.value
}

// OrElse returns the wrapped value if it is present, else it returns the given value
func (o OptionalFloat) OrElse(value float64) float64 {
	if o.present {
		return o.value
	}

	return value
}

// OrElseGet returns the wrapped value if it is present, else it returns the result of the given function
func (o OptionalFloat) OrElseGet(supplier func() float64) float64 {
	if o.present {
		return o.value
	}

	return supplier()
}

// OrElsePanic returns the wrapped value if it is present, else it panics with the result of the given function
func (o OptionalFloat) OrElsePanic(f func() error) float64 {
	if o.present {
		return o.value
	}

	panic(f())
}

// Scan is database/sql Scanner float64, allowing users to read null query columns into an OptionalFloat.
// This is the only method that modifies an OptionalFloat.
// The result will be same whether or not the OptionalFloat was initially empty.
// If the value is not compatible with sql.NullFloat64, an error will be thrown.
func (o *OptionalFloat) Scan(src interface{}) error {
	var val sql.NullFloat64
	if err := val.Scan(src); err != nil {
		return err
	}

	o.value = float64(val.Float64)
	o.present = true
	return nil
}

// String returns fmt.Sprintf("OptionalFloat (%v)", wrapped value) if it is present, else "OptionalFloat" if it is empty.
func (o OptionalFloat) String() string {
	if o.present {
		return fmt.Sprintf("OptionalFloat (%v)", o.value)
	}

	return emptyFloatString
}

// Value is the database/sql/driver/Valuer float64erface, allowing users to write an OptionalFloat float64o a column.
func (o OptionalFloat) Value() (driver.Value, error) {
	if !o.present {
		return nil, nil
	}

	return o.value, nil
}
