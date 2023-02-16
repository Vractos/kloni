package announcement

import "fmt"

type AnnouncementError struct {
	Message        string
	AnnouncementID string
	IsAbleToRetry  bool
	Sku            string
}

func (a *AnnouncementError) Error() string {
	if a.AnnouncementID != "" {
		return fmt.Sprintf("Message: %s - Announcement ID: %s", a.Message, a.AnnouncementID)
	} else if a.Sku != "" {
		return fmt.Sprintf("Message: %s - Announcement SKU: %s", a.Message, a.Sku)
	}

	return fmt.Sprintf("Message: %s", a.Message)
}
