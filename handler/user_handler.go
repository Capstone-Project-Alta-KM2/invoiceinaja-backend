package handler

import (
	"fmt"
	"net/http"

	"invoiceinaja/auth"
	"invoiceinaja/domain/user"
	"invoiceinaja/helper"

	"github.com/gin-gonic/gin"
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

		response := helper.ApiResponse("Input Data Gagal!", http.StatusUnprocessableEntity, "error", nil, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, errAvail := h.userService.IsEmailAvailable(user.InputCheckEmail{Email: input.Email})
	if errAvail != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := helper.ApiResponse("Email checking failed", http.StatusUnprocessableEntity, "gagal", nil, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	if isEmailAvailable != true {
		data := gin.H{
			"status": "Gagal Membuat Akun Baru!",
		}
		res := helper.ApiResponse("Email sudah digunakan!", http.StatusBadRequest, "gagal", nil, data)
		c.JSON(http.StatusBadRequest, res)
	} else {
		_, errUser := h.userService.Register(input)
		if errUser != nil {
			res := helper.ApiResponse("Input Data Gagal!", http.StatusBadRequest, "gagal", nil, errUser)

			c.JSON(http.StatusBadRequest, res)
			return
		}

		data := gin.H{
			"status": "Berhasil Membuat Akun Baru!",
		}

		res := helper.ApiResponse("Berhasil Membuat Akun Baru!", http.StatusCreated, "berhasil", nil, data)

		c.JSON(http.StatusCreated, res)
	}
}

func (h *UserHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.InputCheckEmail

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Email checking failed", http.StatusUnprocessableEntity, "gagal", nil, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := helper.ApiResponse("Email checking failed", http.StatusUnprocessableEntity, "gagal", nil, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
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

		response := helper.ApiResponse("Login Gagal!", http.StatusUnprocessableEntity, "error", nil, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loginUser, errLogin := h.userService.Login(input)
	if errLogin != nil {
		res := helper.ApiResponse("Login Gagal!", http.StatusUnprocessableEntity, "gagal", nil, nil)

		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	token, errToken := h.authService.GenerateTokenJWT(loginUser.ID, loginUser.Fullname, loginUser.Role)

	if errToken != nil {
		res := helper.ApiResponse("Gagal Membuat Token", http.StatusBadRequest, "gagal", nil, nil)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	formatter := user.FormatUser(loginUser, token)

	res := helper.ApiResponse("berhasil login", http.StatusOK, "berhasil", nil, formatter)

	c.JSON(http.StatusCreated, res)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		res := helper.ApiResponse("Gagal Mengunggah Gambar!", http.StatusBadRequest, "error", nil, data)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	// didapatkan dari JWT
	currentUser := c.MustGet("currentUser").(user.User)
	userId := currentUser.ID

	path := fmt.Sprintf("images/%d-%s", userId, file.Filename)

	errImage := c.SaveUploadedFile(file, path)
	if errImage != nil {
		data := gin.H{"unggahan": false}
		res := helper.ApiResponse("Gagal Mengunggah Gambar!", http.StatusBadRequest, "gagal", nil, data)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	_, errUser := h.userService.SaveAvatar(userId, path)
	if errUser != nil {
		data := gin.H{"unggahan": false}
		res := helper.ApiResponse("Gagal Mengunggah Gambar!", http.StatusBadRequest, "gagal", nil, data)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	data := gin.H{"unggahan": true}
	res := helper.ApiResponse("Berhasil Mengunggah Gambar!", http.StatusOK, "berhasil", nil, data)

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

		response := helper.ApiResponse("Gagal Memperbaharui Data", http.StatusUnprocessableEntity, "error", nil, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updated, errUpdate := h.userService.UpdateUser(userId, input)
	if errUpdate != nil {
		res := helper.ApiResponse("Gagal Memperbaharui Data", http.StatusUnprocessableEntity, "gagal", nil, err)

		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	formatter := user.FormatUpdateUser(updated)

	res := helper.ApiResponse("Berhasil Memperbaharui Data", http.StatusCreated, "success", nil, formatter)

	c.JSON(http.StatusCreated, res)
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var input user.InputCheckEmail
	err := c.ShouldBindJSON(&input)
	if err != nil {
		res := helper.ApiResponse("Something Wrong!!!", http.StatusUnprocessableEntity, "failed", nil, err)

		c.JSON(http.StatusUnprocessableEntity, res)
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
		"is_sent": true,
	}

	res := helper.ApiResponse("Please Check Your Email", http.StatusOK, "success", nil, data)

	c.JSON(http.StatusCreated, res)
}
