## Interface

1. In go, interface are not attached to any objects
2. They are independent, they only validate if the passed object has speficied method assosiated with it

1. How to accept any type of data in go
```go
func allData(data interface{}){
    passedVal, ok := data.(int)
    //data.() return the passed value and also bool whether it is the specified data type
}
```

## Generics

Kind of `any` but we can restrict the data types
s
```go
func onlyNumber[T int|float64] (a, b T) T{
    return a+b
}
```
