// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package resharing_test

import (
	"context"
	"fmt"
	"testing"
	comm2 "tss-demo/tss_util/comm"
	"tss-demo/tss_util/comm/elector"
	"tss-demo/tss_util/keyshare"
	"tss-demo/tss_util/tss"
	"tss-demo/tss_util/tss/ecdsa/resharing"
	tsstest2 "tss-demo/tss_util/tss/test"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/sourcegraph/conc/pool"
	"github.com/stretchr/testify/suite"
)

type ResharingTestSuite struct {
	tsstest2.CoordinatorTestSuite
}

func TestRunResharingTestSuite(t *testing.T) {
	suite.Run(t, new(ResharingTestSuite))
}

func (s *ResharingTestSuite) Test_ValidResharingProcess_OldAndNewSubset() {
	communicationMap := make(map[peer.ID]*tsstest2.TestCommunication)
	coordinators := []*tss.Coordinator{}
	processes := []tss.TssProcess{}

	hosts := []host.Host{}
	for i := 0; i < s.PartyNumber+1; i++ {
		host, _ := tsstest2.NewHost(i)
		hosts = append(hosts, host)
	}
	for _, host := range hosts {
		for _, peer := range hosts {
			host.Peerstore().AddAddr(peer.ID(), peer.Addrs()[0], peerstore.PermanentAddrTTL)
		}
	}

	for i, host := range hosts {
		communication := tsstest2.TestCommunication{
			Host:          host,
			Subscriptions: make(map[comm2.SubscriptionID]chan *comm2.WrappedMessage),
		}
		communicationMap[host.ID()] = &communication
		storer := keyshare.NewECDSAKeyshareStore(fmt.Sprintf("../../test/keyshares/%d.keyshare", i))
		share, _ := storer.GetKeyshare()
		s.MockECDSAStorer.EXPECT().LockKeyshare()
		s.MockECDSAStorer.EXPECT().UnlockKeyshare()
		s.MockECDSAStorer.EXPECT().GetKeyshare().Return(share, nil)
		s.MockECDSAStorer.EXPECT().StoreKeyshare(gomock.Any()).Return(nil)
		resharing := resharing.NewResharing("resharing2", 1, host, &communication, s.MockECDSAStorer)
		electorFactory := elector.NewCoordinatorElectorFactory(host, s.BullyConfig)
		coordinators = append(coordinators, tss.NewCoordinator(host, &communication, electorFactory))
		processes = append(processes, resharing)
	}
	tsstest2.SetupCommunication(communicationMap)

	resultChn := make(chan interface{})
	pool := pool.New().WithContext(context.Background()).WithCancelOnError()
	for i, coordinator := range coordinators {
		pool.Go(func(ctx context.Context) error {
			return coordinator.Execute(ctx, []tss.TssProcess{processes[i]}, resultChn)
		})
	}

	err := pool.Wait()
	s.Nil(err)
}

func (s *ResharingTestSuite) Test_ValidResharingProcess_RemovePeer() {
	communicationMap := make(map[peer.ID]*tsstest2.TestCommunication)
	coordinators := []*tss.Coordinator{}
	processes := []tss.TssProcess{}

	hosts := []host.Host{}
	for i := 0; i < s.PartyNumber-1; i++ {
		host, _ := tsstest2.NewHost(i)
		hosts = append(hosts, host)
	}
	for _, host := range hosts {
		for _, peer := range hosts {
			host.Peerstore().AddAddr(peer.ID(), peer.Addrs()[0], peerstore.PermanentAddrTTL)
		}
	}

	for i, host := range hosts {
		communication := tsstest2.TestCommunication{
			Host:          host,
			Subscriptions: make(map[comm2.SubscriptionID]chan *comm2.WrappedMessage),
		}
		communicationMap[host.ID()] = &communication
		storer := keyshare.NewECDSAKeyshareStore(fmt.Sprintf("../../test/keyshares/%d.keyshare", i))
		share, _ := storer.GetKeyshare()
		s.MockECDSAStorer.EXPECT().LockKeyshare()
		s.MockECDSAStorer.EXPECT().UnlockKeyshare()
		s.MockECDSAStorer.EXPECT().GetKeyshare().Return(share, nil)
		s.MockECDSAStorer.EXPECT().StoreKeyshare(gomock.Any()).Return(nil)
		resharing := resharing.NewResharing("resharing2", 1, host, &communication, s.MockECDSAStorer)
		electorFactory := elector.NewCoordinatorElectorFactory(host, s.BullyConfig)
		coordinators = append(coordinators, tss.NewCoordinator(host, &communication, electorFactory))
		processes = append(processes, resharing)
	}
	tsstest2.SetupCommunication(communicationMap)

	resultChn := make(chan interface{})
	pool := pool.New().WithContext(context.Background()).WithCancelOnError()
	for i, coordinator := range coordinators {
		pool.Go(func(ctx context.Context) error {
			return coordinator.Execute(ctx, []tss.TssProcess{processes[i]}, resultChn)
		})
	}

	err := pool.Wait()
	s.Nil(err)
}

