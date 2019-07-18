package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"syscall/js"
)

type RecipesModel struct {
	Uuid uuid.UUID
	Name string
	Time string
	Type string
}

type RecipeModel struct {
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

var recipesModels []RecipesModel
var recipesModel map[uuid.UUID]RecipeModel

func main() {
	c := make(chan struct{}, 0)
	registerCallbacks()
	startup()
	<-c
}

func registerCallbacks() {
	js.Global().Set("search", js.FuncOf(search))
	js.Global().Set("clearSearch", js.FuncOf(clearSearch))
	js.Global().Set("createRecipe", js.FuncOf(createRecipe))
	js.Global().Set("recipeBack", js.FuncOf(recipeBack))
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
	for _, recipe := range recipesModels {
		if strings.Contains(strings.ToLower(recipe.Name), searchTerm) ||
			strings.Contains(strings.ToLower(recipe.Type), searchTerm) ||
			strings.Contains(strings.ToLower(recipe.Time), searchTerm) {
			htmlListRecipe := createRecipeListItem(&recipe)
			list.Call("appendChild", *htmlListRecipe)
		}
	}
}

func startup() {
	hash := js.Global().Get("window").Get("location").Get("hash")
	if hash.String() != "" {
		hashStr := hash.String()[1:]
		recipeUuid, err := uuid.Parse(hashStr)
		if err != nil {
			js.Global().Get("history").Call("pushState", nil, nil, " ")
		} else {
			openRecipe(recipeUuid)
		}
	}

	recipesModel = make(map[uuid.UUID]RecipeModel)
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

	var recipes []RecipesModel
	err = json.Unmarshal(body, &recipes)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, recipe := range recipes {
		htmlListRecipe := createRecipeListItem(&recipe)
		list.Call("appendChild", *htmlListRecipe)
	}

	recipesModels = recipes
}

func createRecipeListItem(recipe *RecipesModel) *js.Value {
	li := js.Global().Get("document").Call("createElement", "li")
	innerHtml := "<div class=\"recipe-list-item\" onClick=\"createRecipe('" + recipe.Uuid.String() + "');\">" +
		"<h3>" + recipe.Name + "</h3>" +
		"<p>Type: " + recipe.Type + "</p>" +
		"<p>Time: " + recipe.Time + "</p>" +
		"</div>"
	li.Set("innerHTML", innerHtml)
	return &li
}

func createRecipe(this js.Value, i []js.Value) interface{} {
	js.Global().Get("history").Call("pushState", nil, nil, "#"+i[0].String())
	recipeUuid, err := uuid.Parse(i[0].String())
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	openRecipe(recipeUuid)
	return nil
}

func openRecipe(recipeUuid uuid.UUID) {
	outerDiv := js.Global().Get("document").Call("createElement", "div")
	outerDiv.Set("className", "recipe-item")
	outerDiv.Set("id", "recipe-item")
	backButton := "<div class=\"recipe-item-back\" onClick=\"recipeBack();\">"
	outerDiv.Set("innerHTML", backButton)
	js.Global().Get("document").Get("body").Call("appendChild", outerDiv)

	go func() {

		var recipe RecipeModel

		if _, model := recipesModel[recipeUuid]; model {
			recipe = recipesModel[recipeUuid]
		} else {
			req, err := http.NewRequest("GET", serverURL+"/recipe"+"/"+recipeUuid.String(), nil)
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

			err = json.Unmarshal(body, &recipe)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			recipesModel[recipeUuid] = recipe
		}

		outerDiv := js.Global().Get("document").Call("getElementById", "recipe-item")
		innerDiv := js.Global().Get("document").Call("createElement", "div")
		innerDiv.Set("className", "recipe-item-inner")
		// Recipe Metadata
		title := js.Global().Get("document").Call("createElement", "h1")
		title.Set("innerHTML", recipe.Name)
		rType := js.Global().Get("document").Call("createElement", "h3")
		rType.Set("innerHTML", "Type: "+recipe.Type)
		time := js.Global().Get("document").Call("createElement", "h3")
		time.Set("innerHTML", "Time: "+recipe.Time)
		innerDiv.Call("appendChild", title)
		innerDiv.Call("appendChild", rType)
		innerDiv.Call("appendChild", time)
		// Equipment
		equipmentTitle := js.Global().Get("document").Call("createElement", "h3")
		equipmentTitle.Set("innerHTML", "Equipment:")
		equipmentList := js.Global().Get("document").Call("createElement", "ul")
		equipmentList.Set("className", "recipe-item-equipment")
		for _, eq := range *recipe.Equipments {
			equipment := js.Global().Get("document").Call("createElement", "li")
			equipment.Set("innerHTML", strconv.Itoa(eq.Quantity)+" - "+eq.Name)
			equipmentList.Call("appendChild", equipment)
		}
		innerDiv.Call("appendChild", equipmentTitle)
		innerDiv.Call("appendChild", equipmentList)
		// Ingredients
		ingredientTitle := js.Global().Get("document").Call("createElement", "h3")
		ingredientTitle.Set("innerHTML", "Ingredients:")
		ingredientList := js.Global().Get("document").Call("createElement", "ul")
		ingredientList.Set("className", "recipe-item-ingredient")
		for _, ing := range *recipe.Ingredients {
			ingredient := js.Global().Get("document").Call("createElement", "li")
			ingredient.Set("innerHTML", strconv.Itoa(ing.Quantity)+" - "+ing.Name)
			ingredientList.Call("appendChild", ingredient)
		}
		innerDiv.Call("appendChild", ingredientTitle)
		innerDiv.Call("appendChild", ingredientList)

		// Steps
		stepTitle := js.Global().Get("document").Call("createElement", "h3")
		stepTitle.Set("innerHTML", "Steps:")
		stepList := js.Global().Get("document").Call("createElement", "ul")
		stepList.Set("className", "recipe-item-step")
		for _, st := range *recipe.Steps {
			step := js.Global().Get("document").Call("createElement", "li")
			step.Set("innerHTML", st.Content)
			stepList.Call("appendChild", step)
		}
		innerDiv.Call("appendChild", stepTitle)
		innerDiv.Call("appendChild", stepList)

		// Pictures
		images := js.Global().Get("document").Call("createElement", "div")
		images.Set("className", "recipe-item-image")
		imgHtml := ""
		for _, img := range *recipe.Pictures {
			imgHtml += "<img src=\"" + img.ImageSource + "\" alt\"" + img.Uuid.String() + "\"/>"
		}
		images.Set("innerHTML", imgHtml)
		innerDiv.Call("appendChild", images)

		outerDiv.Call("appendChild", innerDiv)
	}()
}

func recipeBack(this js.Value, i []js.Value) interface{} {
	js.Global().Get("history").Call("pushState", nil, nil, " ")
	recipeDiv := js.Global().Get("document").Call("getElementById", "recipe-item")
	js.Global().Get("document").Get("body").Call("removeChild", recipeDiv.JSValue())
	return nil
}
