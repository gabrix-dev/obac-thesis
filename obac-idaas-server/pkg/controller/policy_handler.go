package controller

import (
	"encoding/json"
	"entrust.com/iat/obac-idaas-server/pkg/config"
	"entrust.com/iat/obac-idaas-server/pkg/logging"
	"entrust.com/iat/obac-idaas-server/pkg/models"
	"entrust.com/iat/obac-idaas-server/pkg/repository"
	service "entrust.com/iat/obac-idaas-server/pkg/service/policy"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type PolicyHandlerController interface {
	Edit(response http.ResponseWriter, request *http.Request)
	GetPolicy(response http.ResponseWriter, request *http.Request)
	GetPolicyContracts(response http.ResponseWriter, request *http.Request)
}

type policyHandlerControllerImpl struct {
	service service.PolicyHandlerService
	appUri  string
}

func NewPolicyHandlerController(appUriConfig config.UriConfig, repository repository.Repository, bcRepository repository.BlockchainRepository) PolicyHandlerController {
	p := policyHandlerControllerImpl{
		service: service.NewPolicyHandlerService(repository, bcRepository),
		appUri:  appUriConfig.Protocol + "://" + appUriConfig.Host + ":" + appUriConfig.Port + appUriConfig.Path,
	}
	p.service.InitializeContracts()
	return &p
}

func (p policyHandlerControllerImpl) Edit(response http.ResponseWriter, request *http.Request) {
	var policy models.Policy
	err := json.NewDecoder(request.Body).Decode(&policy)
	if err != nil {
		HandleError(request, response, fmt.Errorf("error decoding the received policy to json: %s", err.Error()))
		return
	}
	err = p.service.Edit(policy)
	if err != nil {
		HandleError(request, response, fmt.Errorf("error writing the policy to the repo: %s", err.Error()))
		return
	}
	if err = p.service.InitializeContracts(); err != nil {
		HandleError(request, response, fmt.Errorf("error initializing the contracts: %s", err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	logging.InfoEvent(logging.EditHandled, logrus.Fields{"policy": policy})
}

func (p policyHandlerControllerImpl) GetPolicy(response http.ResponseWriter, request *http.Request) {
	policy, err := p.service.GetPolicy()
	response.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(response).Encode(policy); err != nil {
		HandleError(request, response, errors.New("Error encoding policy to http response"))
		return
	}
	logging.InfoEvent(logging.GetPolictHandled, nil)
}

func (p policyHandlerControllerImpl) GetPolicyContracts(response http.ResponseWriter, request *http.Request) {
	currentPolicy, err := p.service.GetPolicy()
	if err != nil {
		HandleError(request, response, fmt.Errorf("error retreiving the current policy: %s", err.Error()))
		return
	}
	currentPolicySCs := currentPolicy.GetSCs()
	policyScInfo, err := p.service.GetPolicySCInfo(currentPolicySCs)
	if err != nil {
		HandleErrorWithFields(request, response, fmt.Errorf("error getting the sc info: %s", err.Error()), &logrus.Fields{"smart_contracts": currentPolicySCs})
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(policyScInfo)
	logging.InfoEvent(logging.GetPolicyContractsHandled, nil)
}
