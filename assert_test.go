package assert_test

import (
	"fmt"
	"io/fs"
	"testing"
	"time"

	"github.com/andreyvit/assert"
)

func TestOK(t *testing.T) {
	assert.OK(noerr(t), true)
	assert.OK(fake(t, "** got false, wanted true"), false)
}

func TestOK_msg(t *testing.T) {
	assert.OK(noerr(t), true, "with %v", panicOnString{})
	assert.OK(fake(t, "** with 42: got false, wanted true"), false, "with %v", 42)
}

func TestFalse(t *testing.T) {
	assert.False(noerr(t), false)
	assert.False(fake(t, "** got true, wanted false"), true)
}

func TestEq(t *testing.T) {
	assert.Eq(noerr(t), 42, 42)
	assert.Eq(fake(t, "** got 10, wanted 42"), 10, 42)
}

func TestEq_floats(t *testing.T) {
	assert.Eq(noerr(t), 42.2, 42.2)
	assert.Eq(fake(t, "** got 42, wanted 42.2"), 42.0, 42.2)
}

func TestNotEq(t *testing.T) {
	assert.NotEq(noerr(t), 12, 42)
	assert.NotEq(fake(t, "** got 42, wanted anything else"), 42, 42)
}

func TestDeepEqual(t *testing.T) {
	assert.DeepEqual(noerr(t), []int{12, 34}, []int{12, 34})
	assert.DeepEqual(fake(t, "** got [42], wanted [12 34]"), []int{42}, []int{12, 34})
}

func TestNotDeepEqual(t *testing.T) {
	assert.NotDeepEqual(noerr(t), []int{12, 34}, []int{12})
	assert.NotDeepEqual(fake(t, "** got [12 34], wanted anything else"), []int{12, 34}, []int{12, 34})
}

func TestNil(t *testing.T) {
	assert.Nil(noerr(t), (*int)(nil))
	v := 42
	assert.Nil(fake(t, "** got &42, wanted nil"), &v)
}

func TestNonNil(t *testing.T) {
	v := 42
	assert.NonNil(noerr(t), &v)
	assert.NonNil(fake(t, "** got nil *int, wanted non-nil"), (*int)(nil))
}

func TestZero(t *testing.T) {
	assert.Zero(noerr(t), 0)
	assert.Zero(noerr(t), 0.0)
	assert.Zero(noerr(t), false)
	assert.Zero(noerr(t), "")
	assert.Zero(noerr(t), time.Time{})
	assert.Zero(noerr(t), any(nil))
	assert.Zero(fake(t, "** got 42, wanted zero value 0"), 42)
	assert.Zero(fake(t, "** got 42.2, wanted zero value 0"), 42.2)
	assert.Zero(fake(t, "** got true, wanted zero value false"), true)
	assert.Zero(fake(t, "** got abc, wanted zero value "), "abc")
	assert.Zero(fake(t, "** got 2023-02-05 00:00:00 +0000 UTC, wanted zero value 0001-01-01 00:00:00 +0000 UTC"), time.Date(2023, 2, 5, 0, 0, 0, 0, time.UTC))
	assert.Zero(fake(t, "** got 42, wanted zero value <nil>"), any(42))
}

func TestNonZero(t *testing.T) {
	assert.NonZero(noerr(t), 42)
	assert.NonZero(noerr(t), 42.2)
	assert.NonZero(noerr(t), true)
	assert.NonZero(noerr(t), "abc")
	assert.NonZero(noerr(t), time.Date(2023, 2, 5, 0, 0, 0, 0, time.UTC))
	assert.NonZero(noerr(t), any(42))
	assert.NonZero(fake(t, "** got zero value 0, wanted non-zero"), 0)
	assert.NonZero(fake(t, "** got zero value 0, wanted non-zero"), 0.0)
	assert.NonZero(fake(t, "** got zero value false, wanted non-zero"), false)
	assert.NonZero(fake(t, "** got zero value , wanted non-zero"), "")
	assert.NonZero(fake(t, "** got zero value 0001-01-01 00:00:00 +0000 UTC, wanted non-zero"), time.Time{})
	assert.NonZero(fake(t, "** got zero value <nil>, wanted non-zero"), any(nil))
}

func TestEmptySlice(t *testing.T) {
	assert.EmptySlice(noerr(t), []int(nil))
	assert.EmptySlice(noerr(t), []int{})
	assert.EmptySlice(fake(t, "** got [42], wanted empty slice"), []int{42})
}

func TestNonEmptySlice(t *testing.T) {
	assert.NonEmptySlice(noerr(t), []int{42})
	assert.NonEmptySlice(fake(t, "** got empty []int, wanted non-empty"), []int{})
}

func TestEmptyMap(t *testing.T) {
	assert.EmptyMap(noerr(t), map[int]string(nil))
	assert.EmptyMap(noerr(t), map[int]string{})
	assert.EmptyMap(fake(t, "** got map[42:x], wanted empty map"), map[int]string{42: "x"})
}

func TestNonEmptyMap(t *testing.T) {
	assert.NonEmptyMap(noerr(t), map[int]string{42: "x"})
	assert.NonEmptyMap(fake(t, "** got empty map[int]string, wanted non-empty"), map[int]string{})
}

