package eth

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// TODO(stgleb): Add circuit breaker for ethclient.Client
type EthService struct {
	client *ethclient.Client
}

func NewEthService(nodeUrl string) (*EthService, error) {
	c, err := rpc.Dial(nodeUrl)

	if err != nil {
		return nil, err
	}

	client := ethclient.NewClient(c)

	return &EthService{
		client: client,
	}, nil

}

func (s *EthService) GenerateKeyPair() (string, string, error) {
	key, err := crypto.GenerateKey()

	if err != nil {
		return "", "", err
	}

	address := crypto.PubkeyToAddress(key.PublicKey)

	return hex.EncodeToString(key.D.Bytes()), address.String(), nil
}

func (s *EthService) GetBalance(ctx context.Context, address common.Address) (int64, error) {
	balance, err := s.client.BalanceAt(ctx, address, nil)

	return balance.Int64(), err
}

func (s *EthService) GetTxStatus(ctx context.Context, txid string) (*types.Transaction, bool, error) {
	hash := common.HexToHash(txid)

	txStatus, isPending, err := s.client.TransactionByHash(ctx, hash)

	if err != nil {
		return nil, false, err
	}

	return txStatus, isPending, nil
}
