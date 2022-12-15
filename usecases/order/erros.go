package order

import "fmt"

type OrderError struct {
	Message            string
	AnnouncementsError []struct {
		id       string
		quantity int
	}
}

func (o *OrderError) Error() string {
	return fmt.Sprintf("Message: %s", o.Message)
}
