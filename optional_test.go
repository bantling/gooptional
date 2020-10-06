package gooptional

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	"github.com/bantling/goiter"
	"github.com/stretchr/testify/assert"
)

func TestOptionalOfEmptyPresentGet(t *testing.T) {
	var (
		c chan bool
		f func()
		m map[string]int
		p *int
		s []int
		e error

		val   interface{}
		valid bool
	)

	for _, opt := range []Optional{Of(), Of(c), Of(f), Of(m), Of(p), Of(s), Of(e)} {
		assert.Nil(t, opt.value)
		assert.False(t, opt.present)
		assert.True(t, opt.IsEmpty())
		assert.False(t, opt.IsPresent())
		called := false
		opt.IfPresent(func(interface{}) { called = true })
		assert.False(t, called)
		opt.IfEmpty(func() { called = true })
		assert.True(t, called)
		called = false
		opt.IfPresentOrElse(func(interface{}) {}, func() { called = true })
		assert.True(t, called)

		func() {
			defer func() {
				assert.True(t, errNotPresent == recover())
			}()

			opt.MustGet()
			assert.Fail(t, "Expected Panic")
		}()
	}

	opt := Of(0)
	assert.Equal(t, 0, opt.value)
	assert.True(t, opt.present)
	assert.False(t, opt.IsEmpty())
	assert.True(t, opt.IsPresent())
	val = 1
	opt.IfPresent(func(v interface{}) { val = v })
	assert.Equal(t, 0, val)
	val = 1
	opt.IfEmpty(func() { val = 0 })
	assert.Equal(t, 1, val)
	val = 1
	opt.IfPresentOrElse(func(v interface{}) { val = 2 }, func() { val = 3 })
	assert.Equal(t, 2, val)

	val, valid = opt.Get()
	assert.Equal(t, 0, val)
	assert.True(t, valid)
	assert.Equal(t, 0, opt.MustGet())

	opt = Of("")
	assert.Equal(t, "", opt.value)
	assert.True(t, opt.present)

	// Test zero value
	var zval Optional
	assert.Nil(t, zval.value)
	assert.False(t, zval.present)
	assert.True(t, zval.IsEmpty())
	assert.False(t, zval.IsPresent())
	called := false
	zval.IfPresent(func(interface{}) { called = true })
	assert.False(t, called)
	zval.IfEmpty(func() { called = true })
	assert.True(t, called)
	called = false
	zval.IfPresentOrElse(func(interface{}) {}, func() { called = true })
	assert.True(t, called)
	func() {
		defer func() {
			assert.True(t, errNotPresent == recover())
		}()

		zval.MustGet()
		assert.Fail(t, "Expected Panic")
	}()
}

func TestOptionalIter(t *testing.T) {
	var (
		opt      Optional        = Of(1)
		iterable goiter.Iterable = opt
		iter                     = iterable.Iter()
	)
	assert.True(t, iter.Next())
	assert.Equal(t, 1, iter.Value())
	assert.False(t, iter.Next())
}

func TestOptionalEqual(t *testing.T) {
	// Not present optional == not present optional
	assert.True(t, Of().Equal(Of()))

	opt1 := Of(0)
	opt2 := Of(0)
	// Present optional != not present optional
	assert.False(t, opt1.Equal(Of()))

	// Present optional == present optional if values equal
	assert.True(t, opt1.Equal(opt1))
	assert.True(t, opt1.Equal(opt2))

	// Present optional != present optional if values differ
	assert.False(t, opt1.Equal(Of(1)))

	// Not present optional never equals any value
	assert.False(t, Of().EqualValue(nil))

	// Not present optional != non-nil value
	assert.False(t, Of().EqualValue(0))

	// Present optional == value if values equal
	assert.True(t, opt1.EqualValue(0))

	// Present optional != value if values differ
	assert.False(t, opt1.Equal(Of(1)))
}

func TestOptionalNotEqual(t *testing.T) {
	// Not present optional == not present optional
	assert.False(t, Of().NotEqual(Of()))

	opt1 := Of(0)
	opt2 := Of(0)
	// Present optional != not present optional
	assert.True(t, opt1.NotEqual(Of()))

	// Present optional == present optional if values equal
	assert.False(t, opt1.NotEqual(opt1))
	assert.False(t, opt1.NotEqual(opt2))

	// Present optional != present optional if values differ
	assert.True(t, opt1.NotEqual(Of(1)))

	// Not present optional never equals any value
	assert.True(t, Of().NotEqualValue(nil))

	// Not present optional != non-nil value
	assert.True(t, Of().NotEqualValue(0))

	// Present optional == value if values equal
	assert.False(t, opt1.NotEqualValue(0))

	// Present optional != value if values differ
	assert.True(t, opt1.NotEqual(Of(1)))
}

