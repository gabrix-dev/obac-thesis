package models

import (
	"github.com/ethereum/go-ethereum/common"
	signer "github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type SignedMessage struct {
	Signature string           `json:"signature"`
	Message   signer.TypedData `json:"message"`
}

type CodeVerify struct {
	CodeVerifier string `json:"code_verifier"`
	Code         string `json:"code"`
	Nonce        string `json:"nonce"`
}

type AssociationValues struct {
	AssetAttribute []string `json:"asset_attribute"`
	AssetRelated   []string `json:"asset_related"`
}

type TokenResponse struct {
	TokenType      string `json:"token_type"`
	ExpiresIn      int    `json:"expires_in"`
	AccessToken    string `json:"access_token,omitempty"`
	IDToken        string `json:"id_token,omitempty"`
	OwnershipToken string `json:"ownership_token"`
	Scope          string `json:"scope"`
}

type CodeResponse struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type SelectedAssetInfo struct {
	Id        int64
	SCAddress common.Address
	BCNetwork string
}

type SCInfo struct {
	SCId
	Name string `json:"name"`
}

type ScopeMapping struct {
	Name   string  `json:"name"`
	Claims []Claim `json:"claims"`
}

type Claim struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Oidc  string `json:"oidc"`
}
