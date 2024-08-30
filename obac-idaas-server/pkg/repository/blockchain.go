package repository

import (
	"context"
	"entrust.com/iat/obac-idaas-server/pkg/config"
	"entrust.com/iat/obac-idaas-server/pkg/models"
	"entrust.com/iat/obac-idaas-server/pkg/utils"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strconv"
	"time"
)

type EthereumBasedBlockchainRepository struct {
	contractGw         models.ContractGateway
	providers          map[string]*ethclient.Client
	tokenIdMap         map[models.SCId]map[common.Address][]int64
	tokenUriMap        map[models.SCId]map[int64]string
	assetRelatedValues []string
}

// Todo: providers[int]*ethClient.Client --> chainId (int) not chain's name (string))
func NewEthereumBasedBlockchainRepository(providerConfig *config.ProviderConfig) (*EthereumBasedBlockchainRepository, error) {
	providers, err := initProviders(providerConfig)
	if err != nil {
		return &EthereumBasedBlockchainRepository{}, err
	}
	rateConfig := models.RateConfig{NumThreads: 14, MaxRequestsPerInterval: 3, IntervalDuration: 130 * time.Millisecond}
	cgw := models.NewContractGateway(rateConfig)
	return &EthereumBasedBlockchainRepository{
		providers:          providers,
		contractGw:         *cgw,
		tokenIdMap:         map[models.SCId]map[common.Address][]int64{},
		tokenUriMap:        map[models.SCId]map[int64]string{},
		assetRelatedValues: []string{"", "Smart contract address", "Owner address", "Blockchain", "Asset id"},
	}, nil
}

func initProviders(providerConfig *config.ProviderConfig) (map[string]*ethclient.Client, error) {
	providers := make(map[string]*ethclient.Client)
	mumbaiProvider, err := ethclient.Dial(providerConfig.MumbaiApiUrl)
	if err != nil {
		return nil, errors.New("error initializing polygon provider, " + err.Error())
	}
	goerliProvider, err := ethclient.Dial(providerConfig.GoerliApiUrl)
	if err != nil {
		return nil, errors.New("Error initializing goerli provider, " + err.Error())
	}
	binanceProvider, err := ethclient.Dial(providerConfig.BinanceApiUrl)
	if err != nil {
		return nil, errors.New("Error initializing binance provider, " + err.Error())
	}
	goerliWebSocketProvider, err := ethclient.Dial(providerConfig.GoerliApiWebSocketUrl)
	if err != nil {
		return nil, errors.New("Error initializing goerli wss provider: " + err.Error())
	}
	providers["polygon"] = mumbaiProvider
	providers["ethereum"] = goerliProvider
	providers["bsc"] = binanceProvider
	providers["ethereumWebSocket"] = goerliWebSocketProvider
	return providers, nil
}

func (e *EthereumBasedBlockchainRepository) UpdateContracts(contracts []models.SCId) error {
	for _, contract := range contracts {
		if _, err := e.contractGw.GetContract(contract); err != nil {
			instance, err := NewMyContract(contract.Address, e.providers[contract.Blockchain])
			if err != nil {
				return err
			}
			e.contractGw.SetContract(contract, instance)
		}
	}
	return nil
}

func (e *EthereumBasedBlockchainRepository) GetAssetRelatedValues() ([]string, error) {
	return e.assetRelatedValues, nil
}

const (
	goerliChainId = 5
	mumbaiChainId = 80001
)

func (e *EthereumBasedBlockchainRepository) GetMetadataUris(contracts []models.SCId) ([]string, error) {
	var metadataUris []string
	for _, contract := range contracts {
		uri, err := e.GetTokenMetadataUri(int64(1), contract)
		if err == nil {
			metadataUris = append(metadataUris, uri)
		}
	}
	return metadataUris, nil
}

func (e *EthereumBasedBlockchainRepository) GetTokenMetadataUri(id int64, sc models.SCId) (string, error) {
	uri, exist := e.tokenUriMap[sc][id]
	if exist {
		return uri, nil
	}
	resultchan := make(chan models.TaskResult)
	defer close(resultchan)
	task := models.Task{Contract: sc, TokenId: id, Method: "uri", ResultChan: &resultchan}
	e.contractGw.AddTask(&task)
	result := <-resultchan
	if result.Error != nil {
		return "", result.Error
	}
	return result.Result.(models.TokenIdInfo).TokenUri, nil
}

func (e *EthereumBasedBlockchainRepository) GetTokenIds(ownerAddress common.Address, sc models.SCId) ([]int64, error) {
	if err := e.UpdateTokenIdMap2(sc); err != nil {
		return []int64{}, err
	}
	return e.tokenIdMap[sc][ownerAddress], nil
}

