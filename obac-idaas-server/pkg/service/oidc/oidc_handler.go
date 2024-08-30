package oidc

import (
	"encoding/hex"
	"encoding/json"
	"entrust.com/iat/obac-idaas-server/pkg/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"

	"entrust.com/iat/obac-idaas-server/pkg/models"
	"entrust.com/iat/obac-idaas-server/pkg/repository"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	signer "github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/pkg/errors"
)

type OIDCHandlerService interface {
	GetState(nonce string) (*models.State, error)
	WriteState(nonce string, state *models.State) error
	GetTokenMetadata(address common.Address) (*models.OpenSeaMetadata, error)
	WriteTokenMetadata(address common.Address, metadata *models.OpenSeaMetadata)
	GetAuthenticatedTokenInfo(address common.Address) (*models.SelectedAssetInfo, error)
	WriteAuthenticatedTokenInfo(address common.Address, assetInfo *models.SelectedAssetInfo)
	Enforce(userAddress common.Address, selectedAssetInfo *models.SelectedAssetInfo) error
	VerifySignature(signedMessage models.SignedMessage) (common.Address, error)
	GetTokenIds(address string, smartContract models.SCId) ([]int64, error)
	GetOwnerAssets(address string, smartContracts []models.SCId) ([]*models.UserAssets, error)
	GetMessageNonce(message *models.SignedMessage) string
	GetMessageSelectedAssetInfo(message *models.SignedMessage) *models.SelectedAssetInfo
	GetPolicy() (*models.Policy, error)
	SetImageUris(scTokenIds []*models.UserAssets) error
	WriteMapping(mapping []models.ScopeMapping) error
	GetMapping() ([]models.ScopeMapping, error)
	GenerateJWT(nftMetadata *models.OpenSeaMetadata, wallet common.Address, scopes string) (tokenString string, expiration time.Duration, err error)
	GetOwnershipClaims(metadata *models.OpenSeaMetadata, wallet common.Address, scopes string) (map[string]interface{}, error)
	GetAssociations() (models.AssociationValues, error)
}

func NewOIDCHandlerService(repo repository.Repository, bcRepo repository.BlockchainRepository, cadRepo repository.ContentAddressableRepository) OIDCHandlerService {
	return &oidcHandlerServiceImpl{
		repo:    repo,
		bcRepo:  bcRepo,
		cadRepo: cadRepo,
	}
}

type oidcHandlerServiceImpl struct {
	repo    repository.Repository
	bcRepo  repository.BlockchainRepository
	cadRepo repository.ContentAddressableRepository
}

func (h *oidcHandlerServiceImpl) GetAssociations() (models.AssociationValues, error) {
	policy, err := h.GetPolicy()
	if err != nil {
		return models.AssociationValues{}, err
	}
	policyContracts := policy.GetSCs()
	uris, err := h.bcRepo.GetMetadataUris(policyContracts)
	if err != nil {
		return models.AssociationValues{}, err
	}
	files, err := h.cadRepo.GetFilesBatch(uris)
	if err != nil {
		return models.AssociationValues{}, err
	}
	assetAttributeValues := []string{""}
	for _, file := range files {
		err = addAssetAttributes(&assetAttributeValues, file)
		if err != nil {
			return models.AssociationValues{}, err
		}
	}
	utils.RemoveDuplicates(&assetAttributeValues)
	assetRelatedValues, err := h.bcRepo.GetAssetRelatedValues()
	if err != nil {
		return models.AssociationValues{}, err
	}
	return models.AssociationValues{AssetAttribute: assetAttributeValues, AssetRelated: assetRelatedValues}, nil
}

func addAssetAttributes(attributes *[]string, file []byte) error {
	var openseaMetadata models.OpenSeaMetadata
	if err := json.Unmarshal(file, &openseaMetadata); err != nil {
		return err
	}
	for _, attribute := range openseaMetadata.Attributes {
		*attributes = append(*attributes, attribute.TraitType)
	}
	return nil
}

func (h *oidcHandlerServiceImpl) GetMapping() ([]models.ScopeMapping, error) {
	return h.repo.GetMapping()
}

