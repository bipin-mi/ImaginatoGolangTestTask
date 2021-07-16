package repository

import (
	"fmt"
	"time"

	"ImaginatoGolangTestTask/model"
	"ImaginatoGolangTestTask/resources/request"
	"ImaginatoGolangTestTask/resources/response"
	"ImaginatoGolangTestTask/shared/common"
	"ImaginatoGolangTestTask/shared/database"
	"ImaginatoGolangTestTask/shared/log"
	"ImaginatoGolangTestTask/shared/utils"
	"ImaginatoGolangTestTask/shared/utils/const"
)

const (
	updateResetToken       = "UPDATE admin SET reset_token = ?,reset_token_expiry=?  WHERE id = ? AND deleted_at IS NULL"
	updateEmailVerifyAdmin = "UPDATE admin SET reset_token = ?, status = ?, verified_status = ? WHERE reset_token = ? AND deleted_at IS NULL"
)

type IAdminRepository interface {
	AdminSignUp(conn database.IConnection, admin *request.CreateAdminRequest) error
	GetAdminByEmail(conn database.IConnection, email string) (model.Admin, error)
	AdminLogin(conn database.IConnection, admin *request.LoginRequest) (response.Login, error)
	AdminForgotPassword(conn database.IConnection, admin *request.ForgotPasswordData) error
	GetAdminByResetToken(conn database.IConnection, token string) (model.Admin, error)
	UpdateAdmin(conn database.IConnection, admin *request.ResetPasswordData) (response.LoginResponse, error)
	AdminList(conn database.IConnection, req *utils.PageAttr) ([]model.Admin, uint)
	AdminDelete(conn database.IConnection, admin *request.DeleteAdminRequest) error
	UpdateEmailVerifyAdmin(conn database.IConnection, user *request.VerifyEmailRequest) error
}

type adminRepo struct{}

func NewAdminWriter() IAdminRepository {
	return &adminRepo{}
}

//AdminSignUp insert the admin user details
func (ar *adminRepo) AdminSignUp(conn database.IConnection, req *request.CreateAdminRequest) error {
	log.GetLog().Info("INFO : ", "Admin Repo Called(SignUp).")

	admin := model.Admin{
		Name:             req.Name,
		Email:            req.Email,
		Password:         req.Password,
		ResetToken:       req.ResetToken,
		ResetTokenExpiry: req.ResetTokenExpiry,
	}

	err := conn.GetDB().Create(&admin).Error
	if err != nil {
		log.GetLog().Info("ERROR(query) : ", err.Error())
		return err
	}
	return nil
}

//GetAdminByEmail retrieve and returns admin details based on email
func (ar *adminRepo) GetAdminByEmail(conn database.IConnection, email string) (model.Admin, error) {
	log.GetLog().Info("INFO : ", "Admin Repo Called(GetAdminByEmail).")
	admin := model.Admin{}
	err := conn.GetDB().Where("email = ?", email).FirstOrInit(&admin).Error
	if err != nil {
		log.GetLog().Info("ERROR(query) : ", err.Error())
		return model.Admin{}, err
	}
	return admin, nil
}

//getOneAdmin retrieve and returns single admin user details
func (ar *adminRepo) getOneAdmin(id uint64) (response.LoginResponse, error) {
	conn := database.NewConnection()
	admin := response.LoginResponse{}
	err := conn.GetDB().Table("admin").Where("id =?", id).First(&admin).Error
	if err != nil {
		log.GetLog().Info("ERROR(query) : ", err.Error())
		return response.LoginResponse{}, err
	}
	return admin, nil
}

//AdminLogin retrieve and returns login users details
func (ar *adminRepo) AdminLogin(conn database.IConnection, req *request.LoginRequest) (response.Login, error) {
	log.GetLog().Info("INFO : ", "Admin Repo Called(Login).")
	admin := response.Login{}

	err := conn.GetDB().Table("admin").Where("email = ? AND status = ? AND verified_status = ?", req.Email, _const.Active, _const.Verified).First(&admin).Error

	if err != nil {
		log.GetLog().Error(err, "ERROR(query) : ")
		return admin, err
	}
	return admin, nil
}

