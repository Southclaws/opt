package opt

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	a := assert.New(t)
	v := "value"

	a.Empty(NewEmpty[string]())
	a.NotEmpty(New(v))
	a.Equal("value", New(v).String())
	a.Equal("Optional[value]", New(v).GoString())
	a.Equal("Optional[]", NewEmpty[string]().GoString())
	a.Equal("VALUE", NewMap(v, strings.ToUpper).String())
	a.Equal("value", NewSafe(v, true).String())
	a.Equal("", NewSafe(v, false).String())
	a.Equal("value", NewIf(v, func(v string) bool { return true }).String())
	a.Equal("", NewIf(v, func(v string) bool { return false }).String())
	a.Equal("value", NewPtr(&v).String())
	a.Equal("", NewPtr[string](nil).String())
	a.Equal("value", NewPtrOr(&v, "fallback").String())
	a.Equal("fallback", NewPtrOr[string](nil, "fallback").String())
	a.Equal("VALUE", NewPtrMap(&v, strings.ToUpper).String())
	a.Equal("", NewPtrMap(nil, strings.ToUpper).String())
	a.Equal("value", NewPtrIf(&v, func(v string) bool { return true }).String())
	a.Equal("", NewPtrIf(nil, func(v string) bool { return true }).String())
	a.Equal("", NewPtrIf(&v, func(v string) bool { return false }).String())
	a.Equal("", NewPtrIf(nil, func(v string) bool { return false }).String())
}

func TestOk(t *testing.T) {
	s := "ptr to string"
	tests := []struct {
		Optional          Optional[string]
		ExpectedIsPresent bool
	}{
		{NewEmpty[string](), false},
		{New(""), true},
		{New("string"), true},
		{NewPtr((*string)(nil)), false},
		{NewPtr((*string)(&s)), true},
	}

	for _, test := range tests {
		isPresent := test.Optional.Ok()

		if isPresent != test.ExpectedIsPresent {
			t.Errorf("%#v Ok got %#v, want %#v", test.Optional, isPresent, test.ExpectedIsPresent)
		}
	}
}

func TestGet(t *testing.T) {
	s := "ptr to string"
	tests := []struct {
		Optional      Optional[string]
		ExpectedValue string
		ExpectedOk    bool
	}{
		{NewEmpty[string](), "", false},
		{New(""), "", true},
		{New("string"), "string", true},
		{NewPtr((*string)(nil)), "", false},
		{NewPtr((*string)(&s)), "ptr to string", true},
	}

	for _, test := range tests {
		value, ok := test.Optional.Get()

		if value != test.ExpectedValue || ok != test.ExpectedOk {
			t.Errorf("%#v Get got %#v, %#v, want %#v, %#v", test.Optional, ok, test.ExpectedOk, value, test.ExpectedValue)
		}
	}
}

func TestGetMap(t *testing.T) {
	a := assert.New(t)

	in := New("value")
	out, ok := GetMap(in, strings.ToUpper)
	a.True(ok)
	a.Equal("VALUE", out)

	in = NewEmpty[string]()
	out, ok = GetMap(in, strings.ToUpper)
	a.False(ok)
	a.Equal("", out)
}

func TestPtr(t *testing.T) {
	a := assert.New(t)

	in := New("value")
	out := in.Ptr()
	a.NotNil(a)
	a.Equal("value", *out)

	in = NewEmpty[string]()
	out = in.Ptr()
	a.Nil(out)
}

func TestPtrMap(t *testing.T) {
	a := assert.New(t)

	in := New("value")

	out := PtrMap(in, strings.ToUpper)
	a.Equal("VALUE", *out)

	in = NewEmpty[string]()
	out = PtrMap(in, strings.ToUpper)
	a.Nil(out)
}

func TestOr(t *testing.T) {
	s := "ptr to string"
	const orElse = "orelse"
	tests := []struct {
		Optional       Optional[string]
		ExpectedResult string
	}{
		{NewEmpty[string](), orElse},
		{New(""), ""},
		{New("string"), "string"},
		{NewPtr((*string)(nil)), orElse},
		{NewPtr((*string)(&s)), "ptr to string"},
	}

	for _, test := range tests {
		result := test.Optional.Or(orElse)

		if result != test.ExpectedResult {
			t.Errorf("%#v OrElse(%#v) got %#v, want %#v", test.Optional, orElse, result, test.ExpectedResult)
		}
	}
}

