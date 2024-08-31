package coat

import (
	"context"
	"fmt"
	"github.com/Henus321/boney-james-go-backend/pkg/logging"
	"github.com/Henus321/boney-james-go-backend/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	coatListURL = "/coat"
	coatURL     = "/coat/:id"
)

type Handler interface {
	Register(router *httprouter.Router)
}

type handler struct {
	coatService Service
	logger      *logging.Logger
}

func NewHandler(coatService Service, logger *logging.Logger) Handler {
	return &handler{
		coatService: coatService,
		logger:      logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(coatListURL, h.GetAllCoats)
	router.POST(coatListURL, h.CreateCoat)
	router.GET(coatURL, h.GetCoatByID)
	router.DELETE(coatURL, h.DeleteCoat)
}

func (h *handler) GetAllCoats(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	const op = "coat.handler.GetAllCoats"

	coats, err := h.coatService.GetAllCoats(context.Background())
	if err != nil {
		h.logger.Error(op, err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
		return
	}

	// TODO error handle
	utils.WriteJSON(w, http.StatusOK, coats)
	h.logger.Infof("%s: success", op)
}

func (h *handler) GetCoatByID(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	const op = "coat.handler.GetCoatByID"

	coat, err := h.coatService.GetCoatByID(context.Background(), params.ByName("id"))

	if err != nil {
		h.logger.Error(op, err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
		return
	}

	// TODO error handle
	utils.WriteJSON(w, http.StatusOK, coat)
	h.logger.Infof("%s: success", op)
}

func (h *handler) CreateCoat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	const op = "coat.handler.CreateCoat"

	var input CreateCoatInput
	if err := utils.ParseJSON(r, &input); err != nil {
		h.logger.Error(op, err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("%s: %w", op, err))
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %s: %v", op, errors))
		return
	}

	err := h.coatService.CreateCoat(context.Background(), input)
	if err != nil {
		h.logger.Error(op, err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
		return
	}

	// TODO error handle
	utils.WriteJSON(w, http.StatusCreated, nil)
	h.logger.Infof("%s: success", op)
}

func (h *handler) DeleteCoat(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	const op = "coat.handler.DeleteCoat"

	err := h.coatService.DeleteCoat(context.Background(), params.ByName("id"))
	if err != nil {
		h.logger.Error(op, err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
		return
	}

	// TODO error handle
	utils.WriteJSON(w, http.StatusOK, nil)
	h.logger.Infof("%s: success", op)
}
