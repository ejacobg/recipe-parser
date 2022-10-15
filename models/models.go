package models

import (
	"encoding/json"
	"os"
)

type Ingredient struct {
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
	Name   string `json:"name"`
	Notes  string `json:"notes"`
}

type Recipe struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	URL          string       `json:"url"`
	Image        string       `json:"image"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions []string     `json:"instructions"`
}

func (r *Recipe) SaveAs(path string) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path+".json", data, 00666)
	if err != nil {
		return err
	}
	return nil
}

func (r *Recipe) ToJSON() (string, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
