package Sort

import (
	"github.com/alistairfink/Steak/Backend/Data/Models"
	"sort"
)

func SortPicturesBySortOrder(picture *[]Models.PictureModel) {
	sortOrder := func(pic1, pic2 Models.PictureModel) bool {
		return pic1.SortOrder < pic2.SortOrder
	}

	sortPictures(sortOrder).Sort(picture)
}

type sortPictures func(pic1, pic2 Models.PictureModel) bool

func (this sortPictures) Sort(picture *[]Models.PictureModel) {
	pictureSorter := &pictureSorter{
		picture: picture,
		by:      this,
	}

	sort.Sort(pictureSorter)
}

type pictureSorter struct {
	picture *[]Models.PictureModel
	by      func(pic1, pic2 Models.PictureModel) bool
}

func (this *pictureSorter) Len() int {
	return len(*this.picture)
}

func (this *pictureSorter) Swap(i, j int) {
	(*this.picture)[i], (*this.picture)[j] = (*this.picture)[j], (*this.picture)[i]
}

func (this *pictureSorter) Less(i, j int) bool {
	return this.by((*this.picture)[i], (*this.picture)[j])
}
