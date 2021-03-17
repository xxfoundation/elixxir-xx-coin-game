////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Package cmd initializes the CLI and config parsers as well as the logger

package cmd

import (
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/xx-coin-game/io"
	"os"
)

var (
	logPath      string
	filePath     string
	logLevel     uint
	session      string
	writeContact string
	password     string
	ndfPath      string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xx-coin-game",
	Short: "Runs the xx coin game",
	Long:  `This binary provides a bot wrapping client`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		// Main program initialization here
		_, _ = io.StartIo(filePath)

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

	rootCmd.Flags().StringVarP(&writeContact, "writeContact", "w",
		"-", "Write contact information, if any, to this file, "+
			" defaults to stdout")

	rootCmd.Flags().StringVarP(&password, "password", "p", "",
		"Password to the session file")

	rootCmd.Flags().StringVarP(&ndfPath, "ndf", "n", "ndf.json",
		"Path to the network definition JSON file")
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
