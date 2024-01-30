// Package opt provides an optional type for expressing the possibility of a
// value being present or not. It's more ergonomic than pointers as it forces
// you to use accessor methods instead of dereferencing, removing the risk of a
// nil pointer dereference.
//
// There are also utilities for mapping the optional type to another type if it
// is present, constructing optional types from pointers and getting a pointer
// to the optional value or nil for easy usage with APIs that accept pointers.
//
// Note that some APIs are not implemented as methods due to the way that Go's
// generics are designed. This may change in future Go versions if it becomes
// possible to write a method of some type T which takes additional type params.
package opt

import (
	"encoding/json"
	"fmt"
)

// Optional wraps the type `T` within a container which provides abstractions
// to make conditional access and transformation of optional types easier.
type Optional[T any] container[T]

// -
// Constructors
// -

// NewEmpty creates an empty optional of the specified type `T`.
func NewEmpty[T any]() Optional[T] { return nil }

// New wraps the input value in an optional type.
func New[T any](value T) Optional[T] { return Optional[T]{value} }

// NewMap wraps `v` after applying `fn` and producing a new type `R`.
func NewMap[T, R any](v T, fn func(T) R) Optional[R] {
	return New(fn(v))
}

// NewSafe works with common "safe" APIs that return (T, boolean)
func NewSafe[T any](value T, ok bool) Optional[T] {
	if ok {
		return Optional[T]{value}
	}
	return nil
}

// NewIf wraps `v` if `fn` returns true. Useful for sanitisation of input such
// as trimming spaces and treating empty strings as none.
func NewIf[T any](v T, fn func(T) bool) Optional[T] {
	if fn(v) {
		return New(v)
	}
	return NewEmpty[T]()
}

// NewPtr wraps the input if it's non-nil, otherwise returns an empty optional.
func NewPtr[T any](ptr *T) Optional[T] {
	if ptr == nil {
		return NewEmpty[T]()
	}
	return New(*ptr)
}

// NewPtrMap wraps the input if it's non-nil and applies the transformer
// function, otherwise returns an empty optional.
func NewPtrMap[T, R any](ptr *T, fn func(T) R) Optional[R] {
	if ptr == nil {
		return NewEmpty[R]()
	}
	return New(fn(*ptr))
}

// NewPtrIf is the same as `NewIf` except will return empty if `ptr` is nil.
func NewPtrIf[T any](ptr *T, fn func(T) bool) Optional[T] {
	if ptr == nil {
		return NewEmpty[T]()
	}
	if !fn(*ptr) {
		return NewEmpty[T]()
	}
	return New(*ptr)
}

// NewPtrOr wraps the input if it's non-nil, otherwise returns a fallback value.
func NewPtrOr[T any](ptr *T, fallback T) Optional[T] {
	if ptr == nil {
		return New(fallback)
	}
	return New(*ptr)
}

// -
// Accessors
// -

// Ok returns true if there's a value inside.
func (o Optional[T]) Ok() bool {
	return o != nil
}

// Get returns the wrapped value if it's present, `ok` signals existence.
func (o Optional[T]) Get() (value T, ok bool) {
	if o == nil {
		return
	}
	return o[0], true
}

// GetMap returns the wrapped value if it's present and applies `fn`.
func GetMap[In, Out any](in Optional[In], fn func(In) Out) (v Out, ok bool) {
	if in == nil {
		return
	}
	return fn(in[0]), true
}

// GetMapC is the curried version of GetMap. See `Map` for an example.
func GetMapC[In, Out any](fn func(In) Out) func(in Optional[In]) (v Out, ok bool) {
	return func(in Optional[In]) (v Out, ok bool) {
		return GetMap(in, fn)
	}
}

// Ptr turns an optional value into a pointer to that value or nil.
func (o Optional[T]) Ptr() *T {
	if o == nil {
		return nil
	}
	return &o[0]
}

// PtrMap turns an optional value into a pointer to that value then transforms
// the value to a new type using the given transformer function.
func PtrMap[In any, Out any](o Optional[In], fn func(In) Out) *Out {
	if val, ok := o.Get(); ok {
		v := fn(val)
		return &v
	}
	return nil
}

// PtrMapC is the curried version of PtrMap. See `Map` for an example.
func PtrMapC[In any, Out any](fn func(In) Out) func(o Optional[In]) *Out {
	return func(in Optional[In]) *Out {
		return PtrMap(in, fn)
	}
}

// Or returns the underlying value or `v`.
func (o Optional[T]) Or(v T) (value T) {
	if o == nil {
		return v
	}

	return o[0]
}

// Call calls `fn` if there is a value wrapped by this optional.
func (o Optional[T]) Call(f func(value T)) {
	if o != nil {
		f(o[0])
	}
}

// OrCall calls `fn` if the optional is empty.
func (o Optional[T]) OrCall(fn func() T) (value T) {
	if o != nil {
		return o[0]
	}

	return fn()
}

// OrZero returns the zero value of `T` if it's not present.
func (o Optional[T]) OrZero() (value T) {
	if o == nil {
		var zero T
		return zero
	}

	return o[0]
}

// -
// Conditional pipelines
// -

// Map calls `fn` on `in` if it's present and returns the new optional value.
func Map[In, Out any](in Optional[In], fn func(In) Out) (v Optional[Out]) {
	if in == nil {
		return
	}

	return Optional[Out]{fn(in[0])}
}

// MapC is the curried version of Map. It's useful for applying the same mapping
// function to many values either manually or using a functor library.
//
// First, construct the function with the transformer.
//
//	fn := Map(transformer)
//
// Now you can use `fn` to transform any optional value using the transformer.
//
//	b := fn(a)
//	dt.Map(array, fn)
func MapC[In, Out any](fn func(In) Out) func(in Optional[In]) (v Optional[Out]) {
	return func(in Optional[In]) (v Optional[Out]) {
		return Map(in, fn)
	}
}

// MapErr calls `fn` on `in` if it's present and returns the result or an error.
func MapErr[In, Out any](in Optional[In], fn func(In) (Out, error)) (v Optional[Out], err error) {
	if in == nil {
		return
	}

	out, err := fn(in[0])
	if err != nil {
		return nil, err
	}

	return Optional[Out]{out}, nil
}

// MapErrC calls `fn` on `in` if it's present and returns the result or an error.
func MapErrC[In, Out any](fn func(In) (Out, error)) func(in Optional[In]) (v Optional[Out], err error) {
	return func(in Optional[In]) (v Optional[Out], err error) {
		return MapErr(in, fn)
	}
}

// -
// Utilities
// -

// String returns the string representation of the value or an empty string.
func (o Optional[T]) String() string {
	if v, ok := o.Get(); ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// GoString is only used for verbose printing.
func (o Optional[T]) GoString() string {
	if v, ok := o.Get(); ok {
		return fmt.Sprintf("Optional[%v]", v)
	}
	return "Optional[]"
}

// MarshalJSON marshals the value being wrapped to JSON. If there is no vale
// being wrapped, the zero value of its type is marshaled.
func (o Optional[T]) MarshalJSON() (data []byte, err error) {
	if v, ok := o.Get(); ok {
		return json.Marshal(v)
	}
	return []byte("null"), nil
}

// UnmarshalJSON unmarshals the JSON into a value wrapped by this optional.
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*o = NewEmpty[T]()
		return nil
	}

	var v T
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	*o = New(v)
	return nil
}

// container is an internal type that hides and holds the underlying value.
type container[T any] []T
