package search_sorter

import (
	"github.com/lowl11/lazy-collection/array"
	"tok-core/src/data/models"
)

const (
	maxLength = 10

	topUserLength     = 2
	topCategoryLength = 2

	bottomUserLength     = 3
	bottomCategoryLength = 3
)

/*
	1. <user>
	2. <user>
	3. <category>
	4. <category>
	5. <user>
	6. <user>
	7. <user>
	8. <category>
	9. <category>
	10. <category>
*/
func combineLists(userList []models.SearchUserGet, categoryList []models.SearchCategoryGet) []models.SearchItemGet {
	userArray := array.NewWithList[models.SearchUserGet](userList...)
	categoryArray := array.NewWithList[models.SearchCategoryGet](categoryList...)
	list := array.New[models.SearchItemGet]()

	var firstUserCounter int
	for userArray.Size() != 0 && firstUserCounter < topUserLength {
		user := userArray.Get(0)
		userArray.PopForward()

		list.Push(models.SearchItemGet{
			User: &user,
		})
		firstUserCounter++
	}

	var firstCategoryCounter int
	for categoryArray.Size() != 0 && firstCategoryCounter < topCategoryLength {
		category := categoryArray.Get(0)
		categoryArray.PopForward()

		list.Push(models.SearchItemGet{
			Category: &category,
		})
		firstCategoryCounter++
	}

	var bottomUserCounter int
	for userArray.Size() != 0 && bottomUserCounter < bottomUserLength {
		user := userArray.Get(0)
		userArray.PopForward()

		list.Push(models.SearchItemGet{
			User: &user,
		})
		bottomUserCounter++
	}

	var bottomCategoryCounter int
	for categoryArray.Size() != 0 && bottomCategoryCounter < bottomCategoryLength {
		category := categoryArray.Get(0)
		categoryArray.PopForward()

		list.Push(models.SearchItemGet{
			Category: &category,
		})
		bottomCategoryCounter++
	}

	if list.Size() < maxLength && (userArray.Size() > 0 || categoryArray.Size() > 0) {
		difference := maxLength - list.Size()

		var differenceCounter int
		for differenceCounter < difference {
			if userArray.Size() > 0 {
				user := userArray.Get(0)
				userArray.PopForward()

				list.Push(models.SearchItemGet{
					User: &user,
				})
				differenceCounter++
			} else if categoryArray.Size() > 0 {
				category := categoryArray.Get(0)
				categoryArray.PopForward()

				list.Push(models.SearchItemGet{
					Category: &category,
				})
				differenceCounter++
			}
		}
	}

	return list.Slice()
}
