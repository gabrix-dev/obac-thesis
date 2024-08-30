package models

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type OpenSeaMetadata struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Image       string           `json:"image"`
	Attributes  []AttributeValue `json:"attributes"`
}

type AttributeValue struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

func (o *OpenSeaMetadata) ToMap() map[string]string {
	metadataMap := make(map[string]string)
	for _, attribute := range o.Attributes {
		metadataMap[attribute.TraitType] = attribute.Value
	}
	return metadataMap
}

type ERC721 interface {
	TokenURI(opts *bind.CallOpts, _tokenId *big.Int) (string, error)
	OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error)
	Name(opts *bind.CallOpts) (string, error)
}
