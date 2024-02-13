package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// apiRequest sends a GET request to the specified URL and returns the response body.
// It includes a delay before sending the request to avoid rate limiting and sets
// common HTTP headers for the request. (headers are set to avoid Cloudflare's rate limiting & blocks)
func apiRequest(url string) ([]byte, error) {
	// Delay to mitigate rate limiting issues
	time.Sleep(500 * time.Millisecond)

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set common headers for the request
	req.Header.Set("authority", "neal.fun")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("dnt", "1")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("referer", "https://neal.fun/infinite-craft/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="121", "Not A(Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	
	// Initialize and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Read and return the response body
	return ioutil.ReadAll(resp.Body)
}

// Recipe represents a combination of two items to produce a new item.
type Recipe struct {
	First  string `json:"first"`  // ID of the first item
	Second string `json:"second"` // ID of the second item
	Result string `json:"result"` // ID of the resulting item

	HasResult bool `json:"-"` // Indicates whether the recipe has been processed

	// Data is a pointer to the CraftingData containing this recipe.
	Data *CraftingData `json:"-"`

	// Pointers to the items involved in the recipe.
	FirstItem  *Item `json:"-"`
	SecondItem *Item `json:"-"`
	ResultItem *Item `json:"-"`
}

// equals checks if two recipes are equal, ignoring the order of items.
func (r *Recipe) equals(other *Recipe) bool {
	return (r.First == other.First && r.Second == other.Second) || (r.First == other.Second && r.Second == other.First)
}

// equalsAny checks if the recipe matches any recipe in a slice.
func (r *Recipe) equalsAny(recipes []Recipe) bool {
	for _, recipe := range recipes {
		if r.equals(&recipe) {
			return true
		}
	}
	return false
}

// craft performs the crafting operation by sending a request to an external API
// and updates the recipe with the result.
func (r *Recipe) craft() (string, error) {
	url := fmt.Sprintf("https://neal.fun/api/infinite-craft/pair?first=%s&second=%s", r.FirstItem.Name, r.SecondItem.Name)
	body, err := apiRequest(url)
	if err != nil {
		return "", err
	}

	// Unmarshal the API response
	var result struct {
		Result string `json:"result"`
		Emoji  string `json:"emoji"`
		IsNew  bool   `json:"isNew"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// Process the result
	if result.Result != "" {
		newItem := r.Data.addItem(result.Result, result.Emoji, r.First+"+"+r.Second)

		// Update the recipe with the result
		r.Result = newItem.Name
		r.ResultItem = newItem
		r.HasResult = true

		log.Printf("Found new item: %s %s\n", newItem.Emoji, newItem.Name)
		return newItem.Name, nil
	}

	return "", nil
}

// Item represents a craftable or basic item in the crafting system.
type Item struct {
	Emoji   string   `json:"emoji"`   // Emoji representation of the item
	Name    string   `json:"name"`    // Name of the item
	Recipes []string `json:"recipes"` // IDs of recipes that produce this item
}

// CraftingData holds the entire dataset for the crafting simulation,
// including items, recipes, and identifiers for newly generated items and recipes.
type CraftingData struct {
	Items    map[string]*Item   `json:"items"`                    // Map of items by their names
	Recipes  map[string]*Recipe `json:"recipes,omitempty"` // Map of processed recipes
	ItemID   int                `json:"-"`                   // Counter for generating unique item IDs
	RecipeID int                `json:"-"`                 // Counter for generating unique recipe IDs
}

// addItem adds a new item to the CraftingData or updates an existing one with a new recipe.
func (c *CraftingData) addItem(name, emoji, recipe string) *Item {
	// Check if item already exists and update it
	if item, ok := c.Items[name]; ok {
		for _, r := range item.Recipes {
			if r == recipe {
				return item
			}
		}
		item.Recipes = append(item.Recipes, recipe)
		return item
	}

	// Create and add a new item
	item := &Item{
		Emoji:   emoji,
		Name:    name,
		Recipes: []string{recipe},
	}
	c.Items[name] = item
	c.ItemID++
	return item
}

// addRecipe creates and adds a new recipe to the CraftingData.
func (c *CraftingData) addRecipe(first, second string) *Recipe {
	recipe := &Recipe{
		First:     first,
		Second:    second,
		HasResult: false,
		Result:    "",
		Data:      c,
	}
	recipe.FirstItem = c.Items[first]
	recipe.SecondItem = c.Items[second]

	c.Recipes[first+"+"+second] = recipe
	c.RecipeID++
	return recipe
}

// generateValidRecipes generates all possible unique combinations of items as recipes.
func (c *CraftingData) generateValidRecipes() {
	for _, itemA := range c.Items {
		for _, itemB := range c.Items {
			if itemA.Name != itemB.Name {
				recipe := Recipe{First: itemA.Name, Second: itemB.Name}
				if !recipe.equalsAny(c.getAllRecipes()) {
					c.addRecipe(itemA.Name, itemB.Name)
				}
			}
		}
	}
}

// getAllRecipes returns a slice of all recipes in CraftingData.
func (c *CraftingData) getAllRecipes() []Recipe {
	var recipes []Recipe
	for _, recipe := range c.Recipes {
		recipes = append(recipes, *recipe)
	}
	return recipes
}

// processRecipes attempts to craft each unprocessed recipe until no new items are found.
func (c *CraftingData) processRecipes() {
	for {
		newItemsFound := false
		for _, recipe := range c.Recipes {
			if !recipe.HasResult {
				_, err := recipe.craft()
				if err == nil && recipe.HasResult {
					newItemsFound = true
				}
			}
		}
		if !newItemsFound {
			break
		}
		time.Sleep(500 * time.Millisecond) // Delay to avoid rate limiting
	}
}

// save serializes the CraftingData to a JSON file.
func (c *CraftingData) save() {
	jsonData, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("Error marshalling data: %v", err)
	}
	err = ioutil.WriteFile("crafting_data.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
}

// main initializes the crafting simulation with basic items and processes recipes in iterations.
func main() {
	data := CraftingData{
		Items:   make(map[string]*Item),
		Recipes: make(map[string]*Recipe),
	}

	// Initialize with basic items
	basicItems := []Item{
		{"🌎", "Earth", nil},
		{"💧", "Water", nil},
		{"🔥", "Fire", nil},
		{"🌬️", "Wind", nil},
	}

	for _, item := range basicItems {
		data.addItem(item.Name, item.Emoji, "")
	}

	// Generate and process recipes in iterations
	for i := 0; i < 3; i++ {
		data.generateValidRecipes()
		data.processRecipes()
		data.save()
		fmt.Printf("Iteration %d: Items: %d, Recipes: %d\n", i+1, len(data.Items), len(data.Recipes))
	}
}
