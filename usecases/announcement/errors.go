package announcement

import "fmt"

type AnnouncementError struct {
	Message        string
	AnnouncementID string
	IsAbleToRetry  bool
	Sku            string
}

func (a *AnnouncementError) Error() string {
	return fmt.Sprintf("Message: %s - Announcement ID/SKU: %s", a.Message, a.AnnouncementID)
}
