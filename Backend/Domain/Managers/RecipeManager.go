package Managers

import (
	"github.com/alistairfink/Steak/Backend/Data/Commands"
	"github.com/alistairfink/Steak/Backend/Data/Models"
	"github.com/go-pg/pg"
	"github.com/google/uuid"
)

type RecipeManger struct {
	recipeCommand     *Commands.RecipeCommand
	ingredientCommand *Commands.IngredientCommand
	pictureCommand    *Commands.PictureCommand
	equipmentCommand  *Commands.EquipmentCommand
	stepCommand       *Commands.StepCommand
}

func (this *RecipeManger) Init(db *pg.DB) {
	this.recipeCommand = &Commands.RecipeCommand{DB: db}
	this.ingredientCommand = &Commands.IngredientCommand{DB: db}
	this.pictureCommand = &Commands.PictureCommand{DB: db}
	this.equipmentCommand = &Commands.EquipmentCommand{DB: db}
	this.stepCommand = &Commands.StepCommand{DB: db}
}

func (this *RecipeManger) Get(uuid uuid.UUID) *Models.RecipeDomainModel {
	panic("Not Implemented")
}

func (this *RecipeManger) GetAll() *[]Models.RecipeModel {
	recipeModels := this.recipeCommand.GetAll()
	return recipeModels
}

func (this *RecipeManger) Update(model *Models.RecipeDomainModel) *Models.RecipeDomainModel {
	if !this.recipeCommand.Exists(model.Uuid) {
		return nil
	}

	existingIngredients := this.ingredientCommand.GetByRecipeUuid(model.Uuid)
	keepIngredients := make(map[uuid.UUID]bool)
	for i, ingredient := range *model.Ingredients {
		keepIngredients[ingredient.Uuid] = true
		(*model.Ingredients)[i].RecipeUuid = model.Uuid
		if ingredient.Uuid != uuid.Nil && !this.ingredientCommand.Exists(ingredient.Uuid) {
			return nil
		}
	}

	existingPictures := this.pictureCommand.GetByRecipeUuid(model.Uuid)
	keepPictures := make(map[uuid.UUID]bool)
	for i, picture := range *model.Pictures {
		keepPictures[picture.Uuid] = true
		(*model.Pictures)[i].RecipeUuid = model.Uuid
		if picture.Uuid != uuid.Nil && !this.pictureCommand.Exists(picture.Uuid) {
			return nil
		}
	}

	existingEquipment := this.equipmentCommand.GetByRecipeUuid(model.Uuid)
	keepEquipment := make(map[uuid.UUID]bool)
	for i, equipment := range *model.Equipments {
		keepEquipment[equipment.Uuid] = true
		(*model.Equipments)[i].RecipeUuid = model.Uuid
		if equipment.Uuid != uuid.Nil && !this.equipmentCommand.Exists(equipment.Uuid) {
			return nil
		}
	}

	existingStep := this.stepCommand.GetByRecipeUuid(model.Uuid)
	keepStep := make(map[uuid.UUID]bool)
	for i, step := range *model.Steps {
		keepStep[step.Uuid] = true
		(*model.Steps)[i].RecipeUuid = model.Uuid
		if step.Uuid != uuid.Nil && !this.stepCommand.Exists(step.Uuid) {
			return nil
		}
	}

	for _, ingredient := range *existingIngredients {
		if !keepIngredients[ingredient.Uuid] {
			this.ingredientCommand.Delete(ingredient.Uuid)
		}
	}

	for _, picture := range *existingPictures {
		if !keepPictures[picture.Uuid] {
			this.pictureCommand.Delete(picture.Uuid)
		}
	}

	for _, equipment := range *existingEquipment {
		if !keepEquipment[equipment.Uuid] {
			this.equipmentCommand.Delete(equipment.Uuid)
		}
	}

	for _, step := range *existingStep {
		if !keepStep[step.Uuid] {
			this.stepCommand.Delete(step.Uuid)
		}
	}

	return this.Create(model)
}

func (this *RecipeManger) Create(model *Models.RecipeDomainModel) *Models.RecipeDomainModel {
	preUpsertModel := model.ToRecipeModel()
	postUpsertModel := this.recipeCommand.Upsert(preUpsertModel)

	for _, ingredient := range *model.Ingredients {
		ingredient.RecipeUuid = postUpsertModel.Uuid
		this.ingredientCommand.Upsert(&ingredient)
	}

	for _, picture := range *model.Pictures {
		picture.RecipeUuid = postUpsertModel.Uuid
		this.pictureCommand.Upsert(&picture)
	}

	for _, equipment := range *model.Equipments {
		equipment.RecipeUuid = postUpsertModel.Uuid
		this.equipmentCommand.Upsert(&equipment)
	}

	for _, step := range *model.Steps {
		step.RecipeUuid = postUpsertModel.Uuid
		this.stepCommand.Upsert(&step)
	}

	return this.Get(postUpsertModel.Uuid)
}

func (this *RecipeManger) Delete(uuid uuid.UUID) bool {
	if !this.recipeCommand.Exists(uuid) {
		return false
	}

	ingredients := this.ingredientCommand.GetByRecipeUuid(uuid)
	pictures := this.pictureCommand.GetByRecipeUuid(uuid)
	equipments := this.equipmentCommand.GetByRecipeUuid(uuid)
	steps := this.stepCommand.GetByRecipeUuid(uuid)

	for _, ingredient := range *ingredients {
		this.ingredientCommand.Delete(ingredient.Uuid)
	}

	for _, picture := range *pictures {
		this.pictureCommand.Delete(picture.Uuid)
	}

	for _, equipment := range *equipments {
		this.equipmentCommand.Delete(equipment.Uuid)
	}

	for _, step := range *steps {
		this.stepCommand.Delete(step.Uuid)
	}

	return this.recipeCommand.Delete(uuid)
}
