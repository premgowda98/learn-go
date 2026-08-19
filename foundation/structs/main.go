package main

import (
	"fmt"
	"learn/structs/user"
)

func main() {
	userfirstName := getUserData("Please enter your first name: ")
	userlastName := getUserData("Please enter your last name: ")
	userbirthdate := getUserData("Please enter your birthdate (MM/DD/YYYY): ")

	// var appUser user

	// appUser := user{
	// 	firstName: userfirstName,
	// 	lastName:  userlastName,
	// 	birthDate: userbirthdate,
	// 	createdAt: time.Now(),
	// }

	appUser, _ := user.New(userfirstName, userlastName, userbirthdate)

	// printName(appUser)
	appUser.PrintNameMethod()
	appUser.MutateUser()
	appUser.PrintNameMethod()

	appUserAdmin := user.NewAdmin("Prem@gmail.com", "123")
	appUserAdmin.PrintNameMethod()

}

func getUserData(promptText string) string {
	fmt.Print(promptText)
	var value string
	fmt.Scan(&value)
	return value
}
