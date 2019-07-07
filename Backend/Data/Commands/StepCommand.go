package Commands

import (
	"github.com/alistairfink/Steak/Backend/Data/Models"
	"github.com/go-pg/pg"
	"github.com/google/uuid"
)

type StepCommand struct {
	DB *pg.DB
}

func (this *StepCommand) Get(uuid uuid.UUID) *Models.StepModel {
	if !this.Exists(uuid) {
		return nil
	}

	var models []Models.StepModel
	err := this.DB.Model(&models).Where("id = ?", uuid).Select()
	if err != nil {
		panic(err)
	}

	return &models[0]
}

func (this *StepCommand) GetAll() *[]Models.StepModel {
	var models []Models.StepModel
	err := this.DB.Model(&models).Select()
	if err != nil {
		panic(err)
	}

	return &models
}

func (this *StepCommand) GetByRecipeUuid(recipeUuid uuid.UUID) *[]Models.StepModel {
	var models []Models.StepModel
	err := this.DB.Model(&models).Where("recipe_id = ?", recipeUuid).Select()
	if err != nil {
		panic(err)
	}

	return &models
}

func (this *StepCommand) Upsert(model *Models.StepModel) *Models.StepModel {
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

func (this *StepCommand) Delete(uuid uuid.UUID) bool {
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

func (this *StepCommand) Exists(uuid uuid.UUID) bool {
	var models []Models.StepModel
	exists, err := this.DB.Model(&models).Where("id = ?", uuid).Exists()
	if err != nil {
		panic(err)
	}

	return exists
}
