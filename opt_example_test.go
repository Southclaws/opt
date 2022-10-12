package opt_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Southclaws/opt"
)

func Example_get() {
	i := 1001
	values := []opt.Optional[int]{
		opt.NewEmpty[int](),
		opt.New(1000),
		opt.NewPtr[int](nil),
		opt.NewPtr(&i),
	}

	for _, v := range values {
		if i, ok := v.Get(); ok {
			fmt.Println(i)
		}
	}

	// Output:
	// 1000
	// 1001
}

func Example_if() {
	i := 1001
	values := []opt.Optional[int]{
		opt.NewEmpty[int](),
		opt.New(1000),
		opt.NewPtr[int](nil),
		opt.NewPtr(&i),
	}

	for _, v := range values {
		v.If(func(i int) {
			fmt.Println(i)
		})
	}

	// Output:
	// 1000
	// 1001
}

func Example_else() {
	i := 1001
	values := []opt.Optional[int]{
		opt.NewEmpty[int](),
		opt.New(1000),
		opt.NewPtr[int](nil),
		opt.NewPtr(&i),
	}

	for _, v := range values {
		fmt.Println(v.Or(1))
	}

	// Output:
	// 1
	// 1000
	// 1
	// 1001
}

func Example_elseZero() {
	i := 1001
	values := []opt.Optional[int]{
		opt.NewEmpty[int](),
		opt.New(1000),
		opt.NewPtr[int](nil),
		opt.NewPtr(&i),
	}

	for _, v := range values {
		fmt.Println(v.OrZero())
	}

	// Output:
	// 0
	// 1000
	// 0
	// 1001
}

func Example_elseFunc() {
	i := 1001
	values := []opt.Optional[int]{
		opt.NewEmpty[int](),
		opt.New(1000),
		opt.NewPtr[int](nil),
		opt.NewPtr(&i),
	}

	for _, v := range values {
		fmt.Println(v.OrCall(func() int {
			return 2
		}))
	}

	// Output:
	// 2
	// 1000
	// 2
	// 1001
}

func Example_jsonMarshalOmitEmpty() {
	s := struct {
		Bool    opt.Optional[bool]      `json:"bool,omitempty"`
		Byte    opt.Optional[byte]      `json:"byte,omitempty"`
		Float32 opt.Optional[float32]   `json:"float32,omitempty"`
		Float64 opt.Optional[float64]   `json:"float64,omitempty"`
		Int16   opt.Optional[int16]     `json:"int16,omitempty"`
		Int32   opt.Optional[int32]     `json:"int32,omitempty"`
		Int64   opt.Optional[int64]     `json:"int64,omitempty"`
		Int     opt.Optional[int]       `json:"int,omitempty"`
		Rune    opt.Optional[rune]      `json:"rune,omitempty"`
		String  opt.Optional[string]    `json:"string,omitempty"`
		Time    opt.Optional[time.Time] `json:"time,omitempty"`
		Uint16  opt.Optional[uint16]    `json:"uint16,omitempty"`
		Uint32  opt.Optional[uint32]    `json:"uint32,omitempty"`
		Uint64  opt.Optional[uint64]    `json:"uint64,omitempty"`
		Uint    opt.Optional[uint]      `json:"uint,omitempty"`
		Uintptr opt.Optional[uintptr]   `json:"uintptr,omitempty"`
	}{
		Bool:    opt.NewEmpty[bool](),
		Byte:    opt.NewEmpty[byte](),
		Float32: opt.NewEmpty[float32](),
		Float64: opt.NewEmpty[float64](),
		Int16:   opt.NewEmpty[int16](),
		Int32:   opt.NewEmpty[int32](),
		Int64:   opt.NewEmpty[int64](),
		Int:     opt.NewEmpty[int](),
		Rune:    opt.NewEmpty[rune](),
		String:  opt.NewEmpty[string](),
		Time:    opt.NewEmpty[time.Time](),
		Uint16:  opt.NewEmpty[uint16](),
		Uint32:  opt.NewEmpty[uint32](),
		Uint64:  opt.NewEmpty[uint64](),
		Uint:    opt.NewEmpty[uint](),
		Uintptr: opt.NewEmpty[uintptr](),
	}

	output, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(output))

	// Output:
	// {}
}

