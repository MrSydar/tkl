package client

import "fmt"

type CustomerNotFoundError struct {
	Text string
}

func (e *CustomerNotFoundError) Error() string {
	return fmt.Sprintf(e.Text)
}
