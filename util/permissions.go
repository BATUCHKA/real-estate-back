package util

import (
	"net/http"
	"reflect"

	"github.com/BATUCHKA/real-estate-back/database"
	"github.com/BATUCHKA/real-estate-back/database/models"
	"gorm.io/gorm/clause"
)

type PermissionKeys struct {
	// User management permissions
	UserEdit   string `key:"user_edit" description:"Edit user information"`
	UserCreate string `key:"user_create" description:"Create new users"`
	UserRead   string `key:"user_read" description:"View user information"`
	UserDelete string `key:"user_delete" description:"Delete users"`

	// Project management permissions (apartment complexes)
	ProjectEdit   string `key:"project_edit" description:"Edit project information"`
	ProjectCreate string `key:"project_create" description:"Create new projects"`
	ProjectRead   string `key:"project_read" description:"View project information"`
	ProjectDelete string `key:"project_delete" description:"Delete projects"`

	// Apartment management permissions
	ApartmentEdit   string `key:"apartment_edit" description:"Edit apartment information"`
	ApartmentCreate string `key:"apartment_create" description:"Create new apartments"`
	ApartmentRead   string `key:"apartment_read" description:"View apartment information"`
	ApartmentDelete string `key:"apartment_delete" description:"Delete apartments"`

	// Purchase management permissions
	PurchaseEdit   string `key:"purchase_edit" description:"Edit purchase information"`
	PurchaseCreate string `key:"purchase_create" description:"Create new purchases"`
	PurchaseRead   string `key:"purchase_read" description:"View purchase information"`
	PurchaseDelete string `key:"purchase_delete" description:"Delete purchases"`

	// Agent management permissions
	AgentEdit   string `key:"agent_edit" description:"Edit agent information"`
	AgentCreate string `key:"agent_create" description:"Create new agents"`
	AgentRead   string `key:"agent_read" description:"View agent information"`
	AgentDelete string `key:"agent_delete" description:"Delete agents"`

	// Client management permissions
	ClientEdit   string `key:"client_edit" description:"Edit client information"`
	ClientCreate string `key:"client_create" description:"Create new clients"`
	ClientRead   string `key:"client_read" description:"View client information"`
	ClientDelete string `key:"client_delete" description:"Delete clients"`

	// Finance management permissions
	FinanceEdit   string `key:"finance_edit" description:"Edit financial information"`
	FinanceCreate string `key:"finance_create" description:"Create financial records"`
	FinanceRead   string `key:"finance_read" description:"View financial information"`
	FinanceDelete string `key:"finance_delete" description:"Delete financial records"`

	// Document management permissions
	DocumentEdit   string `key:"document_edit" description:"Edit documents"`
	DocumentCreate string `key:"document_create" description:"Create new documents"`
	DocumentRead   string `key:"document_read" description:"View documents"`
	DocumentDelete string `key:"document_delete" description:"Delete documents"`

	// Settings permissions
	SettingsEdit string `key:"settings_edit" description:"Edit system settings"`
	SettingsRead string `key:"settings_read" description:"View system settings"`
}

var Permissions *PermissionKeys

func init() {
	Permissions = &PermissionKeys{}
	permissions := reflect.TypeOf(Permissions).Elem()
	for i := 0; i < permissions.NumField(); i++ {
		field := permissions.Field(i)
		fieldValue := reflect.ValueOf(Permissions).Elem().Field(i)
		key := field.Tag.Get("key")
		fieldValue.SetString(key)
	}
}

func (k *PermissionKeys) Flush() {
	db := database.Database
	permissions := reflect.TypeOf(k).Elem()
	for i := 0; i < permissions.NumField(); i++ {
		field := permissions.Field(i)
		key := field.Tag.Get("key")
		description := field.Tag.Get("description")
		permission := &models.Permission{
			Key:         key,
			Description: description,
		}
		db.GormDB.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "key"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"key":         key,
				"description": description,
			}),
		}).Create(&permission)
	}
}

type rolePermission struct {
	models.Permission
	PermissionExist int
}

