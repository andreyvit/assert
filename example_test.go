package assert_test

import (
	"io/fs"
	"math"
	"testing"

	"github.com/andreyvit/assert"
)

type intPtr *int

func Example() {
	var t *testing.T // this is an argument of your test method
	var (
		found        bool
		count        int
		pointer      *int
		typedPointer intPtr
		someSlice    []int
		someMap      map[int]int
		value        any
		err          error
	)

	assert.OK(t, found)
	assert.False(t, found)

	assert.Eq(t, count, 42)
	assert.NotEq(t, count, 42)

	assert.DeepEqual(t, someSlice, []int{42})
	assert.NotDeepEqual(t, someSlice, []int{42})

	assert.Nil(t, pointer)
	assert.NonNil(t, pointer)
	assert.Nil(t, typedPointer)
	assert.NonNil(t, typedPointer)

	assert.Zero(t, value)
	assert.NonZero(t, value)

	assert.EmptySlice(t, someSlice)
	assert.NonEmptySlice(t, someSlice)
	assert.EmptyMap(t, someMap)
	assert.NonEmptyMap(t, someMap)

	assert.Success(t, err)
	assert.Error(t, err, fs.ErrNotExist)
	assert.ErrorMsg(t, err, "file does not exist")
	assert.PanicMsg(t, panickyFunc, "runtime error: index out of range [2] with length 2")
}

func TestApproxEq(t *testing.T) {
	approxEq(noerr(t), 1, 1.0000001)
	approxEq(fake(t, "** got 1, wanted 1.0001 ± 1e-06"), 1, 1.0001)
}

func approxEq(t assert.TB, a, e float64, messageAndArgs ...any) bool {
	const eps = 1e-6
	if math.Abs(a-e) > eps {
		t.Helper()
		t.Errorf("** %sgot %v, wanted %v ± %v", assert.FormatPrefix(messageAndArgs), a, e, eps)
		return false
	}
	return true
}
