package main

import (
	"bufio"
	"dictionaryApp/dictionary"
	"fmt"
	"os"
	"strings"
)

func main() {
	const filename = "db_dictionary.json"
	dict := dictionary.New()
	reader := bufio.NewReader(os.Stdin)
	dict.LoadFromFile()

	for {
		fmt.Printf("\nEnter action ->  add | define | remove | list | exit  : ")
		action, _ := reader.ReadString('\n')
		action = strings.TrimSpace(action)

		switch action {
		case "add":
			actionAdd(dict, reader)
		case "define":
			actionDefine(dict, reader)
		case "remove":
			actionRemove(dict, reader)
		case "list":
			actionList(dict)
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("!!! Invalid action !!!")
		}
	}
}

func actionAdd(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Printf("\nEnter word: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	fmt.Printf("\nEnter definition: ")
	definition, _ := reader.ReadString('\n')
	definition = strings.TrimSpace(definition)

	d.Add(word, definition)
	d.SaveToFile()
	fmt.Printf("%s added.\n", word)
}

func actionDefine(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter word to define: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	entry, err := d.Get(word)
	if err != nil {
		fmt.Println("Not found.")
		return
	}

	fmt.Printf("- %s: %s\n", word, entry.Definition)
}

func actionRemove(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter word to remove: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	d.Remove(word)
	d.SaveToFile()
	fmt.Printf("%s removed.\n", word)
}

func actionList(d *dictionary.Dictionary) {
	entries := d.List()
	if len(entries) == 0 {
		fmt.Println("Dictionary is empty.")
		return
	}

	fmt.Println("\n--------------- list ---------------")
	for _, entry := range entries {
		fmt.Printf("- %s:\n \t%s \n\n", entry.Word, entry.Definition)
	}
}

/*
laptop
compact and portable personal computer.

telepheric
term typically used in the context of cable cars and transportation systems.

web
system of interconnected public webpages accessible through the Internet.

*/
