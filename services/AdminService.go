package services

import (
	"ImaginatoGolangTestTask/repository"
	"ImaginatoGolangTestTask/resources/request"
	"golang.org/x/crypto/bcrypt"

	"errors"
	"net/http"
	"os"
	"time"

	"ImaginatoGolangTestTask/model"
	"ImaginatoGolangTestTask/shared/common"
	"ImaginatoGolangTestTask/shared/database"
	"ImaginatoGolangTestTask/shared/log"
	"ImaginatoGolangTestTask/shared/utils"
	"ImaginatoGolangTestTask/shared/utils/const"
	msg "ImaginatoGolangTestTask/shared/utils/message"
	"ImaginatoGolangTestTask/shared/utils/middleware"
	"ImaginatoGolangTestTask/thirdparty/email"
)

type IAdminService interface {
	Create(req request.CreateAdminRequest) error
	Login(req request.LoginRequest) map[string]interface{}
	ForgotPassword(req request.ForgotPasswordRequest) map[string]interface{}
	ResetPassword(req request.ResetPasswordRequest) map[string]interface{}
	List(page utils.PageAttr) map[string]interface{}
	ServiceDelete(req request.DeleteAdminRequest) map[string]interface{}
	VerifyEmail(req request.VerifyEmailRequest) map[string]interface{}
}

type AdminService struct {
	Admin     model.Admin
	AdminRepo repository.IAdminRepository
}

func NewAdminService() IAdminService {
	adminRepo := repository.NewAdminWriter()
	return &AdminService{
		AdminRepo: adminRepo,
	}
}

//Create this service create an admin
func (as *AdminService) Create(req request.CreateAdminRequest) error {
	log.GetLog().Info("INFO : ", "Admin Service Called(SignUp).")
	conn := database.NewConnection()

	// Duplicate Admin validation
	adminResponse, err := as.AdminRepo.GetAdminByEmail(conn, req.Email)
	if err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return err
	}

	if adminResponse.Id > 0 {
		log.GetLog().Info("WARN : ", "Duplicate Admin...")
		return errors.New(msg.EmailInUse)
	}

	req.Password = utils.HashedPassword(req.Password)
	req.ResetToken = utils.RandomKeyGenerator(32, "alphaNum")
	req.ResetTokenExpiry = time.Now().Add(time.Hour * 24)

	//Call admin repository
	repoErr := as.AdminRepo.AdminSignUp(conn, &req)
	if repoErr != nil {
		log.GetLog().Info("ERROR(from repo) : ", repoErr.Error())
		return repoErr
	}

	// Send email for verify email
	from := os.Getenv("Sender")
	emailRequest := email.Request{
		From:    from,
		To:      []string{req.Email},
		Subject: "Verify Email",
		Body:    os.Getenv("VerifyEmailURI") + req.ResetToken,
	}

	go emailRequest.SendSmtpEmail(req.Name, "/templates/adminSignUp.html")

	return nil
}

//Login this service check the login process
func (as *AdminService) Login(req request.LoginRequest) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Admin Service Called(Login).")
	conn := database.NewConnection()

	resp, err := as.AdminRepo.AdminLogin(conn, &req)
	if err != nil {
		log.GetLog().Info("ERROR(from repo) : ", err.Error())
		return common.ResponseErrorWithCode(http.StatusBadRequest, msg.InvalidCredential)
	}

	if resp.Email != req.Email {
		log.GetLog().Info("ERROR(from repo) : ", "Email is not registered with us.")
		return common.ResponseErrorWithCode(http.StatusBadRequest, msg.InvalidCredential)
	}

	credErr := bcrypt.CompareHashAndPassword([]byte(resp.Password), []byte(req.Password))
	if credErr != nil && credErr == bcrypt.ErrMismatchedHashAndPassword {
		log.GetLog().Info("ERROR(hash) : ", credErr.Error())
		return common.ResponseErrorWithCode(http.StatusBadRequest, msg.InvalidCredential)
	}

	accessSecret := os.Getenv("AccessSecret")
	userData := request.TokenDataRequest{
		ID: resp.Id,
	}

	token, tokenErr := middleware.GenerateToken([]byte(accessSecret), userData)

	if tokenErr != nil {
		log.GetLog().Info("ERROR(from repo) : ", tokenErr.Error())
		return common.ResponseErrorWithCode(http.StatusBadRequest, msg.InvalidCredential)
	}
	resp.Password = ""
	loginResponse := common.ResponseSuccessWithToken(msg.LoginSuccess, resp, token)
	return loginResponse
}

