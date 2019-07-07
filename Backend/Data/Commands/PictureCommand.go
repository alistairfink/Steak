package Commands

import (
	"github.com/alistairfink/Steak/Backend/Data/Models"
	"github.com/go-pg/pg"
	"github.com/google/uuid"
)

type PictureCommand struct {
	DB *pg.DB
}

func (this *PictureCommand) Get(uuid uuid.UUID) *Models.PictureModel {
	if !this.Exists(uuid) {
		return nil
	}

	var models []Models.PictureModel
	err := this.DB.Model(&models).Where("id = ?", uuid).Select()
	if err != nil {
		panic(err)
	}

	return &models[0]
}

func (this *PictureCommand) GetAll() *[]Models.PictureModel {
	var models []Models.PictureModel
	err := this.DB.Model(&models).Select()
	if err != nil {
		panic(err)
	}

	return &models
}

func (this *PictureCommand) GetByRecipeUuid(recipeUuid uuid.UUID) *[]Models.PictureModel {
	var models []Models.PictureModel
	err := this.DB.Model(&models).Where("recipe_id = ?", recipeUuid).Select()
	if err != nil {
		panic(err)
	}

	return &models
}

func (this *PictureCommand) Upsert(model *Models.PictureModel) *Models.PictureModel {
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

func (this *PictureCommand) Delete(uuid uuid.UUID) bool {
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

func (this *PictureCommand) Exists(uuid uuid.UUID) bool {
	var models []Models.PictureModel
	exists, err := this.DB.Model(&models).Where("id = ?", uuid).Exists()
	if err != nil {
		panic(err)
	}

	return exists
}
