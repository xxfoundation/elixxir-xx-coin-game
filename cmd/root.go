////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Package cmd initializes the CLI and config parsers as well as the logger

package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/interfaces/contact"
	"gitlab.com/elixxir/client/interfaces/message"
	"gitlab.com/elixxir/client/interfaces/params"
	"gitlab.com/elixxir/client/single"
	"gitlab.com/elixxir/client/switchboard"
	"gitlab.com/elixxir/xx-coin-game/crypto"
	"gitlab.com/elixxir/xx-coin-game/game"
	"gitlab.com/elixxir/xx-coin-game/io"
	"io/ioutil"
	"os"
	"time"
)

var (
	logPath     string
	filePath    string
	logLevel    uint
	session     string
	contactPath string
	password    string
	ndfPath     string
	salt        []byte
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xx-coin-game",
	Short: "Runs the xx coin game",
	Long:  `This binary provides a bot wrapping client`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		// Main program initialization here
		addressMap, addressWriteCh := io.StartIo(filePath)

		gameMap := game.New(addressMap, salt, crypto.NewRng())

		client := initClient()

		// Write user contact to file
		user := client.GetUser()
		jww.INFO.Printf("User: %s", user.ReceptionID)
		jww.INFO.Printf("User Transmission: %s", user.TransmissionID)
		writeContact(user.GetContact())

		// Set up reception handler
		swBoard := client.GetSwitchboard()
		recvCh := make(chan message.Receive, 10000)
		listenerID := swBoard.RegisterChannel("DefaultCLIReceiver",
			switchboard.AnyUser(), message.Text, recvCh)
		jww.INFO.Printf("Message ListenerID: %v", listenerID)

		// Set up auth request handler, which simply prints the user ID of the
		// requester
		//authMgr := client.GetAuthRegistrar()

		// NOTE: We would only need this to support E2E messages.
		// authMgr.AddGeneralRequestCallback(func(
		// 	requester contact.Contact, message string) {
		// 		jww.INFO.Printf("Got request: %s", requester.ID)
		// 		err := client.ConfirmAuthenticatedChannel(requester)
		// 		if err != nil {
		// 			jww.FATAL.Panicf("%+v", err)
		// 		}
		// 	})

		_, err := client.StartNetworkFollower()
		if err != nil {
			jww.FATAL.Panicf("%+v", err)
		}

		// Wait until connected or crash on timeout
		connected := make(chan bool, 10)
		client.GetHealth().AddChannel(connected)
		waitUntilConnected(connected)

		// Make single-use manager and start receiving process
		singleMng := single.NewManager(client)

		// Register the callback

		callback := func(payload []byte, c single.Contact) {
			if payload == nil {
				jww.WARN.Printf("Empty payload from %s",
					c.GetPartner())
				return
			}

			//process the message
			address, text, err := crypto.HandleMessage(string(payload))
			if err != nil {
				jww.WARN.Printf("Payload %s from %s failed handling: %s",
					string(payload), c.GetPartner(), err.Error())
				err := singleMng.RespondSingleUse(c, []byte(err.Error()), 30*time.Second)
				if err != nil {
					jww.WARN.Printf("Failed to transmit resonce to %s: %+v",
						c.GetPartner(), err)
				}
			}

			new, value, err := gameMap.Play(address, string(payload))

			if err != nil {
				jww.WARN.Printf("Address %s from %s could nto be found: %s",
					address, c.GetPartner(), err.Error())
				err = singleMng.RespondSingleUse(c, []byte(err.Error()), 30*time.Second)
				if err != nil {
					jww.WARN.Printf("Failed to transmit resonce to %s: %+v",
						c.GetPartner(), err)
				}
			}

			message := fmt.Sprintf("Address %s said %s and won %d xx coins!", address, text, value)

			if new {
				addressWriteCh <- io.AddressUpdate{
					Address: address,
					Value:   uint64(value),
				}
				jww.INFO.Println(message)
			}

			err = singleMng.RespondSingleUse(c, []byte(message), 30*time.Second)
			if err != nil {
				jww.WARN.Printf("Failed to transmit resonce to %s: %+v",
					c.GetPartner(), err)
			}
		}
		singleMng.RegisterCallback("xxCoinGame", callback)
		client.AddService(singleMng.StartProcesses)

		// Wait to receive a message or stop after timeout occurs
		fmt.Println("Bot Started...")
		select {}
	},
}

// Execute adds all child commands to the root command and sets flags
// appropriately.  This is called by main.main(). It only needs to
// happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		jww.ERROR.Println(err)
		os.Exit(1)
	}
}

// init is the initialization function for Cobra which defines commands
// and flags.
func init() {
	// NOTE: The point of init() is to be declarative.
	// There is one init in each sub command. Do not put variable declarations
	// here, and ensure all the Flags are of the *P variety, unless there's a
	// very good reason not to have them as local params to sub command."
	cobra.OnInitialize(initLog)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.Flags().UintVarP(&logLevel, "logLevel", "", 1,
		"Level of debugging to display. 0 = info, 1 = debug, >1 = trace")

	rootCmd.Flags().StringVarP(&filePath, "filePath", "f",
		"", "Sets the address file path")

	rootCmd.Flags().StringVarP(&logPath, "logPath", "l",
		"-", "Sets the log file path")

	rootCmd.Flags().StringVarP(&session, "session", "s",
		"", "Sets the initial storage directory for "+
			"client session data")

	rootCmd.Flags().StringVarP(&contactPath, "contactPath", "w",
		"-", "Write contact information, if any, to this file, "+
			" defaults to stdout")

	rootCmd.Flags().StringVarP(&password, "password", "p", "",
		"Password to the session file")

	rootCmd.Flags().StringVarP(&ndfPath, "ndf", "n", "ndf.json",
		"Path to the network definition JSON file")

	rootCmd.Flags().BytesHexVar(&salt, "salt", make([]byte, 32), "Default value of salt")
}

