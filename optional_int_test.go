package gooptional

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionalIntOfEmptyPresentGet(t *testing.T) {
	opt := OfInt()
	assert.Equal(t, 0, opt.value)
	assert.False(t, opt.present)
	assert.True(t, opt.IsEmpty())
	assert.False(t, opt.IsPresent())
	called := false
	opt.IfPresent(func(int) { called = true })
	assert.False(t, called)
	opt.IfEmpty(func() { called = true })
	assert.True(t, called)
	called = false
	opt.IfPresentOrElse(func(int) {}, func() { called = true })
	assert.True(t, called)

	func() {
		defer func() {
			assert.True(t, errNotPresent == recover())
		}()

		opt.MustGet()
		assert.Fail(t, "Expected Panic")
	}()

	opt = OfInt(0)
	assert.Equal(t, 0, opt.value)
	assert.True(t, opt.present)
	assert.False(t, opt.IsEmpty())
	assert.True(t, opt.IsPresent())
	val := 1
	opt.IfPresent(func(v int) { val = v })
	assert.Equal(t, 0, val)
	val = 1
	opt.IfEmpty(func() { val = 0 })
	assert.Equal(t, 1, val)
	val = 1
	opt.IfPresentOrElse(func(v int) { val = 2 }, func() { val = 3 })
	assert.Equal(t, 2, val)

	val, valid := opt.Get()
	assert.Equal(t, 0, val)
	assert.True(t, valid)
	assert.Equal(t, 0, opt.MustGet())
}

func TestOptionalIntIter(t *testing.T) {
	var opt OptionalInt
	iter := opt.Iter()
	assert.False(t, iter.Next())

	opt = OfInt(1)
	iter = opt.Iter()
	assert.True(t, iter.Next())
	assert.Equal(t, 1, iter.Value())
	assert.False(t, iter.Next())
}

func TestOptionalIntEqual(t *testing.T) {
	// Not present optional == not present optional
	assert.True(t, OfInt().Equal(OfInt()))

	opt1 := OfInt(0)
	opt2 := OfInt(0)
	// Present optional != not present optional
	assert.False(t, opt1.Equal(OfInt()))

	// Present optional == present optional if values equal
	assert.True(t, opt1.Equal(opt1))
	assert.True(t, opt1.Equal(opt2))

	// Present optional != present optional if values differ
	assert.False(t, opt1.Equal(OfInt(1)))

	// Not present optional never equals any value
	assert.False(t, OfInt().EqualValue(0))

	// Not present optional != non-nil value
	assert.False(t, OfInt().EqualValue(0))

	// Present optional == value if values equal
	assert.True(t, opt1.EqualValue(0))

	// Present optional != value if values differ
	assert.False(t, opt1.Equal(OfInt(1)))
}

func TestOptionalIntNotEqual(t *testing.T) {
	// Not present optional == not present optional
	assert.False(t, OfInt().NotEqual(OfInt()))

	opt1 := OfInt(0)
	opt2 := OfInt(0)
	// Present optional != not present optional
	assert.True(t, opt1.NotEqual(OfInt()))

	// Present optional == present optional if values equal
	assert.False(t, opt1.NotEqual(opt1))
	assert.False(t, opt1.NotEqual(opt2))

	// Present optional != present optional if values differ
	assert.True(t, opt1.NotEqual(OfInt(1)))

	// Not present optional never equals any value
	assert.True(t, OfInt().NotEqualValue(0))

	// Not present optional != non-nil value
	assert.True(t, OfInt().NotEqualValue(0))

	// Present optional == value if values equal
	assert.False(t, opt1.NotEqualValue(0))

	// Present optional != value if values differ
	assert.True(t, opt1.NotEqual(OfInt(1)))
}

func TestOptionalIntFilter(t *testing.T) {
	opt := OfInt(1)
	assert.True(t, opt == opt.Filter(func(val int) bool { return true }))
	assert.True(t, opt.Filter(func(int) bool { return false }).IsEmpty())

	assert.True(t, OfInt().Filter(func(int) bool { return true }).IsEmpty())
}

