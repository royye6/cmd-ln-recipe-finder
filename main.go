package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ApiResponse struct {
	Results []Recipe `json:"results"`
}

type Recipe struct {
	Title       string   `json:"title"`
	PrepTime    *int     `json:"readyInMinutes"`
	Servings    *int     `json:"servings"`
	SourceUrl   string   `json:"sourceUrl"`
	Image       string   `json:"image"`
	DishType    []string `json:"dishTypes"`
	DietType    []string `json:"diets"`
	Description string   `json:"summary"`
	RecipeScore float64  `json:"spoonacularScore"`
}

func cleanDescription(htmlText string) string {
	cleanText := strings.ReplaceAll(htmlText, "<[^>]*", "")
	cleanText = strings.TrimSpace(cleanText)

	return cleanText
}

func getData() (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	apiKey := os.Getenv("SPOONACULAR_API_KEY")

	url := fmt.Sprintf("https://api.spoonacular.com/recipes/complexSearch?apiKey=%s&type=breakfast&query=bacon&addRecipeInformation=True&number=1&sort=random", apiKey)

	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return "", fmt.Errorf("error reading response body: %v", err)
		}

		var apiResponse ApiResponse
		err = json.Unmarshal(body, &apiResponse)

		if err != nil {
			return "", fmt.Errorf("error unmarshalling(deserializing) JSON response: %v", err)
		}

		randRecipe := apiResponse.Results[0]

		cleanDescription := cleanDescription(randRecipe.Description)

		recipeDetails := fmt.Sprintf("\nTitle: %s\n\nPrep Time (minutes): %d\n\nServings: %d\n\nSource URL: %s\n\nImage: %s\n\nDescription:\n\n%s\n\nRecipe Score: %.2f\n",
			randRecipe.Title, *randRecipe.PrepTime, *randRecipe.Servings, randRecipe.SourceUrl, randRecipe.Image, cleanDescription, randRecipe.RecipeScore)

		return recipeDetails, nil
	}

	return "", fmt.Errorf("error fetching recipe from URL: %s", url)

}

func main() {
	details, err := getData()
	if err != nil {
		fmt.Println("error fetching recipe details", err)
		return
	}
	fmt.Println(details)
}
