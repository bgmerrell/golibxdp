package main

import (
	"github.com/mdlayher/netlink"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/unix"
)

func status(ifname string) {
	if ifname == "" {
		log.Debug().Msgf("Getting status for all interfaces")
	} else {
		log.Debug().Msgf("Getting status for: %s", ifname)
	}

	conn, err := netlink.Dial(unix.NETLINK_ROUTE, nil)
	if err != nil {
		log.Fatal().Msgf("Failed to Dial NETLINK_ROUTE: %v", err)
	}
	defer conn.Close()

	if err = conn.SetOption(netlink.ExtendedAcknowledge, true); err != nil {
		log.Fatal().Msgf("Failed to set NETLINK_EXT_ACK option: %v", err)
	}

	// TODO: Finish implementing
	log.Warn().Msg("Command not fully implemented")
}
