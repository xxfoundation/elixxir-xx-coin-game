////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Package cmd initializes the CLI and config parsers as well as the logger

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"gitlab.com/xx_network/primitives/utils"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xx-coin-game",
	Short: "Runs the xx coin game",
	Long:  `This binary provides a bot wrapping client`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		// Main program initialization here
		fmt.Printf("Hello, World!")
		//FIXME: Dump config here..
		//fmt.Printf("%+v",
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
	rootCmd.Flags().StringP("config", "c", "",
		"Path to load the configuration file from. If not set, this "+
			"file must be named xx-coin-game.yaml and must be "+
			"located in ~/.xxnetwork/, /opt/xxnetwork, "+
			" or /etc/xxnetwork.")

	rootCmd.Flags().UintP("logLevel", "l", 1,
		"Level of debugging to display. "+
			"0 = info, 1 = debug, >1 = trace")
	viper.BindPFlag("logLevel",
		rootCmd.PersistentFlags().Lookup("logLevel"))

	rootCmd.Flags().StringP("filePath", "f",
		"", "Sets the address file path")
	viper.BindPFlag("filePath",
		rootCmd.PersistentFlags().Lookup("filePath"))

	rootCmd.Flags().StringP("logPath", "l",
		"", "Sets the log file path")
	viper.BindPFlag("logPath", rootCmd.PersistentFlags().Lookup("logPath"))

	rootCmd.PersistentFlags().BoolP("verbose", "v", false,
		"Verbose mode for debugging")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().StringP("session", "s",
		"", "Sets the initial storage directory for "+
			"client session data")
	viper.BindPFlag("session", rootCmd.PersistentFlags().Lookup("session"))

	rootCmd.PersistentFlags().StringP("writeContact", "w",
		"-", "Write contact information, if any, to this file, "+
			" defaults to stdout")
	viper.BindPFlag("writeContact", rootCmd.PersistentFlags().Lookup(
		"writeContact"))

	rootCmd.PersistentFlags().StringP("password", "p", "",
		"Password to the session file")
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup(
		"password"))

	rootCmd.PersistentFlags().StringP("ndf", "n", "ndf.json",
		"Path to the network definition JSON file")
	viper.BindPFlag("ndf", rootCmd.PersistentFlags().Lookup("ndf"))

	rootCmd.PersistentFlags().StringP("log", "l", "-",
		"Path to the log output path (- is stdout)")
	viper.BindPFlag("log", rootCmd.PersistentFlags().Lookup("log"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	cfgFile := viper.GetString("config")
	if cfgFile == "" {
		var err error
		cfgFile, err = utils.SearchDefaultLocations("xx-coin-game.yaml",
			"xxnetwork")
		if err != nil {
			jww.FATAL.Panicf("Failed to find config file: %+v", err)
		}
	}
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		jww.FATAL.Panicf("Unable to read config file (%s): %+v",
			cfgFile, err.Error())
	}

}

// initLog initializes logging thresholds and the log path.
func initLog() {
	logPath := viper.Get("logPath")
	if len(logPath) > 0 {
		logLevel := viper.Get("logLevel")
		// Check the level of logs to display
		if logLevel > 1 {
			// Set the GRPC log level
			err := os.Setenv("GRPC_GO_LOG_SEVERITY_LEVEL", "info")
			if err != nil {
				jww.ERROR.Printf("Could not set "+
					"GRPC_GO_LOG_SEVERITY_LEVEL: %+v", err)
			}

			err = os.Setenv("GRPC_GO_LOG_VERBOSITY_LEVEL", "99")
			if err != nil {
				jww.ERROR.Printf("Could not set "+
					"GRPC_GO_LOG_VERBOSITY_LEVEL: %+v", err)
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
			jww.WARN.Println("Invalid or missing log path, " +
				"default path used.")
		} else {
			jww.SetLogOutput(logFile)
		}
	}
}
