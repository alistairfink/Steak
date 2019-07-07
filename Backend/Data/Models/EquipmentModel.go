package Models

import (
	"github.com/google/uuid"
)

type EquipmentModel struct {
	tableName  struct{}  `sql:"equipment"`
	Uuid       uuid.UUID `sql:"id, pk"`
	RecipeUuid uuid.UUID `sql:"recipe_id, fk:recipe.id, notnull"`
	Name       string    `sql:"name, notnull"`
	Quantity   int       `sql:"quantity, notnull"`
}
