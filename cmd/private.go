/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/spf13/cobra"
)

// privateCmd represents the private command
var privateCmd = &cobra.Command{
	Use:   "private",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: privateChannel,
}

func init() {
	transmitterCmd.AddCommand(privateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// privateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// privateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var MessageHistory = make([]string, 0)

func privateChannel(cmd *cobra.Command, args []string) {
	sk1, pk1, _, _ := generateKey(User)
	UserPubKey = append(UserPubKey, []string{pk1, sk1})

	_, pk2, _, npub2 := generateKey((User + 1) % 2)

	// Channel creation
	messages := make(chan string)

	go func() {
		subscribeToRelayPrivateChannel(sk1, npub2, messages)
	}()

	// go func() {
	// 	for {
	// 		// Wait for message to be received on channel
	// 		message := <-messages
	// 		fmt.Println(message)
	// 		MessageHistory = append(MessageHistory, message)
	// 	}
	// }()

	// for {
	// 	var message string
	// 	fmt.Scanln(&message)
	// 	fmt.Println("You: " + message)
	// 	if message == "exit" {
	// 		break
	// 	}
	// 	sendPrivateMessage(sk1, pk1, pk2, message)
	// }

	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()

	// Create chat history widget
	history := widgets.NewList()
	history.Title = "Chat History"
	history.Rows = []string{}

	// Create user input widget
	input := widgets.NewParagraph()
	input.Title = "User Input"
	input.Text = "> "

	// Create chat layout
	// layout := termui.NewVBox(history, input) This is the old version
	layout := termui.NewGrid()
	layout.Set(
		termui.NewRow(1.0/10*9, history),
		termui.NewRow(1.0/10, input),
	)
	layout.SetRect(0, 0, 50, 30)

	// Create channel for receiving messages

	// Start goroutine for receiving messages
	go func() {
		for {
			// Wait for message to be received on channel
			message := <-messages

			// Add message to chat history
			history.Rows = append(history.Rows, message)
			history.ScrollBottom()

			// Update chat history widget
			termui.Render(history)
		}
	}()

	// Handle user input
	termui.Render(layout)
	termuiEvents := termui.PollEvents()

	for {
		e := <-termuiEvents
		if e.Type == termui.KeyboardEvent {
			switch e.ID {
			case "<C-c>":
				// Quit application on Ctrl+C
				return
			case "<Enter>":
				// Send message on Enter
				message := input.Text[2:]

				// Clear user input
				input.Text = "> "

				go sendPrivateMessage(sk1, pk1, pk2, message)

				// Update user input widget
				termui.Render(input)
			case "<Backspace>":
				// Remove last character from user input
				input.Text = input.Text[:len(input.Text)-1]
				termui.Render(input)
			case "<Space>":
				// Add space to user input
				input.Text += " "
				termui.Render(input)
			case "<Tab>":
				// Add tab to user input
				input.Text += "\t"
				termui.Render(input)
			case "<Escape>":
				// Clear user input
				input.Text = "> "
				termui.Render(input)

			default:
				// Add character to user input
				input.Text += string(e.ID)
				termui.Render(input)
			}

		}
	}

}

func sendPrivateMessage(sk1 string, pk1 string, pk2 string, message string) {
	sharedKey, err := nip04.ComputeSharedSecret(pk2, sk1)
	if err != nil {
		fmt.Println(err)
	}

	encryptedMessage, err := nip04.Encrypt(message, sharedKey)

	if err != nil {
		fmt.Println(err)
	}

	publishToRelayPrivateMessage(sk1, pk1, encryptedMessage, pk2)

}

func publishToRelayPrivateMessage(sk string, pub string, encryptedMessage string, to string) {

	ev := nostr.Event{
		PubKey:    pub,
		CreatedAt: time.Now(),
		Kind:      4,
		Tags:      nostr.Tags{{"p", to}},
		Content:   encryptedMessage,
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

		relay.Publish(context.Background(), ev)

		count++
	}

}

func receivePrivateMessage(sk string, pub string, encryptedMessage string) string {

	sharedKey, err := nip04.ComputeSharedSecret(pub, sk)
	if err != nil {
		fmt.Println(err)
	}

	decryptedMessage, err := nip04.Decrypt(encryptedMessage, sharedKey)

	if err != nil {
		fmt.Println("Error decrypting message: ", err)
		fmt.Println("pub: ", pub)
	}

	return decryptedMessage

}

func subscribeToRelayPrivateChannel(sk string, npub string, eventChan chan string) {
	relay, err := nostr.RelayConnect(context.Background(), "wss://knostr.neutrine.com")
	if err != nil {
		panic(err)
	}

	var filters nostr.Filters
	if _, _, err := nip19.Decode(npub); err == nil {

		filters = []nostr.Filter{{
			Kinds: []int{4},
			Limit: 5,
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

		// Write the received event to the channel
		isMe := false
		for potentialKeyIndex := range UserPubKey {
			if UserPubKey[potentialKeyIndex][0] == ev.PubKey {
				isMe = true
				decodedMessage := receivePrivateMessage(UserPubKey[potentialKeyIndex][1], ev.Tags.GetFirst([]string{"p"}).Value(), ev.Content)

				eventChan <- "Sent: " + decodedMessage
			}
		}
		if !isMe {
			decodedMessage := receivePrivateMessage(sk, ev.PubKey, ev.Content)

			eventChan <- "Received: " + decodedMessage
		}
	}
}
