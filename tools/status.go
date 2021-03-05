package main

import (
	"fmt"

	"github.com/jsimonetti/rtnetlink"
	"github.com/rs/zerolog/log"
)

func status(ifname string) error {
	// outTemplate := "%-16s %-5v %-17s %-8s %-4v %-17s %s\n"
	if ifname == "" {
		log.Debug().Msgf("Getting status for all interfaces")
	} else {
		log.Debug().Msgf("Getting status for: %s", ifname)
	}

	// Dial a connection to the rtnetlink socket
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	defer conn.Close()

	// Request a list of interfaces
	msgs, err := conn.Link.List()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	for _, msg := range msgs {
		if msg.Attributes.XDP.Attached > 0 {
			fmt.Printf("%-17sXDP Program loaded with ID: %v\n", msg.Attributes.Name, msg.Attributes.XDP.ProgID)
			// fmt.Printf(outTemplate, msg.Attributes.Name, 10, "xdp_dispatcher", "skb", 50, "d51e469e988d81da", "XDP_PASS")
			// fmt.Printf("XDP: %#v\n", *msg.Attributes.XDP)
		} else {
			fmt.Printf("%-17s<No XDP program loaded!>\n", msg.Attributes.Name)
		}
	}

	// TODO: Finish implementing
	log.Warn().Msg("Command not fully implemented")

	return nil
}
