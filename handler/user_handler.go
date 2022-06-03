package handler

import (
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
		res := helper.ApiResponse("Input Data Gagal!", http.StatusUnprocessableEntity, "gagal", err)

		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	_, errUser := h.userService.Register(input)
	if errUser != nil {
		res := helper.ApiResponse("Input Data Gagal!", http.StatusBadRequest, "gagal", errUser)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	data := gin.H{
		"status": "Berhasil Membuat Akun Baru!",
	}

	res := helper.ApiResponse("Berhasil Membuat Akun Baru!", http.StatusCreated, "berhasil", data)

	c.JSON(http.StatusCreated, res)
}

func (h *UserHandler) Login(c *gin.Context) {
	var input user.InputLogin

	err := c.ShouldBindJSON(&input)
	if err != nil {
		res := helper.ApiResponse("Login Gagal!", http.StatusUnprocessableEntity, "gagal", nil)
		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	loginUser, errLogin := h.userService.Login(input)
	if errLogin != nil {
		res := helper.ApiResponse("Login Gagal!", http.StatusUnprocessableEntity, "gagal", nil)

		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	token, errToken := h.authService.GenerateTokenJWT(loginUser.ID, loginUser.Fullname, loginUser.Role)

	if errToken != nil {
		res := helper.ApiResponse("Gagal Membuat Token", http.StatusBadRequest, "gagal", nil)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	formatter := user.FormatUser(loginUser, token)

	res := helper.ApiResponse("berhasil login", http.StatusOK, "berhasil", formatter)

	c.JSON(http.StatusCreated, res)
}
