package order

import (
	"fmt"

	"github.com/Vractos/kloni/usecases/announcement"
)

type OrderError struct {
	Message            string
	AnnouncementsError []announcement.Announcements
}

func (o *OrderError) Error() string {
	return fmt.Sprintf("%s", o.Message)
}
