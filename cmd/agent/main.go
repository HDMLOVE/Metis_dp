package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	_ "net"
	"os"
	_ "strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// CapInfo 配置文件信息
type CapInfoSt struct {
	deviceName string
	filter     string
	ip         string
	port       string
	cpu        string
	mem        string
	dir        string
	size       string
}

var CapInfo CapInfoSt

// loadConfig 加载配置问价
func loadConfig(filename string) {
	cfg, err := ini.Load(filename)
	if err != nil {
		log.Fatal("fail to read the file: \n", err)
	}
	CapInfo.deviceName = cfg.Section("captrue").Key("interface").String()
	CapInfo.filter = cfg.Section("captrue").Key("filter").String()
	CapInfo.ip = cfg.Section("server").Key("ip").String()
	CapInfo.port = cfg.Section("port").Key("port").String()
	fmt.Println(CapInfo)
	os.Exit(0)
}

func main() {

	fmt.Println("packet start...")
	loadConfig("agent.ini")

	deviceName := "eth0"
	snapLen := int32(65535)
	port := uint16(22)
	filter := getFilter(port)
	fmt.Printf("device:%v, snapLen:%v, port:%v\n", deviceName, snapLen, port)
	fmt.Println("filter:", filter)

	//打开网络接口，抓取在线数据
	handle, err := pcap.OpenLive(deviceName, snapLen, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("pcap open live failed: %v", err)
		return
	}

	// 设置过滤器
	if err := handle.SetBPFFilter(filter); err != nil {
		fmt.Printf("set bpf filter failed: %v", err)
		return
	}
	defer handle.Close()

	// 抓包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.NoCopy = true
	for packet := range packetSource.Packets() {
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			fmt.Println("unexpected packet")
			continue
		}

		fmt.Printf("packet:%v\n", packet)

		// tcp 层
		tcp := packet.TransportLayer().(*layers.TCP)
		fmt.Printf("tcp:%v\n", tcp)
		// tcp payload，也即是tcp传输的数据
		fmt.Printf("tcp payload:%v\n", tcp.Payload)
	}
}

//定义过滤器
func getFilter(port uint16) string {
	filter := fmt.Sprintf("tcp and ((src port %v) or (dst port %v))", port, port)
	return filter
}
