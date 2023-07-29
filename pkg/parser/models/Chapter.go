package models

type Chapter struct {
	Title       string    `json:"title"`
	Raw         string    `json:"raw"`
	Recording   string    `json:"recording"`
	Translation string    `json:"translation"`
	Navigation  string    `json:"navigation"`
	Children    []Chapter `json:"children"`
}
