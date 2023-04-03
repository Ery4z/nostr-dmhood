module github.com/Ery4z/nostr-dmhood

go 1.19

// q: How can i know where my go modules are stored?
// a: https://stackoverflow.com/questions/58033366/how-can-i-know-where-my-go-modules-are-stored
require (
	github.com/SaveTheRbtz/generic-sync-map-go v0.0.0-20220414055132-a37292614db8 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/btcsuite/btcd/btcutil v1.1.3 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.1 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.0.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/gizak/termui/v3 v3.1.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-runewidth v0.0.2 // indirect
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/nbd-wtf/go-nostr v0.15.3 // indirect
	github.com/nsf/termbox-go v0.0.0-20190121233118-02980233997d // indirect
	github.com/spf13/cobra v1.6.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stackerstan/go-nostr v0.0.0 // indirect
	github.com/valyala/fastjson v1.6.3 // indirect
	golang.org/x/exp v0.0.0-20221106115401-f9659909a136 // indirect
)

replace github.com/nbd-wtf/go-nostr => ./local/go-nostr
