// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package tsstest

import (
	"fmt"
	"os"
	"time"
	"tss-demo/tss_util/comm/mock"
	mock_tss2 "tss-demo/tss_util/tss/mock"
	"tss-demo/tss_util/tss_config/relayer"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/stretchr/testify/suite"
)

type CoordinatorTestSuite struct {
	suite.Suite
	GomockController  *gomock.Controller
	MockECDSAStorer   *mock_tss2.MockECDSAKeyshareStorer
	MockFrostStorer   *mock_tss2.MockFrostKeyshareStorer
	MockCommunication *mock_comm.MockCommunication
	MockTssProcess    *mock_tss2.MockTssProcess

	Hosts       []host.Host
	Threshold   int
	PartyNumber int
	BullyConfig relayer.BullyConfig
}

func (s *CoordinatorTestSuite) SetupTest() {
	s.GomockController = gomock.NewController(s.T())
	s.MockECDSAStorer = mock_tss2.NewMockECDSAKeyshareStorer(s.GomockController)
	s.MockFrostStorer = mock_tss2.NewMockFrostKeyshareStorer(s.GomockController)
	s.MockCommunication = mock_comm.NewMockCommunication(s.GomockController)
	s.MockTssProcess = mock_tss2.NewMockTssProcess(s.GomockController)
	s.PartyNumber = 3
	s.Threshold = 1

	hosts := []host.Host{}
	for i := 0; i < s.PartyNumber; i++ {
		host, _ := NewHost(i)
		hosts = append(hosts, host)
	}
	for _, host := range hosts {
		for _, peer := range hosts {
			host.Peerstore().AddAddr(peer.ID(), peer.Addrs()[0], peerstore.PermanentAddrTTL)
		}
	}
	s.Hosts = hosts
	s.BullyConfig = relayer.BullyConfig{
		PingWaitTime:     1 * time.Second,
		PingBackOff:      1 * time.Second,
		PingInterval:     1 * time.Second,
		ElectionWaitTime: 2 * time.Second,
		BullyWaitTime:    25 * time.Second,
	}
}

func NewHost(i int) (host.Host, error) {
	privBytes, err := os.ReadFile(fmt.Sprintf("../../test/pks/%d.pk", i))
	if err != nil {
		return nil, err
	}

	priv, err := crypto.UnmarshalPrivateKey(privBytes)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.Identity(priv),
		libp2p.DisableRelay(),
	}
	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func SetupCommunication(commMap map[peer.ID]*TestCommunication) {
	for self, comm := range commMap {
		peerComms := make(map[string]Receiver)
		for p, otherComm := range commMap {
			if self.Pretty() == p.Pretty() {
				continue
			}
			peerComms[p.Pretty()] = otherComm
		}
		comm.PeerCommunications = peerComms
	}
}