//AdminForgotPassword update the admins reset password token and expiry date
func (ar *adminRepo) AdminForgotPassword(conn database.IConnection, admin *request.ForgotPasswordData) error {
	log.GetLog().Info("INFO : ", "Admin Repo Called(ForgotPassword).")

	err := conn.GetDB().Exec(updateResetToken, admin.ResetToken, admin.ResetTokenExpiry, admin.Id).Error
	if err != nil {
		log.GetLog().Info("ERROR(query) : ", err.Error())
		return err
	}

	return nil
}

//GetAdminByResetToken retrieve and returns admin user based on reset password token
func (ar *adminRepo) GetAdminByResetToken(conn database.IConnection, token string) (model.Admin, error) {
	log.GetLog().Info("INFO : ", "Admin Repo Called(GetAdminByResetToken).")
	admin := model.Admin{}

	err := conn.GetDB().Where("reset_token = ?", token).First(&admin).Scan(&admin).Error
	if err != nil {
		log.GetLog().Info("ERROR(query) : ", err.Error())
		return admin, err
	}
	return admin, nil
}

//UpdateAdmin update the admin user
func (ar *adminRepo) UpdateAdmin(conn database.IConnection, req *request.ResetPasswordData) (response.LoginResponse, error) {
	log.GetLog().Info("INFO : ", "Admin Repo Called(UpdateAdmin).")

	//err := conn.GetDB().Exec(updatePassword, utils.HashedPassword(req.Password), "", nil, req.Id).Error
	err := conn.GetDB().Table("admin").Where("id = ?", req.Id).Update(model.Admin{
		Password:         utils.HashedPassword(req.Password),
		ResetToken:       "",
		ResetTokenExpiry: time.Time{},
	}).Error
	if err != nil {
		log.GetLog().Info("ERROR(query) : ", err.Error())
		return response.LoginResponse{}, err
	}

	adminData, adminDataErr := ar.getOneAdmin(req.Id)
	if adminDataErr != nil {
		log.GetLog().Info("ERROR(data) : ", adminDataErr.Error())
		return response.LoginResponse{}, adminDataErr
	}

	return adminData, nil
}

//AdminList retrieves all admin users in DESC order of last updated
func (ar *adminRepo) AdminList(conn database.IConnection, req *utils.PageAttr) ([]model.Admin, uint) {
	log.GetLog().Info("INFO : ", "Admin Repo Called(AdminList).")
	var objs []model.Admin
	var count uint

	query := conn.GetDB().Table("admin").Select("*")
	// filtering
	for k, v := range req.Filter {
		if v != string(common.DefaultEmpty) {
			query = query.Where(fmt.Sprintf("%s LIKE ?", utils.SnakeCase(k)), fmt.Sprintf("%s%%", v))
		}
	}

	// count
	query.Count(&count)

	// sorting
	if req.SortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", utils.SnakeCase(req.SortBy), req.SortDir))
	}

	// paging
	query = query.Offset((req.Page - 1) * req.Size).Limit(req.Size)

	query.Find(&objs)
	return objs, count
}

//AdminDelete soft delete the admin user
func (ar *adminRepo) AdminDelete(conn database.IConnection, admin *request.DeleteAdminRequest) error {
	log.GetLog().Info("INFO : ", "Admin Repo Called(RepoAdminDelete).")

	//err := conn.GetDB().Exec(deleteAdmin, _const.InActive, now, admin.IdMsb, admin.IdLsb).Error
	err := conn.GetDB().Table("admin").Delete(model.Admin{
		Model: model.Model{
			Id: admin.ID,
		},
	}).Update(model.Admin{
		Status: _const.InActive,
	}).Error
	if err != nil {
		log.GetLog().Info("ERROR(query) : ", err.Error())
		return err
	}
	return nil
}

func (ar *adminRepo) UpdateEmailVerifyAdmin(conn database.IConnection, req *request.VerifyEmailRequest) error {
	log.GetLog().Info("INFO : ", "User Repo Called(UpdateEmailVerifyUser).")

	err := conn.GetDB().Exec(updateEmailVerifyAdmin, "", _const.Active, _const.Verified, req.ResetToken).Error
	if err != nil {
		log.GetLog().Info("ERROR(query) : ", err.Error())
		return err
	}
	return nil
}
