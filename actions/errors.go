package actions

import "fmt"

type NotFoundError struct {
	CustomerName string
	OrderID      int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %d not found", e.CustomerName, e.OrderID)
}

func NewNotFoundError(customerName string, orderID int) *NotFoundError {
	return &NotFoundError{
		CustomerName: customerName,
		OrderID:      orderID,
	}
}