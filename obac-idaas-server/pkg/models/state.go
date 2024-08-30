package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type State struct {
	CodeChallenge string
	Code          string
	State         string
	RedirectUri   string
	Scope         string
	Wallet        common.Address
}

func (s *State) SetCodeChallenge(cc string) {
	s.CodeChallenge = cc
}

func (s *State) SetCode(c string) {
	s.Code = c
}

func (s *State) SetState(st string) {
	s.State = st
}

func (s *State) SetWallet(w common.Address) {
	s.Wallet = w
}
