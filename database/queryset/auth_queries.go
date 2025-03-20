package queryset

import (
	"github.com/BATUCHKA/real-estate-back/database"
	"gorm.io/gorm"
)

type AuthIQuerySet interface {
	AuthMeRoles(userID string) func(db *gorm.DB) *gorm.DB
	AuthLoginScope(email string) func(db *gorm.DB) *gorm.DB
}

type authQuerySet struct {
	db *database.Postgres
}

var AuthQuerySet AuthIQuerySet

func init() {
	AuthQuerySet = &authQuerySet{
		db: database.Database,
	}
}


func (qset *authQuerySet) AuthMeScope(userID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(
			"user_content_types uct",
		).Select(
			"uct.id, uct.code, COUNT(DISTINCT uc.content_id)",
		).Joins(
			"LEFT JOIN user_contents uc ON uc.type_id = uct.id",
		).Joins(
			"LEFT JOIN contents c ON c.id = uc.content_id",
		).Where(
			"c.deleted_at IS NULL",
		).Group("uct.id")
	}
}


func (qset *authQuerySet) AuthMeRoles(userID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(
			"user_roles",
		).Select(
			"roles.*",
		).Joins(
			"INNER JOIN roles ON roles.id = user_roles.role_id",
		).Where(
			"user_roles.user_id = ?", userID,
		)
	}
}


func (qset *authQuerySet) AuthLoginScope(email string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(
			"users",
		).Select(
			"users.*",
		).Joins(
			"LEFT JOIN user_roles ON user_roles.user_id = users.id",
		).Joins(
			"LEFT JOIN roles ON roles.id = user_roles.role_id",
		).Where(
			"users.email = ?", email, email,
		).Limit(1)
	}
}
