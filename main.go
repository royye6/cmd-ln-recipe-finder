package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
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

type Trivia struct {
	Text string `json:"text"`
}

func getUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')

	return text
}

func cleanDescription(htmlText string) string {
	re := regexp.MustCompile("<.*?>")
	cleanText := re.ReplaceAllString(htmlText, "")
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

		recipeDetails := fmt.Sprintf("*****\n\nTitle: %s\n\nPrep Time (minutes): %d\n\nServings: %d\n\nSource URL: %s\n\nImage: %s\n\nDescription:\n\n%s\n\nRecipe Score: %.2f\n",
			randRecipe.Title, *randRecipe.PrepTime, *randRecipe.Servings, randRecipe.SourceUrl, randRecipe.Image, cleanDescription, randRecipe.RecipeScore)

		return recipeDetails, nil
	}

	return "", fmt.Errorf("error fetching recipe from URL: %s", url)

}

func saveData(r string) error {
	file, err := os.OpenFile("recipe.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(r + "\n\n")
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

func getFoodFact() (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	apiKey := os.Getenv("SPOONACULAR_API_KEY")

	url := fmt.Sprintf("https://api.spoonacular.com/food/trivia/random?apiKey=%s", apiKey)
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

		var foodFact Trivia
		err = json.Unmarshal(body, &foodFact)

		if err != nil {
			return "", fmt.Errorf("error unmarshalling(deserializing) JSON response: %v", err)
		}

		fact := foodFact.Text

		randFact := fmt.Sprintf("\n%s\n", fact)

		return randFact, nil
	}

	return "", fmt.Errorf("error fetching trivia from URL: %s", url)
}

func main() {
	for {
		mainMenu := []string {
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
					for {
						save_data := strings.TrimSpace(getUserInput("\nDo you want to save this recipe to recipe.txt? y/n: "))
						switch save_data {
						case "y":
							saveData(recipe)
							time.Sleep(1 * time.Second)
							fmt.Println("Recipe saved successfully...")
							fmt.Println("\nReturned to main menu\n ")
							time.Sleep(1 * time.Second)
							main()

							return
						case "n":
							fmt.Println("Exiting...")
							time.Sleep(1 * time.Second)
							return
						default:
							fmt.Println("\nInvalid choice")
						}
					}
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
					for {
						save_data := strings.TrimSpace(getUserInput("\nDo you want to save this recipe to recipe.txt? y/n: "))
						switch save_data {
						case "y":
							saveData(recipe)
							time.Sleep(1 * time.Second)
							fmt.Println("Recipe saved successfully...")
							fmt.Println("\nReturned to main menu\n ")
							time.Sleep(1 * time.Second)
							main()

							return
						case "n":
							fmt.Println("Exiting...")
							time.Sleep(1 * time.Second)
							return
						default:
							fmt.Println("\nInvalid choice")
						}
					}
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
					for {
						save_data := strings.TrimSpace(getUserInput("\nDo you want to save this recipe to recipe.txt? y/n: "))
						switch save_data {
						case "y":
							saveData(recipe)
							time.Sleep(1 * time.Second)
							fmt.Println("Recipe saved successfully...")
							fmt.Println("\nReturned to main menu\n ")
							time.Sleep(1 * time.Second)
							main()

							return
						case "n":
							fmt.Println("Exiting...")
							time.Sleep(1 * time.Second)
							return
						default:
							fmt.Println("\nInvalid choice")
						}
					}
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
					for {
						save_data := strings.TrimSpace(getUserInput("\nDo you want to save this recipe to recipe.txt? y/n: "))
						switch save_data {
						case "y":
							saveData(recipe)
							time.Sleep(1 * time.Second)
							fmt.Println("Recipe saved successfully...")
							fmt.Println("\nReturned to main menu\n ")
							time.Sleep(1 * time.Second)
							main()

							return
						case "n":
							fmt.Println("Exiting...")
							time.Sleep(1 * time.Second)
							return
						default:
							fmt.Println("\nInvalid choice")
						}
					}
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
					for {
						save_data := strings.TrimSpace(getUserInput("\nDo you want to save this recipe to recipe.txt? y/n: "))
						switch save_data {
						case "y":
							saveData(recipe)
							time.Sleep(1 * time.Second)
							fmt.Println("Recipe saved successfully...")
							fmt.Println("\nReturned to main menu\n ")
							time.Sleep(1 * time.Second)
							main()

							return
						case "n":
							fmt.Println("Exiting...")
							time.Sleep(1 * time.Second)
							return
						default:
							fmt.Println("\nInvalid choice")
						}
					}
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
					for {
						save_data := strings.TrimSpace(getUserInput("\nDo you want to save this recipe to recipe.txt? y/n: "))
						switch save_data {
						case "y":
							saveData(recipe)
							time.Sleep(1 * time.Second)
							fmt.Println("Recipe saved successfully...")
							fmt.Println("\nReturned to main menu\n ")
							time.Sleep(1 * time.Second)
							main()

							return
						case "n":
							fmt.Println("Exiting...")
							time.Sleep(1 * time.Second)
							return
						default:
							fmt.Println("\nInvalid choice")
						}
					}
				case 7:
					main()
				default:
					fmt.Println("\nInvalid selection. Enter (1-7).")
				}
			}
		case 2:
			trivia, err := getFoodFact()
			if err != nil {
				fmt.Println("error fetching trivia")
				return
			}
			fmt.Printf(trivia + "\n")
			time.Sleep(3 * time.Second)
			return
		case 3:
			fmt.Println("Exiting...")
			time.Sleep(1 * time.Second)
			return
		default:
			fmt.Println("\nInvalid selection. Enter 1, 2 or 3.")
		}
	}
}