func (h *oidcHandlerServiceImpl) WriteMapping(mapping []models.ScopeMapping) error {
	return h.repo.WriteMapping(mapping)
}

func (h *oidcHandlerServiceImpl) SetImageUris(scTokenIds []*models.UserAssets) error {
	for _, scUserInfo := range scTokenIds {
		for _, asset := range scUserInfo.Assets {
			imageUrl, err := h.cadRepo.GetImage(asset.MetadataUri)
			if err != nil {
				return err
			}
			asset.ImageUrl = imageUrl
		}
	}
	return nil
}

func (h *oidcHandlerServiceImpl) GetState(nonce string) (*models.State, error) {
	state, ok := h.repo.GetState(nonce)
	if !ok {
		return nil, errors.New("state not found")
	}
	return state, nil
}

func (h *oidcHandlerServiceImpl) WriteState(nonce string, state *models.State) error {
	h.repo.WriteState(nonce, state)
	return nil
}

func (h *oidcHandlerServiceImpl) GetTokenMetadata(address common.Address) (*models.OpenSeaMetadata, error) {
	metadata, ok := h.repo.GetTokenMetadata(address)
	if !ok {
		return nil, errors.New("token metadata not found")
	}
	return metadata, nil
}

func (h *oidcHandlerServiceImpl) WriteTokenMetadata(address common.Address, metadata *models.OpenSeaMetadata) {
	h.repo.WriteTokenMetadata(address, metadata)
}

func (h *oidcHandlerServiceImpl) GetAuthenticatedTokenInfo(address common.Address) (*models.SelectedAssetInfo, error) {
	assetInfo, ok := h.repo.GetAuthenticatedTokenInfo(address)
	if !ok {
		return &models.SelectedAssetInfo{}, errors.New("asset info not found")
	}
	return assetInfo, nil
}

func (h *oidcHandlerServiceImpl) WriteAuthenticatedTokenInfo(address common.Address, tokenInfo *models.SelectedAssetInfo) {
	h.repo.WriteAuthenticatedTokenInfo(address, tokenInfo)
}

func (h *oidcHandlerServiceImpl) Enforce(userAddress common.Address, selectedAssetInfo *models.SelectedAssetInfo) error {
	policy, err := h.repo.GetPolicy()
	if err != nil {
		return err
	}
	for _, rule := range policy.AccessControlPolicy.AccessControlRules {
		err = h.enforceRule(&rule, userAddress, selectedAssetInfo)
		if err == nil {
			return nil
		}
	}
	return err
}

func (h *oidcHandlerServiceImpl) enforceRule(rule *models.AccessControlRule, userAddress common.Address, selectedAssetInfo *models.SelectedAssetInfo) error {
	if rule.AccessAddresses != nil {
		if isWhitelisted := checkIfWhitelisted(rule, userAddress); !isWhitelisted {
			return errors.New("user not whitelisted")
		}
	}
	if rule.OwnershipConstraints != nil {
		for _, ownershipConstraint := range rule.OwnershipConstraints {
			if ownershipConstraint.SmartContract.Blockchain != strings.ToLower(selectedAssetInfo.BCNetwork) || ownershipConstraint.SmartContract.Address != selectedAssetInfo.SCAddress {
				continue
			}
			err := h.checkOwnershipConstraint(&ownershipConstraint, userAddress, selectedAssetInfo)
			if err == nil {
				return nil
			}
		}
		return errors.New("the asset did not pass the access control policy ")
	}
	return errors.New("ownership constraints missing")
}

