package Models

import (
	"github.com/google/uuid"
)

type PictureModel struct {
	tableName   struct{}  `sql:"picture"`
	Uuid        uuid.UUID `sql:"id, pk"`
	RecipeUuid  uuid.UUID `sql:"recipe_id, fk:recipe.id, notnull"`
	ImageSource string    `sql:"image_source, notnull"`
	SortOrder   int       `sql:"sort_order, notnull"`
}
