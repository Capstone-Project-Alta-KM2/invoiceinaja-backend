package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"invoiceinaja/auth"
	"invoiceinaja/domain/user"
	"invoiceinaja/helper"
	"invoiceinaja/utils"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-password/password"
)

type UserHandler struct {
	userService user.IService
	authService auth.Service
}

func NewUserHandler(userService user.IService, authService auth.Service) *UserHandler {
	return &UserHandler{userService, authService}
}

func (h *UserHandler) UserRegister(c *gin.Context) {
	var input user.InputRegister
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Data Input Failed!", http.StatusBadRequest, "error", nil, errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	isEmailAvailable, errAvail := h.userService.IsEmailAvailable(user.InputCheckEmail{Email: input.Email})
	if errAvail != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := helper.ApiResponse("Email checking failed", http.StatusBadRequest, "failed", nil, errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if !isEmailAvailable {
		data := gin.H{
			"status": "Failed to Create New Account!",
		}
		res := helper.ApiResponse("Email already used!", http.StatusBadRequest, "failed", nil, data)
		c.JSON(http.StatusBadRequest, res)
	} else {
		user, errUser := h.userService.Register(input)
		if errUser != nil {
			res := helper.ApiResponse("Data Input Failed!", http.StatusBadRequest, "failed", nil, errUser)

			c.JSON(http.StatusBadRequest, res)
			return
		}

		token, errToken := h.authService.GenerateTokenJWT(user.ID, user.Fullname, user.Role)
		if errToken != nil {
			res := helper.ApiResponse("Failed to generate Token", http.StatusBadRequest, "failed", nil, nil)

			c.JSON(http.StatusBadRequest, res)
			return
		}

		otp, _ := password.Generate(4, 4, 0, true, true)

		var a string
		message := utils.SendMailOTP(input.Email, otp)
		for _, v := range message.ResultsV31 {
			a = v.Status
		}
		if a != "success" {
			res := errors.New("failed send email")
			c.JSON(http.StatusCreated, res)
			return
		}

		data := gin.H{
			"status":   "Successfully Created New Account!",
			"token":    token,
			"code_otp": otp,
		}

		res := helper.ApiResponse("Successfully Created New Account!", http.StatusCreated, "success", nil, data)

		c.JSON(http.StatusCreated, res)
	}
}

func (h *UserHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.InputCheckEmail

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Email checking failed", http.StatusBadRequest, "failed", nil, errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := helper.ApiResponse("Email checking failed", http.StatusBadRequest, "failed", nil, errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registered"

	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.ApiResponse(metaMessage, http.StatusOK, "success", nil, data)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Login(c *gin.Context) {
	var input user.InputLogin

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Login Failed!", http.StatusBadRequest, "error", nil, errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loginUser, errLogin := h.userService.Login(input)
	if errLogin != nil {
		res := helper.ApiResponse("Login Failed!", http.StatusBadRequest, "failed", nil, nil)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	token, errToken := h.authService.GenerateTokenJWT(loginUser.ID, loginUser.Fullname, loginUser.Role)
	if errToken != nil {
		res := helper.ApiResponse("Failed to generate Token", http.StatusBadRequest, "failed", nil, nil)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	formatter := user.FormatUser(loginUser, token)

	res := helper.ApiResponse("Login Successfully", http.StatusOK, "success", nil, formatter)

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		res := helper.ApiResponse("Failed to Upload Avatar!", http.StatusBadRequest, "error", nil, data)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	if !strings.Contains(file.Filename, "jpg") && !strings.Contains(file.Filename, "png") {
		data := gin.H{"is_uploaded": false}
		res := helper.ApiResponse("Failed to Upload File!", http.StatusBadRequest, "failed", nil, data)

		c.JSON(http.StatusBadRequest, res)
		return

	}

	// didapatkan dari JWT
	currentUser := c.MustGet("currentUser").(user.User)
	userId := currentUser.ID

	path := fmt.Sprintf("images/%d-%s", userId, file.Filename)

	errImage := c.SaveUploadedFile(file, path)
	if errImage != nil {
		data := gin.H{"is_uploaded": false}
		res := helper.ApiResponse("Failed to Upload File!", http.StatusBadRequest, "failed", nil, data)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	_, errUser := h.userService.SaveAvatar(userId, path)
	if errUser != nil {
		data := gin.H{"is_uploaded": false}
		res := helper.ApiResponse("Failed to Upload Avatar!", http.StatusBadRequest, "failed", nil, data)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	data := gin.H{"is_uploaded": true}
	res := helper.ApiResponse("Successfully Uploaded Avatar!", http.StatusOK, "success", nil, data)

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	// cek yg akses login
	currentUser := c.MustGet("currentUser").(user.User)
	userId := currentUser.ID

	var input user.InputUpdate
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Update Data Failed", http.StatusBadRequest, "error", nil, errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	updated, errUpdate := h.userService.UpdateUser(userId, input)
	if errUpdate != nil {
		res := helper.ApiResponse("Update Data Failed", http.StatusBadRequest, "failed", nil, err)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	formatter := user.FormatUpdateUser(updated)

	res := helper.ApiResponse("Update Data Successfully", http.StatusOK, "success", nil, formatter)

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	// cek yg akses login
	currentUser := c.MustGet("currentUser").(user.User)
	userId := currentUser.ID

	var input user.InputChangePassword
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Change Password Failed", http.StatusBadRequest, "error", nil, errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, errUpdate := h.userService.ChangePassword(input, userId)
	if errUpdate != nil {
		res := helper.ApiResponse("Change Password Failed", http.StatusBadRequest, "failed", nil, err)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	data := gin.H{"is_changed": true}
	res := helper.ApiResponse("Successfully Change Password!", http.StatusCreated, "success", nil, data)

	c.JSON(http.StatusCreated, res)

}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var input user.InputCheckEmail
	err := c.ShouldBindJSON(&input)
	if err != nil {
		res := helper.ApiResponse("Something Wrong!!!", http.StatusBadRequest, "failed", nil, err)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	_, errData := h.userService.IsEmailAvailable(input)
	if errData != nil {
		res := helper.ApiResponse("Something Wrong!!!", http.StatusBadRequest, "failed", nil, errData)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	_, errUser := h.userService.ResetPassword(input)
	if errUser != nil {
		res := helper.ApiResponse("Something Wrong!!!", http.StatusBadRequest, "failed", nil, errUser)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	data := gin.H{
		"is_send": true,
	}

	res := helper.ApiResponse("Please Check Your Email", http.StatusOK, "success", nil, data)

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) ResendOTP(c *gin.Context) {
	// didapatkan dari JWT
	currentUser := c.MustGet("currentUser").(user.User)
	otp, _ := password.Generate(4, 4, 0, true, true)

	var a string
	message := utils.SendMailOTP(currentUser.Email, otp)
	for _, v := range message.ResultsV31 {
		a = v.Status
	}
	if a != "success" {
		res := errors.New("failed send email")
		c.JSON(http.StatusCreated, res)
		return
	}

	data := gin.H{
		"is_send":  true,
		"code_otp": otp,
	}

	res := helper.ApiResponse("Please Check Your Email", http.StatusOK, "success", nil, data)

	c.JSON(http.StatusCreated, res)
}
