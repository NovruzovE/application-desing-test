package order

import (
	"errors"
	"time"
)

type orderRequest struct {
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

func (r *orderRequest) Validate() error {
	var errs []error
	if r.HotelID == "" {
		errs = append(errs, errors.New("hotel_id is required"))
	}
	if r.RoomID == "" {
		errs = append(errs, errors.New("room_id is required"))
	}
	if r.UserEmail == "" {
		errs = append(errs, errors.New("email is required"))
	}
	if r.From.IsZero() {
		errs = append(errs, errors.New("start date is required"))
	}
	if r.To.IsZero() {
		errs = append(errs, errors.New("end date is required"))
	}
	if r.From.After(r.To) {
		errs = append(errs, errors.New("start date should be before end date"))
	}

	return errors.Join(errs...)
}

type status string

const (
	statusSuccess status = "success"
	statusError   status = "error"
)

type orderResponse struct {
	Status  status `json:"status"`
	Message string `json:"message"`
}
