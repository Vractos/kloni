package order

import (
	"fmt"

	"github.com/Vractos/kloni/usecases/common"
)

type OrderError struct {
	Message            string
	AnnouncementsError []common.MeliAnnouncement
}

func (o *OrderError) Error() string {
	return fmt.Sprintf("%s", o.Message)
}
