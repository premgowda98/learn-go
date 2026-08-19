package main // simmilar to __init__.py.
// main is a special package indicating go that this is the main package to be ran
// this main works more like a __name__=='__main__' in python
// anothing analogy is docker CMD and Entrypoint commands

import "fmt" //standard package

func main() { // even this function must be named main
	fmt.Print("Hello World") //use only "", ''will not work or try ``
}

// there should be only one main function in single package
