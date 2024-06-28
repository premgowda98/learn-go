package user

import (
	"errors"
	"fmt"
	"time"
)

// Struct
type User struct {
	firstName string
	lastName  string
	birthDate string
	createdAt time.Time
}

// regular function
func PrintName(u User) {
	fmt.Println(u.firstName, u.lastName, u.birthDate)
}

// method => function attached to struct. () before function name is called receiver argument
func (u User) PrintNameMethod() {
	fmt.Println(u.firstName, u.lastName, u.birthDate)
}

// mutation methods
// In mutation function we should pointer as recivere argument, since we are
// changing the struct data
// If not sent as pointer, then only the copy will be edited and not the original object
func (u *User) MutateUser() {
	u.firstName = "Hi " + u.firstName
}

// constructor functions
func New(userfirstName, userlastName, userbirthdate string) (*User, error) {
	if userfirstName == "" || userfirstName == "" || userbirthdate == "" {
		return nil, errors.New("Invalid")
	}

	return &User{
		firstName: userfirstName,
		lastName:  userlastName,
		birthDate: userbirthdate,
		createdAt: time.Now(),
	}, nil
}

// Struct Embedding

type Admin struct {
	email    string
	password string
	User
}

func NewAdmin(email, password string) Admin {
	return Admin{
		email:    email,
		password: password,
		User: User{
			firstName: "Admin",
			lastName:  "Admin",
			birthDate: "-",
			createdAt: time.Now(),
		},
	}
}
