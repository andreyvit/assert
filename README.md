Go assertion library in stdlib style
====================================

[![Go reference](https://pkg.go.dev/badge/github.com/andreyvit/assert.svg)](https://pkg.go.dev/github.com/andreyvit/assert) ![Zero dependencies](https://img.shields.io/badge/deps-zero-brightgreen) ![Zero magic](https://img.shields.io/badge/magic-none-brightgreen) ![under 200 LOC](https://img.shields.io/badge/size-%3C200%20LOC-green) ![100% coverage](https://img.shields.io/badge/coverage-100%25-green) [![Go Report Card](https://goreportcard.com/badge/github.com/andreyvit/assert)](https://goreportcard.com/report/github.com/andreyvit/assert)


Why?
----

Make tests easier to understand by replacing:

```go
if a, e := Foo(), 123; a != e {
    t.Errorf("** Foo() = %v, wanted %v", a, e)
}
if !Bar() {
    t.Errorf("** Bar is false")
}
```

with:

```go
assert.Eq(t, Foo(), 123)
assert.OK(t, Bar())
```


Usage
-----

```go
assert.OK(t, found)
assert.False(t, found)

assert.Eq(t, count, 42)
assert.NotEq(t, count, 42)

assert.DeepEqual(t, someSlice, []int{42})
assert.NotDeepEqual(t, someSlice, []int{42})

assert.Nil(t, pointer)
assert.NonNil(t, pointer)

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
```

All of these functions return true/false to chain assertions:

```go
if assert.NonNil(t, obj) {
    assert.Eq(t, obj.Foo, 42)
}
```

You can pass an extra message with Printf-style arguments to each:

```go
assert.Eq(t, items[i].Foo, 42, "items[%d].Foo", i)
assert.Eq(t, items[i].Bar, 10, "items[%d].Bar", i)
```


Why this library?
-----------------

There are many alternatives. Here's why this is the best one.

1. Tiny. A small set of straightforward assertions with simple code, not a port of some crazy BDD framework or assert.IsValidEmail-million-assertions library. We include what we think Go team would choose to include if this was a standard library package.
2. Zero magic. We use Go generics for type safety instead of the typical fuzzy comparison logic. (E.g. consider `assert.Eq(t, value, 42)` when value is int64. If the arguments are typed as `any`, they won't be equal. Either you have to write `int64(42)`, or the library needs runtime fuzzy logic to treat all integer types as equal.)
3. Allows adding extra message prefix for test failures.
4. Allows writing your custom assertions in the same style.
5. Follows (actual, expected) argument order customary for Go (and uses got/wanted message wording, too). Duh, but some otherwise reasonable libs mess that up.
6. Zero dependencies, 100% coverage.


What if I want an assertion that you don't have?
------------------------------------------------

Use our code structure and `assert.FormatPrefix` to write your own assertion as part of your test suite.

Example 1. Floating-point equality with tolerance:

```go
func approxEq(t testing.TB, a, e float64, messageAndArgs ...any) bool {
    const eps = 1e-6  // NOT a gospel or magic value; adjust for your project!
    if math.Abs(a-e) > eps {
        t.Helper()
        t.Errorf("** %sgot %v, wanted %v ± %v", assert.FormatPrefix(messageAndArgs), a, e, eps)
        return false
    }
    return true
}
```

Example 2. Smart equality check via [`github.com/google/go-cmp`](https://github.com/google/go-cmp), for when you want to compare large structs and ignore some fields, etc:

```go
func smartEq(t testing.TB, a, e any, messageAndArgsAndCmpOptions ...any) bool {
    var cmpOptions []cmp.Option
    var messageAndArgs []any
    for _, v := range messageAndArgsAndCmpOptions {
        if opt, ok := v.(cmp.Option); ok {
            cmpOptions = append(cmpOptions, opt)
        } else if opts, ok := v.([]cmp.Option); ok {
            cmpOptions = append(cmpOptions, opts...)
        } else {
            messageAndArgs = append(messageAndArgs, v)
        }
    }

    t.Helper()
    if diff := cmp.Diff(e, a, cmpOptions...); diff != "" {
        t.Errorf("** %snot equal (-want +got): %s", assert.FormatPrefix(messageAndArgs), diff)
        return false
    }
    return true
}
```


Contributing
------------

“We include what we think Go team would choose to include if this was a standard library package.”

We accept contributions that:

* add better documentation and examples;
* fix bugs;
* add extra assertions WITH a great argument for why they need to be included.

Out of scope:

* assertions that don't add much value for a typical test suite;
* assertions that don't have a single obvious best implementation or might need slightly different takes for different projects (e.g. floating-point assertions or string diffing);
* anything that requires third-party dependencies;
* exported helpers and utility methods unless they add a LOT of value for typical test suites.

We recommend [modd](https://github.com/cortesi/modd) (`go install github.com/cortesi/modd/cmd/modd@latest`) for continuous testing during development.

Maintain 100% coverage. It's not often the right choice, but it is for this library.

Ensure that all assertions produce clear, meaningful and helpful error messages.

TODO:

- [ ] Figure out if there's a way to build universal `assert.Empty` to replace `EmptyMap` and `EmptySlice` that meets all expectations that users would immediately assign to such function (i.e. works on strings, numbers, structs, etc). Or maybe this can be built without reflection to just work on the types we explicitly support? Do we even want a func like that, or is that too much magic? Alternatively, choose to extend `assert.Zero` to consider empty slices and maps to be zero values too.


FAQ
---

### Why do you require Go 1.20?

It's generics: we need `any` values passed in to be treated as `comparable`.


MIT license
-----------

Copyright (c) 2023 Andrey Tarantsov. Published under the terms of the [MIT license](LICENSE).
