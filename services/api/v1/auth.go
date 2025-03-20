package v1

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/BATUCHKA/real-estate-back/database"
	"github.com/BATUCHKA/real-estate-back/database/models"
	"github.com/BATUCHKA/real-estate-back/database/queryset"
	"github.com/BATUCHKA/real-estate-back/util"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authLoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authLoginResponse struct {
	BearerToken string        `json:"bearer_token"`
	ID          uuid.UUID     `json:"id"`
	Email       string        `json:"email"`
	FirstName   string        `json:"first_name"`
	LastName    string        `json:"last_name"`
	PhoneNumber string        `json:"phone_number"`
	Roles       []models.Role `json:"roles"`
}

func AuthLogin(w http.ResponseWriter, r *http.Request) {
	var body authLoginBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		util.JsonErrorResponse("Json body is not valid").Write(w)
		return
	}

	db := database.Database
	body.Username = strings.ToLower(body.Username)
	var user models.User
	result := db.GormDB.First(&user, "(username = ? OR email = ?)", body.Username, body.Username)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if result.Error != nil || err != nil {
		util.JsonErrorResponse("user is invalid").Write(w)
		return
	}

	session, err := util.CreateUserSession(&user, time.Now().AddDate(0, 0, 7))
	if err != nil {
		util.JsonErrorResponse("failed to create user session").Write(w)
		return
	}
	resBody := authLoginResponse{
		BearerToken: session.Hash,
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
	}

	var roles []models.Role
	roleQuery := db.GormDB.Scopes(queryset.AuthQuerySet.AuthMeRoles(user.ID.String()))
	if result := roleQuery.Find(&roles); result.Error != nil {
		util.JsonErrorResponse("user role not found").Write(w)
		return
	}
	resBody.Roles = roles

	util.JsonResponse(resBody).Write(w)
}

type authSignUpBody struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	RoleID   string `json:"role_id"`
}

func AuthSignUp(w http.ResponseWriter, r *http.Request) {
	var body authSignUpBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.JsonErrorResponse(err.Error()).Write(w)
		return
	}
	if err := util.NewValidator().Validate(body); err != nil {
		util.JsonErrorResponse(err.Error()).Write(w)
		return
	}
	db := database.Database
	var user *models.User

	body.Email = strings.ToLower(body.Email)
	body.Username = strings.ToLower(body.Username)

	if len(body.Email) > 0 && len(body.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			util.JsonErrorResponse(err.Error()).Write(w)
			return
		}

		user = &models.User{
			Email:    body.Email,
			Password: string(hashedPassword),
		}

		var role *models.Role
		if result := db.GormDB.First(&role, "id = ?", body.RoleID); result.Error != nil {
			util.JsonErrorResponse("role not found").Write(w)
			return
		} else if role != nil {
			user.Role = *role
		}

		if result := db.GormDB.Create(&user); result.Error != nil {
			util.JsonErrorResponse("user create error").Write(w)
			return
		}

		util.JsonResponse(&user).Write(w)
	} else {
		util.JsonErrorResponse("empty password!").Write(w)
		return
	}
}

type authLogoutResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

func AuthLogout(w http.ResponseWriter, r *http.Request) {
	user := util.GetUserFromRequestContext(r)
	session := util.GetSessionFromRequestContext(r)
	err := util.ReleaseUserSession(session)

	if err != nil {
		util.JsonErrorResponse(err.Error()).Write(w)
	}
	resBody := authLogoutResponse{
		ID:    user.ID,
		Email: user.Email,
	}
	util.JsonResponse(resBody).Write(w)
}

type authMeResponse struct {
	models.User
	EmailVerified   bool        `json:"email_verified"`
	ProfileImageUrl string      `json:"profile_image_url"`
	Role            models.Role `json:"roles"`
}

func AuthMeGet(w http.ResponseWriter, r *http.Request) {
	user := util.GetUserFromRequestContext(r)
	db := database.Database

	var userInfo models.User
	if result := db.GormDB.First(&userInfo, "id = ?", user.ID); result.Error != nil {
		util.JsonErrorResponse("user info not found").Write(w)
		return
	}

	var role models.Role
	if result := db.GormDB.Find(&role, "id = ?", userInfo.RoleID); result.Error != nil {
		util.JsonErrorResponse("user role not found").Write(w)
		return
	}

	meResponse := &authMeResponse{
		User: userInfo,
		Role: role,
	}

	util.JsonResponse(&meResponse).Write(w)
}

