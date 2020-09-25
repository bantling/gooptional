package gooptional

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/bantling/goiter"
)

// Optional is a mostly immutable generic wrapper for any kind of value with a present flag.
// The only mutable operation is the implementation of the sql.Scanner interface.
// The zero value is ready to use.
type Optional struct {
	value   interface{}
	present bool
}

var (
	errNotPresent = fmt.Errorf("No value present")
	emptyString   = "Optional"
)

// IsNil returns true if a value is nil.
// When a nil value v is received as type interface{}, v == nil will usually be false.
// Need to check if either v == nil or Sprintf("%p") == "0x0".
// Of course, a string value could be "0x0", so check for string type first.
func IsNil(value interface{}) bool {
	if _, ok := value.(string); ok {
		return false
	}

	if (value == nil) || (fmt.Sprintf("%p", value) == "0x0") {
		return true
	}

	return false
}

// Of returns an Optional.
// If no value or a nil value is provided, a new empty Optional is returned.
// Otherwise a new Optional that wraps the value is returned.
func Of(value ...interface{}) Optional {
	opt := Optional{}
	if len(value) == 0 {
		return opt
	}

	v := value[0]
	if IsNil(v) {
		return opt
	}

	opt.value = v
	opt.present = true
	return opt
}

// Equal returns true if:
// 1. This Optional is empty and the Optional passed is empty.
// 2. This Optional is present and the Optional passed is present and contains the same value using reflect.DeepEqual.
func (o Optional) Equal(opt Optional) bool {
	if !o.present {
		return !opt.present
	}

	if !opt.present {
		return false
	}

	return reflect.DeepEqual(o.value, opt.value)
}

// NotEqual returns the opposite of Equal
func (o Optional) NotEqual(opt Optional) bool {
	if !o.present {
		return opt.present
	}

	if !opt.present {
		return true
	}

	return !reflect.DeepEqual(o.value, opt.value)
}

// EqualValue returns true if this Optional is present and contains the value passed.
// Note that an empty Optional never equals any value, including nil.
func (o Optional) EqualValue(val interface{}) bool {
	if !o.present {
		return false
	}

	return reflect.DeepEqual(o.value, val)
}

// NotEqualValue returns the opposite of EqualValue
func (o Optional) NotEqualValue(val interface{}) bool {
	if !o.present {
		return true
	}

	return !reflect.DeepEqual(o.value, val)
}

// Filter applies the predicate to the value of this Optional.
// Returns this Optional only if this Optional is present and the filter returns true for the value.
// Otherwise an empty Optional is returned.
func (o Optional) Filter(predicate func(interface{}) bool) Optional {
	if o.present && predicate(o.value) {
		return o
	}

	return Optional{}
}

// FilterNot applies the inverted predicate to the value of this Optional.
// Returns this Optional only if this Optional is present and the filter returns false for the value.
// Otherwise an empty Optional is returned.
func (o Optional) FilterNot(predicate func(interface{}) bool) Optional {
	if o.present && (!predicate(o.value)) {
		return o
	}

	return Optional{}
}

// Get returns the wrapped value and whether or not it is present.
// The value is only valid if the boolean is true.
func (o Optional) Get() (interface{}, bool) {
	return o.value, o.present
}

// IfPresent executes the consumer function with the wrapped value only if the value is present.
func (o Optional) IfPresent(consumer func(interface{})) {
	if o.present {
		consumer(o.value)
	}
}

// IfEmpty executes the function only if the value is not present.
func (o Optional) IfEmpty(f func()) {
	if !o.present {
		f()
	}
}

// IfPresentOrElse executes the consumer function with the wrapped value if the value is present, otherwise executes the function of no args.
func (o Optional) IfPresentOrElse(consumer func(interface{}), f func()) {
	if o.present {
		consumer(o.value)
	} else {
		f()
	}
}

// IsEmpty returns true if this Optional is not present
func (o Optional) IsEmpty() bool {
	return !o.present
}

// IsPresent returns true if this Optional is present
func (o Optional) IsPresent() bool {
	return o.present
}

// Iter returns an *Iter of one element containing the wrapped value if present, else an empty Iter
func (o Optional) Iter() *goiter.Iter {
	if o.present {
		return goiter.Of(o.value)
	}

	return goiter.Of()
}

// FlatMap operates like Map, except that the mapping function already returns an Optional, which is returned as is.
func (o Optional) FlatMap(f func(interface{}) Optional) Optional {
	if o.present {
		return f(o.value)
	}

	return Optional{}
}

