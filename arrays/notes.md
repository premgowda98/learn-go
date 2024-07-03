## Lists/Arrays

1. **cap** vs **len** in go
    1. cap gets the capacity defined when initiating te array
    2. len gives the current length of the array
2. To create an array of items with different types in Go, you can use the interface{} type.

```go
func main() {
    //Array containing 5 items of different type
    var randomsArray = [5]interface{}{"Hello world!", 35, false, 33.33, 'A'}
    fmt.Println(randomsArray) // => [Hello world! 35 false 33.33 65]

}
```

[Slices vs Array](https://dev.to/dawkaka/go-arrays-and-slices-a-deep-dive-dp8)
[Imp](https://youtu.be/2oXBOaVowYY)

## Maps
1. Similar to dict in python `map [string][string]{}`

### Maps vs Struct
Both can be used as key value pairs, but
1. In **map**, a key can be string, int
2. **Struct** is data type and it is static