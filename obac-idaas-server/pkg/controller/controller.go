package controller

import (
	"entrust.com/iat/obac-idaas-server/pkg/config"
	Errors "entrust.com/iat/obac-idaas-server/pkg/error"
	"entrust.com/iat/obac-idaas-server/pkg/logging"
	"entrust.com/iat/obac-idaas-server/pkg/repository"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Manager struct {
	FIFAHandlerController
	OIDCHandlerController
	PolicyHandlerController
}

func NewControllerManager(appUriConfig, idpUriConfig config.UriConfig, repo repository.Repository, bcRepo repository.BlockchainRepository, cadRepo repository.ContentAddressableRepository) Manager {
	fifaHandlerController := NewFIFAHandlerController(appUriConfig, repo)
	oidcHandlerController := NewOIDCHandlerController(appUriConfig, idpUriConfig, repo, bcRepo, cadRepo)
	policyHandlerController := NewPolicyHandlerController(appUriConfig, repo, bcRepo)
	return Manager{
		fifaHandlerController,
		oidcHandlerController,
		policyHandlerController,
	}
}

func HandleError(request *http.Request, response http.ResponseWriter, err error) {
	event := logging.GetFunctionName()
	errorType, errorDesc, info := GetErrorDetails(err)
	typedErrorController(event, info, request, response, errorType, err, errorDesc)
}

func HandleErrorWithFields(request *http.Request, response http.ResponseWriter, err error, fields *log.Fields) {
	event := logging.GetFunctionName()
	errorType, errorDesc, info := GetErrorDetails(err)
	typedFieldErrorController(event, info, request, response, errorType, err, errorDesc, fields)
}

func GetErrorDetails(err error) (errorType int, errorDesc string, info bool) {
	if errors.Is(err, Errors.BadRequestError{}) {
		errorType = http.StatusBadRequest
		errorDesc = err.(Errors.BadRequestError).GetDescription()
		info = true
	} else if errors.Is(err, Errors.NotAuthorizedError{}) {
		errorType = http.StatusUnauthorized
		errorDesc = err.(Errors.NotAuthorizedError).GetDescription()
	} else if errors.Is(err, Errors.NotFoundError{}) {
		errorType = http.StatusNotFound
		errorDesc = err.(Errors.NotFoundError).GetDescription()
		info = true
	} else {
		errorType = http.StatusInternalServerError
		errorDesc = "Internal server error"
	}
	return
}

func typedErrorController(event string, info bool, request *http.Request, response http.ResponseWriter, errorType int, err error, description string) {
	fields := &log.Fields{
		"httpStatus": errorType,
		"clientIP":   request.RemoteAddr,
	}
	logError(event, info, fields, err)
	response.WriteHeader(errorType)
	_, err = response.Write([]byte(description))
	if err != nil {
		log.Errorf("Error while writing the error description at error controller: " + err.Error())
	}
}

func typedFieldErrorController(event string, info bool, request *http.Request, response http.ResponseWriter, errorType int, err error, description string, fields *log.Fields) {
	if *fields == nil {
		fields = &log.Fields{}
	}
	(*fields)["httpStatus"] = errorType
	(*fields)["clientIP"] = request.RemoteAddr
	logError(event, info, fields, err)
	response.WriteHeader(errorType)
	_, err = response.Write([]byte(description))
	if err != nil {
		log.Errorf("Error while writing the error description at error controller: " + err.Error())
	}
}

func logError(event string, info bool, fields *log.Fields, err error) {
	if info {
		logging.InfoErrorEvent(event, *fields, err)
	} else {
		logging.ErrorEvent(event, *fields, err)
	}
}