func TestOptionalIntFilterNot(t *testing.T) {
	opt := OfInt(1)
	assert.True(t, opt == opt.FilterNot(func(val int) bool { return false }))
	assert.True(t, opt.FilterNot(func(int) bool { return true }).IsEmpty())

	assert.True(t, OfInt().FilterNot(func(int) bool { return false }).IsEmpty())
}

func TestOptionalFlafIntMapInterfaceFloatString(t *testing.T) {
	toi := func(val int) OptionalInt {
		return OfInt(val + 1)
	}
	assert.True(t, OfInt().FlatMap(toi).IsEmpty())
	assert.Equal(t, 2, OfInt(1).FlatMap(toi).MustGet())

	tof := func(val int) OptionalFloat {
		return OfFloat(float64(val + 1))
	}
	assert.True(t, OfInt().FlatMapToFloat(tof).IsEmpty())
	assert.Equal(t, 2.0, OfInt(1).FlatMapToFloat(tof).MustGet())

	too := func(val int) Optional {
		return Of(val + 1)
	}
	assert.True(t, OfInt().FlatMapTo(too).IsEmpty())
	assert.Equal(t, 2, OfInt(1).FlatMapTo(too).MustGet())

	tos := func(val int) OptionalString {
		return OfString(strconv.Itoa(val + 1))
	}
	assert.True(t, OfInt().FlatMapToString(tos).IsEmpty())
	assert.Equal(t, "2", OfInt(1).FlatMapToString(tos).MustGet())
}

func TestOptionalIntMapInterfaceFloatString(t *testing.T) {
	toi := func(val int) int {
		return val + 1
	}
	assert.True(t, OfInt().Map(toi).IsEmpty())
	assert.Equal(t, 2, OfInt(1).Map(toi).MustGet())

	tof := func(val int) float64 {
		return float64(val + 1)
	}
	assert.True(t, OfInt().MapToFloat(tof).IsEmpty())
	assert.Equal(t, 2.0, OfInt(1).MapToFloat(tof).MustGet())

	too := func(val int) interface{} {
		return val + 1
	}
	assert.True(t, OfInt().MapTo(too).IsEmpty())
	assert.Equal(t, 2, OfInt(1).MapTo(too).MustGet())

	tos := func(val int) string {
		return strconv.Itoa(val + 1)
	}
	assert.True(t, OfInt().MapToString(tos).IsEmpty())
	assert.Equal(t, "2", OfInt(1).MapToString(tos).MustGet())
}

func TestOptionalIntOrElseGetPanic(t *testing.T) {
	f := func() int { return 2 }
	assert.Equal(t, 1, OfInt().OrElse(1))
	assert.Equal(t, 2, OfInt().OrElseGet(f))

	err := fmt.Errorf("")
	errf := func() error { return err }
	func() {
		defer func() {
			assert.True(t, err == recover())
		}()
		OfInt().OrElsePanic(errf)
		assert.Fail(t, "Expected Panic")
	}()

	assert.Equal(t, 3, OfInt(3).OrElse(1))
	assert.Equal(t, 3, OfInt(3).OrElseGet(f))
	assert.Equal(t, 3, OfInt(3).OrElsePanic(errf))
}

func TestOptionalIntScan(t *testing.T) {
	var opt OptionalInt
	assert.Nil(t, opt.Scan(0))
	assert.Equal(t, 0, opt.MustGet())
	assert.Nil(t, opt.Scan(1))
	assert.Equal(t, 1, opt.MustGet())

	sc := (sql.Scanner)(&opt)
	assert.NotNil(t, &sc)
}

func TestOptionalIntString(t *testing.T) {
	assert.Equal(t, emptyIntString, fmt.Sprintf("%s", OfInt()))
	assert.Equal(t, "OptionalInt (1)", fmt.Sprintf("%s", OfInt(1)))
}

func TestOptionalIntValue(t *testing.T) {
	val, err := OfInt().Value()
	assert.Nil(t, val)
	assert.Nil(t, err)

	val, err = OfInt(0).Value()
	assert.Equal(t, 0, val)
	assert.Nil(t, err)
}
