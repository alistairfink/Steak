package Models

import (
	"github.com/google/uuid"
)

type StepModel struct {
	tableName struct{} `sql:"step"`
	Uuid uuid.UUID `sql:"id, pk"`
	RecipeUuid uuid.UUID `sql:"recipe_id, fk:recipe.id, notnull"`
	Content string `sql:"content, notnull"`
	StepNumber int `sql:"step_number, notnull`
}