func TestSuccess(t *testing.T) {
	assert.Success(noerr(t), error(nil))
	assert.Success(fake(t, "** failed: file does not exist"), fs.ErrNotExist)
}

func TestError_ok(t *testing.T) {
	assert.Error(noerr(t), fs.ErrNotExist, fs.ErrNotExist)
}
func TestError_assert_success(t *testing.T) {
	assert.Error(noerr(t), error(nil), nil)
	assert.Error(fake(t, "** failed: file does not exist"), fs.ErrNotExist, nil)
}
func TestError_wrong_error(t *testing.T) {
	assert.Error(fake(t, "** failed with: file already exists, wanted: file does not exist"), fs.ErrExist, fs.ErrNotExist)
}
func TestError_unexpected_success(t *testing.T) {
	assert.Error(fake(t, "** succeeded, wanted to fail with: file does not exist"), nil, fs.ErrNotExist)
}

func TestErrorMsg_ok(t *testing.T) {
	assert.ErrorMsg(noerr(t), fs.ErrNotExist, "file does not exist")
}
func TestErrorMsg_assert_success(t *testing.T) {
	assert.ErrorMsg(noerr(t), error(nil), "")
	assert.ErrorMsg(fake(t, "** failed: file does not exist"), fs.ErrNotExist, "")
}
func TestErrorMsg_wrong_error(t *testing.T) {
	assert.ErrorMsg(fake(t, "** failed with: file already exists, wanted: file does not exist"), fs.ErrExist, "file does not exist")
}
func TestErrorMsg_unexpected_success(t *testing.T) {
	assert.ErrorMsg(fake(t, "** succeeded, wanted to fail with: file does not exist"), nil, "file does not exist")
}

func TestPanicMsg_ok(t *testing.T) {
	assert.PanicMsg(noerr(t), panickyFunc, "runtime error: index out of range [2] with length 2")
}
func TestPanicMsg_unexpected_success(t *testing.T) {
	assert.PanicMsg(fake(t, "** succeeded, wanted to panic with: foo"), func() {}, "foo")
}
func TestPanicMsg_wrong_error(t *testing.T) {
	assert.PanicMsg(fake(t, "** paniced with: runtime error: index out of range [2] with length 2, wanted: foo"), panickyFunc, "foo")
}

func TestFormatPrefix_empty(t *testing.T) {
	assert.Eq(t, assert.FormatPrefix(nil), "")
}
func TestFormatPrefix_basic(t *testing.T) {
	assert.Eq(t, assert.FormatPrefix([]any{"foo"}), "foo: ")
	assert.Eq(t, assert.FormatPrefix([]any{"foo v=%v", 42}), "foo v=42: ")
}
func TestFormatPrefix_panic(t *testing.T) {
	assert.PanicMsg(t, func() {
		assert.FormatPrefix([]any{42})
	}, "when passing messageAndArgs to assertion funcs, the first extra argument must be a format string, got int 42")
}

func TestAddPrefix_empty(t *testing.T) {
	assert.DeepEqual(t, assert.AddPrefix(nil, "foo.%d", 42), []any{"foo.%d", 42})
}
func TestAddPrefix_existing(t *testing.T) {
	assert.DeepEqual(t, assert.AddPrefix([]any{"bar %s", "boz"}, "foo.%d", 42), []any{"foo.%d: bar %s", 42, "boz"})
}
func TestAddPrefix_notstr(t *testing.T) {
	assert.DeepEqual(t, assert.AddPrefix([]any{1, 2, 3}, "foo.%d", 42), []any{"foo.%d", 42, 1, 2, 3})
}

func panickyFunc() {
	a := make([]int, 2)
	a[2] = 42
}

type noErrTB struct {
	t testing.TB
}

func noerr(t testing.TB) assert.TB {
	return &noErrTB{t}
}
func (f *noErrTB) Helper() {
}
func (f *noErrTB) Errorf(format string, args ...any) {
	f.t.Helper()
	f.t.Fatalf("unexpected error: %s", fmt.Sprintf(format, args...))
}

type fakeTB struct {
	t            testing.TB
	helperCalled bool
	errorCalled  bool
	expected     string
}

func fake(t testing.TB, expected string) assert.TB {
	f := &fakeTB{
		t:        t,
		expected: expected,
	}
	t.Cleanup(f.verify)
	return f
}

func (f *fakeTB) verify() {
	f.t.Helper()
	if !f.errorCalled {
		f.t.Fatal("Errorf not called")
	}
	if !f.helperCalled {
		f.t.Fatal("Helper not called")
	}
}

func (f *fakeTB) Helper() {
	f.t.Helper()
	f.helperCalled = true
}

func (f *fakeTB) Errorf(format string, args ...any) {
	f.t.Helper()
	actual := fmt.Sprintf(format, args...)
	if actual != f.expected {
		f.t.Fatalf("incorrect error message, got:\n\t%s\nwanted:\n\t%s", actual, f.expected)
	}
	f.errorCalled = true
}

type panicOnString struct{}

func (_ panicOnString) String() string {
	panic("String() should not be called")
}
