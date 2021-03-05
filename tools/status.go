package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/jsimonetti/rtnetlink"
	"github.com/rs/zerolog/log"
)

// Attached modes
const (
	XDP_MODE_UNSPEC = iota
	XDP_MODE_NATIVE
	XDP_MODE_SKB
	XDP_MODE_HW
)

var attachModeToName = map[uint8]string{
	XDP_MODE_UNSPEC: "unspecified",
	XDP_MODE_NATIVE: "native",
	XDP_MODE_SKB:    "skb",
	XDP_MODE_HW:     "hw",
}

type xdpProgram struct {
	iface       string
	fd          uint32
	bpfProgInfo bpfProgramInformation
	netlinkInfo rtnetlink.LinkXDP
}

func (p xdpProgram) Mode() string {
	mode, ok := attachModeToName[p.netlinkInfo.Attached]
	if !ok {
		mode = "UNKNOWN"
	}
	return mode
}

func (p xdpProgram) Exists() bool {
	return p.netlinkInfo.ProgID > 0
}

func (p xdpProgram) Name() string {
	return string(bytes.Trim(p.bpfProgInfo.name[:], "\x00"))
}

func status(ifname string) error {
	outTemplate := "%-16s %-5v %-17s %-8s %-4v %-17s %-13s\n"
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
		log.Fatal().Err(err).Msg("Failed to get interfaces")
	}

	xdpProgs := make([]xdpProgram, len(msgs))
	for i, msg := range msgs {
		progID := msg.Attributes.XDP.ProgID
		if progID <= 0 {
			xdpProgs[i] = xdpProgram{iface: msg.Attributes.Name}
			continue
		}
		fd, err := bpfGetFDByID(progID)
		if err != nil {
			log.Fatal().Err(err).Msg(
				"BPF_PROG_GET_FD_BY_ID syscall returned error")
		}
		bpfProgInfo, err := bpfGetInfoByFD(fd)
		if err != nil {
			log.Fatal().Err(err).Msg(
				"BPF_OBJ_GET_INFO_BY_FD syscall returned error")
		}
		xdpProgs[i] = xdpProgram{
			msg.Attributes.Name,
			fd,
			*bpfProgInfo,
			*msg.Attributes.XDP,
		}
	}
	fmt.Printf(outTemplate,
		"Interface", "Prio", "Program name", "Mode", "ID", "Tag",
		"Chain actions")
	fmt.Println(strings.Repeat("-", 86))
	for _, xdpProg := range xdpProgs {
		if xdpProg.Exists() {
			fmt.Printf(outTemplate,
				xdpProg.iface, 0, xdpProg.Name(),
				xdpProg.Mode(), xdpProg.netlinkInfo.ProgID,
				hex.EncodeToString(xdpProg.bpfProgInfo.tag[:]),
				"XDP_FOO")
			// fmt.Printf("%#v\n", xdpProg.bpfProgInfo)
		} else {
			fmt.Printf("%-17s<No XDP program loaded!>\n", xdpProg.iface)
		}
	}

	// TODO: Finish implementing
	log.Warn().Msg("Command not fully implemented")

	return nil
}
