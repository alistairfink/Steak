package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
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

var testData []RecipeModel

func main() {
	c := make(chan struct{}, 0)
	registerCallbacks()
	startup()
	<-c
}

func registerCallbacks() {
	js.Global().Set("search", js.FuncOf(search))
}

func search(this js.Value, i []js.Value) interface{} {
	list := js.Global().Get("document").Call("getElementById", "recipe-list")
	li := js.Global().Get("document").Call("createElement", "li")
	li.Set("innerHTML", "test")
	list.Call("appendChild", li)
	return nil
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

	testData = recipes
}

func recipeListItem() {

}
