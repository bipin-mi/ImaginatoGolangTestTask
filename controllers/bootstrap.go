package controllers

import (
	"ImaginatoGolangTestTask/services"
	"ImaginatoGolangTestTask/validator"
)

//InitController Initialize the controller
func InitController(validatorService validator.IValidatorService, adminService services.IAdminService) *AdminController {
	adminCtl := AdminController{
		Validator:    validatorService,
		AdminService: adminService,
	}

	return &adminCtl
}