// Map the wrapped value with the given mapping function, which may return a different type.
// If this optional is not present, the function is not invoked and an empty Optional is returned.
// If this optional is present and the map function returns a nil pointer value, an empty Optional is returned.
// Otherwise, an Optional wrapping the mapped value is returned.
func (o Optional) Map(f func(interface{}) interface{}) Optional {
	if o.present {
		v := f(o.value)
		rv := reflect.ValueOf(v)
		if (rv.Kind() == reflect.Ptr) && (rv.IsNil()) {
			return Optional{}
		}

		return Of(v)
	}

	return Optional{}
}

// FlatMapToFloat operates like MapToFloat, except that the mapping function already returns an OptionalFloat, which is returned as is.
func (o Optional) FlatMapToFloat(f func(interface{}) OptionalFloat) OptionalFloat {
	if o.present {
		return f(o.value)
	}

	return OptionalFloat{}
}

// MapToFloat maps the wrapped value to a float64 with the given mapping function.
// If this optional is not present, the function is not invoked and an empty OptionalFloat is returned.
// Otherwise, an OptionalFloat wrapping the mapped value is returned.
func (o Optional) MapToFloat(f func(interface{}) float64) OptionalFloat {
	if o.present {
		return OfFloat(f(o.value))
	}

	return OptionalFloat{}
}

// FlatMapToInt operates like MapToInt, except that the mapping function already returns an OptionalInt, which is returned as is.
func (o Optional) FlatMapToInt(f func(interface{}) OptionalInt) OptionalInt {
	if o.present {
		return f(o.value)
	}

	return OptionalInt{}
}

// MapToInt the wrapped value to an int with the given mapping function.
// If this optional is not present, the function is not invoked and an empty OptionalInt is returned.
// Otherwise, an OptionalInt wrapping the mapped value is returned.
func (o Optional) MapToInt(f func(interface{}) int) OptionalInt {
	if o.present {
		return OfInt(f(o.value))
	}

	return OptionalInt{}
}

// FlatMapToString operates like MapToString, except that the mapping function already returns an OptionalString, which is returned as is.
func (o Optional) FlatMapToString(f func(interface{}) OptionalString) OptionalString {
	if o.present {
		return f(o.value)
	}

	return OptionalString{}
}

// MapToString the wrapped value to a string with the given mapping function.
// If this optional is not present, the function is not invoked and an empty OptionalString is returned.
// Otherwise, an OptionalString wrapping the mapped value is returned.
func (o Optional) MapToString(f func(interface{}) string) OptionalString {
	if o.present {
		return OfString(f(o.value))
	}

	return OptionalString{}
}

// MustGet returns the unwrapped value and panics if it is not present
func (o Optional) MustGet() interface{} {
	if !o.present {
		panic(errNotPresent)
	}

	return o.value
}

// OrElse returns the wrapped value if it is present, else it returns the given value.
func (o Optional) OrElse(value interface{}) interface{} {
	if o.present {
		return o.value
	}

	return value
}

// OrElseGet returns the wrapped value if it is present, else it returns the result of the given function.
func (o Optional) OrElseGet(supplier func() interface{}) interface{} {
	if o.present {
		return o.value
	}

	return supplier()
}

// OrElsePanic returns the wrapped value if it is present, else it panics with the result of the given function
func (o Optional) OrElsePanic(f func() error) interface{} {
	if o.present {
		return o.value
	}

	panic(f())
}

// Scan is database/sql Scanner interface, allowing users to read null query columns into an Optional.
// This is the only method that modifies an Optional.
// The result will be same whether or not the Optional was initially empty.
// The provided value is just stored, so if it is a reference type it must be copied before the next call to Scan.
// Since any value can be stored, the result is always a nil error.
// It is up to the caller to ensure the correct type is being read.
func (o *Optional) Scan(src interface{}) error {
	o.value = src
	o.present = src != nil
	return nil
}

// Value is the database/sql/driver/Valuer interface, allowing users to write an Optional into a column.
// If a present optional does not contain an allowed type, the operation will fail.
// It is up to the caller to ensure the correct type is being written.
func (o Optional) Value() (driver.Value, error) {
	if !o.present {
		return nil, nil
	}

	return o.value, nil
}

// String returns fmt.Sprintf("Optional (%v)", wrapped value) if it is present, else "Optional" if it is empty.
func (o Optional) String() string {
	if o.present {
		return fmt.Sprintf("Optional (%v)", o.value)
	}

	return emptyString
}
