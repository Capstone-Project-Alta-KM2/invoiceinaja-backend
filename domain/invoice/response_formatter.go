package invoice

type InvoiceFormatter struct {
	Id      int    `json:"id"`
	Client  string `json:"client"`
	Date    string `json:"date"`
	PostDue string `json:"post_due"`
	Amount  int    `json:"Amount"`
	Status  string `json:"status"`
}

func FormatInvoice(Invoice Invoice) InvoiceFormatter {
	formatter := InvoiceFormatter{
		Id:      Invoice.ID,
		Client:  Invoice.Client.Fullname,
		Date:    Invoice.InvoiceDate,
		PostDue: Invoice.InvoiceDue,
		Amount:  Invoice.TotalAmount,
		Status:  Invoice.Status,
	}

	return formatter
}

func FormatInvoices(invoice []Invoice) []InvoiceFormatter {
	if len(invoice) == 0 {
		return []InvoiceFormatter{}
	}

	var invoicesFormatter []InvoiceFormatter

	for _, data := range invoice {
		formatter := FormatInvoice(data)
		invoicesFormatter = append(invoicesFormatter, formatter)
	}

	return invoicesFormatter
}
