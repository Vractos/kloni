package order

import (
	"fmt"

	"github.com/Vractos/dolly/usecases/common"
)

type OrderError struct {
	Message            string
	AnnouncementsError []common.OrderItem
}

func (o *OrderError) Error() string {
	return fmt.Sprintf("%s", o.Message)
}
