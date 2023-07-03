package main

import (
	"github.com/sqweek/dialog"
)

func main() {
	dialog.Message("%s", "Hello World!").Title("Info").Info()
}
