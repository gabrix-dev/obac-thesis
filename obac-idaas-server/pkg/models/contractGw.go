package models

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

type ContractGateway struct {
	contracts map[SCId]ERC721
	callOpts  *bind.CallOpts
	RateConfig
	ticker       *time.Ticker
	taskChan     chan Task
	resourceChan chan bool
}

type RateConfig struct {
	NumThreads             int
	MaxRequestsPerInterval int
	IntervalDuration       time.Duration
}

type TaskResult struct {
	Result interface{}
	Error  error
}

type Task struct {
	Method     string
	TokenId    int64
	Contract   SCId
	ResultChan *chan TaskResult
}

func (c *ContractGateway) AddTask(task *Task) {
	c.taskChan <- *task
}

func NewContractGateway(config RateConfig) *ContractGateway {
	c := ContractGateway{RateConfig: config}
	c.contracts = make(map[SCId]ERC721)
	c.callOpts = &bind.CallOpts{Pending: false, Context: context.Background()}
	c.taskChan = make(chan Task, c.MaxRequestsPerInterval)
	c.resourceChan = make(chan bool, c.MaxRequestsPerInterval)
	c.ticker = time.NewTicker(c.IntervalDuration)
	for i := 0; i < c.NumThreads; i++ {
		go c.StartWorker()
	}
	go c.StartResourceManager()
	return &c
}

func (c *ContractGateway) StartResourceManager() {
	for {
		<-c.ticker.C
		nResources := len(c.resourceChan)
		for i := 0; i < nResources; i++ {
			<-c.resourceChan
		}
	}
}

func (c *ContractGateway) StartWorker() (interface{}, error) {
	for {
		task := <-c.taskChan
		instance, ok := c.contracts[task.Contract]
		if !ok {
			*task.ResultChan <- TaskResult{nil, errors.New("contract's instance not found")}
		}
		c.resourceChan <- true
		var result interface{}
		var err error
		switch task.Method {
		case "owner":
			result, err = instance.OwnerOf(c.callOpts, big.NewInt(task.TokenId))
			result = TokenIdInfo{TokenId: task.TokenId, Address: result.(common.Address)}
		case "name":
			result, err = instance.Name(c.callOpts)
			result = SCInfo{SCId: task.Contract, Name: result.(string)}
		case "uri":
			result, err = instance.TokenURI(c.callOpts, big.NewInt(task.TokenId))
			result = TokenIdInfo{TokenId: task.TokenId, TokenUri: result.(string)}
		}
		*task.ResultChan <- TaskResult{result, err}
	}
}

func (c *ContractGateway) GetContract(id SCId) (ERC721, error) {
	contract, ok := c.contracts[id]
	if !ok {
		return nil, errors.New("contract not found")
	}
	return contract, nil
}

func (c *ContractGateway) SetContract(id SCId, instance ERC721) {
	c.contracts[id] = instance
}
