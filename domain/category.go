package domain

func GetAllCategories() []Category {
	categoryNames := [...]string{
		"japanese",
		"geography",
		"history",
		"civics",
		"mathematics",
		"scientific",
		"healthAndPE",
		"art",
		"english",
		"homeEconomics",
		"information",
	}

	var categories []Category
	for i, name := range categoryNames {
		category := Category{
			ID:   int64(i) + 1, // ID starts from 1.
			Name: name,
		}
		categories = append(categories, category)
	}

	categories = append(categories, Category{ID: 999, Name: "other"})

	return categories
}
