package Commands

import (
	"github.com/alistairfink/Steak/Backend/Data/Models"
	"github.com/go-pg/pg"
	"github.com/google/uuid"
)

type RecipeCommand struct {
	DB *pg.DB
}

func (this *RecipeCommand) Get(uuid uuid.UUID) *Models.RecipeModel {
	if !this.Exists(uuid) {
		return nil
	}

	var models []Models.RecipeModel
	err := this.DB.Model(&models).Where("id = ?", uuid).Select()
	if err != nil {
		panic(err)
	}

	return &models[0]
}

func (this *RecipeCommand) GetAll() *[]Models.RecipeModel {
	var models []Models.RecipeModel
	err := this.DB.Model(&models).Select()
	if err != nil {
		panic(err)
	}

	return &models
}

func (this *RecipeCommand) Upsert(model *Models.RecipeModel) *Models.RecipeModel {
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

func (this *RecipeCommand) Delete(uuid uuid.UUID) bool {
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

func (this *RecipeCommand) Exists(uuid uuid.UUID) bool {
	var models []Models.RecipeModel
	exists, err := this.DB.Model(&models).Where("id = ?", uuid).Exists()
	if err != nil {
		panic(err)
	}

	return exists
}
