package notifications

import "fmt"

type Notification struct {
	User    User
	Message string
}

// NotificationService provides methods for sending notifications.
type NotificationService struct{}

func SendNotification(ns NotificationService, notification Notification) {
	emailSubject := "Notification"
	emailBody := notification.Message

	err := ns.SendEmail(notification.User.Email, emailSubject, emailBody)
	if err != nil {
		fmt.Println(err)
		return
		// Handle error
	}

	// smsMessage := fmt.Sprintf("Notification: %s", notification.Message)
	// err = ns.SendSMS(notification.User.Phone, smsMessage)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// 	// Handle error
	// }
}
