package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type loggerConfig struct {
	isVerbose bool
	isDebug   bool
}

func setupLogger(loggerCfg loggerConfig) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano})

	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	if loggerCfg.isDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else if loggerCfg.isVerbose {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func main() {

	var (
		isVerbose bool
		isDebug   bool
	)

	var cmdLoad = &cobra.Command{
		Use:   "load [flags] <ifname> <filenames>",
		Short: "Load an XDP program on an interface",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: implement
			log.Error().Msg("Command not implemented")
		},
	}

	var cmdUnload = &cobra.Command{
		Use:   "unload [flags] <ifname>",
		Short: "Unload an XDP program from an interface",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: implement
			log.Error().Msg("Command not implemented")
		},
	}

	var (
		ifname string
	)
	var cmdStatus = &cobra.Command{
		Use:   "status [flags]",
		Short: "Unload an XDP program from an interface",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			status(ifname)
		},
	}
	cmdStatus.Flags().StringVar(&ifname, "ifname", "", "interface name")

	var rootCmd = &cobra.Command{
		Use: "xdp-loader [flags]",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setupLogger(loggerConfig{
				isVerbose,
				isDebug,
			})
		},
	}
	rootCmd.PersistentFlags().BoolVarP(
		&isVerbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(
		&isDebug, "debug", "d", false, "debug (more verbose) output")
	rootCmd.AddCommand(cmdLoad, cmdUnload, cmdStatus)
	rootCmd.Execute()
}
