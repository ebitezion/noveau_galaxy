package data

import (
	"fmt"
)

// Insert method for inserting a new record in the  table.
func (m Models) Insert(Query string, Args []interface{}) error {
	fmt.Println("function reached")
	// Create a context with a 3-second timeout.
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// Execute the database query using the provided context and arguments.
	err := m.AccountModel.DB.Ping()
	fmt.Println(err)
	// _, err := m.AccountModel.DB.ExecContext(ctx, Query, Args...)
	// if err != nil {
	// 	return err
	// }
	fmt.Println("function broke")
	return nil
}