func (k *PermissionKeys) FlushPreRoles() {
	db := database.Database
	roles := map[string]map[string]any{
		"Admin": {
			"Key":         "admin",
			"Description": "Administrative role with full access",
			"Permissions": "*",
		},
		"Agent": {
			"Key":         "agent",
			"Description": "Real estate agent who purchases properties in bulk",
			"Permissions": []string{
				Permissions.ProjectRead,
				Permissions.ApartmentRead,
				Permissions.PurchaseRead, Permissions.PurchaseCreate, Permissions.PurchaseEdit,
				Permissions.DocumentRead, Permissions.DocumentCreate,
			},
		},
		"Accountant": {
			"Key":         "accountant",
			"Description": "Financial manager with access to transactions and finances",
			"Permissions": []string{
				Permissions.PurchaseRead, Permissions.PurchaseCreate, Permissions.PurchaseEdit,
				Permissions.FinanceRead, Permissions.FinanceCreate, Permissions.FinanceEdit,
				Permissions.DocumentRead,
			},
		},
		"Client": {
			"Key":         "client",
			"Description": "Regular client with minimal access",
			"Permissions": []string{
				Permissions.ProjectRead,
				Permissions.ApartmentRead,
				Permissions.DocumentRead,
				Permissions.PurchaseCreate, Permissions.PurchaseRead,
			},
		},
	}
	for _, v := range roles {
		var role models.Role
		if result := db.GormDB.First(&role, "key = ?", v["Key"].(string)); result.RowsAffected == 0 {
			role = models.Role{
				Key: models.RoleKeyType(v["Key"].(string)),
			}
			db.GormDB.Create(&role)
		}

		if reflect.TypeOf(v["Permissions"]).String() == "string" && v["Permissions"].(string) == "*" {
			var permissions []rolePermission
			db.GormDB.Raw(`
        SELECT t1.id, t1.key, SUM(CASE WHEN t1.role_id IS NULL THEN 0 ELSE 1 END) as permission_exist FROM 
        (
          SELECT id, key, NULL as role_id FROM permissions
          UNION
          SELECT permissions.id, permissions.key, roles.id as role_id FROM roles
          INNER JOIN role_permissions ON roles.id = role_permissions.role_id
          INNER JOIN permissions ON permissions.id = role_permissions.permission_id
          WHERE role_id = ?
        ) t1 GROUP BY t1.id, t1.key
      `, role.ID).Scan(&permissions)
			for _, p := range permissions {
				if p.PermissionExist == 0 {
					role_permission := &models.RolePermission{
						RoleID:       role.ID.String(),
						PermissionID: p.ID.String(),
					}
					db.GormDB.Create(&role_permission)
				}
			}
		} else {
			permission_keys := v["Permissions"].([]string)
			var permissions []rolePermission
			db.GormDB.Raw(`
        SELECT t1.id, t1.key, SUM(CASE WHEN t1.role_id IS NULL THEN 0 ELSE 1 END) as permission_exist FROM 
        (
          SELECT id, key, NULL as role_id FROM permissions
          UNION
          SELECT permissions.id, permissions.key, roles.id as role_id FROM roles
          INNER JOIN role_permissions ON roles.id = role_permissions.role_id
          INNER JOIN permissions ON permissions.id = role_permissions.permission_id
          WHERE role_id = ?
        ) t1 WHERE t1.key IN ? GROUP BY t1.id, t1.key
      `, role.ID, permission_keys).Scan(&permissions)
			for _, p := range permissions {
				if p.PermissionExist == 0 {
					role_permission := &models.RolePermission{
						RoleID:       role.ID.String(),
						PermissionID: p.ID.String(),
					}
					db.GormDB.Create(&role_permission)
				}
			}
		}
	}
}

func PermissionMiddleware(permission_keys ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromRequestContext(r)
			if user == nil {
				JsonErrorResponse("Unauthorized").WithErrorCode(401).Write(w)
				return
			}

			db := database.Database

			// Check if user has admin role (directly via RoleID)
			var role models.Role
			if err := db.GormDB.First(&role, "id = ?", user.RoleID).Error; err == nil {
				// If user has admin role, allow all actions
				if role.Key == "admin" {
					next.ServeHTTP(w, r)
					return
				}
			}

			var permissions []models.Permission
			if result := db.GormDB.Raw(`
        SELECT p.* FROM users 
        INNER JOIN roles r ON r.id = users.role_id
        INNER JOIN role_permissions rp ON rp.role_id = r.id
        INNER JOIN permissions p ON p.id = rp.permission_id 
        WHERE users.deleted_at IS NULL AND users.id = ? AND p.key IN ?
      `, user.ID, permission_keys).Scan(&permissions); result.Error != nil {
				JsonErrorResponse("Failed to check permission").WithErrorCode(500).Write(w)
				return
			} else {
				if len(permissions) >= len(permission_keys) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Check user-specific permissions if role permissions are insufficient
			var userPermissions []models.Permission
			if result := db.GormDB.Raw(`
        SELECT p.* FROM users 
        INNER JOIN user_permissions up ON up.user_id = users.id
        INNER JOIN permissions p ON p.id = up.permission_id 
        WHERE user_permissions.deleted_at IS NULL AND users.id = ? AND p.key IN ?
      `, user.ID, permission_keys).Scan(&userPermissions); result.Error != nil {
				JsonErrorResponse("Failed to check permission").WithErrorCode(500).Write(w)
				return
			} else {
				if len(userPermissions) >= len(permission_keys) {
					next.ServeHTTP(w, r)
					return
				}
			}

			JsonErrorResponse("Permission denied").WithErrorCode(403).Write(w)
		})
	}
}
