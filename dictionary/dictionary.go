package dictionary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var ErrWordAlreadyExists = errors.New("word already exists")
var ErrWordNotFound = errors.New("word not found")

type Entry struct {
	Word       string
	Definition string
}

type Dictionary struct {
	entries map[string]Entry
}

func New() *Dictionary {
	return &Dictionary{
		entries: make(map[string]Entry),
	}
}

const filename = "db_dictionary.json"

func (d *Dictionary) LoadFromFile() error {
	// Check if the file exists. If it doesn't exist, create a new file.
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return err
		}
		defer file.Close()
		return nil
	} else if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// Read data from the file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error loading file:", err)
		return err
	}
	// return if no data in file
	if len(data) == 0 {
		return nil
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		fmt.Println("Error in data:", err)
		return err
	}

	for _, entry := range entries {
		d.Add(entry.Word, entry.Definition)
	}

	return nil
}

func (d *Dictionary) SaveToFile() error {
	entries := make([]Entry, 0, len(d.entries))
	for _, entry := range d.entries {
		entries = append(entries, entry)
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func (d *Dictionary) Add(word string, definition string) error {
	_, found := d.entries[word]
	if found {
		return ErrWordAlreadyExists
	}

	entry := Entry{Word: word, Definition: definition}
	d.entries[word] = entry
	return nil
}

func (d *Dictionary) Get(word string) (*Entry, error) {
	entry, found := d.entries[word]
	if !found {
		return nil, ErrWordNotFound
	}
	return &entry, nil
}

func (d *Dictionary) Remove(word string) error {
	delete(d.entries, word)
	return nil
}

func (d *Dictionary) List() []Entry {
	var entryList []Entry
	for _, entry := range d.entries {
		entryList = append(entryList, entry)
	}
	return entryList
}
