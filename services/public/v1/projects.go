package v1

import (
	"net/http"
	"os"

	"github.com/BATUCHKA/real-estate-back/database"
	"github.com/BATUCHKA/real-estate-back/database/models"
	"github.com/BATUCHKA/real-estate-back/database/queryset"
	"github.com/BATUCHKA/real-estate-back/util"
	"github.com/go-chi/chi/v5"
)

type apartmentDetail struct {
	models.Apartment
	Rooms      []models.Room `json:"rooms"`
	ImageNames []string      `json:"image_names"`
	ImageUrls  []string      `json:"image_urls"`
}

type projectDetail struct {
	models.Project
	MainImageName   *string           `json:"main_image_name"`
	MainImageUrl    *string           `json:"main_image_url"`
	OtherImageNames []string          `json:"other_image_names"`
	OtherImageUrls  []string          `json:"other_image_urls"`
	ApartmentInfo   []apartmentDetail `json:"apartment_info"`
}

func ProjectGetByID(w http.ResponseWriter, r *http.Request) {
	db := database.Database
	projectID := chi.URLParam(r, "id")

	var project projectDetail
	if result := db.GormDB.First(&project, "id = ?", projectID); result.Error != nil {
		util.JsonErrorResponse("project not found").Write(w)
	}

	if project.MainImageName != nil {
		url := os.Getenv("MEDIA_URL") + *project.MainImageName
		project.MainImageUrl = &url
	}

	if project.OtherImageNames != nil {
		for i, v := range project.OtherImageNames {
			project.OtherImageUrls[i] = os.Getenv("MEDIA_URL") + v
		}
	}
}

type projectListFilter struct {
	Title       *string `filter:"title:title;field:title"`
	Description *string `filter:"description:description;field:description"`
	OrderBy     *string `order_by:"title:title"`
}

type projectListBody struct {
	models.Project
	MainImageFileName *string `json:"main_image_file_name"`
	MainImageFileURL  *string `json:"main_image_file_url"`
}

func ProjectListGet(w http.ResponseWriter, r *http.Request) {
	r_model := util.NewRequestModel(r)
	var filter projectListFilter
	r_model.ParseFilter(&filter)
	var projects []projectListBody
	db := database.Database

	query := db.GormDB.Scopes(queryset.ProjectQuerySet.ProjectListQuery())
	if result := query.Find(&projects); result.Error != nil {
		util.JsonErrorResponse("project list get error").Write(w)
		return
	}

	for i, v := range projects {
		if *projects[i].MainImageFileName != "" {
			url := os.Getenv("MEDIA_URL") + *v.MainImageFileName
			projects[i].MainImageFileURL = &url
		}
	}
	util.JsonResponse(projects).WithPagingScopeDistinct(query, r_model, "projects.id").Write((w))
}
