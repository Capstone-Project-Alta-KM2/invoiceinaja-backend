package payment

type InputCreateTansaction struct {
	InvoiceID   int `json:"invoice_id" binding:"required"`
	TotalAmount int `json:"total_amount" binding:"required"`
}

type InputTransactionNotif struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