// type authAuthChangePasswordBody struct {
// 	OldPassword *string `json:"old_password"`
// 	NewPassword *string `json:"new_password"`
// }

// func AuthChangePassword(w http.ResponseWriter, r *http.Request) {
// 	var body authAuthChangePasswordBody
// 	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 		util.JsonErrorResponse(err.Error()).Write(w)
// 		return
// 	}
// 	if err := util.NewValidator().Validate(body); err != nil {
// 		util.JsonErrorResponse(err.Error()).Write(w)
// 		return
// 	}

// 	user := util.GetUserFromRequestContext(r)

// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*body.OldPassword)); err != nil {
// 		util.JsonErrorResponse(invalidPasswordMessage).Write(w)
// 		return
// 	}

// 	if *body.OldPassword == *body.NewPassword {
// 		util.JsonErrorResponse(newPasswordIsSameAsOldError).Write(w)
// 		return
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*body.NewPassword), bcrypt.DefaultCost)
// 	if err != nil {

// 		util.JsonErrorResponse(err.Error()).Write(w)
// 		return
// 	}

// 	db := database.Database
// 	user.Password = string(hashedPassword)
// 	result := db.GormDB.Save(&user)
// 	if result.Error != nil {
// 		util.JsonErrorResponse(failedUpdateUserPasswordMessage).Write(w)
// 		return
// 	}
// 	util.JsonResponse(user).Write(w)
// }

// type authAuthForgotPasswordBody struct {
// 	Email string `json:"email" validate:"required,email"`
// }

// func AuthForgotPassword(w http.ResponseWriter, r *http.Request) {
// 	var body authAuthForgotPasswordBody
// 	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 		util.JsonErrorResponse(err.Error()).Write(w)
// 		return
// 	}
// 	if err := util.NewValidator().Validate(body); err != nil {
// 		util.JsonErrorResponse(err.Error()).Write(w)
// 		return
// 	}

// 	db := database.Database
// 	var user models.User

// 	body.Email = strings.ToLower(body.Email)

// 	if result := db.GormDB.First(&user, "email = ?", body.Email); result.RowsAffected == 1 {
// 		otpCodeGenerated := util.GenerateRandomCode(6)

// 		user.OTPCode = otpCodeGenerated
// 		if result = db.GormDB.Save(user); result.Error != nil {
// 			util.JsonErrorResponse(failedToSaveDataMessage).Write(w)
// 			return
// 		}
// 		util.SendEmailForgotPassword(user.Email, otpCodeGenerated)
// 		util.JsonResponse(nil).WithMessage(emailSuccessfullySendMessage).Write(w)
// 		return
// 	}
// 	util.JsonErrorResponse(emailNotFoundMessage).Write(w)
// }

// type authAuthResetPasswordBody struct {
// 	// Secret      *string `json:"secret"`
// 	NewPassword *string `json:"new_password"`
// 	Email       *string `json:"email"`
// }

// func AuthResetPassword(w http.ResponseWriter, r *http.Request) {
// 	var body authAuthResetPasswordBody
// 	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 		util.JsonErrorResponse(err.Error()).Write(w)
// 		return
// 	}
// 	if err := util.NewValidator().Validate(body); err != nil {
// 		util.JsonErrorResponse(err.Error()).Write(w)
// 		return
// 	}

// 	db := database.Database
// 	var user models.User

// 	// if result := db.GormDB.First(&user, "password_reset_secret = ? AND password_reset_secret_expire_at >= ?", body.Secret, time.Now()); result.Error != nil {
// 	// 	util.JsonErrorResponse("Not found password reset code or code expired").Write(w)
// 	// 	return
// 	// }
// 	if result := db.GormDB.First(&user, "email = ?", body.Email); result.Error != nil {
// 		util.JsonErrorResponse(userNotFound).Write(w)
// 		return
// 	}
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*body.NewPassword), bcrypt.DefaultCost)
// 	if err != nil {
// 		util.JsonErrorResponse(err.Error()).Write(w)
// 		return
// 	}
// 	user.Password = string(hashedPassword)
// 	if result := db.GormDB.Save(&user); result.Error != nil {
// 		util.JsonErrorResponse(failedToUpdateUserPassword).Write(w)
// 		return
// 	}
// 	util.JsonResponse(nil).WithMessage(passwordChangedSuccessfully).Write(w)
// }

