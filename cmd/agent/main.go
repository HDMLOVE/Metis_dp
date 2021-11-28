package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	_ "net"
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

var capInfo CapInfoSt

// loadConfig 加载配置文件
func loadConfig(filename string) {
	cfg, err := ini.Load(filename)
	if err != nil {
		log.Fatal("fail to read the file: \n", err)
	}
	capInfo.deviceName = cfg.Section("captrue").Key("deviceName").String()
	capInfo.filter = cfg.Section("captrue").Key("filter").String()
	capInfo.ip = cfg.Section("server").Key("ip").String()
	capInfo.port = cfg.Section("port").Key("port").String()
	fmt.Println(capInfo)
}

func packetHandle(c <-chan layers.TCP) {
	select {
	case <-c:
		fmt.Println("test")
	default:
		fmt.Println("aaaa")
	}
}

// 主函数入口
func main() {

	fmt.Println("packet start...")
	// 加载配置文件
	loadConfig("agent.ini")

	deviceName := capInfo.deviceName
	snapLen := int32(65535)
	filter := capInfo.filter
	fmt.Printf("device:%v, snapLen:%v, filter:%v\n", deviceName, snapLen, filter)

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

	//c := make(chan layers.TCP, 100)
	//go packetHandle(c)

	// 抓包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.NoCopy = true
	for packet := range packetSource.Packets() {
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			fmt.Println("unexpected packet")
			continue
		}

		//fmt.Printf("packet:%v\n", packet.Data())

		// tcp 层
		tcp := packet.TransportLayer().(*layers.TCP)
		fmt.Printf("\ntcp:%t\n", tcp)
		// tcp payload，也即是tcp传输的数据
		//fmt.Printf("tcp payload:%v\n", tcp.Payload)
	}
}
