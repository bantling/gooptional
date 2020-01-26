package gooptional

import (
	"database/sql"
	"fmt"
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

	func() {
		defer func() {
			assert.True(t, notPresentError == recover())
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

	val, valid := opt.Get()
	assert.Equal(t, 0, val)
	assert.True(t, valid)
	assert.Equal(t, 0, opt.MustGet())
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

func TestOptionalIntFilter(t *testing.T) {
	opt := OfInt(1)
	assert.True(t, opt == opt.Filter(func(val int) bool { return true }))
	assert.True(t, opt.Filter(func(int) bool { return false }).IsEmpty())

	assert.True(t, OfInt().Filter(func(int) bool { return true }).IsEmpty())
}

func TestOptionalIntMap(t *testing.T) {
	m := func(val int) int {
		return val + 1
	}
	assert.True(t, OfInt().Map(m).IsEmpty())
	assert.Equal(t, 2, OfInt(1).Map(m).MustGet())
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
