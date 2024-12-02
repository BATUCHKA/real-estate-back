package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gitlab.com/steppelink/odin/odin-backend/database"
	"gitlab.com/steppelink/odin/odin-backend/database/models"
)

const (
	SessionKey     = "session-key"
	SessionUserKey = "session-user-data"
)

func CreateUserSession(user *models.User, expireAt time.Time) (*models.UserSession, error) {
	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	session_id := uuid.New()
	sig := hmac.New(sha256.New, []byte(session_id.String()))
	sig.Write(data)
	session := models.UserSession{
		Key:      session_id,
		Data:     base64.StdEncoding.EncodeToString(data),
		Hash:     hex.EncodeToString(sig.Sum(nil)),
		ExpireAt: expireAt,
	}
	db := database.Database
	result := db.GormDB.Create(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func ReleaseUserSession(session *models.UserSession) error {
	db := database.Database
	result := db.GormDB.Delete(&session)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func ParseUserSession(key string) (*models.User, *models.UserSession, error) {
	var session models.UserSession

	db := database.Database
	result := db.GormDB.First(&session, "hash = ? AND expire_at >= ?", key, time.Now())
	if result.Error != nil {
		return nil, nil, result.Error
	}
	data, err := base64.StdEncoding.DecodeString(session.Data)
	if result.Error != nil {
		log.Fatal(err)
	}
	var user models.AdminUser
	var userInfo models.User

	json.Unmarshal(data, &user)
	result = db.GormDB.First(&userInfo, "id = ?", user.ID)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	return &userInfo, &session, nil
}

func GetUserFromRequestContext(r *http.Request) *models.User {
	if val := r.Context().Value(SessionUserKey); val != nil {
		return r.Context().Value(SessionUserKey).(*models.User)
	}
	return nil
}

func GetUserSessionFromRequestContext(r *http.Request) *models.UserSession {
	session := r.Context().Value(SessionKey).(*models.UserSession)
	return session
}
func CreateAdminSession(user *models.AdminUser, expireAt time.Time) (*models.Session, error) {
	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	session_id := uuid.New()
	sig := hmac.New(sha256.New, []byte(session_id.String()))
	sig.Write(data)
	session := models.Session{
		Key:      session_id,
		Data:     base64.StdEncoding.EncodeToString(data),
		Hash:     hex.EncodeToString(sig.Sum(nil)),
		ExpireAt: expireAt,
		// LastActivityTime: time.Now(),
	}
	db := database.Database
	result := db.GormDB.Create(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func ReleaseAdminSession(session *models.Session) error {
	db := database.Database
	result := db.GormDB.Delete(&session)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func ParseAdminSession(key string) (*models.AdminUser, *models.Session, error) {
	var session models.Session

	db := database.Database
	result := db.GormDB.First(&session, "hash = ? AND expire_at >= ?", key, time.Now())
	if result.Error != nil {
		return nil, nil, result.Error
	}
	data, err := base64.StdEncoding.DecodeString(session.Data)
	if result.Error != nil {
		log.Fatal(err)
	}
	var user models.AdminUser
	var userInfo models.AdminUser
	json.Unmarshal(data, &user)
	result = db.GormDB.First(&userInfo, "id = ?", user.ID)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	return &userInfo, &session, nil
}

func GetAdminFromRequestContext(r *http.Request) *models.AdminUser {
	if val := r.Context().Value(SessionUserKey); val != nil {
		return r.Context().Value(SessionUserKey).(*models.AdminUser)
	}
	return nil
}

func GetAdminSessionFromRequestContext(r *http.Request) *models.Session {
	session := r.Context().Value(SessionKey).(*models.Session)
	return session
}

// func UpdateLastActivityTime(*models.Session) error {

// 	session := models.Session{}

// 	session.LastActivityTime = time.Now()

// 	db := database.Database
// 	return db.GormDB.Save(&session).Error
// }

// func SessionTimeoutMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		session := GetSessionFromRequestContext(r)

// 		if session != nil {
// 			// Check for session timeout based on last activity time
// 			if session.LastActivityTime.Add(20 * time.Minute).Before(time.Now()) {
// 				// Session has timed out
// 				if err := ReleaseUserSession(session); err != nil {
// 					JsonErrorResponse(err.Error()).Write(w)
// 					return
// 				}

// 				// Redirect to login page
// 				// http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
// 				JsonErrorResponse("Session has timed out").Write(w)
// 				return
// 			} else {
// 				// Session is still active, update last activity time
// 				session.LastActivityTime = time.Now()
// 				if err := UpdateLastActivityTime(session); err != nil {
// 					JsonErrorResponse(err.Error()).Write(w)
// 					return
// 				}
// 			}
// 		}

// 		// Proceed to the actual handler
// 		next.ServeHTTP(w, r)
// 	})
// }
