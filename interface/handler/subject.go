package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/super-dog-human/teraconnectgo/domain"
	"github.com/super-dog-human/teraconnectgo/usecase"
)

type getSubjectResponse struct {
	Subject    domain.Subject         `json:"subject"`
	Categories []domain.ShortCategory `json:"categories"`
}

func getSubjects(c echo.Context) error {
	subjects, err := usecase.GetSubjects(c.Request())
	if err != nil {
		fatalLog(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	withCategories := c.QueryParam("with_categories") == "true"
	if withCategories {
		categories, err := usecase.GetAllCategories(c.Request())
		if err != nil {
			fatalLog(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		var preSubjectID int64
		categoryIDs := make(map[int64][]domain.ShortCategory)
		for _, category := range categories {
			if preSubjectID == 0 {
				preSubjectID = category.SubjectID
			}
			categoryIDs[category.SubjectID] = append(categoryIDs[category.SubjectID], category)
		}

		var responses []getSubjectResponse
		for _, subject := range subjects {
			response := getSubjectResponse{Subject: subject, Categories: categoryIDs[subject.ID]}
			responses = append(responses, response)
		}
		return c.JSON(http.StatusOK, responses)
	} else {
		return c.JSON(http.StatusOK, subjects)
	}
}
