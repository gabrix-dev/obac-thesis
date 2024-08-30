package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	error2 "entrust.com/iat/obac-idaas-server/pkg/error"
	"entrust.com/iat/obac-idaas-server/pkg/logging"
	"github.com/sirupsen/logrus"

	"entrust.com/iat/obac-idaas-server/pkg/config"
	"entrust.com/iat/obac-idaas-server/pkg/models"
	"entrust.com/iat/obac-idaas-server/pkg/repository"
	service "entrust.com/iat/obac-idaas-server/pkg/service/oidc"
	"entrust.com/iat/obac-idaas-server/pkg/utils"
	"github.com/pkg/errors"
)

type OIDCHandlerController interface {
	Authorize(response http.ResponseWriter, request *http.Request)
	Enforce(response http.ResponseWriter, request *http.Request)
	GetToken(response http.ResponseWriter, request *http.Request)
	GetAssets(response http.ResponseWriter, request *http.Request)
	PostMapping(response http.ResponseWriter, request *http.Request)
	GetMapping(response http.ResponseWriter, request *http.Request)
	GetAssociations(response http.ResponseWriter, request *http.Request)
}

type oidcHandlerControllerImpl struct {
	service service.OIDCHandlerService
	idpUri  string
}

func NewOIDCHandlerController(appUriConfig, idpUriConfig config.UriConfig, repository repository.Repository, bcRepo repository.BlockchainRepository, cadRepo repository.ContentAddressableRepository) OIDCHandlerController {
	return &oidcHandlerControllerImpl{
		service: service.NewOIDCHandlerService(repository, bcRepo, cadRepo),
		idpUri:  idpUriConfig.Protocol + "://" + idpUriConfig.Host + ":" + idpUriConfig.Port + idpUriConfig.Path,
	}
}

