package util

import (
	"errors"
	"fmt"
	"github.com/BATUCHKA/real-estate-back/database"
	"github.com/BATUCHKA/real-estate-back/database/models"
	"os"
)

func GetSurveyShortLink(surveyID string) (*string, error) {
	db := database.Database

	var surveyLink models.SurveyLink
	if result := db.GormDB.First(&surveyLink, "survey_id = ?", surveyID); result.Error != nil && result.RowsAffected == 0 {
		return nil, errors.New("survey link not found")
	}

	baseURL := os.Getenv("SURVEY_BASE_URL")
	surveyLinkURL := fmt.Sprintf("%s/surveys/%s", baseURL, surveyLink.ShortKey)
	return &surveyLinkURL, nil
}
