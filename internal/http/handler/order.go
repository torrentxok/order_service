package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/torrentxok/order_service/internal/service"
	"go.uber.org/zap"
)

type OrderHandler struct {
	service *service.OrderService
	logger  *zap.Logger
}

func NewOrderHandler(service *service.OrderService, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{
		service: service,
		logger:  logger,
	}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderUID := chi.URLParam(r, "order_uid")
	if orderUID == "" {
		http.Error(w, "order_uid is required", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(ctx, orderUID)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		h.logger.Error("failed to get order", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
