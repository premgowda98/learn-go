## Basics

1. In Go, Module contains multiple packages, unlike in python where multiple modules form a package
2. In order to create module, use `go mod init <module_name>` this creates **go.mod** file in the directory
3. `go build` will create the compiled code with **module_name**
4. The main file must contain `package main` just like __init__==__main__
5. Even the entry function should be named **main**