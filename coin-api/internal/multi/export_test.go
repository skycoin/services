package multi

import (
	"github.com/skycoin/skycoin/src/api/cli"
)

func (s *SkyСoinService) InjectRPCAPIMock(mockRPCAPI WebRPCClientAPI) {
	s.client = mockRPCAPI
}

func (s *SkyСoinService) InjectCheckBalanceMock(mockCheckBalance func(client WebRPCClientAPI, addresses []string) (*cli.BalanceResult, error)) {
	s.checkBalance = mockCheckBalance
}
