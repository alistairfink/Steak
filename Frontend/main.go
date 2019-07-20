package main

import (
	"bytes"
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
var ApiKey string

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
	js.Global().Set("setApiKey", js.FuncOf(setApiKey))
	js.Global().Set("deleteRecipe", js.FuncOf(deleteRecipe))
	js.Global().Set("openEdit", js.FuncOf(openEdit))
	js.Global().Set("sendEdit", js.FuncOf(sendEdit))
	js.Global().Set("editBack", js.FuncOf(editBack))
	js.Global().Set("openAdd", js.FuncOf(openAdd))
	js.Global().Set("sendNew", js.FuncOf(sendNew))
}

func clearSearch(this js.Value, i []js.Value) interface{} {
	js.Global().Get("document").Call("getElementById", "search").Set("value", "")
	searchLogic()
	return nil
}

func search(this js.Value, i []js.Value) interface{} {
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

	hash := js.Global().Get("window").Get("location").Get("hash")
	if hash.String() == "#admin" {
		openAdmin()
	} else if hash.String() != "" {
		hashStr := hash.String()[1:]
		recipeUuid, err := uuid.Parse(hashStr)
		if err != nil {
			js.Global().Get("history").Call("pushState", nil, nil, " ")
		} else {
			openRecipe(recipeUuid)
		}
	}
}

func createRecipeListItem(recipe *RecipesModel) *js.Value {
	li := js.Global().Get("document").Call("createElement", "li")
	innerHtml := "<div class=\"recipe-list-item\" onClick=\"createRecipe('" + recipe.Uuid.String() + "');\">" +
		"<h3>" + recipe.Name + "</h3>" +
		"<p>Type: " + recipe.Type + "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Time: " + recipe.Time + "</p>" +
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
	backButton := "<div class=\"recipe-item-back\" onClick=\"recipeBack();\"><img src=\"Back.png\" title=\"Back\" alt=\"Back\"/></div>"
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
		stepList := js.Global().Get("document").Call("createElement", "ol")
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
			imgHtml += "<img src=\"" + img.ImageSource + "\" alt=\"" + img.Uuid.String() + "\" class=\"recipe-image\"/>"
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

func openAdmin() {
	outerDiv := js.Global().Get("document").Call("createElement", "div")
	outerDiv.Set("className", "admin")
	js.Global().Get("document").Get("body").Call("appendChild", outerDiv)
	title := js.Global().Get("document").Call("createElement", "h1")
	title.Set("innerHTML", "Admin")
	outerDiv.Call("appendChild", title)

	apiKeyOuter := js.Global().Get("document").Call("createElement", "div")
	apiKeyOuter.Set("innerHTML", "<input type=\"text\" id=\"apiKey\" placeholder=\"Api Key\">"+
		"<button id=\"apiKey_set\" onClick=\"setApiKey('apiKey');\">Set ApiKey</button>"+
		"<button id=\"add\" onClick=\"openAdd();\">Add</button>")
	outerDiv.Call("appendChild", apiKeyOuter)

	for _, rec := range recipesModels {
		elementOuter := js.Global().Get("document").Call("createElement", "div")
		elementOuter.Set("className", "admin-element-outer")
		elementOuter.Set("innerHTML",
			"<p>"+rec.Uuid.String()+"&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</p>"+
				"<p>"+rec.Name+"</p>"+
				"<button id=\""+rec.Uuid.String()+"_edit\" onClick=\"openEdit('"+rec.Uuid.String()+"')\">Edit</button>"+
				"<button id=\""+rec.Uuid.String()+"_delete\" onClick=\"deleteRecipe('"+rec.Uuid.String()+"');\">Delete</button>")
		outerDiv.Call("appendChild", elementOuter)
	}
}

func setApiKey(this js.Value, i []js.Value) interface{} {
	apiKeyId := i[0].String()
	apiKey := js.Global().Get("document").Call("getElementById", apiKeyId)
	ApiKey = apiKey.Get("value").String()
	return nil
}

func deleteRecipe(this js.Value, i []js.Value) interface{} {
	recipeUuid, err := uuid.Parse(i[0].String())
	if err != nil {
		println(err.Error())
		return nil
	}

	go func() {
		req, _ := http.NewRequest("DELETE", serverURL+"/recipe"+"/"+recipeUuid.String(), nil)
		req.Header.Set("APIKey", ApiKey)
		http.DefaultClient.Do(req)
		js.Global().Get("location").Call("reload")
	}()

	return nil
}

func openEdit(this js.Value, i []js.Value) interface{} {
	recipeUuid, _ := uuid.Parse(i[0].String())
	outerDiv := js.Global().Get("document").Call("createElement", "div")
	outerDiv.Set("id", "admin-edit")
	outerDiv.Set("className", "admin")
	backButton := "<div class=\"recipe-item-back\" onClick=\"editBack();\"><img src=\"Back.png\" title=\"Back\" alt=\"Back\"/></div>"
	outerDiv.Set("innerHTML", backButton)
	js.Global().Get("document").Get("body").Call("appendChild", outerDiv)
	title := js.Global().Get("document").Call("createElement", "h1")
	title.Set("innerHTML", "Edit")
	outerDiv.Call("appendChild", title)

	apiKeyOuter := js.Global().Get("document").Call("createElement", "div")
	apiKeyOuter.Set("innerHTML", "<input type=\"text\" id=\"apiKey_edit\" placeholder=\"Api Key\">"+
		"<button id=\"apiKey_set\" onClick=\"setApiKey('apiKey_edit');\">Set ApiKey</button>"+
		"<button id=\"apiKey_set\" onClick=\"sendEdit();\">Send</button>")
	outerDiv.Call("appendChild", apiKeyOuter)

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

		json, _ := json.MarshalIndent(recipe, "", "    ")
		textArea := js.Global().Get("document").Call("createElement", "textarea")
		textArea.Get("value")
		textArea.Set("id", "edit_json")
		textArea.Set("name", "edit_json")
		textArea.Set("className", "json-editor")

		outerDiv.Call("appendChild", textArea)

		js.Global().Get("document").Call("getElementById", "edit_json").Set("value", string(json))
	}()

	return nil
}

