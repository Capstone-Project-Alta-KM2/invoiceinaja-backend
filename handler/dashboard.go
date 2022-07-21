package handler

import (
	"invoiceinaja/auth"
	"invoiceinaja/domain/client"
	"invoiceinaja/domain/invoice"
	"invoiceinaja/domain/user"
	"invoiceinaja/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	invoiceService invoice.IService
	clientService  client.IService
	authService    auth.Service
}

func NewDashboardHandler(invoiceService invoice.IService, clientService client.IService, authService auth.Service) *DashboardHandler {
	return &DashboardHandler{invoiceService, clientService, authService}
}

func (h *DashboardHandler) GetDataOverall(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)

	totalPaid, err := h.invoiceService.GetSumStatus(currentUser.ID)
	if err != nil {
		res := helper.ApiResponse("Any Error", http.StatusBadRequest, "failed", nil, err)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	totalCustomer := h.clientService.TotalCustomer(currentUser.ID)

	formatter := helper.FormatOverall(totalPaid["paid"], totalPaid["unpaid"], totalCustomer)
	res := helper.ApiResponse("invoices", http.StatusOK, "success", nil, formatter)
	c.JSON(http.StatusOK, res)
}

func (h *DashboardHandler) GetDataGraphic(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)

	invoices, err := h.invoiceService.GetInvoices(currentUser.ID)
	if err != nil {
		res := helper.ApiResponse("Any Error", http.StatusBadRequest, "failed", nil, err)

		c.JSON(http.StatusBadRequest, res)
		return
	}

	h.invoiceService.CheckDate(invoices)

	///////////////////////////////////
	s := invoice.GraphicInvoice(invoices)
	///////////////////////////////////

	//formatter := invoice.FormatInvoices(invoices)
	res := helper.ApiResponse("invoices", http.StatusOK, "success", nil, s)
	c.JSON(http.StatusOK, res)
}
