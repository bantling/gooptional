package gooptional

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionalStringOfEmptyPresentGet(t *testing.T) {
	opt := OfString()
	assert.Equal(t, "", opt.value)
	assert.False(t, opt.present)
	assert.True(t, opt.IsEmpty())
	assert.False(t, opt.IsPresent())
	called := false
	opt.IfPresent(func(string) { called = true })
	assert.False(t, called)

	func() {
		defer func() {
			assert.True(t, notPresentError == recover())
		}()

		opt.MustGet()
		assert.Fail(t, "Expected Panic")
	}()

	opt = OfString("0")
	assert.Equal(t, "0", opt.value)
	assert.True(t, opt.present)
	assert.False(t, opt.IsEmpty())
	assert.True(t, opt.IsPresent())
	val := "1"
	opt.IfPresent(func(v string) { val = v })
	assert.Equal(t, "0", val)

	val, valid := opt.Get()
	assert.Equal(t, "0", val)
	assert.True(t, valid)
	assert.Equal(t, "0", opt.MustGet())
}

func TestOptionalStringEqual(t *testing.T) {
	// Not present optional == not present optional
	assert.True(t, OfString().Equal(OfString()))

	opt1 := OfString("0")
	opt2 := OfString("0")
	// Present optional != not present optional
	assert.False(t, opt1.Equal(OfString()))

	// Present optional == present optional if values equal
	assert.True(t, opt1.Equal(opt1))
	assert.True(t, opt1.Equal(opt2))

	// Present optional != present optional if values differ
	assert.False(t, opt1.Equal(OfString("1")))

	// Not present optional never equals any value
	assert.False(t, OfString().EqualValue("0"))

	// Not present optional != non-nil value
	assert.False(t, OfString().EqualValue("0"))

	// Present optional == value if values equal
	assert.True(t, opt1.EqualValue("0"))

	// Present optional != value if values differ
	assert.False(t, opt1.Equal(OfString("1")))
}

func TestOptionalStringFilter(t *testing.T) {
	opt := OfString("1")
	assert.True(t, opt == opt.Filter(func(val string) bool { return true }))
	assert.True(t, opt.Filter(func(string) bool { return false }).IsEmpty())

	assert.True(t, OfString().Filter(func(string) bool { return true }).IsEmpty())
}

func TestOptionalStringMapFloatIntInterface(t *testing.T) {
	m := func(val string) string {
		return val + "1"
	}
	assert.True(t, OfString().Map(m).IsEmpty())
	assert.Equal(t, "11", OfString("1").Map(m).MustGet())

	too := func(val string) interface{} {
		return val + "1"
	}
	assert.True(t, OfString().MapTo(too).IsEmpty())
	assert.Equal(t, "11", OfString("1").MapTo(too).MustGet())

	tof := func(val string) float64 {
		v, _ := strconv.ParseFloat(val+"1", 64)
		return v
	}
	assert.True(t, OfString().MapToFloat(tof).IsEmpty())
	assert.Equal(t, 11.0, OfString("1").MapToFloat(tof).MustGet())

	toi := func(val string) int {
		v, _ := strconv.Atoi(val + "1")
		return v
	}
	assert.True(t, OfString().MapToInt(toi).IsEmpty())
	assert.Equal(t, 11, OfString("1").MapToInt(toi).MustGet())
}

func TestOptionalStringOrElseGetPanic(t *testing.T) {
	f := func() string { return "2" }
	assert.Equal(t, "1", OfString().OrElse("1"))
	assert.Equal(t, "2", OfString().OrElseGet(f))

	err := fmt.Errorf("")
	errf := func() error { return err }
	func() {
		defer func() {
			assert.True(t, err == recover())
		}()
		OfString().OrElsePanic(errf)
		assert.Fail(t, "Expected Panic")
	}()

	assert.Equal(t, "3", OfString("3").OrElse("1"))
	assert.Equal(t, "3", OfString("3").OrElseGet(f))
	assert.Equal(t, "3", OfString("3").OrElsePanic(errf))
}

func TestOptionalStringScan(t *testing.T) {
	var opt OptionalString
	assert.Nil(t, opt.Scan(0))
	assert.Equal(t, "0", opt.MustGet())
	assert.Nil(t, opt.Scan("1"))
	assert.Equal(t, "1", opt.MustGet())

	sc := (sql.Scanner)(&opt)
	assert.NotNil(t, &sc)
}

func TestOptionalStringString(t *testing.T) {
	assert.Equal(t, emptyStringString, fmt.Sprintf("%s", OfString()))
	assert.Equal(t, "OptionalString (1)", fmt.Sprintf("%s", OfString("1")))
}

func TestOptionalStringValue(t *testing.T) {
	val, err := OfString().Value()
	assert.Nil(t, val)
	assert.Nil(t, err)

	val, err = OfString("0").Value()
	assert.Equal(t, "0", val)
	assert.Nil(t, err)
}
