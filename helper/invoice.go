package helper

import (
	"fmt"
	"time"
)

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

type MonthReport struct {
	Paid   int `json:"paid"`
	Unpaid int `json:"unpaid"`
}

type Month struct {
	Jan, Feb, Mar, Apr, May, Jun, Jul, Agt, Sep, Okt, Nov, Des MonthReport
}

func ConvStingDate(date string) time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println(err)
	}

	return t
}