// initLog initializes logging thresholds and the log path.
func initLog() {
	if len(logPath) > 0 {

		// Check the level of logs to display
		if logLevel > 1 {
			// Set the GRPC log level
			err := os.Setenv("GRPC_GO_LOG_SEVERITY_LEVEL", "info")
			if err != nil {
				jww.ERROR.Printf("Could not set GRPC_GO_LOG_SEVERITY_LEVEL: %+v", err)
			}

			err = os.Setenv("GRPC_GO_LOG_VERBOSITY_LEVEL", "99")
			if err != nil {
				jww.ERROR.Printf("Could not set GRPC_GO_LOG_VERBOSITY_LEVEL: %+v", err)
			}
			// Turn on trace logs
			jww.SetLogThreshold(jww.LevelTrace)
			jww.SetStdoutThreshold(jww.LevelTrace)
		} else if logLevel == 1 {
			// Turn on debugging logs
			jww.SetLogThreshold(jww.LevelDebug)
			jww.SetStdoutThreshold(jww.LevelDebug)
		} else {
			// Turn on info logs
			jww.SetLogThreshold(jww.LevelInfo)
			jww.SetStdoutThreshold(jww.LevelInfo)
		}

		// Create log file, overwrites if existing
		logFile, err := os.OpenFile(logPath,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0644)
		if err != nil {
			jww.WARN.Println("Invalid or missing log path, default path used.")
		} else {
			jww.SetLogOutput(logFile)
		}
	}
}

func createClient() *api.Client {
	pass := password
	storeDir := session

	//create a new client if none exist
	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		// Load NDF
		ndfJSON, err := ioutil.ReadFile(ndfPath)
		if err != nil {
			jww.FATAL.Panicf(err.Error())
		}

		err = api.NewClient(string(ndfJSON), storeDir,
			[]byte(pass), "")
		if err != nil {
			jww.FATAL.Panicf("%+v", err)
		}
	}

	netParams := params.GetDefaultNetwork()
	client, err := api.OpenClient(storeDir, []byte(pass), netParams)
	if err != nil {
		jww.FATAL.Panicf("%+v", err)
	}
	return client
}

func initClient() *api.Client {
	createClient()

	pass := password
	storeDir := session

	netParams := params.GetDefaultNetwork()
	client, err := api.Login(storeDir, []byte(pass), netParams)
	if err != nil {
		jww.FATAL.Panicf("%+v", err)
	}

	return client
}

func writeContact(c contact.Contact) {
	outfilePath := contactPath
	if outfilePath == "" {
		return
	}
	err := ioutil.WriteFile(outfilePath, c.Marshal(), 0644)
	if err != nil {
		jww.FATAL.Panicf("%+v", err)
	}
}

func waitUntilConnected(connected chan bool) {
	timeoutTimer := time.NewTimer(90 * time.Second)
	isConnected := false
	//Wait until we connect or panic if we can't by a timeout
	for !isConnected {
		select {
		case isConnected = <-connected:
			jww.INFO.Printf("Network Status: %v\n",
				isConnected)
			break
		case <-timeoutTimer.C:
			jww.FATAL.Panic("timeout on connection")
		}
	}

	// Now start a thread to empty this channel and update us
	// on connection changes for debugging purposes.
	go func() {
		prev := true
		for {
			select {
			case isConnected = <-connected:
				if isConnected != prev {
					prev = isConnected
					jww.INFO.Printf(
						"Network Status Changed: %v\n",
						isConnected)
				}
				break
			}
		}
	}()
}

// responseCallbackChan structure used to collect information sent to the
// response callback.
type responseCallbackChan struct {
	payload []byte
	c       single.Contact
}

// makeResponsePayload generates a new payload that will span the max number of
// message parts in the contact. Each resulting message payload will contain a
// copy of the supplied payload with spaces taking up any remaining data.
func makeResponsePayload(m *single.Manager, payload []byte, maxParts uint8) []byte {
	payloads := make([][]byte, maxParts)
	payloadPart := makeResponsePayloadPart(m, payload)
	for i := range payloads {
		payloads[i] = make([]byte, m.GetMaxResponsePayloadSize())
		copy(payloads[i], payloadPart)
	}
	return bytes.Join(payloads, []byte{})
}

// makeResponsePayloadPart creates a single response payload by coping the given
// payload and filling the rest with spaces.
func makeResponsePayloadPart(m *single.Manager, payload []byte) []byte {
	payloadPart := make([]byte, m.GetMaxResponsePayloadSize())
	for i := range payloadPart {
		payloadPart[i] = ' '
	}
	copy(payloadPart, payload)

	return payloadPart
}
