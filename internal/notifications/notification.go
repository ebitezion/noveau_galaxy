package notifications

import "fmt"

type Notification struct {
	User    User
	Message string
}

// NotificationService provides methods for sending notifications.
type NotificationService struct{}

func SendNotification(ns NotificationService, notification Notification) {
	emailSubject := "Verification Notification"
	emailBody := notification.Message

	err := ns.SendEmail(notification.User.Email, emailSubject, emailBody)
	if err != nil {
		fmt.Println(err)
		return
		// Handle error
	}

}
