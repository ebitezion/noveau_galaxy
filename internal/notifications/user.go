package notifications

type User struct {
	ID       int
	Username string
	Email    string
	Phone    string
}

func GetUserByID(userID int) *User {
	return &User{
		ID:       userID,
		Username: "john_doe",
		Email:    "john@example.com",
		Phone:    "+1234567890",
	}
}
