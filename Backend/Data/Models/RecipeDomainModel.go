package Models

import (
	"github.com/google/uuid"
)

type RecipeDomainModel struct {
	Uuid        uuid.UUID
	Time        string
	Type        string
	Pictures    *[]PictureModel
	Equipments  *[]EquipmentModel
	Ingredients *[]IngredientModel
	Steps       *[]StepModel
}

func (this *RecipeDomainModel) ToRecipeModel() *RecipeModel {
	return &RecipeModel{
		Uuid: this.Uuid,
		Time: this.Time,
		Type: this.Type,
	}
}