func (h *oidcHandlerServiceImpl) checkOwnershipConstraint(ownershipConstraint *models.OwnershipConstraint, userAddress common.Address, selectedAssetInfo *models.SelectedAssetInfo) error {
	tokenIds, err := h.bcRepo.GetTokenIds(userAddress, ownershipConstraint.SmartContract)
	if err != nil {
		return err
	}
	if len(tokenIds) == 0 {
		return errors.New("the user doesn't own a NFT of the collection")
	}
	if !checkIfTokenOwned(tokenIds, selectedAssetInfo.Id) {
		return err
	}
	metadataUrl, err := h.bcRepo.GetTokenMetadataUri(selectedAssetInfo.Id, ownershipConstraint.SmartContract)
	if err != nil {
		return err
	}
	metadataBytes, err := h.cadRepo.GetFile(metadataUrl)
	if err != nil {
		return err
	}
	var openseaMetadata models.OpenSeaMetadata
	if err = json.Unmarshal(metadataBytes, &openseaMetadata); err != nil {
		return err
	}
	if ownershipConstraint.MetadataConstraints != nil {
		metadataMap := openseaMetadata.ToMap()
		for _, metadataConstraint := range ownershipConstraint.MetadataConstraints {
			ok := checkMetadataConstraint(&metadataConstraint, metadataMap)
			if !ok {
				return errors.New("asset doesn't pass the metadata policy")
			}
		}
	}
	h.repo.WriteAuthenticatedTokenInfo(userAddress, selectedAssetInfo)
	h.repo.WriteTokenMetadata(userAddress, &openseaMetadata)
	return nil
}

func checkMetadataConstraint(metadataConstraint *models.MetadataConstraint, assetClaims map[string]string) bool {
	for _, targetValue := range metadataConstraint.Value {
		userValue := assetClaims[metadataConstraint.Claim]
		if metadataConstraint.Operator == "Equals" {
			if strings.ToLower(userValue) == strings.ToLower(targetValue) {
				return true
			}
		} else {
			if strings.ToLower(userValue) == strings.ToLower(targetValue) {
				return false
			}
		}
	}
	return !(metadataConstraint.Operator == "Equals") //If it arrives here returns false for "Equals" and true for "Not Equals"
}

func checkIfTokenOwned(tokenIds []int64, selectedAsset int64) bool {
	for _, id := range tokenIds {
		if id == selectedAsset {
			return true
		}
	}
	return false
}

func checkIfWhitelisted(rule *models.AccessControlRule, userAddress common.Address) bool {
	for _, address := range rule.AccessAddresses {
		if userAddress.String() == address.String() {
			return true
		}
	}
	return false
}

func (h *oidcHandlerServiceImpl) VerifySignature(signedMessage models.SignedMessage) (common.Address, error) {

	var signingAddress common.Address
	incomingMetamaskSignature, typedData := signedMessage.Signature, signedMessage.Message

	// Decode the hex-encoded signature from metamask
	signature, err := hex.DecodeString(incomingMetamaskSignature[2:]) //incomingMetamaskSignature starts with "0x", we want to decode the rest
	if err != nil {
		return signingAddress, err
	}
	if len(signature) != 65 {
		return signingAddress, fmt.Errorf("invalid signature length: %d", len(signature))
	}
	if signature[64] != 27 && signature[64] != 28 {
		return signingAddress, fmt.Errorf("invalid recovery id: %d", signature[64])
	}
	signature[64] -= 27 //0 o 1
	hashedPayload := hashPayload(typedData)
	sigPublicKeyECDSA, err := crypto.SigToPub(hashedPayload.Bytes(), signature)
	if err != nil {
		return signingAddress, fmt.Errorf("invalid signature: %s", err.Error())
	}
	signingAddress = crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	err = h.checkParameters(typedData, signingAddress)
	if err != nil {
		return signingAddress, fmt.Errorf("invalid struct: %s", err.Error())
	}
	return signingAddress, nil
}

func hashPayload(typedData signer.TypedData) common.Hash {
	typedDataHash, _ := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	domainSeparator, _ := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hashedPayload := crypto.Keccak256Hash(rawData)
	return hashedPayload
}

func (h *oidcHandlerServiceImpl) checkParameters(typedData signer.TypedData, recoveredAddress common.Address) error {
	domain := typedData.Domain.Map()
	nonce := typedData.Message["nonce"].(string)
	if domain["name"] != "NFT ownership verification service" {
		return fmt.Errorf("invalid domain: " + domain["name"].(string))
	}

	//We check if the db contains the nonce
	_, ok := h.repo.GetState(nonce)
	if !ok {
		return errors.New("Invalid nonce " + nonce)
	}
	return nil
}

