package payment

type InputCreateTansaction struct {
	InvoiceID   int `json:"invoice_id" binding:"required"`
	TotalAmount int `json:"total_amount" binding:"required"`
}
