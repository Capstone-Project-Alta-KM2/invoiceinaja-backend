package invoice

type InvoiceFormatter struct {
	Id      int    `json:"id"`
	Client  string `json:"client"`
	Date    string `json:"date"`
	PostDue string `json:"post_due"`
	Amount  int    `json:"Amount"`
	Status  string `json:"status"`
	Items   []ItemFormatter
}

type ItemFormatter struct {
	ItemName string `json:"item_name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

func FormatInvoice(Invoice Invoice) InvoiceFormatter {
	formatter := InvoiceFormatter{
		Id:      Invoice.ID,
		Client:  Invoice.Client.Fullname,
		Date:    Invoice.InvoiceDate,
		PostDue: Invoice.InvoiceDue,
		Amount:  Invoice.TotalAmount,
		Status:  Invoice.Status,
		Items:   FormatItem(Invoice),
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

func FormatItem(item Invoice) []ItemFormatter {
	var detail []ItemFormatter

	for _, v := range item.Items {
		detail = append(detail, ItemFormatter{
			ItemName: v.ItemName,
			Price:    v.Price,
			Quantity: v.Quantity,
		})
	}

	return detail
}
