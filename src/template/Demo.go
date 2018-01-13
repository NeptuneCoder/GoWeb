package main

import (
	"text/template"
	"os"
)

func main() {
	t := template.New("fieldname example")
	t, _ = t.Parse("hello {{.UserName}}! {{.Email}} is you email?")
	p := Person{UserName: "yanghai",Email:"1454025171@qq.com"}
	t.Execute(os.Stdout, p)
}
type Person struct {
	UserName string
	Email string
}
