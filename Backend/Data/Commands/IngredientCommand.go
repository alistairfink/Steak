package Commands

import (
	"github.com/alistairfink/Steak/Backend/Data/Models"
	"github.com/go-pg/pg"
	"github.com/google/uuid"
)

type IngredientCommand struct {
	DB *pg.DB
}

func (this *IngredientCommand) Get(uuid uuid.UUID) *Models.IngredientModel {
	if !this.Exists(uuid) {
		return nil
	}

	var models []Models.IngredientModel
	err := this.DB.Model(&models).Where("id = ?", uuid).Select()
	if err != nil {
		panic(err)
	}

	return &models[0]
}

func (this *IngredientCommand) GetAll() *[]Models.IngredientModel {
	var models []Models.IngredientModel
	err := this.DB.Model(&models).Select()
	if err != nil {
		panic(err)
	}

	return &models
}

func (this *IngredientCommand) GetByRecipeUuid(recipeUuid uuid.UUID) *[]Models.IngredientModel {
	var models []Models.IngredientModel
	err := this.DB.Model(&models).Where("recipe_id = ?", recipeUuid).Select()
	if err != nil {
		panic(err)
	}

	return &models
}

func (this *IngredientCommand) Upsert(model *Models.IngredientModel) *Models.IngredientModel {
	if this.Exists(model.Uuid) {
		_, err := this.DB.Model(model).Where("id = ?", model.Uuid).Update(model)
		if err != nil {
			panic(err)
		}
	} else {
		err := this.DB.Insert(model)
		if err != nil {
			panic(err)
		}
	}

	return this.Get(model.Uuid)
}

func (this *IngredientCommand) Delete(uuid uuid.UUID) bool {
	model := this.Get(uuid)
	if model == nil {
		return false
	}

	err := this.DB.Delete(model)
	if err != nil {
		panic(err)
	}

	return true
}

func (this *IngredientCommand) Exists(uuid uuid.UUID) bool {
	var models []Models.IngredientModel
	exists, err := this.DB.Model(&models).Where("id = ?", uuid).Exists()
	if err != nil {
		panic(err)
	}

	return exists
}
