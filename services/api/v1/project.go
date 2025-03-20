package v1

import (
	"encoding/json"
	"net/http"

	"github.com/BATUCHKA/real-estate-back/database"
	"github.com/BATUCHKA/real-estate-back/database/models"
	"github.com/BATUCHKA/real-estate-back/util"
	"github.com/go-chi/chi/v5"
)

type projectCreateBody struct {
	Title              string        `json:"title"`
	Description        string        `json:"description"`
	TotalApartments    int           `json:"total_apartments"`
	RemaningApartments *int          `json:"remaining_apartments"`
	SoldApartments     *int          `json:"sold_apartments"`
	AdvantagesHTML     *string       `json:"advantages_html" gorm:"type:text"`
	Longitude          *float64      `json:"longitude"`
	Latitude           *float64      `json:"latitude"`
	StartedAt          *util.Iso8601 `json:"started_at"`
	EndedAt            *util.Iso8601 `json:"ended_at"`
	// PlanImageIDs       []string      `json:"plan_image_ids"`
	// MainImageID        *string       `json:"main_image_id"`
}

func ProjectCreate(w http.ResponseWriter, r *http.Request) {
	var body projectCreateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.JsonErrorResponse(err.Error()).Write(w)
		return
	}
	if err := util.NewValidator().Validate(body); err != nil {
		util.JsonErrorResponse(err.Error()).Write(w)
		return
	}

	db := database.Database

	var newProject = &models.Project{
		Title:           body.Title,
		Description:     body.Description,
		TotalApartments: body.TotalApartments,
	}

	if body.RemaningApartments != nil {
		newProject.RemaningApartments = *body.RemaningApartments
	}
	if body.SoldApartments != nil {
		newProject.SoldApartments = *body.SoldApartments
	}
	if body.AdvantagesHTML != nil {
		newProject.AdvantagesHTML = *body.AdvantagesHTML
	}
	if body.StartedAt != nil {
		newProject.StartedAt = body.StartedAt.GetTime()
	}
	if body.EndedAt != nil {
		newProject.EndedAt = body.EndedAt.GetTime()
	}

	if body.RemaningApartments != nil {
		newProject.RemaningApartments = *body.RemaningApartments
	}
	if body.RemaningApartments != nil {
		newProject.RemaningApartments = *body.RemaningApartments
	}
	if body.RemaningApartments != nil {
		newProject.RemaningApartments = *body.RemaningApartments
	}

	if err := db.GormDB.Save(&newProject).Error; err != nil {
		util.JsonErrorResponse("new project create error").Write(w)
		return
	}

	util.JsonResponse(&newProject).Write(w)
}

type projectUpdateBody struct {
	Title              *string       `json:"title"`
	Description        *string       `json:"description"`
	TotalApartments    *int          `json:"total_apartments"`
	RemaningApartments *int          `json:"remaining_apartments"`
	SoldApartments     *int          `json:"sold_apartments"`
	AdvantagesHTML     *string       `json:"advantages_html" gorm:"type:text"`
	Longitude          *float64      `json:"longitude"`
	Latitude           *float64      `json:"latitude"`
	StartedAt          *util.Iso8601 `json:"started_at"`
	EndedAt            *util.Iso8601 `json:"ended_at"`
}

func ProjectUpdateByID(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	var body projectUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.JsonErrorResponse(err.Error()).Write(w)
		return
	}

	db := database.Database

	var project models.Project
	if result := db.GormDB.First(&project, "id = ?", projectID); result.Error != nil {
		util.JsonErrorResponse("Project not found").WithErrorCode(http.StatusNotFound).Write(w)
		return
	}

	if body.Title != nil {
		project.Title = *body.Title
	}
	if body.Description != nil {
		project.Description = *body.Description
	}
	if body.TotalApartments != nil {
		project.TotalApartments = *body.TotalApartments
	}
	if body.RemaningApartments != nil {
		project.RemaningApartments = *body.RemaningApartments
	}
	if body.SoldApartments != nil {
		project.SoldApartments = *body.SoldApartments
	}
	if body.AdvantagesHTML != nil {
		project.AdvantagesHTML = *body.AdvantagesHTML
	}
	if body.Longitude != nil {
		project.Longitude = body.Longitude
	}
	if body.Latitude != nil {
		project.Latitude = body.Latitude
	}
	if body.StartedAt != nil {
		project.StartedAt = body.StartedAt.GetTime()
	}
	if body.EndedAt != nil {
		project.EndedAt = body.EndedAt.GetTime()
	}

	if err := db.GormDB.Save(&project); err.Error != nil {
		util.JsonErrorResponse("Failed to update project").Write(w)
		return
	}

	util.JsonResponse(&project).Write(w)
}

func ProjectDeleteByID(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	db := database.Database

	// Check if project exists
	var project models.Project
	if result := db.GormDB.First(&project, "id = ?", projectID); result.Error != nil {
		util.JsonErrorResponse("Project not found").Write(w)
		return
	}

	// Begin transaction
	tx := db.GormDB.Begin()

	if err := tx.Where("project_id = ?", projectID).Delete(&models.ImageFile{}).Error; err != nil {
		tx.Rollback()
		util.JsonErrorResponse("Failed to delete project images").WithErrorCode(http.StatusInternalServerError).Write(w)
		return
	}

	var apartmentIDs []string
	tx.Model(&models.Apartment{}).Where("project_id = ?", projectID).Pluck("id", &apartmentIDs)

	// Delete apartments
	if len(apartmentIDs) > 0 {
		if err := tx.Where("project_id = ?", projectID).Delete(&models.Apartment{}).Error; err != nil {
			tx.Rollback()
			util.JsonErrorResponse("Failed to delete project apartments").WithErrorCode(http.StatusInternalServerError).Write(w)
			return
		}

		// 3. Delete any rooms associated with these apartments
		if err := tx.Where("apartment_id IN ?", apartmentIDs).Delete(&models.Room{}).Error; err != nil {
			tx.Rollback()
			util.JsonErrorResponse("Failed to delete apartment rooms").WithErrorCode(http.StatusInternalServerError).Write(w)
			return
		}
	}

	// 4. Finally delete the project
	if err := tx.Delete(&project).Error; err != nil {
		tx.Rollback()
		util.JsonErrorResponse("Failed to delete project").WithErrorCode(http.StatusInternalServerError).Write(w)
		return
	}

	// Commit the transaction
	tx.Commit()

	util.JsonResponse(map[string]interface{}{
		"status":  "success",
		"message": "Project deleted successfully",
	}).Write(w)
}