func (e *EthereumBasedBlockchainRepository) GetUserAssetInfo(ownerAddress common.Address, scs []models.SCId) ([]*models.UserAssets, error) {
	var tokenIds []*models.UserAssets
	if err := e.UpdateTokenIdMap2(scs...); err != nil {
		return nil, err
	}
	scNames, err := e.GetSCNameBatch(scs...)
	if err != nil {
		return nil, err
	}
	for _, sc := range scs {
		scName, ok := scNames[sc]
		if !ok {
			scName = sc.Address.Hex()
		}
		ids, exist := e.tokenIdMap[sc][ownerAddress]
		if exist {
			scTokenIds := &models.UserAssets{SCInfo: models.SCInfo{Name: scName, SCId: sc}}
			scTokenIds.SetAssetIds(ids) //We create a list of Assets and initialize each with an id
			e.SetAssetUris(scTokenIds)
			tokenIds = append(tokenIds, scTokenIds)
		}
	}
	return tokenIds, nil
}

func (e *EthereumBasedBlockchainRepository) SetAssetUris(collection *models.UserAssets) error {
	var err error
	uriMap, exists := e.tokenUriMap[collection.SCId]
	if !exists {
		return err
	}
	for _, asset := range collection.Assets {
		metadataUri, exist := uriMap[asset.Id]
		if !exist {
			metadataUri, err = e.GetTokenMetadataUri(asset.Id, collection.SCId)
			if err != nil {
				return err
			}
		}
		asset.MetadataUri = metadataUri
	}
	return nil
}

func (e *EthereumBasedBlockchainRepository) UpdateTokenIdMap2(scs ...models.SCId) error {
	for _, sc := range scs {
		var basicErr, plusErr, vipErr bool
		tokenIdMap := make(map[common.Address][]int64)
		finishedTasksCount := models.NewMutexCounter()
		tasksCount := models.NewMutexCounter()

		resultChan := make(chan models.TaskResult, 30)
		taskLimitChan := make(chan bool, e.contractGw.MaxRequestsPerInterval*2) //We limit the amount of tasks to be added before having responses

		go e.InitTaskAdder(1, "owner", &taskLimitChan, &resultChan, &sc, tasksCount, &basicErr)
		go e.InitTaskAdder(1001, "owner", &taskLimitChan, &resultChan, &sc, tasksCount, &plusErr)
		go e.InitTaskAdder(2001, "owner", &taskLimitChan, &resultChan, &sc, tasksCount, &vipErr)

		for {
			taskResult := <-resultChan
			<-taskLimitChan //if handleResult -> break , saveTokenIdMap[sc]
			tokenInfo := taskResult.Result.(models.TokenIdInfo)
			if taskResult.Error != nil {
				finishedTasksCount.Add()
				if tokenInfo.TokenId < 1000 {
					basicErr = true
				} else if tokenInfo.TokenId > 2000 {
					vipErr = true
				} else {
					plusErr = true
				}
				if basicErr && vipErr && plusErr && finishedTasksCount.Number == tasksCount.Number {
					break
				}
			} else {
				finishedTasksCount.Add()
				tokenIdMap[tokenInfo.Address] = append(tokenIdMap[tokenInfo.Address], tokenInfo.TokenId)
			}
			e.tokenIdMap[sc] = tokenIdMap
		}
	}
	return nil
}

func (e *EthereumBasedBlockchainRepository) UpdateTokenUriMap(scs ...models.SCId) error {
	resultChan := make(chan models.TaskResult, 30)
	taskLimitChan := make(chan bool, e.contractGw.MaxRequestsPerInterval*2)
	defer close(resultChan)
	defer close(taskLimitChan)

	for _, sc := range scs {
		var basicErr, plusErr, vipErr bool
		tokenUriMap := make(map[int64]string)
		finishedTasksCount := models.NewMutexCounter()
		tasksCount := models.NewMutexCounter()

		go e.InitTaskAdder(1, "uri", &taskLimitChan, &resultChan, &sc, tasksCount, &basicErr)
		go e.InitTaskAdder(1001, "uri", &taskLimitChan, &resultChan, &sc, tasksCount, &plusErr)
		go e.InitTaskAdder(2001, "uri", &taskLimitChan, &resultChan, &sc, tasksCount, &vipErr)

		for {
			taskResult := <-resultChan
			<-taskLimitChan
			tokenInfo := taskResult.Result.(models.TokenIdInfo)
			if taskResult.Error != nil {
				finishedTasksCount.Add()
				if tokenInfo.TokenId < 1000 {
					basicErr = true
				} else if tokenInfo.TokenId > 2000 {
					vipErr = true
				} else {
					plusErr = true
				}
				if basicErr && vipErr && plusErr && finishedTasksCount.Number == tasksCount.Number {
					break
				}
			} else {
				finishedTasksCount.Add()
				tokenUriMap[tokenInfo.TokenId] = tokenInfo.TokenUri
			}
			e.tokenUriMap[sc] = tokenUriMap
		}
		utils.EmptyChannel(&taskLimitChan)
	}
	return nil
}

