package search_sorter

import "tok-core/src/data/models"

/*
	SortSmart сортировка и комбинирование разных списков и объединение в одно
*/
func SortSmart(userList []models.SearchUserGet, categoryList []models.SearchCategoryGet) []models.SearchItemGet {
	return combineLists(userList, categoryList)
}
