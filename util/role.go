package util

import (
	"net/http"

	"github.com/BATUCHKA/real-estate-back/database"
	"github.com/BATUCHKA/real-estate-back/database/models"
)

func RoleMiddleware(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetAdminFromRequestContext(r)
			if user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// if user.IsSuperAdmin {
			// 	next.ServeHTTP(w, r)
			// 	return
			// }

			role, err := getRoleFromDatabase(user)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			authorized := false
			for _, authorizedRole := range roles {
				if role == authorizedRole {
					authorized = true
					break
				}
			}

			if !authorized {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getRoleFromDatabase(user *models.AdminUser) (string, error) {
	db := database.Database.GormDB
	var role models.Role
	err := db.Model(&user).Association("Role").Find(&role)

	if err != nil {
		return "", err
	}
	return string(role.Key), nil
}
