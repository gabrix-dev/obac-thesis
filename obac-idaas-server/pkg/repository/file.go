package repository

import (
	"encoding/json"
	"entrust.com/iat/obac-idaas-server/pkg/models"
	"github.com/ethereum/go-ethereum/common"
	"os"
)

type FileRepository struct {
	stateMap                map[string]*models.State
	tokenMetadataMap        map[common.Address]*models.OpenSeaMetadata
	authenticatedTokenIdMap map[common.Address]models.SelectedAssetInfo
	bookingMap              map[float64]models.BookingInfo
}

func NewFileRepository() (*FileRepository, error) {
	return &FileRepository{
		stateMap:                make(map[string]*models.State),
		tokenMetadataMap:        make(map[common.Address]*models.OpenSeaMetadata),
		authenticatedTokenIdMap: make(map[common.Address]models.SelectedAssetInfo),
		bookingMap:              make(map[float64]models.BookingInfo),
	}, nil
}

func (f *FileRepository) GetMapping() ([]models.ScopeMapping, error) {
	data, err := os.ReadFile("resources/scopes/current.json")
	if err != nil {
		return []models.ScopeMapping{}, err
	}
	var scopeMapping []models.ScopeMapping
	err = json.Unmarshal(data, &scopeMapping)
	if err != nil {
		return []models.ScopeMapping{}, err
	}
	return scopeMapping, nil
}

func (f *FileRepository) WriteMapping(mappings []models.ScopeMapping) error {
	data, err := json.Marshal(mappings)
	if err != nil {
		return err
	}
	return os.WriteFile("resources/scopes/current.json", data, 0644)
}

func (f *FileRepository) GetPolicy() (*models.Policy, error) {
	data, err := os.ReadFile("resources/policies/current.json")
	if err != nil {
		return &models.Policy{}, err
	}
	var policy models.Policy
	err = json.Unmarshal(data, &policy)
	if err != nil {
		return &models.Policy{}, err
	}
	return &policy, nil
}

func (f *FileRepository) WritePolicy(policy models.Policy) error {
	data, err := json.Marshal(policy)
	if err != nil {
		return err
	}
	return os.WriteFile("resources/policies/current.json", data, 0644)
}

func (f *FileRepository) GetState(nonce string) (*models.State, bool) {
	state, ok := f.stateMap[nonce]
	return state, ok
}

func (f *FileRepository) WriteState(nonce string, state *models.State) error {
	f.stateMap[nonce] = state
	return nil
}

func (f *FileRepository) GetTokenMetadata(address common.Address) (*models.OpenSeaMetadata, bool) {
	metadata, ok := f.tokenMetadataMap[address]
	return metadata, ok
}

func (f *FileRepository) WriteTokenMetadata(address common.Address, metadata *models.OpenSeaMetadata) {
	f.tokenMetadataMap[address] = metadata
}

func (f *FileRepository) GetAuthenticatedTokenInfo(address common.Address) (*models.SelectedAssetInfo, bool) {
	assetInfo, ok := f.authenticatedTokenIdMap[address]
	return &assetInfo, ok
}

func (f *FileRepository) WriteAuthenticatedTokenInfo(address common.Address, assetInfo *models.SelectedAssetInfo) {
	f.authenticatedTokenIdMap[address] = *assetInfo
}

func (f *FileRepository) GetBookingInfo(tokenId float64) (models.BookingInfo, bool) {
	bookingInfo, ok := f.bookingMap[tokenId]
	return bookingInfo, ok
}

func (f *FileRepository) WriteBookingInfo(tokenId float64, bookingInfo models.BookingInfo) {
	f.bookingMap[tokenId] = bookingInfo
}
