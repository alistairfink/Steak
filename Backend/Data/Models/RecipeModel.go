package Models

import (
	"github.com/google/uuid"
)

type RecipeModel struct {
	tablename struct{}  `sql:"recipe"`
	Uuid      uuid.UUID `sql:"id, pk"`
	Time      string    `sql:"time, notnull"`
	Type      string    `sql:"type, notnull"`
}
