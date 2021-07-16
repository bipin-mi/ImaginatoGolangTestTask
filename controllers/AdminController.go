package controllers

import "C"
import (
	"ImaginatoGolangTestTask/resources/request"
	"ImaginatoGolangTestTask/services"
	"ImaginatoGolangTestTask/validator"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/leebenson/conform"
	"strconv"

	"net/http"

	"ImaginatoGolangTestTask/shared/common"
	"ImaginatoGolangTestTask/shared/log"
	"ImaginatoGolangTestTask/shared/utils"
	msg "ImaginatoGolangTestTask/shared/utils/message"
)

type AdminController struct {
	AdminService services.IAdminService
	Validator    validator.IValidatorService
}

//Create is validate request and call the service
func (ac *AdminController) Create(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Admin Controller Called(SignUp).")
	var req request.CreateAdminRequest

	//decode the request body into struct and failed if any error occur
	if err := c.BindJSON(&req); err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, msg.InvalidRequest))
		return
	}

	if reqErr := conform.Strings(&req); reqErr != nil {
		log.GetLog().Info("REQUEST ERROR : ", reqErr.Error())
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, msg.InvalidRequest))
		return
	}

	//Struct field validation
	if resp, ok := ac.Validator.ValidateStruct(req, "CreateAdminRequest"); !ok {
		log.GetLog().Info("ERROR : ", "CreateAdminRequest validation errors")
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, resp))
		return
	}

	//call service
	serviceErr := ac.AdminService.Create(req)
	if serviceErr != nil {
		common.ErrorResponse(c.Writer, serviceErr)
		return
	}
	common.SuccessResponse(c.Writer, map[string]interface{}{}, msg.AdminAddedSuccess)
}

//Login is validate request and call the service
func (ac *AdminController) Login(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Admin Controller Called(Login).")
	var req request.LoginRequest

	//decode the request body into struct and failed if any error occur
	if err := c.BindJSON(&req); err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, msg.InvalidRequest))
		return
	}

	if reqErr := conform.Strings(&req); reqErr != nil {
		log.GetLog().Info("REQUEST ERROR : ", reqErr.Error())
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, msg.InvalidRequest))
		return
	}

	if resp, ok := ac.Validator.ValidateStruct(req, "LoginRequest"); !ok {
		log.GetLog().Info("ERROR : ", "Struct validation error")
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, resp))
		return
	}

	resp := ac.AdminService.Login(req)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])

	common.Respond(c.Writer, statusCode, resp)
	log.GetLog().Info("INFO : ", "Login complete...", resp["data"])
	return
}

//ForgotPassword is validate request and call the service
func (ac *AdminController) ForgotPassword(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Admin Controller Called(Forgot Password).")
	var req request.ForgotPasswordRequest

	//decode the request body into struct and failed if any error occur
	if err := c.BindJSON(&req); err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, msg.InvalidRequest))
		return
	}

	if reqErr := conform.Strings(&req); reqErr != nil {
		log.GetLog().Info("REQUEST ERROR : ", reqErr.Error())
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, msg.InvalidRequest))
		return
	}

	if resp, ok := ac.Validator.ValidateStruct(req, "ForgotPasswordRequest"); !ok {
		log.GetLog().Info("ERROR : ", "Struct validation error")
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, resp))
		return
	}

	forgotResp := ac.AdminService.ForgotPassword(req)
	statusCode := common.GetHTTPStatusCode(forgotResp["res_code"])

	common.Respond(c.Writer, statusCode, forgotResp)
	log.GetLog().Info("INFO : ", "Forgot password complete...", forgotResp["data"])
	return
}

//ResetPassword is validate request and call the service
func (ac *AdminController) ResetPassword(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Admin Controller Called(Reset Password).")
	var req request.ResetPasswordRequest

	//decode the request body into struct and failed if any error occur
	if err := c.BindJSON(&req); err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, msg.InvalidRequest))
		return
	}

	if resp, ok := ac.Validator.ValidateStruct(req, "ResetPasswordRequest"); !ok {
		log.GetLog().Info("ERROR : ", "Struct validation error")
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, resp))
		return
	}

	resetResp := ac.AdminService.ResetPassword(req)
	statusCode := common.GetHTTPStatusCode(resetResp["res_code"])

	common.Respond(c.Writer, statusCode, resetResp)
	log.GetLog().Info("INFO : ", "Password reset successfully...")
	return
}

//List is validate request and call the service
func (ac *AdminController) List(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Admin Controller Called(AdminList).")

	pageStr := c.DefaultQuery("page", "1")
	if c.Param("page") != "" {
		pageStr = c.Param("page")
	}
	sizeStr := c.DefaultQuery("size", "10")
	if c.Param("size") != "" {
		sizeStr = c.Param("size")
	}
	fmt.Println(pageStr, sizeStr)
	page, size := utils.PageAttributes(pageStr, sizeStr)
	fmt.Println(page, size)
	filter := map[string]string{
		"name":  c.DefaultQuery("name", ""),
		"email": c.DefaultQuery("email", ""),
	}

	sortBy := c.DefaultQuery("sort_by", "updated_at")
	sortDir := c.DefaultQuery("sort_dir", "desc")

	// attributes for page
	pageAttr := utils.PageAttr{
		Page:    page,
		Size:    size,
		Filter:  filter,
		SortBy:  sortBy,
		SortDir: sortDir,
	}

	//call service
	resp := ac.AdminService.List(pageAttr)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])

	//return response using api helper
	common.Respond(c.Writer, statusCode, resp)
	return
}

//AdminDelete is validate request and call the service
func (ac *AdminController) AdminDelete(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Admin Controller Called(Delete Admin).")
	var req request.DeleteAdminRequest

	id := c.Param("id")
	if id == "" {
		log.GetLog().Info("ERROR : ", "Invalid Request")
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, "Invalid Request"))
		return
	}

	req.ID, _ = strconv.ParseInt(id, 10, 32)

	// Struct field validation
	if resp, ok := ac.Validator.ValidateStruct(req, "DeleteAdminRequest"); !ok {
		log.GetLog().Info("ERROR : ", "Struct validation error")
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, resp))
		return
	}

	//call service
	deleteResp := ac.AdminService.ServiceDelete(req)
	statusCode := common.GetHTTPStatusCode(deleteResp["res_code"])

	//return response using api helper
	common.Respond(c.Writer, statusCode, deleteResp)
	log.GetLog().Info("INFO : ", "Admin deleted successfully...")
	return

}

func (ac *AdminController) VerifyEmail(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Admin Controller Called(Verify Email).")
	var req request.VerifyEmailRequest

	resetToken := c.Param("reset_token")
	if resetToken == "" {
		log.GetLog().Info("ERROR : ", "Invalid request")
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, msg.InvalidRequest))
		return
	}

	req.ResetToken = resetToken

	// Struct field validation
	if resp, ok := ac.Validator.ValidateStruct(req, "VerifyEmailRequest"); !ok {
		log.GetLog().Info("ERROR : ", "Struct validation error")
		common.Respond(c.Writer, http.StatusBadRequest, common.ResponseErrorWithCode(common.CodeBadRequest, resp))
		return
	}

	//call service
	resp := ac.AdminService.VerifyEmail(req)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])

	//return response using api helper
	common.Respond(c.Writer, statusCode, resp)
	log.GetLog().Info("INFO : ", "Email is verified.")
	return
}
