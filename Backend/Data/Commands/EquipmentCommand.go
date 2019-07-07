package Commands

import (
	"github.com/alistairfink/Steak/Backend/Data/Models"
	"github.com/go-pg/pg"
	"github.com/google/uuid"
)

type EquipmentCommand struct {
	DB *pg.DB
}

func (this *EquipmentCommand) Get(uuid uuid.UUID) *Models.EquipmentModel {
	if !this.Exists(uuid) {
		return nil
	}

	var models []Models.EquipmentModel
	err := this.DB.Model(&models).Where("id = ?", uuid).Select()
	if err != nil {
		panic(err)
	}

	return &models[0]
}

func (this *EquipmentCommand) GetAll() *[]Models.EquipmentModel {
	var models []Models.EquipmentModel
	err := this.DB.Model(&models).Select()
	if err != nil {
		panic(err)
	}

	return &models
}

func (this *EquipmentCommand) GetByRecipeUuid(recipeUuid uuid.UUID) *[]Models.EquipmentModel {
	var models []Models.EquipmentModel
	err := this.DB.Model(&models).Where("recipe_id = ?", recipeUuid).Select()
	if err != nil {
		panic(err)
	}

	return &models
}

func (this *EquipmentCommand) Upsert(model *Models.EquipmentModel) *Models.EquipmentModel {
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

func (this *EquipmentCommand) Delete(uuid uuid.UUID) bool {
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

func (this *EquipmentCommand) Exists(uuid uuid.UUID) bool {
	var models []Models.EquipmentModel
	exists, err := this.DB.Model(&models).Where("id = ?", uuid).Exists()
	if err != nil {
		panic(err)
	}

	return exists
}
