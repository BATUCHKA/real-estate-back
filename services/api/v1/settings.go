package v1

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/BATUCHKA/real-estate-back/util"
)

type infoBody struct {
	Name        *string `validate:"required" json:"name"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	Register    *string `json:"register"`
	Website     *string `json:"website"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
	Logo        *string `json:"logo"`
	FavIcon     *string `json:"fav_icon"`
}

func SettingsInfoPut(w http.ResponseWriter, r *http.Request) {
	var body infoBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.JsonErrorResponse(err.Error()).Write(w)
		return
	}
	if err := util.NewValidator().Validate(body); err != nil {
		util.JsonErrorResponse(err.Error()).Write(w)
		return
	}
	if body.Name != nil {
		util.ConfigKeyVal.OrganizationName = *body.Name
	}

	if body.Phone != nil {
		util.ConfigKeyVal.OrganizationPhone = *body.Phone
	}

	if body.Email != nil {
		util.ConfigKeyVal.OrganizationEmail = *body.Email
	}

	if body.Register != nil {
		util.ConfigKeyVal.OrganizationRegister = *body.Register
	}

	if body.Website != nil {
		util.ConfigKeyVal.OrganizationWebsite = *body.Website
	}

	if body.Description != nil {
		util.ConfigKeyVal.OrganizationDescription = *body.Description
	}

	if body.Color != nil {
		util.ConfigKeyVal.OrganizationColor = *body.Color
	}

	if body.Logo != nil {
		util.ConfigKeyVal.OrganizationLogo = *body.Logo
	}

	if body.FavIcon != nil {
		util.ConfigKeyVal.OrganizationFavIcon = *body.FavIcon
	}

	util.ConfigKeyVal.Save()
	util.JsonResponse(body).Write(w)
}

type infoGetBody struct {
	Name        string `validate:"required" json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Register    string `json:"register"`
	Website     string `json:"website"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Logo        string `json:"logo"`
	FavIcon     string `json:"fav_icon"`
}

func SettingsInfoGet(w http.ResponseWriter, r *http.Request) {
	body := &infoGetBody{
		Name:        util.ConfigKeyVal.OrganizationName,
		Phone:       util.ConfigKeyVal.OrganizationPhone,
		Email:       util.ConfigKeyVal.OrganizationEmail,
		Register:    util.ConfigKeyVal.OrganizationRegister,
		Website:     util.ConfigKeyVal.OrganizationWebsite,
		Description: util.ConfigKeyVal.OrganizationDescription,
		Color:       util.ConfigKeyVal.OrganizationColor,
		Logo:        os.Getenv("MEDIA_URL") + util.ConfigKeyVal.OrganizationLogo,
		FavIcon:     os.Getenv("MEDIA_URL") + util.ConfigKeyVal.OrganizationFavIcon,
	}
	util.JsonResponse(body).Write(w)
}
