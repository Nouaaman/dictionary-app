package main

import (
	"bufio"
	"dictionaryApp/dictionary"
	"fmt"
	"os"
	"strings"
)

func main() {
	dict := dictionary.New()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("\nEnter action ->  add | define | remove | list | exit  : ")
		action, _ := reader.ReadString('\n')

		switch action {
		case "add\n":
			actionAdd(dict, reader)
		case "define\n":
			actionDefine(dict, reader)
		case "remove\n":
			actionRemove(dict, reader)
		case "list\n":
			actionList(dict)
		case "exit\n":
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

	fmt.Printf("Definition: %s\n", entry.Definition)
}

func actionRemove(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter word to remove: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	d.Remove(word)
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
		fmt.Printf("- %s: %s \n\n", entry.Word, entry.Definition)
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
