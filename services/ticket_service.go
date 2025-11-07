package ticket

import (
	"fmt"
	"os"

	"resty.dev/v3"
)

func GetTicketByID(ticketID string) ([]byte, error) {
	apiKey := os.Getenv("FRESHDESK_API_KEY")
	domain := os.Getenv("FRESHDESK_DOMAIN")

	client := resty.New()

	response, err := client.R().
		SetBasicAuth(apiKey, "X").
		Get(fmt.Sprintf("%s/api/v2/tickets/%s", domain, ticketID))

	if err != nil {
		return nil, err
	}

	return response.Bytes(), nil
}
