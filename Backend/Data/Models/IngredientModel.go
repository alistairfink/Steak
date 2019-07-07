package Models

import (
	"github.com/google/uuid"
)

type IngredientModel struct {
	tableName struct{} `sql:"ingredient"`
	Uuid uuid.UUID `sql:"id, pk"`
	RecipeUuid uuid.UUID `sql:"recipe_id, fk:recipe.id, notnull"`
	Name string `sql:"name, notnull"`
	Quantity int `sql:"quantity, notnull"`
}