func (o oidcHandlerControllerImpl) GetAssociations(response http.ResponseWriter, request *http.Request) {
	associations, err := o.service.GetAssociations()
	if err != nil {
		HandleError(request, response, err)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(response).Encode(associations); err != nil {
		HandleError(request, response, fmt.Errorf("error encoding the associations: %s", err.Error()))
		return
	}
	logging.InfoEvent(logging.GetAssociationsHandled, nil)
}

func (o oidcHandlerControllerImpl) GetMapping(response http.ResponseWriter, request *http.Request) {
	mapping, err := o.service.GetMapping()
	if err != nil {
		HandleError(request, response, fmt.Errorf("error retreiving the mapping: %s", err.Error()))
		return
	}
	response.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(response).Encode(mapping); err != nil {
		HandleError(request, response, fmt.Errorf("error encoding the mapping: %s", err.Error()))
		return
	}
	logging.InfoEvent(logging.GetMappingHandled, nil)
}

func (o oidcHandlerControllerImpl) PostMapping(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	var mapping []models.ScopeMapping //Todo: use list of pointers
	if err := json.NewDecoder(request.Body).Decode(&mapping); err != nil {
		HandleError(request, response, error2.NewBadRequestError(err, "error parsing the mapping JSON"))
		return
	}
	if err := o.service.WriteMapping(mapping); err != nil {
		HandleError(request, response, err)
		return
	}
	response.WriteHeader(http.StatusOK)
	logging.InfoEvent(logging.PostMappingHandled, logrus.Fields{"mapping": mapping})
}

func (o oidcHandlerControllerImpl) Authorize(response http.ResponseWriter, request *http.Request) {
	err := checkParameters(request.URL.Query())
	if err != nil {
		HandleErrorWithFields(request, response, error2.NewBadRequestError(err, "error processing the request: "+err.Error()),
			&logrus.Fields{"request_parameters": request.URL.Query()})
		return
	}
	nonce, err := utils.GenerateRandomString(8)
	if err != nil {
		HandleError(request, response, fmt.Errorf("error generating the random nonce: %s", err.Error()))
		return
	}
	err = o.service.WriteState(nonce, &models.State{CodeChallenge: request.URL.Query().Get("code_challenge"),
		State: request.URL.Query().Get("state"), RedirectUri: request.URL.Query().Get("redirect_uri"), Scope: request.URL.Query().Get("scope")})
	if err != nil {
		HandleErrorWithFields(request, response, fmt.Errorf("error writing the authorization data: %s", err.Error()), &logrus.Fields{"incoming_url": request.URL.Query()})
		return
	}
	http.Redirect(response, request, o.idpUri+"/?nonce="+nonce+"&redirect_uri="+request.URL.Query().Get("redirect_uri"), http.StatusFound)
	logging.InfoEvent(logging.AuthorizeHandled, logrus.Fields{
		"nonce":          nonce,
		"code_challenge": request.URL.Query().Get("code_challenge"),
		"redirect_uri":   request.URL.Query().Get("redirect_uri"),
	})
}

func (o oidcHandlerControllerImpl) GetAssets(response http.ResponseWriter, request *http.Request) {
	userAddressStr := request.URL.Query().Get("address")
	if userAddressStr == "" {
		HandleError(request, response, error2.NewBadRequestError(errors.New("wrong request format"), "missing address url parameter"))
		return
	}
	currentPolicy, err := o.service.GetPolicy()
	if err != nil {
		HandleError(request, response, fmt.Errorf("error retreiving the current policy: %s", err.Error()))
		return
	}
	currentPolicySCs := currentPolicy.GetSCs()
	ownerAssets, err := o.service.GetOwnerAssets(userAddressStr, currentPolicySCs)
	if err != nil {
		HandleErrorWithFields(request, response, fmt.Errorf("error retreiving the owner's assets: %s", err.Error()), &logrus.Fields{"owner_address": userAddressStr})
		return
	}
	if err = o.service.SetImageUris(ownerAssets); err != nil {
		HandleErrorWithFields(request, response, fmt.Errorf("error setting the asset's images: %s", err.Error()), &logrus.Fields{"owner_assets": ownerAssets})
		return
	}
	for _, collection := range ownerAssets {
		collection.Sort()
	}
	response.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(response).Encode(ownerAssets); err != nil {
		HandleError(request, response, fmt.Errorf("error encoding the user assets to the response: %s", err.Error()))
		return
	}
	logging.InfoEvent(logging.GetTokenIdsHandled, logrus.Fields{"owner_address": userAddressStr})
}

func (o oidcHandlerControllerImpl) Enforce(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	var receivedMessage models.SignedMessage
	err := json.NewDecoder(request.Body).Decode(&receivedMessage)
	if err != nil {
		HandleError(request, response, error2.NewBadRequestError(err, "error decoding the body parameters"))
		return
	}
	nonce := o.service.GetMessageNonce(&receivedMessage)
	state, err := o.service.GetState(nonce)
	if err != nil {
		HandleErrorWithFields(request, response, err, &logrus.Fields{"nonce": nonce})
		return
	}
	address, err := o.service.VerifySignature(receivedMessage)
	if err != nil {
		HandleError(request, response, error2.NewNotAuthorizedError(err, "invalid signature"))
		return
	}
	selectedAssetInfo := o.service.GetMessageSelectedAssetInfo(&receivedMessage)
	if err = o.service.Enforce(address, selectedAssetInfo); err != nil {
		HandleErrorWithFields(request, response, error2.NewNotAuthorizedError(err, "The asset did not pass the access control policy"), &logrus.Fields{"selected_asset": selectedAssetInfo})
		return
	}
	state.SetWallet(address)
	code, err := utils.GenerateRandomString(32)
	if err != nil {
		HandleError(request, response, fmt.Errorf("error generating the code: %s", err.Error()))
		return
	}
	state.SetCode(code)
	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(response).Encode(models.CodeResponse{Code: code, State: state.State})
	if err != nil {
		HandleError(request, response, fmt.Errorf("error encoding the code response: %s", err.Error()))
	}
	logging.InfoEvent(logging.EnforceHandled, logrus.Fields{
		"nonce":          nonce,
		"owner_address":  address,
		"token_code":     code,
		"selected_asset": selectedAssetInfo,
	})
}

func (o oidcHandlerControllerImpl) GetToken(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	var cv models.CodeVerify
	err := json.NewDecoder(request.Body).Decode(&cv)
	if err != nil {
		HandleError(request, response, error2.NewBadRequestError(err, "error decoding the body parameters"))
		return
	}
	state, err := o.service.GetState(cv.Nonce)
	if err != nil {
		HandleErrorWithFields(request, response, error2.NewNotAuthorizedError(err, "error retreiving the state"), &logrus.Fields{"nonce": cv.Nonce, "code": cv.Code})
		return
	}
	if cv.Code != state.Code {
		HandleErrorWithFields(request, response, error2.NewNotAuthorizedError(errors.New("received code doesn't match the state code"), "unauthorized user"),
			&logrus.Fields{"stored_code": state.Code, "received_code": cv.Code, "nonce": cv.Nonce})
		return
	}
	verifierBytes := []byte(cv.CodeVerifier)
	if err = verifyChallenge(verifierBytes, state.CodeChallenge); err != nil {
		HandleError(request, response, error2.NewNotAuthorizedError(err, "code verifier and challenge don't match"))
		return
	}
	nftMetadata, err := o.service.GetTokenMetadata(state.Wallet)
	if err != nil {
		HandleErrorWithFields(request, response, fmt.Errorf("error retreiving the nft metadata: %s", err.Error()), &logrus.Fields{"owner_address": state.Wallet})
		return
	}
	tokenStr, expiresIn, err := o.service.GenerateJWT(nftMetadata, state.Wallet, state.Scope)
	if err != nil {
		HandleErrorWithFields(request, response, fmt.Errorf("error generating the ownership token: %s", err.Error()), &logrus.Fields{"scope": state.Scope, "metadata": nftMetadata})
		return
	}
	tokenResponse := models.TokenResponse{TokenType: "Bearer", ExpiresIn: int(expiresIn.Seconds()), OwnershipToken: tokenStr, Scope: state.Scope}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(tokenResponse)
	logging.InfoEvent(logging.GetTokenHandled, logrus.Fields{
		"nonce":           cv.Nonce,
		"ownership_token": tokenStr,
		"scope":           state.Scope,
	})
}

func checkParameters(values url.Values) error {
	expectedParameters := []string{"scope", "code_challenge", "redirect_uri", "state"}
	for _, s := range expectedParameters {
		if values.Get(s) == "" {
			return errors.New("URL parameter missing: " + s)
		}
	}
	return nil
}

func verifyChallenge(codeVerifier []byte, codeChallenge string) error {
	hash := sha256.Sum256(codeVerifier)
	hashStr := hex.EncodeToString(hash[:])
	if hashStr != codeChallenge {
		return errors.New("invalid code verifier")
	}
	return nil
}
