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

	"github.com/BATUCHKA/real-estate-back/database"
	"github.com/BATUCHKA/real-estate-back/database/models"
	"github.com/google/uuid"
)

const (
	SessionKey     = "session-key"
	SessionUserKey = "session-user-data"
)

func CreateUserSession(user *models.User, expireAt time.Time) (*models.Session, error) {
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
	}
	db := database.Database
	result := db.GormDB.Create(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func ReleaseUserSession(session *models.Session) error {
	db := database.Database
	result := db.GormDB.Delete(&session)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func ParseUserSession(key string) (*models.User, *models.Session, *[]models.Role, error) {
	var session models.Session

	db := database.Database
	result := db.GormDB.First(&session, "hash = ? AND expire_at >= ?", key, time.Now())
	if result.Error != nil {
		return nil, nil, nil, result.Error
	}
	data, err := base64.StdEncoding.DecodeString(session.Data)
	if result.Error != nil {
		log.Fatal(err)
	}
	var user models.User
	var userInfo models.User
	json.Unmarshal(data, &user)
	result = db.GormDB.First(&userInfo, "id = ?", user.ID)
	if result.Error != nil {
		return nil, nil, nil, result.Error
	}
	var roles []models.Role
	db.GormDB.Find(&roles, "id = ?", userInfo.RoleID)

	return &userInfo, &session, &roles, nil
}

func GetUserFromRequestContext(r *http.Request) *models.User {
	if val := r.Context().Value(SessionUserKey); val != nil {
		return r.Context().Value(SessionUserKey).(*models.User)
	}
	return nil
}

func GetSessionFromRequestContext(r *http.Request) *models.Session {
	if val := r.Context().Value(SessionKey); val != nil {
		return r.Context().Value(SessionKey).(*models.Session)
	}
	return nil
}

// func ParseSession(key string) (*models.User, *models.Session, error) {
// 	var session models.Session

// 	db := database.Database
// 	result := db.GormDB.First(&session, "hash = ? AND expire_at >= ?", key, time.Now())
// 	if result.Error != nil {
// 		return nil, nil, result.Error
// 	}
// 	data, err := base64.StdEncoding.DecodeString(session.Data)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var userData models.User
// 	err = json.Unmarshal(data, &userData)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	var user models.User
// 	result = db.GormDB.First(&user, "id = ?", userData.ID)
// 	if result.Error != nil {
// 		return nil, nil, result.Error
// 	}

// 	return &user, &session, nil
// }
