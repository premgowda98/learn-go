## Packages

1. File 1 can have `package main` and some functions
2. File 2 can also have `package main` and this file can use functions in File 1 since both are main package.
3. And in this case, no need to import from file 1
4. If creating different package then must be in different directory and that folder name should be same as package name
5. When importing from other package, we must import from module name, which was initiaized
6. And also only those functions which are starts with upper case is imported

### 3rd Party Packages

1. Packages can be searched in go website
2. To import it `go get github.com/Pallinder/go-randomdata` usually this will be github repo