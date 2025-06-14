// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package keygen_test

import (
	"context"
	"testing"
	"time"
	comm2 "tss-demo/tss_util/comm"
	"tss-demo/tss_util/comm/elector"
	"tss-demo/tss_util/tss"
	"tss-demo/tss_util/tss/ecdsa/keygen"
	tsstest2 "tss-demo/tss_util/tss/test"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sourcegraph/conc/pool"
	"github.com/stretchr/testify/suite"
)

type KeygenTestSuite struct {
	tsstest2.CoordinatorTestSuite
}

func TestRunKeygenTestSuite(t *testing.T) {
	suite.Run(t, new(KeygenTestSuite))
}

func (s *KeygenTestSuite) Test_ValidKeygenProcess() {
	communicationMap := make(map[peer.ID]*tsstest2.TestCommunication)
	coordinators := []*tss.Coordinator{}
	processes := []tss.TssProcess{}

	for _, host := range s.CoordinatorTestSuite.Hosts {
		communication := tsstest2.TestCommunication{
			Host:          host,
			Subscriptions: make(map[comm2.SubscriptionID]chan *comm2.WrappedMessage),
		}
		communicationMap[host.ID()] = &communication
		keygen := keygen.NewKeygen("keygen", s.Threshold, host, &communication, s.MockECDSAStorer)
		electorFactory := elector.NewCoordinatorElectorFactory(host, s.BullyConfig)
		coordinators = append(coordinators, tss.NewCoordinator(host, &communication, electorFactory))
		processes = append(processes, keygen)
	}
	tsstest2.SetupCommunication(communicationMap)

	s.MockECDSAStorer.EXPECT().LockKeyshare().Times(3)
	s.MockECDSAStorer.EXPECT().UnlockKeyshare().Times(3)
	s.MockECDSAStorer.EXPECT().StoreKeyshare(gomock.Any()).Times(3)
	pool := pool.New().WithContext(context.Background()).WithCancelOnError()
	for i, coordinator := range coordinators {
		pool.Go(func(ctx context.Context) error { return coordinator.Execute(ctx, []tss.TssProcess{processes[i]}, nil) })
	}

	err := pool.Wait()
	s.Nil(err)
}

func (s *KeygenTestSuite) Test_KeygenTimeout() {
	communicationMap := make(map[peer.ID]*tsstest2.TestCommunication)
	coordinators := []*tss.Coordinator{}
	processes := []tss.TssProcess{}
	for _, host := range s.CoordinatorTestSuite.Hosts {
		communication := tsstest2.TestCommunication{
			Host:          host,
			Subscriptions: make(map[comm2.SubscriptionID]chan *comm2.WrappedMessage),
		}
		communicationMap[host.ID()] = &communication
		keygen := keygen.NewKeygen("keygen2", s.Threshold, host, &communication, s.MockECDSAStorer)
		electorFactory := elector.NewCoordinatorElectorFactory(host, s.BullyConfig)
		coordinator := tss.NewCoordinator(host, &communication, electorFactory)
		coordinator.TssTimeout = time.Millisecond
		coordinators = append(coordinators, coordinator)
		processes = append(processes, keygen)
	}
	tsstest2.SetupCommunication(communicationMap)

	s.MockECDSAStorer.EXPECT().LockKeyshare().AnyTimes()
	s.MockECDSAStorer.EXPECT().UnlockKeyshare().AnyTimes()
	s.MockECDSAStorer.EXPECT().StoreKeyshare(gomock.Any()).Times(0)
	pool := pool.New().WithContext(context.Background())
	for i, coordinator := range coordinators {
		pool.Go(func(ctx context.Context) error { return coordinator.Execute(ctx, []tss.TssProcess{processes[i]}, nil) })
	}

	err := pool.Wait()
	s.NotNil(err)
}
