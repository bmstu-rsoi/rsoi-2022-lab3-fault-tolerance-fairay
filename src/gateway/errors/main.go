package errors

import (
	"errors"
)

var (
	ForbiddenTicket    = errors.New("Forbidden ticket for this user")
	BonusUnavailable   = errors.New("Bonus Service")
	FlightUnavailable  = errors.New("Flight Service")
	TicketsUnavailable = errors.New("Tickets Service")
)
