package util

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"gitlab.com/steppelink/odin/odin-backend/database"
	"gitlab.com/steppelink/odin/odin-backend/database/models"
)

type ConfigKeyValType struct {
	OrganizationName        string `key:"organization_name"`
	OrganizationPhone       string `key:"organization_phone"`
	OrganizationEmail       string `key:"organization_email"`
	OrganizationRegister    string `key:"organization_register"`
	OrganizationWebsite     string `key:"organization_website"`
	OrganizationDescription string `key:"organization_description"`
	OrganizationColor       string `key:"organization_color"`
	OrganizationLogo        string `key:"organization_logo"`
	OrganizationFavIcon     string `key:"organization_fav_icon"`

	DomainURL string `key:"domain_url" description:"web site domain base url"`

	LoadedKeyValDB map[string]interface{}
}

var ConfigKeyVal *ConfigKeyValType

// run the load func and save the current database's data in LoadedKeyValDB
func init() {
	ConfigKeyVal = &ConfigKeyValType{
		LoadedKeyValDB: make(map[string]interface{}),
	}
	// getting database values here
	ConfigKeyVal.Load()
	// saving the first value here
	configValue := reflect.ValueOf(ConfigKeyVal).Elem()
	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Field(i)
		tag := configValue.Type().Field(i).Tag.Get("key")
		if tag != "" {
			ConfigKeyVal.LoadedKeyValDB[tag] = field
		}
	}
}

// Load the database's data
func (c *ConfigKeyValType) Load() {
	var db = database.Database
	var configs []models.Config
	db.GormDB.Find(&configs)
	configValue := reflect.TypeOf(c).Elem()
	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Field(i)
		for _, v := range configs {
			if v.Key == field.Tag.Get("key") {
				fieldValue := reflect.ValueOf(c).Elem().Field(i)
				newV := &v.Value
				// hasVal := http.Request.URL.Query().Has(["name"])
				switch fieldValue.Kind() {
				case reflect.String:
					if fieldValue.Type().String() == "string" {
						fieldValue.SetString(v.Value)
					}
				case reflect.Ptr:
					if fieldValue.Type().String() == "*string" {
						if len(*newV) > 0 {
							fieldValue.Set(reflect.ValueOf(newV))
						}
					}
				case reflect.Int:
					intVal, err := strconv.Atoi(*newV)
					if err != nil {
						intVal = 0
					}
					fieldValue.SetInt(int64(intVal))
				case reflect.Float32:
					float32Val, err := strconv.ParseFloat(v.Value, 32)
					if err != nil {
						float32Val = 0
					}
					fieldValue.SetFloat(float32Val)
				case reflect.Float64:
					float64Val, err := strconv.ParseFloat(v.Value, 64)
					if err != nil {
						float64Val = 0
					}
					fieldValue.SetFloat(float64Val)
				}
				break
			}
		}
	}
}

// got the given value from hand and compare it to the saved current database's data and if there is difference between those two it will be changed.
func (c *ConfigKeyValType) Save() {
	var db = database.Database
	var configs []models.Config
	db.GormDB.Find(&configs)
	configValue := reflect.ValueOf(c).Elem()
	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Field(i)
		fieldValue := reflect.ValueOf(c).Elem().Field(i)
		tag := configValue.Type().Field(i).Tag.Get("key")
		desc := configValue.Type().Field(i).Tag.Get("description")
		if tag != "" {
			if result := db.GormDB.Model(&models.Config{}).Where("key = ?", tag).Update("value", fieldValue); result.RowsAffected != 0 {
				continue
			}

			log.Println("hello", tag)

			config := &models.Config{
				Key:         tag,
				Description: desc,
			}
			if fieldValue.Kind() == reflect.Ptr && fieldValue.Type().String() == "*string" {
				if !field.IsNil() {
					config.Value = *field.Interface().(*string)
				}
			} else {
				config.Value = fmt.Sprint(field)
			}
			db.GormDB.Save(&config)
		}
	}
}
