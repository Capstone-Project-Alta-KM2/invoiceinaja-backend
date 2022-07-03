package helper

type OverallData struct {
	Paid     int `json:"paid"`
	Unpaid   int `json:"unpaid"`
	Customer int `json:"customer"`
}

func FormatOverall(paid, unpaid, customer int) OverallData {
	formatter := OverallData{
		Paid:     paid,
		Unpaid:   unpaid,
		Customer: customer,
	}

	return formatter
}
