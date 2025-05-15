package recipes

// Represents individual ingredients
type Ingredient struct {
	Name string `json:"name"`
}

// Represents a recipe
type Recipe struct {
	Name        string       `json:"name"`
	Ingredients []Ingredient `json:"ingredients"`
}
