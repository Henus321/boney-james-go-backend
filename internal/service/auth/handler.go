package auth

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
	loginURL    = "/login"
	registerURL = "/register"
	userURL     = "/user/:id"
)

type Handler interface {
	Register(router *httprouter.Router)
}

type handler struct {
	authService Service
	logger      *logging.Logger
}

func NewHandler(authService Service, logger *logging.Logger) Handler {
	return &handler{
		authService: authService,
		logger:      logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(userURL, h.GetUserByID)
	router.POST(loginURL, h.LoginUser)
	router.POST(registerURL, h.RegisterUser)
}

func (h *handler) GetUserByID(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	const op = "auth.handler.GetUserByID"

	user, err := h.authService.GetUserByID(context.Background(), params.ByName("id"))

	if err != nil {
		h.logger.Error(op, err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
		return
	}

	// TODO error handle
	utils.WriteJSON(w, http.StatusOK, user)
	h.logger.Infof("%s: success", op)
}

func (h *handler) LoginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	const op = "auth.handler.LoginUser"

	var input UserLoginInput
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

	user, token, err := h.authService.LoginUser(context.Background(), &input)
	if err != nil {
		h.logger.Error(op, err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
		return
	}

	utils.SetCookie(w, token)

	// TODO error handle
	utils.WriteJSON(w, http.StatusOK, user)
	h.logger.Infof("%s: success", op)
}

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	const op = "auth.handler.RegisterUser"

	var input UserRegisterInput
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

	user, token, err := h.authService.RegisterUser(context.Background(), &input)
	if err != nil {
		h.logger.Error(op, err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", op, err))
		return
	}

	utils.SetCookie(w, token)

	// TODO error handle
	utils.WriteJSON(w, http.StatusOK, user)
	h.logger.Infof("%s: success", op)
}