func (e *EthereumBasedBlockchainRepository) InitTaskAdder(startingId int, method string, taskLimiter *chan bool, taskResultChan *chan models.TaskResult, sc *models.SCId, taskCount *models.MutexCounter, stop *bool) {
	id := int64(startingId)
	for {
		if *stop {
			break
		}
		*taskLimiter <- true
		e.contractGw.AddTask(&models.Task{TokenId: id, Method: method, Contract: *sc, ResultChan: taskResultChan})
		taskCount.Add()
		id++
	}
}

func (e *EthereumBasedBlockchainRepository) StartListeningToSCEventLogs(sc models.SCId) {
	for {
		//transferEventSignature := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
		eventSignature := []byte("Transfer(address,address,uint256)")
		eventSignatureHash := crypto.Keccak256Hash(eventSignature)
		var topicsFilter [][]common.Hash
		topicsFilter = append(topicsFilter, []common.Hash{eventSignatureHash})

		query := ethereum.FilterQuery{
			Addresses: []common.Address{sc.Address},
			FromBlock: big.NewInt(8409793),
			Topics:    topicsFilter,
		}

		//Create the output channel
		logs := make(chan types.Log)

		//We create a subscription
		sub, err := e.providers["ethereumWebSocket"].SubscribeFilterLogs(context.Background(), query, logs)
		if err != nil {
			fmt.Printf("Error: creating the subscription %s\n", err.Error())
		}

		//Loop that either reads a new log or a subscription error
	Loop:
		for {
			select {
			//In case there is an error we unsubscribe from the subscription
			//break out of the second loop and start again
			case err := <-sub.Err():
				fmt.Println("WebSocket error:", err.Error())
				sub.Unsubscribe()
				break Loop
			case vLog := <-logs:
				var topics [4]string
				for i := range vLog.Topics {
					topics[i] = vLog.Topics[i].Hex()
				}
				id, err := strconv.ParseInt(topics[3][2:], 16, 64)
				if err != nil {
					fmt.Printf("Error parsing the token id\n")
				}
				transferEvent := &models.TransferEvent{
					Contract: sc,
					From:     common.HexToAddress(topics[1]),
					To:       common.HexToAddress(topics[2]),
					TokenId:  id,
				}
				err = e.onTokenTransferEvent(transferEvent)
				if err != nil {
					fmt.Printf("Error on token transfer event: %s\n", err)
				}
				fmt.Printf("Transfer event detected! \n")
				fmt.Printf("From: %s\n", transferEvent.From.Hex())
				fmt.Printf("To: %s\n", transferEvent.To.Hex())
				fmt.Printf("Id: %d\n", transferEvent.TokenId)
			}
		}
	}
}

func (e *EthereumBasedBlockchainRepository) onTokenTransferEvent(event *models.TransferEvent) error {
	tokenUri, err := e.GetTokenMetadataUri(event.TokenId, event.Contract)
	if err != nil {
		return err
	}
	e.tokenUriMap[event.Contract][event.TokenId] = tokenUri
	e.tokenIdMap[event.Contract][event.To] = append(e.tokenIdMap[event.Contract][event.To], event.TokenId)
	return nil
}

func (e *EthereumBasedBlockchainRepository) GetSCName(sc models.SCId) (string, error) {
	resultChan := make(chan models.TaskResult)
	e.contractGw.AddTask(&models.Task{Contract: sc, ResultChan: &resultChan, Method: "name"})
	taskResult := <-resultChan
	if taskResult.Error != nil {
		return "", taskResult.Error
	}
	return taskResult.Result.(string), nil
}

func (e *EthereumBasedBlockchainRepository) GetSCNameBatch(scs ...models.SCId) (map[models.SCId]string, error) {
	resultChan := make(chan models.TaskResult)
	defer close(resultChan)
	resultMap := make(map[models.SCId]string)
	var taskResult models.TaskResult
	for _, sc := range scs {
		e.contractGw.AddTask(&models.Task{Contract: sc, ResultChan: &resultChan, Method: "name"})
	}
	for i := 0; i < len(scs); i++ {
		taskResult = <-resultChan
		if taskResult.Error == nil {
			resultMap[taskResult.Result.(models.SCInfo).SCId] = taskResult.Result.(models.SCInfo).Name
		}
	}
	return resultMap, nil
}
