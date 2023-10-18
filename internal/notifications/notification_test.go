package notifications

import "testing"

func TestSendSMS(t *testing.T) {
	// Create a new instance of NotificationService.
	ns := &NotificationService{}

	// Define test input.
	to := "+1234567890"
	message := "Test message"

	// Call the SendSMS function.
	err := ns.SendSMS(to, message)

	// Check if there was an error.
	if err != nil {
		t.Errorf("Error sending SMS: %v", err)
	}
}

// func TestSendNotification(t *testing.T) {
// 	// Create a new instance of NotificationService.
// 	ns := &NotificationService{}

// 	// Define test input.
// 	user := User{
// 		ID:       1,
// 		Username: "testuser",
// 		Email:    "testuser@example.com",
// 		Phone:    "+1234567890",
// 	}

// 	notification := Notification{
// 		User:    user,
// 		Message: "Test notification",
// 	}

// 	// Call the SendNotification function.
// 	//ns.SendNotification(notification)
// }
