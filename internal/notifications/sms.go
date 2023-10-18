package notifications

import (
	"fmt"
	"net/http"
	"net/url"
)

// SendSMS sends an SMS with the provided details.
func (ns *NotificationService) SendSMS(to, message string) error {
	// Implement SMS sending logic here.
	// Use the provided 'to' and 'message' parameters to send an SMS.

	// Replace the following placeholders with your actual SMS service details.
	accountSid := "YOUR_ACCOUNT_SID"
	authToken := "YOUR_AUTH_TOKEN"
	from := "YOUR_PHONE_NUMBER" // This is the phone number provided by your SMS service.

	// Construct the SMS API URL.
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	// Set up the HTTP client.
	client := &http.Client{}

	// Prepare the form data.
	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", from)
	msgData.Set("Body", message)

	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.PostForm = msgData

	// Send the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Check the response status.
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("SMS sent successfully!")
		return nil
	}

	return fmt.Errorf("Failed to send SMS. Status code: %d", resp.StatusCode)
}
