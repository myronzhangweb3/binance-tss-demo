// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package tsstest

import (
	"fmt"
	"sync"
	"time"
	comm2 "tss-demo/tss_util/comm"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Receiver interface {
	ReceiveMessage(msg *comm2.WrappedMessage, topic comm2.MessageType, sessionID string)
}

type TestCommunication struct {
	Host               host.Host
	Subscriptions      map[comm2.SubscriptionID]chan *comm2.WrappedMessage
	PeerCommunications map[string]Receiver
	lock               sync.Mutex
}

func (tc *TestCommunication) Broadcast(
	peers peer.IDSlice,
	msg []byte,
	msgType comm2.MessageType,
	sessionID string,
) error {
	wMsg := comm2.WrappedMessage{
		MessageType: msgType,
		SessionID:   sessionID,
		Payload:     msg,
		From:        tc.Host.ID(),
	}

	time.Sleep(100 * time.Millisecond)
	for _, peer := range peers {
		if tc.PeerCommunications[peer.Pretty()] == nil {
			continue
		}

		go tc.PeerCommunications[peer.Pretty()].ReceiveMessage(&wMsg, msgType, sessionID)
	}

	return nil
}

func (ts *TestCommunication) Subscribe(
	sessionID string,
	topic comm2.MessageType,
	channel chan *comm2.WrappedMessage,
) comm2.SubscriptionID {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	subID := comm2.SubscriptionID(fmt.Sprintf("%s-%s", sessionID, topic))
	ts.Subscriptions[subID] = channel
	return subID
}

func (ts *TestCommunication) UnSubscribe(subscriptionID comm2.SubscriptionID) {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	delete(ts.Subscriptions, subscriptionID)
}

func (ts *TestCommunication) ReceiveMessage(msg *comm2.WrappedMessage, topic comm2.MessageType, sessionID string) {
	// simulate real world conditions
	ts.Subscriptions[comm2.SubscriptionID(fmt.Sprintf("%s-%s", sessionID, topic))] <- msg
}

func (ts *TestCommunication) CloseSession(sessionID string) {}
