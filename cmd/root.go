// Copyright (c) 2021 patrick-ogrady
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/ava-labs/avalanchego/utils/constants"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	stakingKeyFile  = "staker.key"
	stakingCertFile = "staker.crt"
)

// Source: https://github.com/ava-labs/avalanchego/blob/e2944176f9e87562140ecd979cafebb4707578c4/main/params.go#L44-L53
var (
	// homeDir is the $HOME directory to use
	// when creating the .avalanchego folder.
	homeDir = os.ExpandEnv("$HOME")

	// stakingDirectory is the directory containing
	// the staking key and certificate.
	stakingDirectory = filepath.Join(
		fmt.Sprintf(".%s", constants.AppName),
		"staking",
	)

	// stakingKeyPath is filepath containing staking key.
	stakingKeyPath = filepath.Join(stakingDirectory, stakingKeyFile)

	// stakingCertPath is filepath containing staking certificate.
	stakingCertPath = filepath.Join(stakingDirectory, stakingCertFile)
)

var (
	// Context is the context to use for this invocation of the cli.
	Context context.Context

	// Cancel is the context.CancelFunc for this invocation of the cli.
	Cancel context.CancelFunc

	// SignalReceived is set to true when a signal causes us to exit. This makes
	// determining the error message to show on exit much more easy.
	SignalReceived = false
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "avalanche-runner",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.avalanche-runner.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Setup signal handling for context
	ctx := context.Background()
	Context, Cancel = context.WithCancel(ctx)
	go handleSignals([]context.CancelFunc{Cancel})

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".avalanche-runner" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".avalanche-runner")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// handleSignals handles OS signals so we can ensure we close database
// correctly. We call multiple sigListeners because we
// may need to cancel more than 1 context.
func handleSignals(listeners []context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Printf("received signal: %s\n", sig)
		SignalReceived = true
		for _, listener := range listeners {
			listener()
		}
	}()
}