func sendEdit(this js.Value, i []js.Value) interface{} {
	editRecipeJson := js.Global().Get("document").Call("getElementById", "edit_json").Get("value")
	go func() {
		req, _ := http.NewRequest("PUT", serverURL+"/recipe", bytes.NewBuffer([]byte(editRecipeJson.String())))
		req.Header.Set("APIKey", ApiKey)
		req.Header.Set("Content-Type", "application/json")
		http.DefaultClient.Do(req)
		js.Global().Get("location").Call("reload")
	}()

	return nil
}

func editBack(this js.Value, i []js.Value) interface{} {
	recipeDiv := js.Global().Get("document").Call("getElementById", "admin-edit")
	js.Global().Get("document").Get("body").Call("removeChild", recipeDiv.JSValue())
	return nil
}

func openAdd(this js.Value, i []js.Value) interface{} {
	outerDiv := js.Global().Get("document").Call("createElement", "div")
	outerDiv.Set("id", "admin-edit")
	outerDiv.Set("className", "admin")
	backButton := "<div class=\"recipe-item-back\" onClick=\"editBack();\"><img src=\"Back.png\" title=\"Back\" alt=\"Back\"/></div>"
	outerDiv.Set("innerHTML", backButton)
	js.Global().Get("document").Get("body").Call("appendChild", outerDiv)
	title := js.Global().Get("document").Call("createElement", "h1")
	title.Set("innerHTML", "New")
	outerDiv.Call("appendChild", title)

	apiKeyOuter := js.Global().Get("document").Call("createElement", "div")
	apiKeyOuter.Set("innerHTML", "<input type=\"text\" id=\"apiKey_edit\" placeholder=\"Api Key\">"+
		"<button id=\"apiKey_set\" onClick=\"setApiKey('apiKey_edit');\">Set ApiKey</button>"+
		"<button id=\"apiKey_set\" onClick=\"sendNew();\">Send</button>")
	outerDiv.Call("appendChild", apiKeyOuter)

	textArea := js.Global().Get("document").Call("createElement", "textarea")
	textArea.Get("value")
	textArea.Set("id", "edit_json")
	textArea.Set("name", "edit_json")
	textArea.Set("className", "json-editor")

	outerDiv.Call("appendChild", textArea)
	json :=
		"{\n" +
			"\t\"Name\": \"\",\n" +
			"\t\"Time\": \"\",\n" +
			"\t\"Type\": \"\",\n" +
			"\t\"Pictures\": [\n" +
			"\t\t{\n" +
			"\t\t\t\"ImageSource\": \"\",\n" +
			"\t\t\t\"SortOrder\": \n" +
			"\t\t}\n" +
			"\t],\n" +
			"\t\"Equipments\": [\n" +
			"\t\t{\n" +
			"\t\t\t\"Name\": \"\",\n" +
			"\t\t\t\"Quantity\": \n" +
			"\t\t}\n" +
			"\t],\n" +
			"\t\"Ingredients\": [\n" +
			"\t\t{\n" +
			"\t\t\t\"Name\": \"\",\n" +
			"\t\t\t\"Quantity\": \n" +
			"\t\t}\n" +
			"\t],\n" +
			"\t\"Steps\": [\n" +
			"\t\t{\n" +
			"\t\t\t\"Content\": \"\",\n" +
			"\t\t\t\"StepNumber\": \n" +
			"\t\t}\n" +
			"\t]\n" +
			"}\n"

	js.Global().Get("document").Call("getElementById", "edit_json").Set("value", string(json))

	return nil
}

func sendNew(this js.Value, i []js.Value) interface{} {
	newRecipeJson := js.Global().Get("document").Call("getElementById", "edit_json").Get("value")
	go func() {
		req, _ := http.NewRequest("POST", serverURL+"/recipe", bytes.NewBuffer([]byte(newRecipeJson.String())))
		req.Header.Set("APIKey", ApiKey)
		req.Header.Set("Content-Type", "application/json")
		http.DefaultClient.Do(req)
		js.Global().Get("location").Call("reload")
	}()

	return nil
}
