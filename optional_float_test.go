package gooptional

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionalFloatOfEmptyPresentGet(t *testing.T) {
	opt := OfFloat()
	assert.Equal(t, 0.0, opt.value)
	assert.False(t, opt.present)
	assert.True(t, opt.IsEmpty())
	assert.False(t, opt.IsPresent())
	called := false
	opt.IfPresent(func(float64) { called = true })
	assert.False(t, called)

	func() {
		defer func() {
			assert.True(t, notPresentError == recover())
		}()

		opt.MustGet()
		assert.Fail(t, "Expected Panic")
	}()

	opt = OfFloat(0)
	assert.Equal(t, 0.0, opt.value)
	assert.True(t, opt.present)
	assert.False(t, opt.IsEmpty())
	assert.True(t, opt.IsPresent())
	val := 1.0
	opt.IfPresent(func(v float64) { val = v })
	assert.Equal(t, 0.0, val)

	val, valid := opt.Get()
	assert.Equal(t, 0.0, val)
	assert.True(t, valid)
	assert.Equal(t, 0.0, opt.MustGet())
}

func TestOptionalFloatEqual(t *testing.T) {
	// Not present optional == not present optional
	assert.True(t, OfFloat().Equal(OfFloat()))

	opt1 := OfFloat(0)
	opt2 := OfFloat(0)
	// Present optional != not present optional
	assert.False(t, opt1.Equal(OfFloat()))

	// Present optional == present optional if values equal
	assert.True(t, opt1.Equal(opt1))
	assert.True(t, opt1.Equal(opt2))

	// Present optional != present optional if values differ
	assert.False(t, opt1.Equal(OfFloat(1)))

	// Not present optional never equals any value
	assert.False(t, OfFloat().EqualValue(0))

	// Not present optional != non-nil value
	assert.False(t, OfFloat().EqualValue(0))

	// Present optional == value if values equal
	assert.True(t, opt1.EqualValue(0))

	// Present optional != value if values differ
	assert.False(t, opt1.Equal(OfFloat(1)))
}

func TestOptionalFloatFilter(t *testing.T) {
	opt := OfFloat(1)
	assert.True(t, opt == opt.Filter(func(val float64) bool { return true }))
	assert.True(t, opt.Filter(func(float64) bool { return false }).IsEmpty())

	assert.True(t, OfFloat().Filter(func(float64) bool { return true }).IsEmpty())
}

func TestOptionalFloatMap(t *testing.T) {
	m := func(val float64) float64 {
		return val + 1
	}
	assert.True(t, OfFloat().Map(m).IsEmpty())
	assert.Equal(t, 2.0, OfFloat(1).Map(m).MustGet())
}

func TestOptionalFloatOrElseGetPanic(t *testing.T) {
	f := func() float64 { return 2 }
	assert.Equal(t, 1.0, OfFloat().OrElse(1))
	assert.Equal(t, 2.0, OfFloat().OrElseGet(f))

	err := fmt.Errorf("")
	errf := func() error { return err }
	func() {
		defer func() {
			assert.True(t, err == recover())
		}()
		OfFloat().OrElsePanic(errf)
		assert.Fail(t, "Expected Panic")
	}()

	assert.Equal(t, 3.0, OfFloat(3).OrElse(1))
	assert.Equal(t, 3.0, OfFloat(3).OrElseGet(f))
	assert.Equal(t, 3.0, OfFloat(3).OrElsePanic(errf))
}

func TestOptionalFloatScan(t *testing.T) {
	var opt OptionalFloat
	assert.Nil(t, opt.Scan(0))
	assert.Equal(t, 0.0, opt.MustGet())
	assert.Nil(t, opt.Scan(1))
	assert.Equal(t, 1.0, opt.MustGet())

	sc := (sql.Scanner)(&opt)
	assert.NotNil(t, &sc)
}

func TestOptionalFloatString(t *testing.T) {
	assert.Equal(t, emptyFloatString, fmt.Sprintf("%s", OfFloat()))
	assert.Equal(t, "OptionalFloat (1)", fmt.Sprintf("%s", OfFloat(1)))
}

func TestOptionalFloatValue(t *testing.T) {
	val, err := OfFloat().Value()
	assert.Nil(t, val)
	assert.Nil(t, err)

	val, err = OfFloat(0).Value()
	assert.Equal(t, 0.0, val)
	assert.Nil(t, err)
}