type authMeBody struct {
	Username    *string    `json:"username"`
	PhoneNumber *string    `json:"phone_number"`
	Birthdate   *util.Date `json:"birthdate"`
	Email       *string    `json:"email"`
	// FCMToken    *string    `json:"fcm_token"`
}

func AuthMePut(w http.ResponseWriter, r *http.Request) {
	db := database.Database
	var body authMeBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.JsonErrorResponse(err.Error()).Write(w)
		return
	}
	if err := util.NewValidator().Validate(body); err != nil {
		util.JsonErrorResponse(err.Error()).Write(w)
		return
	}

	user := util.GetUserFromRequestContext(r)

	if body.PhoneNumber != nil {
		user.PhoneNumber = *body.PhoneNumber
	}

	var emailCheck models.User
	if result := db.GormDB.First(&emailCheck, "email = ?", body.Email); result.RowsAffected == 0 {
		if body.Email != nil {
			user.Email = *body.Email
		}
	} else {
		util.JsonErrorResponse("duplicated email").Write(w)
		return
	}
	if result := db.GormDB.Save(&user); result.Error != nil {
		util.JsonErrorResponse("failed to update user data").Write(w)
		return
	}

	util.JsonResponse(user).Write(w)
}

// type confirmEmailBody struct {
// 	UserID  string `json:"user_id" validate:"required"`
// 	OtpCode string `json:"otp_code" validate:"required"`
// }

// func AuthConfirmEmail(w http.ResponseWriter, r *http.Request) {
// 	var user models.User

// 	db := database.Database

// 	var body confirmEmailBody
// 	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 		util.JsonErrorResponse(err.Error()).Write(w)
// 		return
// 	}
// 	if err := util.NewValidator().Validate(body); err != nil {
// 		util.JsonErrorResponse(err.Error()).Write(w)
// 		return
// 	}

// 	if result := db.GormDB.First(&user, "id = ? AND otp_code = ?", body.UserID, body.OtpCode); result.Error != nil {
// 		util.JsonErrorResponse("otp code wrong").Write(w)
// 		return
// 	}

// 	// user.EmailVerified = true
// 	if result := db.GormDB.Save(&user); result.Error != nil {
// 		util.JsonErrorResponse("failed to verify email").Write(w)
// 		return
// 	}

// 	util.JsonResponse(nil).WithMessage("email verified successfully").Write(w)
// }

// func OtpSend(w http.ResponseWriter, r *http.Request) {
// 	db := database.Database
// 	userID := chi.URLParam(r, "user_id")
// 	var user models.User
// 	if result := db.GormDB.First(&user, "id = ?", userID); result.RowsAffected == 0 {
// 		util.JsonErrorResponse(userNotFound).Write(w)
// 		return
// 	}

// 	otpCodeGenerated := util.GenerateRandomCode(6)
// 	log.Println(user.Email, otpCodeGenerated)

// 	if r.URL.Query().Get("type") == "deactivate" {
// 		util.SendEmailDeactivateUser(user.Email, otpCodeGenerated)
// 	} else {
// 		util.SendEmailAccountVerify(user.Email, otpCodeGenerated)
// 	}

// 	if result := db.GormDB.Model(&models.User{}).Where("id = ?", userID).Update("otp_code", otpCodeGenerated); result.Error != nil {
// 		util.JsonErrorResponse("User otp code save error: " + result.Error.Error()).Write(w)
// 		return
// 	}

// 	user.OTPCode = ""
// 	util.JsonResponse(user).Write(w)
// }

// type otpConfirmBody struct {
// 	Email   string `json:"email" validate:"required"`
// 	OtpCode string `json:"otp_code" validate:"required"`
// }

// func OtpConfirm(w http.ResponseWriter, r *http.Request) {
// 	var user models.User

// 	db := database.Database

// 	var body otpConfirmBody
// 	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 		util.JsonErrorResponse(err.Error()).Write(w)
// 		return
// 	}
// 	if err := util.NewValidator().Validate(body); err != nil {
// 		util.JsonErrorResponse(err.Error()).Write(w)
// 		return
// 	}

// 	if result := db.GormDB.First(&user, "email = ? AND otp_code = ?", body.Email, body.OtpCode); result.Error != nil {
// 		util.JsonErrorResponse(otpMissMatch).Write(w)
// 		return
// 	}

// 	util.JsonResponse(nil).WithMessage(otpMatched).Write(w)
// }