func (h *oidcHandlerServiceImpl) GetTokenIds(userAddress string, sc models.SCId) ([]int64, error) {
	return h.bcRepo.GetTokenIds(common.HexToAddress(userAddress), sc)
}

func (h *oidcHandlerServiceImpl) GetOwnerAssets(userAddress string, scs []models.SCId) ([]*models.UserAssets, error) {
	return h.bcRepo.GetUserAssetInfo(common.HexToAddress(userAddress), scs)
}

func (h *oidcHandlerServiceImpl) GetMessageNonce(message *models.SignedMessage) string {
	return message.Message.Message["nonce"].(string)
}

func (h *oidcHandlerServiceImpl) GetMessageSelectedAssetInfo(message *models.SignedMessage) *models.SelectedAssetInfo {
	selectedAsset := message.Message.Message["asset ID"]
	collection := message.Message.Message["asset collection"]
	network := message.Message.Message["asset blockchain"]
	return &models.SelectedAssetInfo{
		Id:        int64(selectedAsset.(float64)),
		SCAddress: common.HexToAddress(collection.(string)),
		BCNetwork: network.(string),
	}
}

func (h *oidcHandlerServiceImpl) GenerateJWT(nftMetadata *models.OpenSeaMetadata, wallet common.Address, scopes string) (tokenString string, expiration time.Duration, err error) {
	expiresIn := 1 * time.Hour
	expirationTime := time.Now().Add(expiresIn)
	ownershipClaims, err := h.GetOwnershipClaims(nftMetadata, wallet, scopes)
	if err != nil {
		return "", expiresIn, fmt.Errorf("error getting the ownership claims: %s" + err.Error())
	}
	ownershipClaims["exp"] = expirationTime.Unix()
	claims := jwt.MapClaims(ownershipClaims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte("superdopekey"))
	if err != nil {
		return "", expiresIn, fmt.Errorf("error signing the token: %s" + err.Error())
	}
	return
}

func (h *oidcHandlerServiceImpl) GetOwnershipClaims(metadata *models.OpenSeaMetadata, wallet common.Address, scopes string) (map[string]interface{}, error) {
	assetInfo, exist := h.repo.GetAuthenticatedTokenInfo(wallet)
	if !exist {
		return map[string]interface{}{}, errors.New("asset info not found")
	}
	mapping, err := h.repo.GetMapping()
	if err != nil {
		return map[string]interface{}{}, err
	}
	scopeList := strings.Fields(scopes)
	ownershipClaims := make(map[string]interface{})
	for _, scope := range scopeList {
		if scope != "openid" {
			scopeClaims := getScopeClaims(scope, mapping)
			if len(scopeClaims) == 0 {
				continue
			}
			for _, claim := range scopeClaims {
				err = addClaim(ownershipClaims, claim, assetInfo, wallet, metadata.ToMap())
				if err != nil {
					//Todo: afegir InfoErrorLog
				}
			}
		}
	}
	return ownershipClaims, nil
}

func addClaim(claims map[string]interface{}, claim models.Claim, assetInfo *models.SelectedAssetInfo, wallet common.Address, nftMetadata map[string]string) error {
	if claim.Type == "asset_related" {
		switch claim.Value {
		case "Blockchain":
			claims[claim.Oidc] = assetInfo.BCNetwork
		case "Smart contract address":
			claims[claim.Oidc] = assetInfo.SCAddress
		case "Asset id":
			claims[claim.Oidc] = assetInfo.Id
		case "Owner address":
			claims[claim.Oidc] = wallet
		}
	} else if claim.Type == "asset_attribute" {
		value, exist := nftMetadata[claim.Value]
		if exist {
			claims[claim.Oidc] = value
		}
	} else {
		return errors.New("claim type not found")
	}
	return nil
}

func getScopeClaims(name string, mappings []models.ScopeMapping) []models.Claim {
	for _, scopeMapping := range mappings {
		if scopeMapping.Name == name {
			return scopeMapping.Claims
		}
	}
	return []models.Claim{}
}

func (h *oidcHandlerServiceImpl) GetPolicy() (*models.Policy, error) {
	return h.repo.GetPolicy()
}
