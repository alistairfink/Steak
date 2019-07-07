package Sort

import (
	"github.com/alistairfink/Steak/Backend/Data/Models"
	"sort"
)

func SortStepsByStepNumber(step *[]Models.StepModel) {
	sortOrder := func(step1, step2 Models.StepModel) bool {
		return step1.StepNumber < step2.StepNumber
	}

	sortSteps(sortOrder).Sort(step)
}

type sortSteps func(step1, step2 Models.StepModel) bool

func (this sortSteps) Sort(step *[]Models.StepModel) {
	stepSorter := &stepSorter{
		step: step,
		by:   this,
	}

	sort.Sort(stepSorter)
}

type stepSorter struct {
	step *[]Models.StepModel
	by   func(step1, step2 Models.StepModel) bool
}

func (this *stepSorter) Len() int {
	return len(*this.step)
}

func (this *stepSorter) Swap(i, j int) {
	(*this.step)[i], (*this.step)[j] = (*this.step)[j], (*this.step)[i]
}

func (this *stepSorter) Less(i, j int) bool {
	return this.by((*this.step)[i], (*this.step)[j])
}