func TestCall(t *testing.T) {
	s := "ptr to string"
	tests := []struct {
		Optional       Optional[string]
		ExpectedCalled bool
		IfCalledValue  string
	}{
		{NewEmpty[string](), false, ""},
		{New(""), true, ""},
		{New("string"), true, "string"},
		{NewPtr((*string)(nil)), false, ""},
		{NewPtr((*string)(&s)), true, "ptr to string"},
	}

	for _, test := range tests {
		called := false
		test.Optional.Call(func(v string) {
			called = true
			if v != test.IfCalledValue {
				t.Errorf("%#v IfPresent got %#v, want #%v", test.Optional, v, test.IfCalledValue)
			}
		})

		if called != test.ExpectedCalled {
			t.Errorf("%#v IfPresent called %#v, want %#v", test.Optional, called, test.ExpectedCalled)
		}
	}
}

func TestOrCall(t *testing.T) {
	s := "ptr to string"
	const orElse = "orelse"
	tests := []struct {
		Optional       Optional[string]
		ExpectedResult string
	}{
		{NewEmpty[string](), orElse},
		{New(""), ""},
		{New("string"), "string"},
		{NewPtr((*string)(nil)), orElse},
		{NewPtr((*string)(&s)), "ptr to string"},
	}

	for _, test := range tests {
		result := test.Optional.OrCall(func() string { return orElse })

		if result != test.ExpectedResult {
			t.Errorf("%#v OrElse(%#v) got %#v, want %#v", test.Optional, orElse, result, test.ExpectedResult)
		}
	}
}

func TestOrZero(t *testing.T) {
	s := "ptr to string"
	tests := []struct {
		Optional       Optional[string]
		ExpectedResult string
	}{
		{NewEmpty[string](), ""},
		{New(""), ""},
		{New("string"), "string"},
		{NewPtr((*string)(nil)), ""},
		{NewPtr((*string)(&s)), "ptr to string"},
	}

	for _, test := range tests {
		result := test.Optional.OrZero()

		if result != test.ExpectedResult {
			t.Errorf("%#v OrZero() got %#v, want %#v", test.Optional, result, test.ExpectedResult)
		}
	}
}

func TestMap(t *testing.T) {
	a := assert.New(t)

	in := New("value")
	out := Map(in, strings.ToUpper)
	a.Equal("VALUE", out.String())

	in = NewEmpty[string]()
	out = Map(in, strings.ToUpper)
	a.Empty(out)
}

func TestMapErr(t *testing.T) {
	a := assert.New(t)

	valid := New("69")
	invalid := New("value")

	out, err := MapErr(valid, strconv.Atoi)
	a.NoError(err)
	a.Equal(69, out.OrZero())

	out, err = MapErr(invalid, strconv.Atoi)
	a.Error(err)
	a.Equal(0, out.OrZero())

	in := NewEmpty[string]()
	out, err = MapErr(in, strconv.Atoi)
	a.NoError(err, "because input is optionally empty, fn is never called")
	a.Empty(out)
}

func TestJSON(t *testing.T) {
	a := assert.New(t)

	type Data struct {
		ID   string
		Name Optional[string]
		Age  Optional[int]
	}

	in := Data{
		ID: "southclaws",
	}

	b1, err := json.Marshal(in)
	a.NoError(err)
	a.Equal(`{"ID":"southclaws","Name":null,"Age":null}`, string(b1))

	in.Age = New(69)

	b2, err := json.Marshal(in)
	a.NoError(err)
	a.Equal(`{"ID":"southclaws","Name":null,"Age":69}`, string(b2))

	in.Name = New("Southclaws")

	b3, err := json.Marshal(in)
	a.NoError(err)
	a.Equal(`{"ID":"southclaws","Name":"Southclaws","Age":69}`, string(b3))

	var out Data
	err = json.Unmarshal(b3, &out)
	a.NoError(err)

	a.Equal("southclaws", out.ID)
	a.Equal("Southclaws", out.Name.OrZero())
	a.Equal(69, out.Age.OrZero())

	err = json.Unmarshal(b2, &out)
	a.NoError(err)

	a.Equal("southclaws", out.ID)
	a.Empty(out.Name.OrZero())
	a.Equal(69, out.Age.OrZero())

	err = json.Unmarshal(b1, &out)
	a.NoError(err)

	a.Equal("southclaws", out.ID)
	a.Empty(out.Name)
	a.Empty(out.Age)
}
