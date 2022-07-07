package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"invoiceinaja/auth"
	"invoiceinaja/domain/client"
	"invoiceinaja/domain/invoice"

	"invoiceinaja/domain/payment"
	"invoiceinaja/domain/user"
	"invoiceinaja/helper"
	"invoiceinaja/utils"

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
		res := helper.ApiResponse("Add New Data Has Been Failed", http.StatusUnprocessableEntity, "failed", nil, err)

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
			res := helper.ApiResponse("Add Invoice Has Been Failed", http.StatusBadRequest, "failed", nil, errOrder)

			c.JSON(http.StatusBadRequest, res)
		}

		// record data detail order
		_, errDetails := h.invoiceService.SaveDetail(newInvoice.ID, input)
		if errDetails != nil {
			res := helper.ApiResponse("Save Data Has Been Failed", http.StatusBadRequest, "failed", nil, errDetails)

			c.JSON(http.StatusBadRequest, res)
		}

		_, errData := h.invoiceService.SendMailInvoice(newInvoice.ID, currentUser, client)
		if errData != nil {
			res := helper.ApiResponse("Send Invoice Has Been Failed", http.StatusBadRequest, "failed", nil, errOrder)

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
		res := helper.ApiResponse("Failed to Upload File!", http.StatusBadRequest, "error", nil, data)

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
		res := helper.ApiResponse("Failed to Upload File!", http.StatusBadRequest, "failed", nil, data)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	lines, errRead := utils.ReadCsv(path)
	os.Remove(path)
	if errRead != nil {
		data := gin.H{"unggahan": true}
		res := helper.ApiResponse("Failed to Read File!", http.StatusBadRequest, "gagal", nil, data)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	invoices := invoice.Mapping(lines)
	var input []invoice.InputAddInvoice
	var idClient []int

	// cek email
	for _, v := range invoices {
		client, errClient := h.clientService.GetByEmail(v.Email, currentUser.ID)
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
		v.ID = client.ID
		idClient = append(idClient, client.ID)
		if client.UserID != currentUser.ID {
			res := helper.ApiResponse("This is not your client!", http.StatusUnprocessableEntity, "failed", nil, errClient)

			c.JSON(http.StatusUnprocessableEntity, res)
			return
		}
	}

	for i, x := range invoices {
		totalAmount := 0
		var dtl []invoice.DetailInvoiceData

		for _, total := range x.Items {
			totalAmount += total.Price * total.Quantity
		}
		for _, item := range x.Items {
			dtl = append(dtl, invoice.DetailInvoiceData(item))
		}

		input = append(input, invoice.InputAddInvoice{
			Invoice: invoice.InvoiceData{
				ClientID:    idClient[i],
				TotalAmount: totalAmount,
				InvoiceDate: x.InvoiceDate,
				InvoiceDue:  x.InvoiceDue,
			},
			DetailInvoice: dtl,
		})

		// record data invoice
		newInvoice, errOrder := h.invoiceService.AddInvoice(input[i])
		if errOrder != nil {
			res := helper.ApiResponse("Add Invoice Has Been Failed", http.StatusBadRequest, "failed", nil, errOrder)

			c.JSON(http.StatusBadRequest, res)
			return
		}

		// record data detail order
		_, errDetails := h.invoiceService.SaveDetail(newInvoice.ID, input[i])
		if errDetails != nil {
			res := helper.ApiResponse("Save Data Has Been Failed", http.StatusBadRequest, "failed", nil, errDetails)

			c.JSON(http.StatusBadRequest, res)
			return
		}

		client, _ := h.clientService.GetByID(input[i].Invoice.ClientID)

		_, errData := h.invoiceService.SendMailInvoice(newInvoice.ID, currentUser, client)
		if errData != nil {
			res := helper.ApiResponse("Send Invoice Has Been Failed", http.StatusBadRequest, "failed", nil, errOrder)

			c.JSON(http.StatusBadRequest, res)
			return
		}
	}

	data := gin.H{"is_recorded": true}
	res := helper.ApiResponse("Invoices Has Been Created", http.StatusCreated, "success", nil, data)
	c.JSON(http.StatusCreated, res)
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

func (h *InvoiceHandler) GetInvoicesByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// didapatkan dari JWT
	currentUser := c.MustGet("currentUser").(user.User)

	invoices, err := h.invoiceService.GetByID(id)
	if err != nil {
		res := helper.ApiResponse("Invoice Not Found", http.StatusNotFound, "failed", nil, err)

		c.JSON(http.StatusNotFound, res)
		return
	}
	if invoices.Client.ID != currentUser.ID {
		res := helper.ApiResponse("Invoice Not Found", http.StatusNotFound, "failed", nil, err)

		c.JSON(http.StatusNotFound, res)
		return
	}

	formatter := invoice.FormatInvoice(invoices)
	res := helper.ApiResponse("invoices", http.StatusOK, "success", nil, formatter)
	c.JSON(http.StatusOK, res)
}

func (h *InvoiceHandler) DeleteInvoice(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// cek apakah yg akses adalah user benar
	currentUser := c.MustGet("currentUser").(user.User)

	invoice, err := h.invoiceService.GetByID(id)
	if err != nil {
		res := helper.ApiResponse("Invoice Not Found", http.StatusBadRequest, "failed", nil, err)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	details, errDetail := h.invoiceService.GetDetailByID(invoice.ID)
	if errDetail != nil {
		res := helper.ApiResponse("Invoice Not Found", http.StatusBadRequest, "failed", nil, err)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	for _, v := range details {
		h.invoiceService.DeleteDetailInvoice(v)
	}

	if invoice.ID == 0 {
		res := helper.ApiResponse("Invoice Not Found", http.StatusBadRequest, "failed", nil, err)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	if currentUser.ID != invoice.Client.UserID {
		res := helper.ApiResponse("Failed to Delete Invoice", http.StatusBadRequest, "failed", nil, errors.New("you don't have access"))

		c.JSON(http.StatusBadRequest, res)
		return
	}

	_, errDel := h.invoiceService.DeleteInvoice(invoice.ID)
	if errDel != nil {
		res := helper.ApiResponse("Invoice Not Found", http.StatusBadRequest, "failed", nil, errDel)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	cekInvoice, errCek := h.invoiceService.GetByID(id)
	if errCek != nil {
		res := helper.ApiResponse("Failed to Delete Invoice", http.StatusBadRequest, "failed", nil, errCek)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	if cekInvoice.ID != 0 {
		res := helper.ApiResponse("Failed to Delete Invoice", http.StatusOK, "failed", nil, nil)

		c.JSON(http.StatusOK, res)
		return
	}

	data := gin.H{"is_deleted": true}
	res := helper.ApiResponse("Successfuly Delete Invoice", http.StatusBadRequest, "success", nil, data)

	c.JSON(http.StatusCreated, res)
}

func (h *InvoiceHandler) InvoicePay(c *gin.Context) {
	// get Invoice Id, total amount and
	var input payment.InputCreateTansaction

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errMessage := gin.H{"errors": errors}

		res := helper.ApiResponse("Failed to Create Transaction", http.StatusUnprocessableEntity, "failed", nil, errMessage)
		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	// currentUser := c.MustGet("currentUser").(user.User)

	invoice, errInv := h.invoiceService.GetByID(input.InvoiceID)
	if errInv != nil {
		res := helper.ApiResponse("Invoice Not Found", http.StatusBadRequest, "failed", nil, err)

		c.JSON(http.StatusBadRequest, res)
		return
	}
	// if invoice.Client.UserID != currentUser.ID {
	// 	res := helper.ApiResponse("Invoice Not Found", http.StatusBadRequest, "failed", nil, err)

	// 	c.JSON(http.StatusBadRequest, res)
	// 	return
	// }

	url, errTrans := h.invoiceService.PayInvoice(input, invoice.Client)
	if errTrans != nil {
		res := helper.ApiResponse("Failed to Create Transaction", http.StatusUnprocessableEntity, "failed", nil, errTrans)
		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	_, errUpdate := h.invoiceService.UpdateInvoice(invoice, url)
	if errUpdate != nil {
		res := helper.ApiResponse("Failed to Create Transaction", http.StatusUnprocessableEntity, "failed", nil, errTrans)
		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	data := gin.H{"payment_url": url}
	res := helper.ApiResponse("Successfuly Create Payment URL", http.StatusOK, "success", nil, data)

	c.JSON(http.StatusOK, res)
}
