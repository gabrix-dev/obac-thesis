package models

import (
	"github.com/ethereum/go-ethereum/common"
	"sort"
	"sync"
)

type TokenIdInfo struct {
	Address  common.Address
	TokenId  int64
	TokenUri string
}

type MutexCounter struct {
	mu     *sync.RWMutex
	Number int
}

func NewMutexCounter() *MutexCounter {
	return &MutexCounter{&sync.RWMutex{}, 0}
}

func (c *MutexCounter) Add() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Number = c.Number + 1
}

type TransferEvent struct {
	Contract SCId
	From     common.Address
	To       common.Address
	TokenId  int64
}

type UserAssets struct {
	SCInfo
	Assets []*AssetData `json:"assetData"`
}

type SCId struct {
	Address    common.Address `json:"address"`
	Blockchain string         `json:"blockchain"`
}

type AssetData struct {
	Id          int64  `json:"id"`
	MetadataUri string `json:"-"`
	ImageUrl    string `json:"imageUrl"`
}

func (s *UserAssets) Sort() {
	sort.Slice(s.Assets, func(i, j int) bool {
		return s.Assets[i].Id < s.Assets[j].Id
	})
}

func (s *UserAssets) SetAssetIds(ids []int64) {
	s.Assets = make([]*AssetData, len(ids))
	for i, id := range ids {
		s.Assets[i] = &AssetData{
			Id: id,
		}
	}
}