func Example_jsonMarshalEmpty() {
	s := struct {
		Bool    opt.Optional[bool]      `json:"bool"`
		Byte    opt.Optional[byte]      `json:"byte"`
		Float32 opt.Optional[float32]   `json:"float32"`
		Float64 opt.Optional[float64]   `json:"float64"`
		Int16   opt.Optional[int16]     `json:"int16"`
		Int32   opt.Optional[int32]     `json:"int32"`
		Int64   opt.Optional[int64]     `json:"int64"`
		Int     opt.Optional[int]       `json:"int"`
		Rune    opt.Optional[rune]      `json:"rune"`
		String  opt.Optional[string]    `json:"string"`
		Time    opt.Optional[time.Time] `json:"time"`
		Uint16  opt.Optional[uint16]    `json:"uint16"`
		Uint32  opt.Optional[uint32]    `json:"uint32"`
		Uint64  opt.Optional[uint64]    `json:"uint64"`
		Uint    opt.Optional[uint]      `json:"uint"`
		Uintptr opt.Optional[uintptr]   `json:"uintptr"`
	}{
		Bool:    opt.NewEmpty[bool](),
		Byte:    opt.NewEmpty[byte](),
		Float32: opt.NewEmpty[float32](),
		Float64: opt.NewEmpty[float64](),
		Int16:   opt.NewEmpty[int16](),
		Int32:   opt.NewEmpty[int32](),
		Int64:   opt.NewEmpty[int64](),
		Int:     opt.NewEmpty[int](),
		Rune:    opt.NewEmpty[rune](),
		String:  opt.NewEmpty[string](),
		Time:    opt.NewEmpty[time.Time](),
		Uint16:  opt.NewEmpty[uint16](),
		Uint32:  opt.NewEmpty[uint32](),
		Uint64:  opt.NewEmpty[uint64](),
		Uint:    opt.NewEmpty[uint](),
		Uintptr: opt.NewEmpty[uintptr](),
	}

	output, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(output))

	// Output:
	// {
	//   "bool": null,
	//   "byte": null,
	//   "float32": null,
	//   "float64": null,
	//   "int16": null,
	//   "int32": null,
	//   "int64": null,
	//   "int": null,
	//   "rune": null,
	//   "string": null,
	//   "time": null,
	//   "uint16": null,
	//   "uint32": null,
	//   "uint64": null,
	//   "uint": null,
	//   "uintptr": null
	// }
}

func Example_jsonMarshalPresent() {
	s := struct {
		Bool    opt.Optional[bool]      `json:"bool"`
		Byte    opt.Optional[byte]      `json:"byte"`
		Float32 opt.Optional[float32]   `json:"float32"`
		Float64 opt.Optional[float64]   `json:"float64"`
		Int16   opt.Optional[int16]     `json:"int16"`
		Int32   opt.Optional[int32]     `json:"int32"`
		Int64   opt.Optional[int64]     `json:"int64"`
		Int     opt.Optional[int]       `json:"int"`
		Rune    opt.Optional[rune]      `json:"rune"`
		String  opt.Optional[string]    `json:"string"`
		Time    opt.Optional[time.Time] `json:"time"`
		Uint16  opt.Optional[uint16]    `json:"uint16"`
		Uint32  opt.Optional[uint32]    `json:"uint32"`
		Uint64  opt.Optional[uint64]    `json:"uint64"`
		Uint    opt.Optional[uint]      `json:"uint"`
		Uintptr opt.Optional[uintptr]   `json:"uintptr"`
	}{
		Bool:    opt.New(true),
		Byte:    opt.New[byte](1),
		Float32: opt.New[float32](2.1),
		Float64: opt.New(2.2),
		Int16:   opt.New[int16](3),
		Int32:   opt.New[int32](4),
		Int64:   opt.New[int64](5),
		Int:     opt.New(6),
		Rune:    opt.New[rune](7),
		String:  opt.New("string"),
		Time:    opt.New(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
		Uint16:  opt.New[uint16](8),
		Uint32:  opt.New[uint32](9),
		Uint64:  opt.New[uint64](10),
		Uint:    opt.New[uint](11),
		Uintptr: opt.New[uintptr](12),
	}

	output, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(output))

	// Output:
	// {
	//   "bool": true,
	//   "byte": 1,
	//   "float32": 2.1,
	//   "float64": 2.2,
	//   "int16": 3,
	//   "int32": 4,
	//   "int64": 5,
	//   "int": 6,
	//   "rune": 7,
	//   "string": "string",
	//   "time": "2006-01-02T15:04:05Z",
	//   "uint16": 8,
	//   "uint32": 9,
	//   "uint64": 10,
	//   "uint": 11,
	//   "uintptr": 12
	// }
}

