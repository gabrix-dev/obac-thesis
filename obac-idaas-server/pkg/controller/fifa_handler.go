package controller

import (
	"encoding/json"
	"entrust.com/iat/obac-idaas-server/pkg/config"
	error2 "entrust.com/iat/obac-idaas-server/pkg/error"
	"entrust.com/iat/obac-idaas-server/pkg/models"
	"entrust.com/iat/obac-idaas-server/pkg/repository"
	service "entrust.com/iat/obac-idaas-server/pkg/service/fifa"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type FIFAHandlerController interface {
	GetBookingInfo(response http.ResponseWriter, request *http.Request)
	InsertBookingInfo(response http.ResponseWriter, request *http.Request)
}

type fifaHandlerControllerImpl struct {
	service service.FIFAHandlerService
	appUri  string
}

func NewFIFAHandlerController(appUriConfig config.UriConfig, repository repository.Repository) FIFAHandlerController {
	return &fifaHandlerControllerImpl{
		service: service.NewFIFAHandlerService(repository),
		appUri:  appUriConfig.Protocol + "://" + appUriConfig.Host + ":" + appUriConfig.Port + appUriConfig.Path,
	}
}

func (p fifaHandlerControllerImpl) GetBookingInfo(response http.ResponseWriter, request *http.Request) {
	token := request.Header.Get("Authorization")
	if token == "" {
		HandleError(request, response, error2.NewBadRequestError(errors.New("token missing"), "authorization failed"))
		return
	}
	if err := p.service.ValidateToken(token); err != nil {
		HandleErrorWithFields(request, response, error2.NewNotAuthorizedError(err, "error validating the token"), &logrus.Fields{"token": token})
		return
	}
	claims, err := p.service.GetTokenClaims(token)
	if err != nil {
		HandleError(request, response, fmt.Errorf("error retreiving token's claims: %s", err.Error()))
		return
	}
	bookingInfo, err := p.service.GetBookingInfo((*claims)["id"].(float64))
	if err != nil {
		HandleError(request, response, error2.NewNotFoundError(err, "the user has not booked a seat yet"))
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(bookingInfo)
}

func (p fifaHandlerControllerImpl) InsertBookingInfo(response http.ResponseWriter, request *http.Request) {
	var booking models.BookingInfo
	token := request.Header.Get("Authorization")
	if token == "" {
		HandleError(request, response, error2.NewBadRequestError(errors.New("token missing"), "authorization failed"))
		return
	}
	if err := p.service.ValidateToken(token); err != nil {
		HandleErrorWithFields(request, response, error2.NewNotAuthorizedError(err, "error validating the token"), &logrus.Fields{"token": token})
		return
	}
	claims, err := p.service.GetTokenClaims(token)
	if err != nil {
		HandleError(request, response, fmt.Errorf("error retreiving token's claims: %s", err.Error()))
		return
	}
	err = json.NewDecoder(request.Body).Decode(&booking)
	if err != nil {
		HandleError(request, response, fmt.Errorf("error encoding the booking info: %s", err.Error()))
		return
	}

	if err := booking.ValidateBooking((*claims)["ticket_type"].(string), (*claims)["id"].(float64)); err != nil {
		http.Error(response, "User not authorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	//We add the booking to the db (associated with the token ID)
	p.service.WriteBookingInfo((*claims)["id"].(float64), booking)
	response.WriteHeader(http.StatusOK)
}
