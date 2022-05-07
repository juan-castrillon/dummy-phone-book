package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Entry struct {
	Name        string
	Surname     string
	PhoneNumber string
}

var data []Entry

func search(key string) *Entry {
	for _, entry := range data {
		if entry.Surname == key {
			return &entry
		}
	}
	return nil
}

func list() {
	for _, entry := range data {
		fmt.Println(entry)
	}
}

func main() {
	args := os.Args
	nArgs := len(args)
	if nArgs == 1 {
		exe := filepath.Base(args[0])
		fmt.Printf("Usage: %s search|list <arguments>\n", exe)
		os.Exit(1)
	}
	data = append(data, Entry{"Juan", "Castrillon", "01705257416"})
	data = append(data, Entry{"Benito", "Camelas", "02589748862"})
	data = append(data, Entry{"Socio", "Elbi", "05874962158"})
	switch args[1] {
	case "search":
		if nArgs != 3 {
			fmt.Println("Usage: search Surname")
			os.Exit(1)
		}
		result := search(args[2])
		if result == nil {
			fmt.Println("Entry not found: ", args[2])
			os.Exit(0)
		}
		fmt.Println(*result)
	case "list":
		list()
	default:
		fmt.Println("Invalid option ", args[1])
		os.Exit(1)
	}
}
