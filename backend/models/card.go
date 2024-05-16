package models

// Card represents a tarot card.

type Card struct {
	CardName       string `json:"name"`
	Type           string `json:"type"`
	MeaningUp      string `json:"meaning_up"`
	MeaningReverse string `json:"meaning_rev"`
	Description    string `json:"desc"`
	ShortName      string `json:"name_short"`
	Reversed       bool   `json:"reversed"`
}