func TestOptionalFilter(t *testing.T) {
	opt := Of(1)
	assert.True(t, opt == opt.Filter(func(val interface{}) bool { return true }))
	assert.True(t, opt.Filter(func(interface{}) bool { return false }).IsEmpty())

	assert.True(t, Of().Filter(func(interface{}) bool { return true }).IsEmpty())
}

func TestOptionalFilterNot(t *testing.T) {
	opt := Of(1)
	assert.True(t, opt == opt.FilterNot(func(val interface{}) bool { return false }))
	assert.True(t, opt.FilterNot(func(interface{}) bool { return true }).IsEmpty())

	assert.True(t, Of().FilterNot(func(interface{}) bool { return false }).IsEmpty())
}

func TestOptionalFlatMapFloatIntString(t *testing.T) {
	too := func(val interface{}) Optional {
		return Of(val.(int) + 1)
	}
	assert.True(t, Of().FlatMap(too).IsEmpty())
	assert.Equal(t, 2, Of(1).FlatMap(too).MustGet())

	tof := func(val interface{}) OptionalFloat {
		return OfFloat(float64(val.(int) + 1))
	}
	assert.True(t, Of().FlatMapToFloat(tof).IsEmpty())
	assert.Equal(t, 2.0, Of(1).FlatMapToFloat(tof).MustGet())

	toi := func(val interface{}) OptionalInt {
		return OfInt(val.(int) + 1)
	}
	assert.True(t, Of().FlatMapToInt(toi).IsEmpty())
	assert.Equal(t, 2, Of(1).FlatMapToInt(toi).MustGet())

	tos := func(val interface{}) OptionalString {
		return OfString(strconv.Itoa(val.(int) + 1))
	}
	assert.True(t, Of().FlatMapToString(tos).IsEmpty())
	assert.Equal(t, "2", Of(1).FlatMapToString(tos).MustGet())
}

func TestOptionalMapFloatIntString(t *testing.T) {
	too := func(val interface{}) interface{} {
		return val.(int) + 1
	}
	assert.True(t, Of().Map(too).IsEmpty())
	assert.Equal(t, 2, Of(1).Map(too).MustGet())

	tonp := func(val interface{}) interface{} {
		return nil
	}
	assert.True(t, Of(1).Map(tonp).IsEmpty())

	toz := func(val interface{}) interface{} {
		return 0
	}
	assert.False(t, Of(1).Map(toz).IsEmpty())
	assert.True(t, Of(1).Map(toz, true).IsEmpty())

	tof := func(val interface{}) float64 {
		return float64(val.(int) + 1)
	}
	assert.True(t, Of().MapToFloat(tof).IsEmpty())
	assert.Equal(t, 2.0, Of(1).MapToFloat(tof).MustGet())

	toi := func(val interface{}) int {
		return val.(int) + 1
	}
	assert.True(t, Of().MapToInt(toi).IsEmpty())
	assert.Equal(t, 2, Of(1).MapToInt(toi).MustGet())

	tos := func(val interface{}) string {
		return strconv.Itoa(val.(int) + 1)
	}
	assert.True(t, Of().MapToString(tos).IsEmpty())
	assert.Equal(t, "2", Of(1).MapToString(tos).MustGet())
}

func TestOptionalOrElseGetPanic(t *testing.T) {
	f := func() interface{} { return 2 }
	assert.Equal(t, 1, Of().OrElse(1))
	assert.Equal(t, 2, Of().OrElseGet(f))

	err := fmt.Errorf("")
	errf := func() error { return err }
	func() {
		defer func() {
			assert.True(t, err == recover())
		}()
		Of().OrElsePanic(errf)
		assert.Fail(t, "Expected Panic")
	}()

	assert.Equal(t, 3, Of(3).OrElse(1))
	assert.Equal(t, 3, Of(3).OrElseGet(f))
	assert.Equal(t, 3, Of(3).OrElsePanic(errf))
}

func TestOptionalScan(t *testing.T) {
	var opt Optional
	assert.Nil(t, opt.Scan(0))
	assert.Equal(t, 0, opt.MustGet())
	assert.Nil(t, opt.Scan(1))
	assert.Equal(t, 1, opt.MustGet())

	sc := (sql.Scanner)(&opt)
	assert.NotNil(t, &sc)
}

type OptionalT int

func (t OptionalT) String() string {
	return fmt.Sprintf("%d", t+1)
}

func TestOptionalString(t *testing.T) {
	assert.Equal(t, emptyString, fmt.Sprintf("%s", Of()))
	assert.Equal(t, "Optional (1)", fmt.Sprintf("%s", Of(1)))
	assert.Equal(t, "Optional (2)", fmt.Sprintf("%s", Of((OptionalT)(1))))
}

func TestOptionalValue(t *testing.T) {
	val, err := Of().Value()
	assert.Nil(t, val)
	assert.Nil(t, err)

	val, err = Of(0).Value()
	assert.Equal(t, 0, val)
	assert.Nil(t, err)
}
