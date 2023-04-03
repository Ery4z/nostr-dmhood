/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nostr-dmhood",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: runTransmiter,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var User int
var UserPubKey = make([][]string, 0)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nostr-dmhood.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().IntVarP(&User, "user", "u", 1, "User 1 or 2")
}

func runReceiver(cmd *cobra.Command, args []string) {
	_, _, _, npub := generateKey(User)
	subscribeToRelay(npub)

}

func runTransmiter(cmd *cobra.Command, args []string) {
	sk, pk, _, _ := generateKey(User)

	for true {
		// Get the message from the user
		var message string
		fmt.Print("Enter message: ")
		fmt.Scanln(&message)

		publishToRelay(sk, pk, message)

	}

}

func generateKey(num int) (string, string, string, string) {
	// sk := nostr.GeneratePrivateKey()
	// pk, _ := nostr.GetPublicKey(sk)
	// nsec, _ := nip19.EncodePrivateKey(sk)
	// npub, _ := nip19.EncodePublicKey(pk)
	var sk, pk, nsec, npub string
	if num == 0 {

		sk = "23c3ce90e9dcae090f1da4562a10ed2f9d52994221508ff6ada60d7c1ed9a40d"
		pk = "0a364cbbef49135a150477e5ce0babef5e4475414006cfab3a7795c9058454bb"
		nsec = "nsec1y0puay8fmjhqjrca53tz5y8d97w49x2zy9ggla4d5cxhc8ke5sxsnenz98"
		npub = "npub1pgmyewl0fyf459gywljuuzataa0yga2pgqrvl2e6w72ujpvy2jas6vm8tm"
	} else if num == 1 {
		sk = "a54054cbb180c7f1f067253bbe379c0fc29593d2cf8d7d5902291f113cebf2a0"
		pk = "325b4427de846afef143c85eb03024ac66dbf30b51daa77c82bb86809455f224"
		nsec = "nsec154q9fja3srrlrur8y5amuduuplpfty7je7xh6kgz9y03z08t72sqzl5txq"
		npub = "npub1xfd5gf77s340au2rep0tqvpy43ndhuct28d2wlyzhwrgp9z47gjqflgezs"

	} else {
		sk = nostr.GeneratePrivateKey()
		pk, _ = nostr.GetPublicKey(sk)
		nsec, _ = nip19.EncodePrivateKey(sk)
		npub, _ = nip19.EncodePublicKey(pk)

		fmt.Println("sk: ", sk)
		fmt.Println("pk: ", pk)
		fmt.Println("nsec: ", nsec)
		fmt.Println("npub: ", npub)

	}
	return sk, pk, nsec, npub
}

func subscribeToRelay(npub string) {
	relay, err := nostr.RelayConnect(context.Background(), "wss://knostr.neutrine.com")
	if err != nil {
		panic(err)
	}

	var filters nostr.Filters
	if _, _, err := nip19.Decode(npub); err == nil {
		filters = []nostr.Filter{{
			Kinds: []int{1},
			Limit: 1,
		}}
	} else {
		panic(err)
	}

	ctx, _ := context.WithCancel(context.Background())
	sub := relay.Subscribe(ctx, filters)

	go func() {
		<-sub.EndOfStoredEvents
		// handle end of stored events (EOSE, see NIP-15)
	}()

	for ev := range sub.Events {
		// handle returned event.
		// channel will stay open until the ctx is cancelled (in this case, by calling cancel())

		fmt.Println("Received event: ", ev.Content+" from "+ev.PubKey)
	}

}

func publishToRelay(sk string, pub string, message string) {

	ev := nostr.Event{
		PubKey:    pub,
		CreatedAt: time.Now(),
		Kind:      1,
		Tags:      nil,
		Content:   message,
	}

	// calling Sign sets the event ID field and the event Sig field
	ev.Sign(sk)

	// publish the event to two relays
	count := 0
	for _, url := range []string{"wss://knostr.neutrine.com", "wss://nostr.einundzwanzig.space", "wss://nostr-pub.wellorder.net"} {
		relay, e := nostr.RelayConnect(context.Background(), url)
		if e != nil {
			fmt.Println(e)
			continue
		}

		status, _ := relay.Publish(context.Background(), ev)

		fmt.Println(count, " published to ", url, status)
		count++
	}

}
