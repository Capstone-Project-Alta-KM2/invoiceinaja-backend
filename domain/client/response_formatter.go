package client

type ClientFormatter struct {
	Id       int    `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	City     string `json:"city"`
	ZipCode  string `json:"zip_code"`
	Company  string `json:"company"`
}

func FormatClient(client Client) ClientFormatter {
	formatter := ClientFormatter{
		Id:       client.ID,
		Fullname: client.Fullname,
		Email:    client.Email,
		Address:  client.Address,
		City:     client.City,
		ZipCode:  client.ZipCode,
		Company:  client.Company,
	}

	return formatter
}

func FormatClients(clients []Client) []ClientFormatter {
	if len(clients) == 0 {
		return []ClientFormatter{}
	}

	var clientsFormatter []ClientFormatter

	for _, client := range clients {
		formatter := FormatClient(client)
		clientsFormatter = append(clientsFormatter, formatter)
	}

	return clientsFormatter
}
