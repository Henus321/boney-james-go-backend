package shop

import (
	"context"
	"fmt"
	"github.com/Henus321/boney-james-go-backend/pkg/logging"
	"github.com/Henus321/boney-james-go-backend/pkg/utils"
	//"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	//shopListURL   = "/shop"
	shopURL = "/shop/:id"
	//shopOptionURL = "/shop-option"
)

type Handler interface {
	Register(router *httprouter.Router)
}

type handler struct {
	shopService Service
	logger      *logging.Logger
}

func NewHandler(shopService Service, logger *logging.Logger) Handler {
	return &handler{
		shopService: shopService,
		logger:      logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	//router.GET(coatListURL, h.GetAllCoats)
	router.GET(shopURL, h.GetShopByID)
}

//func (h *handler) GetAllCoats(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
//	const op = "coat.handler.GetAllCoats"
//
//	coats, err := h.coatService.GetAllCoats(context.Background())
//	if err != nil {
//		h.logger.Error(op, err)
//		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
//		return
//	}
//
//	// TODO error handle
//	utils.WriteJSON(w, http.StatusOK, coats)
//	h.logger.Infof("%s: success", op)
//}

func (h *handler) GetShopByID(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	const op = "coat.handler.GetShopByID"

	coat, err := h.shopService.GetShopByID(context.Background(), params.ByName("id"))

	if err != nil {
		h.logger.Error(op, err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
		return
	}

	// TODO error handle
	utils.WriteJSON(w, http.StatusOK, coat)
	h.logger.Infof("%s: success", op)
}
