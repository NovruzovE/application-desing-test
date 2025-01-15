package order

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/NovruzovE/application-design-test/internal/core/entity"
	"github.com/NovruzovE/application-design-test/internal/core/usecase/order"
)

type OrderController struct {
	orderService *order.OrderUseCase
	log          *slog.Logger
}

func New(s *order.OrderUseCase, logger *slog.Logger) *OrderController {
	return &OrderController{
		orderService: s,
		log:          logger,
	}
}

func (o *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq orderRequest

	err := json.NewDecoder(r.Body).Decode(&orderReq)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err = orderReq.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	err = o.orderService.CreateOrder(ctx, entity.Order{
		HotelID:   orderReq.HotelID,
		RoomID:    orderReq.RoomID,
		UserEmail: orderReq.UserEmail,
		From:      orderReq.From,
		To:        orderReq.To,
	})
	if err != nil {
		switch {
		case errors.Is(err, order.ErrNoAvailableRooms):
			orderResp := orderResponse{
				Message: "No available rooms for selected dates",
				Status:  statusError,
			}
			resp, err := json.Marshal(orderResp)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(resp)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		case errors.Is(err, order.ErrCannotBookRoom):
			orderResp := orderResponse{
				Message: "Cannot book room",
				Status:  statusError,
			}
			resp, err := json.Marshal(orderResp)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(resp)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	orderResp := orderResponse{
		Message: "The room has been successfully booked",
		Status:  statusSuccess,
	}
	resp, err := json.Marshal(orderResp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
