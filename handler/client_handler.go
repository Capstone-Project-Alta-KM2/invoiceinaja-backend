package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"invoiceinaja/auth"
	"invoiceinaja/domain/client"
	"invoiceinaja/domain/user"
	"invoiceinaja/helper"
	"invoiceinaja/utils"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	clientService client.IService
	userService   user.IService
	authService   auth.Service
}

func NewClientHandler(ClientService client.IService, userService user.IService, authService auth.Service) *ClientHandler {
	return &ClientHandler{ClientService, userService, authService}
}

func (h *ClientHandler) AddClient(c *gin.Context) {
	//
	var input client.InputAddClient

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Add Client Failed!", http.StatusUnprocessableEntity, "error", nil, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// // didapatkan dari JWT
	currentUser := c.MustGet("currentUser").(user.User)

	isEmailAvailable, errAvail := h.clientService.IsEmailClientAvailable(input.Email, currentUser.ID)
	if errAvail != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := helper.ApiResponse("Email checking failed", http.StatusUnprocessableEntity, "failed", nil, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	if !isEmailAvailable {
		data := gin.H{
			"status": "Failed to Add New Client!",
		}
		res := helper.ApiResponse("Email is already used!", http.StatusBadRequest, "failed", nil, data)
		c.JSON(http.StatusBadRequest, res)
	} else {
		_, errClient := h.clientService.AddClient(currentUser.ID, input)
		if errClient != nil {
			errors := helper.FormatValidationError(errClient)
			errorMessage := gin.H{"errors": errors}

			response := helper.ApiResponse("Failed to added New Client!", http.StatusUnprocessableEntity, "error", nil, errorMessage)
			c.JSON(http.StatusUnprocessableEntity, response)
			return
		}

		data := gin.H{
			"status": "Successfully added New Client!",
		}

		res := helper.ApiResponse("Succsessfully created New Client!", http.StatusCreated, "success", nil, data)

		c.JSON(http.StatusCreated, res)
	}
}

func (h *ClientHandler) AddClientsByCSV(c *gin.Context) {
	file, err := c.FormFile("csv_file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		res := helper.ApiResponse("Gagal Mengunggah File!", http.StatusBadRequest, "error", nil, data)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	// didapatkan dari JWT
	currentUser := c.MustGet("currentUser").(user.User)
	userId := currentUser.ID

	path := fmt.Sprintf("csv_file/%d-%s", userId, file.Filename)

	errImage := c.SaveUploadedFile(file, path)
	if errImage != nil {
		data := gin.H{"unggahan": false}
		res := helper.ApiResponse("Gagal Mengunggah File!", http.StatusBadRequest, "gagal", nil, data)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	// baca isi file yg sudah berhasil di upload
	lines, errRead := utils.ReadCsv(path)
	if errRead != nil {
		data := gin.H{"unggahan": true}
		res := helper.ApiResponse("Gagal Membaca File!", http.StatusBadRequest, "gagal", nil, data)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	clients := client.Mapping(lines)

	for _, v := range clients {
		isEmailAvailable, errAvail := h.clientService.IsEmailClientAvailable(v.Email, currentUser.ID)
		if errAvail != nil {
			os.Remove(path)
			errorMessage := gin.H{"errors": "Server error"}
			response := helper.ApiResponse("Email checking failed", http.StatusUnprocessableEntity, "gagal", nil, errorMessage)
			c.JSON(http.StatusUnprocessableEntity, response)
			return
		}
		if !isEmailAvailable {
			os.Remove(path)
			data := gin.H{
				"status": "Gagal Membuat Akun Baru!",
			}
			res := helper.ApiResponse("Email "+v.Email+" sudah digunakan!", http.StatusBadRequest, "gagal", nil, data)
			c.JSON(http.StatusBadRequest, res)
			return
		}
	}

	for _, v := range clients {
		_, errClient := h.clientService.AddClient(currentUser.ID, v)
		if errClient != nil {
			errors := helper.FormatValidationError(errClient)
			errorMessage := gin.H{"errors": errors}

			response := helper.ApiResponse("Menambahkan Client Baru Gagal!", http.StatusUnprocessableEntity, "error", nil, errorMessage)
			c.JSON(http.StatusUnprocessableEntity, response)
			return
		}

	}

	data := gin.H{
		"status": "Berhasil Menambahkan Client Baru!",
	}

	res := helper.ApiResponse("Berhasil Membuat client Baru!", http.StatusCreated, "berhasil", nil, data)
	os.Remove(path)
	c.JSON(http.StatusCreated, res)
}

func (h *ClientHandler) GetClients(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	currentUser := c.MustGet("currentUser").(user.User)

	name := c.Query("name")

	clients, total, perPage, err := h.clientService.GetAll(name, currentUser.ID, page)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Failed to get Clients Data!", http.StatusBadRequest, "error", nil, errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var lastPage float64
	if total%5 >= 1 && total%5 <= 4 {
		lastPage = float64(total/perPage + 1)
	} else {
		lastPage = float64(total / perPage)
	}

	info := gin.H{
		"total":     total,
		"page":      page,
		"last_page": lastPage,
	}
	formatter := client.FormatClients(clients)
	res := helper.ApiResponse("Succesfully to get Client Data!", http.StatusOK, "success", info, formatter)

	c.JSON(http.StatusOK, res)
}

func (h *ClientHandler) UpdateClient(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// cek apakah yg akses adalah admin
	currentUser := c.MustGet("currentUser").(user.User)

	var input client.InputUpdate
	err := c.ShouldBindJSON(&input)
	if err != nil {
		res := helper.ApiResponse("Update Data Has Been Failed", http.StatusUnprocessableEntity, "failed", nil, err)

		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	_, errUpdate := h.clientService.UpdateClient(input, currentUser.ID, id)
	if errUpdate != nil {
		errorMessage := gin.H{"errors": "Client Data does'nt exist"}
		res := helper.ApiResponse("Update Data Has Been Failed", http.StatusUnprocessableEntity, "failed", nil, errorMessage)

		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	data := gin.H{"is_uploaded": true}
	res := helper.ApiResponse("Update Data Has Been Success", http.StatusCreated, "success", nil, data)

	c.JSON(http.StatusCreated, res)

}

func (h *ClientHandler) DeleteClient(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// cek apakah yg akses adalah admin
	currentUser := c.MustGet("currentUser").(user.User)

	client, err := h.clientService.GetByID(id)
	if err != nil {
		res := helper.ApiResponse("Item Not Found", http.StatusBadRequest, "failed", nil, err)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	if client.ID == 0 {
		res := helper.ApiResponse("User Not Found", http.StatusBadRequest, "failed", nil, err)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	if currentUser.ID != client.UserID {
		res := helper.ApiResponse("Failed to Delete Client", http.StatusBadRequest, "failed", nil, errors.New("you don't have access"))

		c.JSON(http.StatusBadRequest, res)
		return
	}

	_, errDel := h.clientService.DeleteClient(client.ID)
	if errDel != nil {
		res := helper.ApiResponse("User Not Found", http.StatusBadRequest, "failed", nil, errDel)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	cekItem, errCek := h.clientService.GetByID(id)
	if errCek != nil {
		res := helper.ApiResponse("Any Error", http.StatusBadRequest, "failed", nil, errCek)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	if cekItem.ID == 0 {
		res := helper.ApiResponse("Successfuly Delete Item", http.StatusOK, "success", nil, nil)

		c.JSON(http.StatusOK, res)
		return
	}

	data := gin.H{"is_deleted": true}
	res := helper.ApiResponse("Any Error", http.StatusBadRequest, "failed", nil, data)

	c.JSON(http.StatusCreated, res)
}
