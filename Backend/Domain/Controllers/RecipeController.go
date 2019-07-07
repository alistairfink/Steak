package Controllers

import (
	"encoding/json"
	"github.com/alistairfink/Steak/Backend/Data/Models"
	"github.com/alistairfink/Steak/Backend/Domain/Managers"
	"github.com/alistairfink/Steak/Backend/Utilities"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-pg/pg"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
)

type RecipeController struct {
	db            *pg.DB
	config        *Utilities.Config
	recipeManager *Managers.RecipeManger
}

func NewRecipeController(
	db *pg.DB,
	config *Utilities.Config,
) *RecipeController {
	recipeManager := Managers.RecipeManger{}
	recipeManager.Init(db)
	return &RecipeController{
		db:            db,
		config:        config,
		recipeManager: &recipeManager,
	}
}

func (this *RecipeController) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{recipe_uuid}", this.Get)
	router.Get("/", this.GetAll)
	router.Put("/", this.Update)
	router.Post("/", this.Create)
	router.Delete("/{recipe_uuid}", this.Delete)
	return router
}

func (this *RecipeController) Get(w http.ResponseWriter, r *http.Request) {
	uuidUnparsed := chi.URLParam(r, "recipe_uuid")
	uuid, err := uuid.Parse(uuidUnparsed)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Invalid Recipe Uuid", http.StatusBadRequest)
		return
	}

	result := this.recipeManager.Get(uuid)
	if result == nil {
		http.Error(w, "Error Processing Request", http.StatusBadRequest)
		return
	}

	render.JSON(w, r, result)
}

func (this *RecipeController) GetAll(w http.ResponseWriter, r *http.Request) {
	result := this.recipeManager.GetAll()
	render.JSON(w, r, result)
}

func (this *RecipeController) Update(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("APIKey")
	if apiKey == this.config.ApiKey {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error Processing JSON", http.StatusBadRequest)
			return
		}

		var model Models.RecipeDomainModel
		err = json.Unmarshal(body, &model)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error Processing JSON", http.StatusBadRequest)
			return
		}

		result := this.recipeManager.Update(&model)
		if result == nil {
			http.Error(w, "Error Processing Request", http.StatusBadRequest)
			return
		}

		render.JSON(w, r, result)
	} else {
		http.Error(w, "Invalid API Key", http.StatusUnauthorized)
	}
}

func (this *RecipeController) Create(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("APIKey")
	if apiKey == this.config.ApiKey {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error Processing JSON", http.StatusBadRequest)
			return
		}

		var model Models.RecipeDomainModel
		err = json.Unmarshal(body, &model)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error Processing JSON", http.StatusBadRequest)
			return
		}

		result := this.recipeManager.Create(&model)
		if result == nil {
			http.Error(w, "Error Processing Request", http.StatusBadRequest)
			return
		}

		render.JSON(w, r, result)
	} else {
		http.Error(w, "Invalid API Key", http.StatusUnauthorized)
	}
}

func (this *RecipeController) Delete(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("APIKey")
	if apiKey == this.config.ApiKey {
		uuidUnparsed := chi.URLParam(r, "recipe_uuid")
		uuid, err := uuid.Parse(uuidUnparsed)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Invlaid Recipe Uuid", http.StatusBadRequest)
			return
		}

		if !this.recipeManager.Delete(uuid) {
			http.Error(w, "Error Processing Request", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Invalid API Key", http.StatusUnauthorized)
	}
}
