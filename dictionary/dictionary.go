package dictionary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

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

// get data from the file
func (d *Dictionary) LoadFromFile() {
	// check if the file exixst or create new file
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()
		return
	} else if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Read data from the file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error loading file:", err)
		return
	}
	//return if no data in file
	if len(data) == 0 {
		return
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		fmt.Println("Error in data:", err)
		return
	}

	for _, entry := range entries {
		d.Add(entry.Word, entry.Definition)
	}

}

// save data to file
func (d *Dictionary) SaveToFile() error {
	entries := make([]Entry, 0, len(d.entries))
	for _, entry := range d.entries {
		entries = append(entries, entry)
	}

	data, err := json.MarshalIndent(entries, "", "	")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func (d *Dictionary) Add(word string, definition string) {
	entry := Entry{Word: word, Definition: definition}
	d.entries[word] = entry
}

func (d *Dictionary) Get(word string) (Entry, error) {
	entry, found := d.entries[word]
	if !found {
		return Entry{}, errors.New("word not found")
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string) {
	delete(d.entries, word)
}

func (d *Dictionary) List() []Entry {
	var entryList []Entry
	for _, entry := range d.entries {
		entryList = append(entryList, entry)
	}
	return entryList
}
