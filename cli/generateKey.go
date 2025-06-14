// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package cli

import (
	"encoding/base64"
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"
	peer "github.com/libp2p/go-libp2p/core/peer"
)

const (
	KEY_LENGTH = 2048
)

func generateKey() error {
	priv, pub, err := crypto.GenerateKeyPair(crypto.RSA, KEY_LENGTH)
	if err != nil {
		return err
	}

	peerID, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return err
	}

	marshPriv, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return err
	}
	encPriv := base64.StdEncoding.EncodeToString(marshPriv)

	fmt.Printf(`
LibP2P peer identity: %s \n
LibP2P private key: %s
`,
		peerID.Pretty(),
		encPriv,
	)
	return nil
}
