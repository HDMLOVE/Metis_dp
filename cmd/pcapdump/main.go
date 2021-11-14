package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket/dumpcommand"
	"github.com/google/gopacket/examples/util"
	"github.com/google/gopacket/pcap"
	"log"
	"os"
	"strings"
	"time"
)

var iface = flag.String("i", "eth0", "Interface to read packet from")
var fname = flag.String("r", "", "filename to read from, overrides -i")
var snaplen = flag.Int("s", 65536, "Snap length (number of bytes max to read per packet")
var tstype = flag.String("timestamp_type", "", "Type of timestamps to use")
var promisc = flag.Bool("promisc", true, "Set promiscuous mode")

func main() {
	defer util.Run()()
	var handle *pcap.Handle
	var err error
	fmt.Println("start...\n")
	if *fname != "" {
		if handle, err = pcap.OpenOffline(*fname); err != nil {
			log.Fatal("pcap openOffline error:", err)
		}
	} else {
		inactive, err := pcap.NewInactiveHandle(*iface)
		if err != nil {
			log.Fatalf("Could not create:%v", err)
		}
		defer inactive.CleanUp()
		if err = inactive.SetSnapLen(*snaplen); err != nil {
			log.Fatalf("could not set snap lenght:%v", err)
		} else if err = inactive.SetPromisc(*promisc); err != nil {
			log.Fatalf("could not set promisc mode: %v", err)
		} else if err = inactive.SetTimeout(time.Second); err != nil {
			log.Fatalf("could not set timeout %v", err)
		}
		if *tstype != "" {
			if t, err := pcap.TimestampSourceFromString(*tstype); err != nil {
				log.Fatalf("Supported timestamp types %v", inactive.SupportedTimestamps())
			} else if err = inactive.SetTimestampSource(t); err != nil {
				log.Fatalf("Support timestamp types %v")
			}
		}
		if handle, err = inactive.Activate(); err != nil {
			log.Fatalf("pcap active error：%v", err)
		}
		defer handle.Close()
	}
	if len(flag.Args()) > 0 {
		bpffilter := strings.Join(flag.Args(), " ")
		fmt.Fprintf(os.Stderr, "Using BPF filter %q\n", bpffilter)
		if err = handle.SetBPFFilter(bpffilter); err != nil {
			log.Fatalf("BPF filter error:%v", err)
		}
	}
	dumpcommand.Run(handle)
}
