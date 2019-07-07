package Models

import (
	"github.com/google/uuid"
)

type RecipeModel struct {
	tableName struct{}  `sql:"recipe"`
	Uuid      uuid.UUID `sql:"id, pk"`
	Name      string    `sql:"name, notnull`
	Time      string    `sql:"time, notnull"`
	Type      string    `sql:"type, notnull"`
}
