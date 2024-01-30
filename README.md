# opt

> Optional types and utilities for egonomic data transformation.

[![GoDoc](https://pkg.go.dev/badge/github.com/Southclaws/opt)](https://pkg.go.dev/github.com/Southclaws/opt?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/Southclaws/opt)](https://goreportcard.com/report/github.com/Southclaws/opt)

opt provides a simple generic optional type with a variety of utilities for
performing various transformations without the need for explicit branching.

So, while there are many `Optional[T]` style packages out there, this one has a
focus on making _data transformations_ easier to write and easier to read.

It also prevents certain categories of bug such as nil pointer dereferencing. Of
course this comes at a cost and if you're writing performance sensitive code,
this library may not be for you and you may be better off just being explicit.

The status of this library is pre-1.0 but the API is stable and probably won't
change. It has been dogfooded in 3 production codebases for about a year and all
APIs were built to solve some real problems in those projects.

## Basics

Let's get the obvious out of the way first...

```go
func main() {
    maybe := opt.New("I exist!")

    maybe.Ok() // true
    value, exists := maybe.Get() // "I exist!", true
    ptr := maybe.Ptr() // some address

    maybe_not := opt.NewEmpty[string]()

    maybe_not.Ok() // false
    value, exists := maybe_not.Get() // "", false
    ptr := maybe_not.Ptr() // nil
}
```

Optional, generic, yada yada, whatever. Every other optional package does it.

The interesting parts are the construction, mapping and access utilities...

## Accessing

Once you've constructed an optional value, you can access the underlying data in
a few ways. These make it easy to build branching logic without the need for
explicit if statements. Which can be useful for transforming large structures.

The simplest ones are `Ok` and `Get` which have examples above. See the GoDoc
for more info on these, they're fairly simple and do what you'd expect.

One method that isn't mentioned above is `Call`. Which simply lets you call a
function with the value if it's present:

```go
maybe := opt.New("I exist!")
maybe.Call(func(value string) {
    fmt.Println(value)
})
```

These have been handy for some ORM setter APIs:

```go
email.Call(accountQuery.SetEmailAddress)
```

## Mapping

One of the core reasons this library was written was to facilitate easy mapping
of data types that may or may not be present. Without the need for code that
looks like this:

```go
var newValue *T
if oldValue != nil {
    newValue = transform(*oldValue)
}
```

Which is fine on its own, but if you have many values, it can get quite verbose.

opt instead provides a way to map data as an access or map data as a pipeline.
To access the data, you already know about `Get` but if you want to change the
type at the same time as accessing, you can use `GetMap`:

```go
maybe := opt.New("I exist!")
value, exists := opt.GetMap(maybe, strings.ToUpper)
// "I EXIST!", true
```

If your destination is expecting a pointer, you can use `PtrMap`:

```go
maybe := opt.New("I exist!")
value := opt.PtrMap(maybe, strings.ToUpper)
// "I EXIST!" as a `*string`
```

Note how these are functions of the library, not methods on the type. It would
be nice to be able to write `maybe.PtrMap(strings.ToUpper)` but currently, this
is not possible to do in the current version of Go's generics.

If you want to transform the data but keep it wrapped as an optional type, you
can use `Map` or `MapErr` to execute the closure, only if the value is present:

```go
maybe := opt.New("I exist!")
maybe = opt.Map(maybe, strings.ToUpper)
// opt.Optional[string]("I EXIST!")
```

And of course `MapErr` does the same thing but allows you to return an error:

```go
maybe := opt.New("5629")
maybe, err := opt.MapErr(maybe, strconv.Atoi)
// opt.Optional[int](5629)

maybe_not := opt.New("not a number :(")
maybe_not, err := opt.MapErr(maybe, strconv.Atoi)
// Empty optional plus the error from Atoi.
```

And, as an escape hatch, a `.String()` method which is useful for tests:

```go
maybe := opt.New("5629")
maybe.String()
// "5629"
```

If the value exists, it'll use `fmt` to stringify, if not it'll just be empty.

### When there's nothing inside

There are also methods to deal with empty values `Or`, `OrZero` and `OrCall`:

```go
maybe := opt.New("I exist!")
maybe.Or("I don't exist!") // "I exist!"
```

```go
maybe := opt.Empty[string]()
maybe.Or("I don't exist!") // "I don't exist!"
```

The `Or` method simply lets you return a default value if the optional value is
empty. This is handy for providing defaults.

The `OrZero` method simply returns the type's zero-value:

```go
maybe := opt.Empty[time.Time]()
t := maybe.OrZero()
t.IsZero() // true
```

And finally, `OrCall` lets you call a function to provide a default value:

```go
maybe := opt.Empty[string]()
t := maybe.OrCall(func() string {
    return "a default value from somewhere"
})
// "a default value from somewhere"
```

## Curried `C` Functions

Some APIs will have a second version with `C` appended to the name. These are
curried versions of those functions to aid in ergonomic usage.

Say for example you have a function that converts a number to a GBP currency
representation. You want to apply this function to a few values in a struct or
to a slice of items.

```go
// Given: ConvertUSD(value int) string

func Convert(input Table) PriceBreakdown {
    return PriceBreakdown{
        Cost:          ConvertGBP(input.UnitCost),
        ShippingFee:   NewPtrMap(input.ShippingFee, ConvertGBP),
        ServiceCharge: NewPtrMap(input.ServiceCharge, ConvertGBP),
        Discount:      NewPtrMap(input.Discount, ConvertGBP),
    }
}
```

A small example, but you could imagine how much this can get in a larger system.

Using curried APIs, we can make this a little more terse:

```go
func Convert(input Table) PriceBreakdown {
    gbp := NewPtrMapC(ConvertGBP)
    return PriceBreakdown{
        Cost:          ConvertGBP(input.UnitCost),
        ShippingFee:   gbp(input.ShippingFee),
        ServiceCharge: gbp(input.ServiceCharge),
        Discount:      gbp(input.Discount),
    }
}
```

Now this may not seem like much but it can make refactors easier and keep diffs
small. Once you start thinking in curried functions, certain tasks get simpler!

Let's see what this looks like for a slice of items:

```go
func ConvertMany(prices []*int) []Optional[string] {
    output := []Optional[string]{}
    for _, v := range prices {
        output = append(output, NewPtrMap(v, ConvertGBP))
    }
    return output
}
```

If you like to use functional libraries like [lo](https://github.com/samber/lo)
and [fp-go](https://github.com/repeale/fp-go) then this might be useful:

```go
func ConvertMany(prices []*int) []Optional[string] {
    fn := PtrMapC(ConvertGBP)
    mapper := fp.Map(fn)
    return mapper(prices)
}
```

## Construction

There are quite a few places data can come from. opt provides a few helpers to
create optional wrappers from various sources.

We've covered the boring ones already, `New` and `NewEmpty` just create values
from either something or nothing.

### `NewMap`

This tool creates an optional type but facilitates mapping the data type using a
function first. This is similar to `.map( x => y )` in many other languages.

```go
v := opt.NewMap("hello", strings.ToUpper)
```

`v` now contains an optional `string` value set to `"HELLO"`. because, before
storing the data, it passed the input value through `strings.ToUpper`.

### `NewSafe`

A common Go pattern is return values that look like `(T, bool)` where the bool
represents validity. `NewSafe` lets you easily build optional values from this.

```go
// where getThing is: func getThing() (v string, ok bool)
v := opt.NewSafe(getThing())
```

It's also just handy sometimes for simple logic:

```go
v := opt.NewSafe(account.Email, account.IsEmailPublic)
```

Here, we're storing the optional value of the account's email only if the value
of `IsEmailPublic` is true.

Sadly this does not work with built-in operations:

```go
hash := map[string]string{"s": "asd"}
NewSafe(hash["dsf"])
// not enough arguments in call to NewSafe have (string) want (T, bool)

var cast any = "hi"
NewSafe(cast.(string))
// not enough arguments in call to NewSafe have (string) want (T, bool)
```

This is because the bool part of these expressions is optional.

### `NewIf`

This one is another way to encode optionality based on some branching logic. In
this variant, the logic exists within a closure that returns a bool.

```go
v := opt.NewIf(account.Email, isValidEmailAddress)
v := opt.NewIf(company.LegalName, func(s string) bool { return s != "" })
v := opt.NewIf(createdAt, func(t time.Time) bool { return !t.IsZero() })
```

### `NewPtr`, `NewPtrMap`, `NewPtrIf` and `NewPtrOr`

If one area of your application is using pointers already but you want to expose
optionals, you can use this one to easily construct an optional from a pointer.

```go
type Account struct {
    Twitter *string
}

// ...

v := opt.NewPtr(account.Twitter)
v := opt.NewPtrOr(account.Twitter, "@southclaws")
```

## Prior Art

- https://github.com/leighmcculloch/go-optional
- https://github.com/samber/mo
- https://github.com/phelmkamp/valor

## Contributing

Issues and pull requests welcome!
