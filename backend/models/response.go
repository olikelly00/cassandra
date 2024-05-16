package models

// JSONCard represents a tarot card struct that has been converted into JSON so that it can be rendered in the UI.

type JSONCard struct {
	CardName       string `json:"name"`
	Type           string `json:"type"`
	MeaningUp      string `json:"meaning_up"`
	MeaningReverse string `json:"meaning_rev"`
	Description    string `json:"desc"`
	ImageName      string `json:"image_file_name"`
	Reversed       bool   `json:"reversed"`
}
