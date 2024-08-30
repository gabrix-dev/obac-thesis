package repository

import (
	"entrust.com/iat/obac-idaas-server/pkg/config"
	"entrust.com/iat/obac-idaas-server/pkg/models"
	"github.com/ethereum/go-ethereum/common"
)

type Repository interface {
	GetPolicy() (*models.Policy, error)
	WritePolicy(models.Policy) error
	GetState(nonce string) (*models.State, bool)
	WriteState(nonce string, state *models.State) error
	GetTokenMetadata(address common.Address) (*models.OpenSeaMetadata, bool)
	WriteTokenMetadata(address common.Address, metadata *models.OpenSeaMetadata)
	GetAuthenticatedTokenInfo(address common.Address) (*models.SelectedAssetInfo, bool)
	WriteAuthenticatedTokenInfo(address common.Address, assetInfo *models.SelectedAssetInfo)
	GetBookingInfo(tokenId float64) (models.BookingInfo, bool)
	WriteBookingInfo(tokenId float64, bookingInfo models.BookingInfo)
	WriteMapping([]models.ScopeMapping) error
	GetMapping() ([]models.ScopeMapping, error)
}

type BlockchainRepository interface {
	GetTokenMetadataUri(id int64, sc models.SCId) (string, error)
	GetTokenIds(ownerAddress common.Address, sc models.SCId) ([]int64, error)
	GetUserAssetInfo(ownerAddress common.Address, scs []models.SCId) ([]*models.UserAssets, error)
	StartListeningToSCEventLogs(sc models.SCId)
	UpdateTokenUriMap(scs ...models.SCId) error
	GetSCName(sc models.SCId) (string, error)
	GetSCNameBatch(scs ...models.SCId) (map[models.SCId]string, error)
	GetMetadataUris(scs []models.SCId) ([]string, error)
	GetAssetRelatedValues() ([]string, error)
	UpdateContracts([]models.SCId) error
}

type ContentAddressableRepository interface {
	GetFile(ipfsLink string) ([]byte, error)
	GetFilesBatch(ipfsLinks []string) ([][]byte, error)
	GetImage(ipfsLink string) (string, error)
}

func NewContentAddressableRepository() (ContentAddressableRepository, error) {
	return NewIpfsRepositoy()
}

func NewBlockchainRepository(providerConfig *config.ProviderConfig) (BlockchainRepository, error) {
	return NewEthereumBasedBlockchainRepository(providerConfig)
}

func NewRepository() (Repository, error) {
	return NewFileRepository()
}