func Example_jsonUnmarshalEmpty() {
	s := struct {
		Bool    opt.Optional[bool]      `json:"bool"`
		Byte    opt.Optional[byte]      `json:"byte"`
		Float32 opt.Optional[float32]   `json:"float32"`
		Float64 opt.Optional[float64]   `json:"float64"`
		Int16   opt.Optional[int16]     `json:"int16"`
		Int32   opt.Optional[int32]     `json:"int32"`
		Int64   opt.Optional[int64]     `json:"int64"`
		Int     opt.Optional[int]       `json:"int"`
		Rune    opt.Optional[rune]      `json:"rune"`
		String  opt.Optional[string]    `json:"string"`
		Time    opt.Optional[time.Time] `json:"time"`
		Uint16  opt.Optional[uint16]    `json:"uint16"`
		Uint32  opt.Optional[uint32]    `json:"uint32"`
		Uint64  opt.Optional[uint64]    `json:"uint64"`
		Uint    opt.Optional[uint]      `json:"uint"`
		Uintptr opt.Optional[uintptr]   `json:"uintptr"`
	}{}

	x := `{}`
	err := json.Unmarshal([]byte(x), &s)
	fmt.Println("error:", err)
	fmt.Println("Bool:", s.Bool.Ok())
	fmt.Println("Byte:", s.Byte.Ok())
	fmt.Println("Float32:", s.Float32.Ok())
	fmt.Println("Float64:", s.Float64.Ok())
	fmt.Println("Int16:", s.Int16.Ok())
	fmt.Println("Int32:", s.Int32.Ok())
	fmt.Println("Int64:", s.Int64.Ok())
	fmt.Println("Int:", s.Int.Ok())
	fmt.Println("Rune:", s.Rune.Ok())
	fmt.Println("String:", s.String.Ok())
	fmt.Println("Time:", s.Time.Ok())
	fmt.Println("Uint16:", s.Uint16.Ok())
	fmt.Println("Uint32:", s.Uint32.Ok())
	fmt.Println("Uint64:", s.Uint64.Ok())
	fmt.Println("Uint64:", s.Uint64.Ok())
	fmt.Println("Uint:", s.Uint.Ok())
	fmt.Println("Uintptr:", s.Uint.Ok())

	// Output:
	// error: <nil>
	// Bool: false
	// Byte: false
	// Float32: false
	// Float64: false
	// Int16: false
	// Int32: false
	// Int64: false
	// Int: false
	// Rune: false
	// String: false
	// Time: false
	// Uint16: false
	// Uint32: false
	// Uint64: false
	// Uint64: false
	// Uint: false
	// Uintptr: false
}

func Example_jsonUnmarshalPresent() {
	s := struct {
		Bool    opt.Optional[bool]      `json:"bool"`
		Byte    opt.Optional[byte]      `json:"byte"`
		Float32 opt.Optional[float32]   `json:"float32"`
		Float64 opt.Optional[float64]   `json:"float64"`
		Int16   opt.Optional[int16]     `json:"int16"`
		Int32   opt.Optional[int32]     `json:"int32"`
		Int64   opt.Optional[int64]     `json:"int64"`
		Int     opt.Optional[int]       `json:"int"`
		Rune    opt.Optional[rune]      `json:"rune"`
		String  opt.Optional[string]    `json:"string"`
		Time    opt.Optional[time.Time] `json:"time"`
		Uint16  opt.Optional[uint16]    `json:"uint16"`
		Uint32  opt.Optional[uint32]    `json:"uint32"`
		Uint64  opt.Optional[uint64]    `json:"uint64"`
		Uint    opt.Optional[uint]      `json:"uint"`
		Uintptr opt.Optional[uintptr]   `json:"uintptr"`
	}{}

	x := `{
   "bool": false,
   "byte": 0,
   "float32": 0,
   "float64": 0,
   "int16": 0,
   "int32": 0,
   "int64": 0,
   "int": 0,
   "rune": 0,
   "string": "string",
   "time": "0001-01-01T00:00:00Z",
   "uint16": 0,
   "uint32": 0,
   "uint64": 0,
   "uint": 0,
   "uintptr": 0
 }`
	err := json.Unmarshal([]byte(x), &s)
	fmt.Println("error:", err)
	fmt.Println("Bool:", s.Bool.Ok(), s.Bool)
	fmt.Println("Byte:", s.Byte.Ok(), s.Byte)
	fmt.Println("Float32:", s.Float32.Ok(), s.Float32)
	fmt.Println("Float64:", s.Float64.Ok(), s.Float64)
	fmt.Println("Int16:", s.Int16.Ok(), s.Int16)
	fmt.Println("Int32:", s.Int32.Ok(), s.Int32)
	fmt.Println("Int64:", s.Int64.Ok(), s.Int64)
	fmt.Println("Int:", s.Int.Ok(), s.Int)
	fmt.Println("Rune:", s.Rune.Ok(), s.Rune)
	fmt.Println("String:", s.String.Ok(), s.String)
	fmt.Println("Time:", s.Time.Ok(), s.Time)
	fmt.Println("Uint16:", s.Uint16.Ok(), s.Uint16)
	fmt.Println("Uint32:", s.Uint32.Ok(), s.Uint32)
	fmt.Println("Uint64:", s.Uint64.Ok(), s.Uint64)
	fmt.Println("Uint64:", s.Uint64.Ok(), s.Uint64)
	fmt.Println("Uint:", s.Uint.Ok(), s.Uint)
	fmt.Println("Uintptr:", s.Uint.Ok(), s.Uint)

	// Output:
	// error: <nil>
	// Bool: true false
	// Byte: true 0
	// Float32: true 0
	// Float64: true 0
	// Int16: true 0
	// Int32: true 0
	// Int64: true 0
	// Int: true 0
	// Rune: true 0
	// String: true string
	// Time: true 0001-01-01 00:00:00 +0000 UTC
	// Uint16: true 0
	// Uint32: true 0
	// Uint64: true 0
	// Uint64: true 0
	// Uint: true 0
	// Uintptr: true 0
}
