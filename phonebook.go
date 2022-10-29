package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Entry struct {
	Name       string
	Surname    string
	Tel        string
	LastAccess string
}

// CSV Reading and Writing
//------------------------
const CSVFILE = "data.csv"

var data = []Entry{}
var index map[string]int
var telRegexp = regexp.MustCompile(`^\d+$`)

func readOrCreateCSV(filepath string) error {
	fI, err := os.Stat(filepath)
	exists := true
	if err != nil {
		if os.IsNotExist(err) {
			exists = false
		} else {
			return err
		}
	}
	if exists {
		// Check if regular file
		if !fI.Mode().IsRegular() {
			return fmt.Errorf("%s is not a regular file", filepath)
		}
		return readCSV(filepath)
	} else {
		f, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer f.Close()
		return nil
	}
}

func readCSV(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}
	index = make(map[string]int)
	for i, line := range lines {
		entry := Entry{
			Name:       line[0],
			Surname:    line[1],
			Tel:        line[2],
			LastAccess: line[3],
		}
		data = append(data, entry)
		index[entry.Tel] = i
	}
	return nil
}

func saveCSV(filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	csvwriter := csv.NewWriter(f)
	for _, row := range data {
		tmp := []string{row.Name, row.Surname, row.Tel, row.LastAccess}
		_ = csvwriter.Write(tmp)
	}
	csvwriter.Flush()
	return csvwriter.Error()
}

//------------------------

func search(key string) *Entry {
	i, ok := index[key]
	if !ok {
		return nil
	}
	data[i].LastAccess = time.Now().Format(time.RFC3339)
	return &data[i]
}

func list() {
	for _, entry := range data {
		fmt.Printf("%s %s: %s\n", entry.Name, entry.Surname, entry.Tel)
	}
}

func createEntry(name, surname, tel string) (*Entry, error) {
	if tel == "" || surname == "" {
		return nil, fmt.Errorf("surname and tel are mandatory")
	}
	return &Entry{
		Name:       name,
		Surname:    surname,
		Tel:        tel,
		LastAccess: time.Now().Format(time.RFC3339),
	}, nil
}

func isTel(s string) bool {
	return telRegexp.MatchString(s)
}

func insert(e *Entry) error {
	_, ok := index[e.Tel]
	if ok {
		return fmt.Errorf("entry with tel %s already exists", e.Tel)
	}
	data = append(data, *e)
	index[e.Tel] = len(data) - 1
	return saveCSV(CSVFILE)
}

func deleteEntry(key string) error {
	i, ok := index[key]
	if !ok {
		return fmt.Errorf("no entry with key %s", key)
	}
	data = append(data[:i], data[i+1:]...)
	delete(index, key)
	return saveCSV(CSVFILE)
}

func main() {
	args := os.Args
	nArgs := len(args)
	if nArgs == 1 {
		exe := filepath.Base(args[0])
		fmt.Printf("Usage: %s insert|delete|search|list <arguments>\n", exe)
		os.Exit(1)
	}
	err := readOrCreateCSV(CSVFILE)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	switch args[1] {
	case "insert":
		if nArgs != 5 {
			fmt.Println("Usage: insert Name Surname Tel")
			os.Exit(1)
		}
		t := strings.ReplaceAll(args[4], "-", "")
		if !isTel(t) {
			fmt.Println("Tel must be a number")
			os.Exit(1)
		}
		entry, err := createEntry(args[2], args[3], t)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = insert(entry)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "delete":
		if nArgs != 3 {
			fmt.Println("Usage: delete Tel")
			os.Exit(1)
		}
		deleteEntry(args[2])
	case "search":
		if nArgs != 3 {
			fmt.Println("Usage: search Tel")
			os.Exit(1)
		}
		result := search(args[2])
		if result == nil {
			fmt.Println("Entry not found: ", args[2])
			os.Exit(0)
		}
		fmt.Printf("%s %s: %s\n", result.Name, result.Surname, result.Tel)
	case "list":
		list()
	default:
		fmt.Println("Invalid option ", args[1])
		os.Exit(1)
	}
}