func (s *ResharingTestSuite) Test_InvalidResharingProcess_InvalidOldThreshold_LessThenZero() {
	communicationMap := make(map[peer.ID]*tsstest2.TestCommunication)
	coordinators := []*tss.Coordinator{}
	processes := []tss.TssProcess{}

	hosts := []host.Host{}
	for i := 0; i < s.PartyNumber+1; i++ {
		host, _ := tsstest2.NewHost(i)
		hosts = append(hosts, host)
	}
	for _, host := range hosts {
		for _, peer := range hosts {
			host.Peerstore().AddAddr(peer.ID(), peer.Addrs()[0], peerstore.PermanentAddrTTL)
		}
	}

	for i, host := range hosts {
		communication := tsstest2.TestCommunication{
			Host:          host,
			Subscriptions: make(map[comm2.SubscriptionID]chan *comm2.WrappedMessage),
		}
		communicationMap[host.ID()] = &communication
		storer := keyshare.NewECDSAKeyshareStore(fmt.Sprintf("../../test/keyshares/%d.keyshare", i))
		share, _ := storer.GetKeyshare()

		// set old threshold to invalid value
		share.Threshold = -1

		s.MockECDSAStorer.EXPECT().LockKeyshare().AnyTimes()
		s.MockECDSAStorer.EXPECT().UnlockKeyshare().AnyTimes()
		s.MockECDSAStorer.EXPECT().GetKeyshare().Return(share, nil)
		resharing := resharing.NewResharing("resharing3", 1, host, &communication, s.MockECDSAStorer)
		electorFactory := elector.NewCoordinatorElectorFactory(host, s.BullyConfig)
		coordinators = append(coordinators, tss.NewCoordinator(host, &communication, electorFactory))
		processes = append(processes, resharing)
	}
	tsstest2.SetupCommunication(communicationMap)

	resultChn := make(chan interface{})
	pool := pool.New().WithContext(context.Background())
	for i, coordinator := range coordinators {
		pool.Go(func(ctx context.Context) error {
			return coordinator.Execute(ctx, []tss.TssProcess{processes[i]}, resultChn)
		})
	}
	err := pool.Wait()
	s.NotNil(err)
}

func (s *ResharingTestSuite) Test_InvalidResharingProcess_InvalidOldThreshold_BiggerThenSubsetLength() {
	communicationMap := make(map[peer.ID]*tsstest2.TestCommunication)
	coordinators := []*tss.Coordinator{}
	processes := []tss.TssProcess{}

	hosts := []host.Host{}
	for i := 0; i < s.PartyNumber+1; i++ {
		host, _ := tsstest2.NewHost(i)
		hosts = append(hosts, host)
	}
	for _, host := range hosts {
		for _, peer := range hosts {
			host.Peerstore().AddAddr(peer.ID(), peer.Addrs()[0], peerstore.PermanentAddrTTL)
		}
	}

	for i, host := range hosts {
		communication := tsstest2.TestCommunication{
			Host:          host,
			Subscriptions: make(map[comm2.SubscriptionID]chan *comm2.WrappedMessage),
		}
		communicationMap[host.ID()] = &communication
		storer := keyshare.NewECDSAKeyshareStore(fmt.Sprintf("../../test/keyshares/%d.keyshare", i))
		share, _ := storer.GetKeyshare()

		// set old threshold to invalid value
		share.Threshold = 314

		s.MockECDSAStorer.EXPECT().LockKeyshare()
		s.MockECDSAStorer.EXPECT().UnlockKeyshare().AnyTimes()
		s.MockECDSAStorer.EXPECT().GetKeyshare().Return(share, nil)
		resharing := resharing.NewResharing("resharing4", 1, host, &communication, s.MockECDSAStorer)
		electorFactory := elector.NewCoordinatorElectorFactory(host, s.BullyConfig)
		coordinators = append(coordinators, tss.NewCoordinator(host, &communication, electorFactory))
		processes = append(processes, resharing)
	}
	tsstest2.SetupCommunication(communicationMap)

	resultChn := make(chan interface{})
	pool := pool.New().WithContext(context.Background())
	for i, coordinator := range coordinators {
		pool.Go(func(ctx context.Context) error {
			return coordinator.Execute(ctx, []tss.TssProcess{processes[i]}, resultChn)
		})
	}

	err := pool.Wait()
	s.NotNil(err)
}
