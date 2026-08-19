package ui

import "github.com/premgowda/learn-go/wailsAppGo/pkg"

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return pkg.SomeFunc() + " " + name
}
