package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

func getUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')

	return text
}

func cleanDescription(htmlText string) string {
	cleanText := strings.ReplaceAll(htmlText, "<[^>]*", "")
	cleanText = strings.TrimSpace(cleanText)

	return cleanText
}

func getData(c string, i string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	apiKey := os.Getenv("SPOONACULAR_API_KEY")
	category, ingredient := c, i

	url := fmt.Sprintf("https://api.spoonacular.com/recipes/complexSearch?apiKey=%s&type=%s&query=%s&addRecipeInformation=True&number=1&sort=random", apiKey, category, ingredient)

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

		recipeDetails := fmt.Sprintf("\n\n\nTitle: %s\n\nPrep Time (minutes): %d\n\nServings: %d\n\nSource URL: %s\n\nImage: %s\n\nDescription:\n\n%s\n\nRecipe Score: %.2f\n",
			randRecipe.Title, *randRecipe.PrepTime, *randRecipe.Servings, randRecipe.SourceUrl, randRecipe.Image, cleanDescription, randRecipe.RecipeScore)

		return recipeDetails, nil
	}

	return "", fmt.Errorf("error fetching recipe from URL: %s", url)

}

func main() {
	for {
		mainMenu := []string{
			"Main menu:\n",
			"1. Get a recipe",
			"2. Learn a surprising food fact",
			"3. Exit",
		}

		for _, item := range mainMenu {
			fmt.Println(item)
		}

		userInput := strings.TrimSpace(getUserInput("\nEnter a selection (1-3): "))

		selection, err := strconv.Atoi(userInput)
		if err != nil {
			fmt.Println("\nInvalid input")
			continue
		}

		switch selection {
		case 1:
			for {
				categories := []string{
					"Meal Categories:\n",
					"1. Breakfast",
					"2. Snack",
					"3. Salad",
					"4. Main Course",
					"5. Dessert",
					"6. Drink",
					"7. <-- Back",
				}

				for _, item := range categories {
					fmt.Println(item)
				}

				userInput := strings.TrimSpace(getUserInput("\nWhich category of food are you looking for?: "))

				category, err := strconv.Atoi(userInput)
				if err != nil {
					fmt.Println("\nInvalid input")
					continue
				}

				switch category {
				case 1:
					category := "Breakfast"
					fmt.Printf("\nCategory: %s\n", category)

					ingredient := strings.TrimSpace(getUserInput("\nWhat ingredients? (Max: 3) eg. beef egg blueberry tomato: "))
					fmt.Printf("\nIngredients: %s\n", ingredient)

					recipe, err := getData(category, ingredient)
					if err != nil {
						fmt.Println("Case 1: error getting recipe")
						return
					}

					time.Sleep(1 * time.Second)
					fmt.Println("Getting a recipe for you...")
					time.Sleep(1 * time.Second)
					fmt.Println(recipe)
					time.Sleep(1 * time.Second)

					return
				case 2:
					category := "Snack"
					fmt.Printf("\nCategory: %s\n", category)

					ingredient := strings.TrimSpace(getUserInput("\nWhat ingredients? (Max: 3) eg. beef egg blueberry tomato: "))
					fmt.Printf("\nIngredients: %s\n", ingredient)

					recipe, err := getData(category, ingredient)
					if err != nil {
						fmt.Println("Case 1: error getting recipe")
						return
					}

					time.Sleep(1 * time.Second)
					fmt.Println("Getting a recipe for you...")
					time.Sleep(1 * time.Second)
					fmt.Println(recipe)
					time.Sleep(1 * time.Second)

					return
				case 3:
					category := "Salad"
					fmt.Printf("\nCategory: %s\n", category)

					ingredient := strings.TrimSpace(getUserInput("\nWhat ingredients? (Max: 3) eg. beef egg blueberry tomato: "))
					fmt.Printf("\nIngredients: %s\n", ingredient)

					recipe, err := getData(category, ingredient)
					if err != nil {
						fmt.Println("Case 1: error getting recipe")
						return
					}

					time.Sleep(1 * time.Second)
					fmt.Println("Getting a recipe for you...")
					time.Sleep(1 * time.Second)
					fmt.Println(recipe)
					time.Sleep(1 * time.Second)

					return
				case 4:
					category := "Main Course"
					fmt.Printf("\nCategory: %s\n", category)
					ingredient := strings.TrimSpace(getUserInput("\nWhat ingredients? (Max: 3) eg. beef egg blueberry tomato: "))

					fmt.Printf("\nIngredients: %s\n", ingredient)

					recipe, err := getData(category, ingredient)
					if err != nil {
						fmt.Println("Case 1: error getting recipe")
						return
					}

					time.Sleep(1 * time.Second)
					fmt.Println("Getting a recipe for you...")
					time.Sleep(1 * time.Second)
					fmt.Println(recipe)
					time.Sleep(1 * time.Second)

					return
				case 5:
					category := "Dessert"
					fmt.Printf("\nCategory: %s\n", category)

					ingredient := strings.TrimSpace(getUserInput("\nWhat ingredients? (Max: 3) eg. beef egg blueberry tomato: "))
					fmt.Printf("\nIngredients: %s\n", ingredient)

					recipe, err := getData(category, ingredient)
					if err != nil {
						fmt.Println("Case 1: error getting recipe")
						return
					}

					time.Sleep(1 * time.Second)
					fmt.Println("Getting a recipe for you...")
					time.Sleep(1 * time.Second)
					fmt.Println(recipe)
					time.Sleep(1 * time.Second)

					return
				case 6:
					category := "Drink"
					fmt.Printf("\nCategory: %s\n", category)

					ingredient := strings.TrimSpace(getUserInput("\nWhat ingredients? (Max: 3) eg. beef egg blueberry tomato: "))
					fmt.Printf("\nIngredients: %s\n", ingredient)

					recipe, err := getData(category, ingredient)
					if err != nil {
						fmt.Println("Case 1: error getting recipe")
						return
					}

					time.Sleep(1 * time.Second)
					fmt.Println("Getting a recipe for you...")
					time.Sleep(1 * time.Second)
					fmt.Println(recipe)
					time.Sleep(1 * time.Second)

					return
				case 7:
					main()
				default:
					fmt.Println("\nInvalid selection. Enter (1-7).")
				}
			}
		case 2:
			fmt.Println("Option 2")
		case 3:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("\nInvalid selection. Enter 1, 2 or 3.")
		}
	}
}
