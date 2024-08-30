package models

import "github.com/ethereum/go-ethereum/common"

var PolicyContract common.Address = common.HexToAddress("0xfd6b508d7a6253809d796101b2b6ae9164ab5959") //Goerli (ethereum)
//var PolicyContract common.Address = common.HexToAddress("0x3aFB1ad9C538C085F97D9486bB46A462cB8EaEE8") //Mumbai (polygon)

type Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MetadataConstraint struct {
	Claim    string   `json:"key"`
	Value    []string `json:"value"`
	Operator string   `json:"operator"`
}

type OwnershipConstraint struct {
	SmartContract       SCId                 `json:"smartContract"`
	MetadataConstraints []MetadataConstraint `json:"metadataConstraints,omitempty"`
}

type AccessControlRule struct {
	Name                 string                `json:"name"`
	Action               string                `json:"action"`
	AccessAddresses      []common.Address      `json:"accessAddresses,omitempty"`
	OwnershipConstraints []OwnershipConstraint `json:"ownershipConstraints"`
}

type AccessControlPolicy struct {
	AccessControlRules []AccessControlRule `json:"accessControlRules"`
}

type Policies struct {
	AccessControlPolicy AccessControlPolicy `json:"accessControlPolicy"`
}

type Policy struct {
	Metadata Metadata `json:"metadata"`
	Policies `json:"policies"`
}

func (p *Policy) GetSCs() []SCId {
	var scList []SCId
	for _, rule := range p.AccessControlPolicy.AccessControlRules {
		for _, constraint := range rule.OwnershipConstraints {
			scList = append(scList, constraint.SmartContract)
		}
	}
	return scList
}
