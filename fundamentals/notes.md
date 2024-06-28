## Data Types

1. int, float64, string, bool
2. unit => unsigned non negative integer
3. int32, rune, int64, unit32, int8, unit8

## Functions

1. 

## Error Handling
1. Error handling is bit different compared to other lang
2. If any error occurs, then the program will not crash, instead go provides option to specify default values in case of error occurs,
3. Example: If file does'not exist to read, instead of crashing the application, go provides 0 as ouptut.
4. Python **raise** => Go **panic**
5. **errors** package is avaibale to create errors.
Some functions by default returns two values one for output and other for error, if no error then value will be **nil**