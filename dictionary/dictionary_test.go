package dictionary

import (
	"testing"
)

func TestAdd(t *testing.T) {
	dict := New()

	// Test adding a new word successfully
	err := dict.Add("test", "definition")
	if err != nil {
		t.Errorf("Error adding a new word: %v", err)
	}

	// Test adding a word that already exists (should return ErrWordAlreadyExists)
	err = dict.Add("test", "definition")
	if err != ErrWordAlreadyExists {
		t.Errorf("Expected ErrWordAlreadyExists, got %v", err)
	}
}

func TestGet(t *testing.T) {
	dict := New()

	// Test retrieving a definition for an existing word
	dict.Add("test", "definition")
	entry, err := dict.Get("test")
	if err != nil {
		t.Errorf("Error retrieving definition: %v", err)
	}
	if entry.Definition != "definition" {
		t.Errorf("Expected definition 'definition', got '%s'", entry.Definition)
	}

	// Test retrieving a definition for a non-existing word (should return ErrWordNotFound)
	_, err = dict.Get("nonexistent")
	if err != ErrWordNotFound {
		t.Errorf("Expected ErrWordNotFound, got %v", err)
	}
}

func TestRemove(t *testing.T) {
	dict := New()

	// Test removing an existing word successfully
	dict.Add("test", "definition")
	err := dict.Remove("test")
	if err != nil {
		t.Errorf("Error removing an existing word: %v", err)
	}

	// Test removing a non-existing word (should not result in an error)
	err = dict.Remove("nonexistent")
	if err != nil {
		t.Errorf("Error removing a non-existing word: %v", err)
	}
}

func TestList(t *testing.T) {
	dict := New()

	// Test listing all entries in the dictionary
	dict.Add("test1", "definition1")
	dict.Add("test2", "definition2")

	entries := dict.List()
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}
}
