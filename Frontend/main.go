package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strings"
	"syscall/js"
)

type RecipeModel struct {
	Uuid uuid.UUID
	Name string
	Time string
	Type string
}

type RecipesModel struct {
	Uuid        uuid.UUID
	Name        string
	Time        string
	Type        string
	Pictures    *[]PictureModel
	Equipments  *[]EquipmentModel
	Ingredients *[]IngredientModel
	Steps       *[]StepModel
}

type PictureModel struct {
	Uuid        uuid.UUID
	RecipeUuid  uuid.UUID
	ImageSource string
	SortOrder   int
}

type StepModel struct {
	Uuid       uuid.UUID
	RecipeUuid uuid.UUID
	Content    string
	StepNumber int
}

type EquipmentModel struct {
	Uuid       uuid.UUID
	RecipeUuid uuid.UUID
	Name       string
	Quantity   int
}

type IngredientModel struct {
	Uuid       uuid.UUID
	RecipeUuid uuid.UUID
	Name       string
	Quantity   int
}

const serverURL = "http://localhost:41690"

var recipeModels []RecipeModel

func main() {
	c := make(chan struct{}, 0)
	registerCallbacks()
	startup()
	<-c
}

func registerCallbacks() {
	js.Global().Set("search", js.FuncOf(search))
	js.Global().Set("clearSearch", js.FuncOf(clearSearch))
}

func search(this js.Value, i []js.Value) interface{} {
	searchLogic()
	return nil
}

func clearSearch(this js.Value, i []js.Value) interface{} {
	js.Global().Get("document").Call("getElementById", "search").Set("value", "")
	searchLogic()
	return nil
}

func searchLogic() {
	list := js.Global().Get("document").Call("getElementById", "recipe-list")
	list.Set("innerHTML", "")
	searchTerm := strings.ToLower(js.Global().Get("document").Call("getElementById", "search").Get("value").String())
	for _, recipe := range recipeModels {
		if strings.Contains(strings.ToLower(recipe.Name), searchTerm) ||
			strings.Contains(strings.ToLower(recipe.Type), searchTerm) ||
			strings.Contains(strings.ToLower(recipe.Time), searchTerm) {
			htmlListRecipe := createRecipeListItem(&recipe)
			list.Call("appendChild", *htmlListRecipe)
		}
	}
}

func startup() {
	// println(js.Global().Get("window").Get("location").Get("hash").String())
	list := js.Global().Get("document").Call("getElementById", "recipe-list")
	list.Set("innerHTML", "")
	req, err := http.NewRequest("GET", serverURL+"/recipe", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	response, err := http.DefaultClient.Do(req)
	defer response.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var recipes []RecipeModel
	err = json.Unmarshal(body, &recipes)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, recipe := range recipes {
		htmlListRecipe := createRecipeListItem(&recipe)
		list.Call("appendChild", *htmlListRecipe)
	}

	recipeModels = recipes
}

func createRecipeListItem(recipe *RecipeModel) *js.Value {
	// Outter
	li := js.Global().Get("document").Call("createElement", "li")
	outerDiv := js.Global().Get("document").Call("createElement", "div")
	outerDiv.Set("className", "recipe-list-item")
	li.Call("appendChild", outerDiv)
	title := js.Global().Get("document").Call("createElement", "h3")
	title.Set("innerHTML", recipe.Name)
	recipeType := js.Global().Get("document").Call("createElement", "p")
	recipeType.Set("innerHTML", "Type: "+recipe.Type)
	recipeTime := js.Global().Get("document").Call("createElement", "p")
	recipeTime.Set("innerHTML", "Time: "+recipe.Time)
	outerDiv.Call("appendChild", title)
	outerDiv.Call("appendChild", recipeType)
	outerDiv.Call("appendChild", recipeTime)
	return &li
}
