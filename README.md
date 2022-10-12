# opt

Optional types and utilities.

_Inspired by [go-optional by Leigh McCulloch](https://github.com/leighmcculloch/go-optional)_

## Usage

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

    // Control flow

    if value, exists := maybe.Get(); exists {
        // it's there!
    } else {
        // it's not there...
    }

    // Mapping and transformation

    maybe = opt.Map(maybe, strings.ToUpper) // I EXIST!
    maybe_not = opt.Map(maybe_not, strings.ToUpper) // no-op

    // Conditional execution

    maybe.OrCall(func() { fmt.Println("It exists!") }) // It exists!
    maybe_not.OrCall(func() { fmt.Println("It exists?") }) // silence...
}
```
