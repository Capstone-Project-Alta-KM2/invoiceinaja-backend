package handler

import (
	"fmt"
	"net/http"
	"os"

	"invoiceinaja/auth"
	"invoiceinaja/domain/client"
	"invoiceinaja/domain/invoice"
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

	// for _, v := range invoices {
	// 	fmt.Println(v.ID)
	// }

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
	}

	// ipn := invoice.InputAddInvoice{
	// 	Invoice: invoice.InvoiceData{
	// 		ClientID:    49,
	// 		TotalAmount: 1111110,
	// 		InvoiceDate: "12-12-2022",
	// 		InvoiceDue:  "12-12-2022",
	// 	},
	// 	DetailInvoice: []invoice.DetailInvoiceData{
	// 		{
	// 			ItemName: "A",
	// 			Price:    1000,
	// 			Quantity: 1,
	// 		},
	// 	},
	// }

	// for _, z := range input {

	//}

	data := gin.H{"is_recorded": true}
	res := helper.ApiResponse("Invoices Has Been Created", http.StatusCreated, "success", nil, data)
	fmt.Println(invoices[0].ID)
	c.JSON(http.StatusCreated, res)

	// var input []invoice.InputAddInvoice

	// for _, v := range invoices {
	// 	totalAmount := 0
	// 	var dtl []invoice.DetailInvoiceData

	// 	for _, x := range v.Items {
	// 		totalAmount += x.Price * x.Quantity
	// 	}

	// 	input = append(input, invoice.InputAddInvoice{
	// 		Invoice: invoice.InvoiceData{
	// 			ClientID:    v.ID,
	// 			TotalAmount: totalAmount,
	// 			InvoiceDate: v.InvoiceDate,
	// 			InvoiceDue:  v.InvoiceDue,
	// 		},
	// 		DetailInvoice: dtl,
	// 	})
	// }

	// for _, v := range input {
	// 	// record data invoice
	// 	newInvoice, errOrder := h.invoiceService.AddInvoice(v)
	// 	if errOrder != nil {
	// 		res := helper.ApiResponse("Add Invoice Has Been Failed", http.StatusBadRequest, "failed", nil, errOrder)

	// 		c.JSON(http.StatusBadRequest, res)
	// 	}

	// 	// record data detail order
	// 	_, errDetails := h.invoiceService.SaveDetail(newInvoice.ID, v)
	// 	if errDetails != nil {
	// 		res := helper.ApiResponse("Save Data Has Been Failed", http.StatusBadRequest, "failed", nil, errDetails)

	// 		c.JSON(http.StatusBadRequest, res)
	// 	}

	// _, errData := h.invoiceService.SendMailInvoice(newInvoice.ID, currentUser, newInvoice.Client)
	// if errData != nil {
	// 	res := helper.ApiResponse("Send Invoice Has Been Failed", http.StatusBadRequest, "failed", nil, errOrder)

	// 	c.JSON(http.StatusBadRequest, res)
	// }
	//}

	// data := gin.H{"is_recorded": true}
	// res := helper.ApiResponse("Order Has Been Created", http.StatusCreated, "success", nil, data)

	// c.JSON(http.StatusCreated, res)
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
