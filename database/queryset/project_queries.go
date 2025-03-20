package queryset

import (
	"github.com/BATUCHKA/real-estate-back/database"
	"gorm.io/gorm"
)

type ProjectIQuerySet interface {
	ProjectListQuery() func(db *gorm.DB) *gorm.DB
	// ProjectQuery(projectID string) func(db *gorm.DB) *gorm.DB
}

type projectQuerySet struct {
	db *database.Postgres
}

var ProjectQuerySet ProjectIQuerySet

func init() {
	ProjectQuerySet = &projectQuerySet{
		db: database.Database,
	}
}

func (qset *projectQuerySet) ProjectListQuery() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(
			"projects p",
		).Select(
			"p.*, if.file_name as main_image_file_name",
		).Joins(
			"LEFT JOIN image_files if ON if.project_id = p.id AND if.image_type = 'main'",
		).Where(
			"p.deleted_at IS NULL",
		)
	}
}

// func (qset *projectQuerySet) ProjectQuery(projectID string) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		return db.Table(
// 			"project p",
// 		).Select(
// 			"p.*, if.file_name as image_file_name",
// 		).Joins(
// 			"LEFT JOIN image_files if ON if.project_id = p.id",
// 		).Where(
// 			"c.deleted_at IS NULL",
// 		)
// 	}
// }