//ForgotPassword this service used in forgot password process
func (as *AdminService) ForgotPassword(req request.ForgotPasswordRequest) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Admin Service Called(Forgot Password).")
	conn := database.NewConnection()

	forgotReq := request.ForgotPasswordData{}
	forgotReq.ResetToken = utils.RandomKeyGenerator(32, "alphaNum")
	forgotReq.Email = req.Email

	// Check Admin exists
	adminResponse, err := as.AdminRepo.GetAdminByEmail(conn, req.Email)
	if err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return common.ResponseErrorWithCode(http.StatusBadRequest, msg.SomethingWrong)
	}

	if adminResponse.Email != req.Email {
		log.GetLog().Info("WARN : ", "Admin Not Found...")
		return common.ResponseErrorWithCode(http.StatusBadRequest, msg.EmailNotRegistered)
	}

	// Check for inactive req
	if adminResponse.Status == _const.InActive || adminResponse.VerifiedStatus == _const.Unverified {
		log.GetLog().Info("WARN : ", "Admin Not active...")
		return common.ResponseErrorWithCode(http.StatusBadRequest, msg.InActive)
	}

	forgotReq.Id = adminResponse.Id
	forgotReq.ResetTokenExpiry = time.Now().Add(time.Hour * 24)

	from := os.Getenv("Sender")
	emailRequest := email.Request{
		From:    from,
		To:      []string{req.Email},
		Subject: "Reset your password.",
		Body:    forgotReq.ResetToken,
	}

	go emailRequest.SendSmtpEmail(adminResponse.Name, "/templates/email.html")

	err = as.AdminRepo.AdminForgotPassword(conn, &forgotReq)
	if err != nil {
		log.GetLog().Info("ERROR(from repo) : ", err.Error())
		return nil
	}

	forgotResponse := common.ResponseSuccessWithObj(msg.ForgotPasswordSuccess, nil)
	return forgotResponse
}

//ResetPassword this service used in reset password process
func (as *AdminService) ResetPassword(req request.ResetPasswordRequest) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Admin Service Called(Reset Password).")
	conn := database.NewConnection()

	resReq := request.ResetPasswordData{}
	resReq.Password = req.Password
	resReq.ResetPasswordToken = req.ResetToken
	adminResponse, err := as.AdminRepo.GetAdminByResetToken(conn, req.ResetToken)
	if err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return common.ResponseErrorWithCode(http.StatusBadRequest, msg.TokenIsNotValid)
	}

	resReq.Id = uint64(adminResponse.Id)

	//set token expired validation
	diff := time.Since(adminResponse.ResetTokenExpiry)
	if diff.Hours() > 24 {
		log.GetLog().Info("WARN : ", "Admin Token Expired...")
		return common.ResponseErrorWithCode(http.StatusBadRequest, msg.TokenIsExpired)
	}

	_, err = as.AdminRepo.UpdateAdmin(conn, &resReq)
	if err != nil {
		log.GetLog().Info("ERROR(from repo) : ", err.Error())
		return nil
	}

	resetResponse := common.ResponseSuccessWithObj(msg.ResetPasswordSuccess, nil)
	return resetResponse
}

//List this service used to get list of admin users
func (as *AdminService) List(page utils.PageAttr) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Admin Service Called(AdminList).")
	conn := database.NewConnection()

	resp, total := as.AdminRepo.AdminList(conn, &page)

	response := common.ResponseSuccessWithArray(msg.AdminListSuccess, resp)
	response["meta"].(map[string]interface{})["count"] = total
	return response
}

//ServiceDelete this service used to delete the admin user
func (as *AdminService) ServiceDelete(req request.DeleteAdminRequest) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Admin Service Called(ServiceAdminDelete).")
	conn := database.NewConnection()

	err := as.AdminRepo.AdminDelete(conn, &req)
	if err != nil {
		log.GetLog().Info("ERROR(from repo) : ", err.Error())
		return common.ResponseErrorWithCode(http.StatusBadRequest, msg.InvalidCredential)
	}
	deleteResponse := common.ResponseSuccessWithObj(msg.AdminDeletedSuccess, nil)
	return deleteResponse
}

//VerifyEmail this service used to verify the admin user
func (as *AdminService) VerifyEmail(req request.VerifyEmailRequest) map[string]interface{} {
	log.GetLog().Info("INFO : ", "User Service Called(Verify Email).")

	conn := database.NewConnection()
	// Check User exists
	adminRes, err := as.AdminRepo.GetAdminByResetToken(conn, req.ResetToken)
	if err != nil {
		log.GetLog().Info("ERROR(DB) : ", err.Error())
		return common.ResponseForWeb(http.StatusBadRequest, msg.SomethingWrong)
	}
	if adminRes.Id == 0 || adminRes.VerifiedStatus != _const.Unverified || adminRes.Status != _const.InActive {
		log.GetLog().Info("WARN : ", "User Not Found...")
		return common.ResponseForWeb(http.StatusBadRequest, msg.UserNotFoundByToken)
	}

	if adminRes.ResetToken != req.ResetToken {
		log.GetLog().Info("WARN : ", "User Not Found...")
		return common.ResponseForWeb(http.StatusBadRequest, msg.UserNotFoundByToken)
	}

	err = as.AdminRepo.UpdateEmailVerifyAdmin(conn, &req)
	if err != nil {
		log.GetLog().Info("ERROR(from repo) : ", err.Error())
		return nil
	}

	response := common.ResponseForWeb(http.StatusOK, msg.EmailVerificationSuccess)
	return response
}
