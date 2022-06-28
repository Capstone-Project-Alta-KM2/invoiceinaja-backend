package handler

import (
	"fmt"
	"net/http"

	"invoiceinaja/auth"
	"invoiceinaja/domain/client"
	"invoiceinaja/domain/invoice"
	"invoiceinaja/domain/user"
	"invoiceinaja/helper"

	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
	invoiceService invoice.IService
	clientService  client.IService
	authService    auth.Service
}

func NewInvoiceHandler(invoiceService invoice.IService, clientService client.IService, authService auth.Service) *InvoiceHandler {
	return &InvoiceHandler{invoiceService, clientService, authService}
}

func (h *InvoiceHandler) AddInvoice(c *gin.Context) {
	// didapatkan dari JWT
	currentUser := c.MustGet("currentUser").(user.User)

	var input invoice.InputAddInvoice
	// tangkap input body
	err := c.ShouldBindJSON(&input)
	if err != nil {
		res := helper.ApiResponse("New Data Has Been Failed1", http.StatusUnprocessableEntity, "failed", nil, err)

		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	// cek id client relate dengan user'
	client, errClient := h.clientService.GetByID(input.Invoice.ClientID)
	if errClient != nil {
		res := helper.ApiResponse("Any Error", http.StatusBadRequest, "failed", nil, errClient)

		c.JSON(http.StatusBadRequest, res)
		return
	}
	if client.UserID == 0 {
		res := helper.ApiResponse("Client not found!", http.StatusUnprocessableEntity, "failed", nil, errClient)

		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}
	if client.UserID != currentUser.ID {
		res := helper.ApiResponse("This is not your client!", http.StatusUnprocessableEntity, "failed", nil, errClient)

		c.JSON(http.StatusUnprocessableEntity, res)
		return
	} else {
		// record data invoice
		newInvoice, errOrder := h.invoiceService.AddInvoice(input)
		if errOrder != nil {
			res := helper.ApiResponse("Invoice Has Been Failed", http.StatusBadRequest, "failed", nil, errOrder)

			c.JSON(http.StatusBadRequest, res)
		}

		// record data detail order
		_, errDetails := h.invoiceService.SaveDetail(newInvoice.ID, input)
		if errDetails != nil {
			res := helper.ApiResponse("New Data Has Been Failed", http.StatusBadRequest, "failed", nil, errDetails)

			c.JSON(http.StatusBadRequest, res)
		}

		data := gin.H{"is_recorded": true}
		res := helper.ApiResponse("Order Has Been Created", http.StatusCreated, "success", nil, data)

		c.JSON(http.StatusCreated, res)
	}

}

func (h *InvoiceHandler) GenerateByCSV(c *gin.Context) {
	file, err := c.FormFile("csv_file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		res := helper.ApiResponse("Failed to Upload Image!", http.StatusBadRequest, "error", nil, data)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	// didapatkan dari JWT
	currentUser := c.MustGet("currentUser").(user.User)
	userId := currentUser.ID

	path := fmt.Sprintf("csv_file/%d-%s", userId, file.Filename)

	errImage := c.SaveUploadedFile(file, path)
	if errImage != nil {
		data := gin.H{"is_uploaded": false}
		res := helper.ApiResponse("Failed to Upload Image!", http.StatusBadRequest, "failed", nil, data)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	data := gin.H{"is_uploaded": true}
	res := helper.ApiResponse("Succesfully Uploaded Image!", http.StatusOK, "success", nil, data)

	c.JSON(http.StatusOK, res)
}

func (h *InvoiceHandler) GetInvoices(c *gin.Context) {
	// didapatkan dari JWT
	currentUser := c.MustGet("currentUser").(user.User)

	invoices, err := h.invoiceService.GetInvoices(currentUser.ID)
	if err != nil {
		res := helper.ApiResponse("Any Error", http.StatusBadRequest, "failed", nil, err)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	formatter := invoice.FormatInvoices(invoices)
	res := helper.ApiResponse("invoices", http.StatusCreated, "success", nil, formatter)
	c.JSON(http.StatusCreated, res)
}
