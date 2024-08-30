package policy

import (
	"entrust.com/iat/obac-idaas-server/pkg/models"
	"entrust.com/iat/obac-idaas-server/pkg/repository"
	"fmt"
)

type PolicyHandlerService interface {
	Edit(models.Policy) error
	InitPolicyContractIdsUris(smartContracts []models.SCId) error
	GetPolicy() (*models.Policy, error)
	GetPolicySCInfo(policySCs []models.SCId) ([]models.SCInfo, error)
	InitializeContracts() error
}

func NewPolicyHandlerService(repo repository.Repository, bcRepo repository.BlockchainRepository) PolicyHandlerService {
	return &policyHandlerServiceImpl{
		repo:   repo,
		bcRepo: bcRepo,
	}
}

type policyHandlerServiceImpl struct {
	repo   repository.Repository
	bcRepo repository.BlockchainRepository
}

func (h *policyHandlerServiceImpl) InitializeContracts() error {
	policy, err := h.GetPolicy()
	if err != nil {
		return err
	}
	policyContracts := policy.GetSCs()
	err = h.bcRepo.UpdateContracts(policyContracts)
	if err != nil {
		return fmt.Errorf("error updating contracts: %s", err.Error())
	}
	err = h.InitPolicyContractIdsUris(policyContracts)
	return err
}

func (h *policyHandlerServiceImpl) Edit(policy models.Policy) error {
	return h.repo.WritePolicy(policy)
}

func (h *policyHandlerServiceImpl) GetPolicy() (*models.Policy, error) {
	return h.repo.GetPolicy()
}

func (h *policyHandlerServiceImpl) GetPolicySCInfo(policySCs []models.SCId) ([]models.SCInfo, error) {
	contractsInfo := make([]models.SCInfo, len(policySCs))
	scNames, err := h.bcRepo.GetSCNameBatch(policySCs...)
	if err != nil {
		return nil, err
	}
	for i, sc := range policySCs {
		name, ok := scNames[sc]
		if !ok {
			name = sc.Address.Hex()
		}
		contractsInfo[i] = models.SCInfo{SCId: sc, Name: name}
	}
	return contractsInfo, nil
}

func (h *policyHandlerServiceImpl) InitPolicyContractIdsUris(smartContracts []models.SCId) error {
	err := h.bcRepo.UpdateTokenUriMap(smartContracts...) //Todo: only update the tokenUriMap for the new contracts in the policy
	if err != nil {
		return err
	}
	//for _, smartContract := range smartContracts {
	//	go h.bcRepo.StartListeningToSCEventLogs(smartContract)
	//}
	return nil
}
