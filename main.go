package main

import (
	"ssego/Revenge"
)

func main() {
	var RV *Revenge.RevengeRoot = Revenge.New("8080")
	RV.Commence()
}
