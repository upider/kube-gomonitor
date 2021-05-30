package packet

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var (
	snapshot_len int32         = 65536
	promiscuous  bool          = false
	timeout      time.Duration = 30 * time.Second

	ethLayer layers.Ethernet
	ipLayer  layers.IPv4
	tcpLayer layers.TCP
	payload  gopacket.Payload

	LocalIP, RemoteIP     net.IP
	LocalPort, RemotePort layers.TCPPort
	Dir                   string
	pLen                  uint16

	AccIntv int64 = 3
	Start   bool
	PkgAcc  map[string]map[string]float64
	Ctx     context.Context
)

func init() {
	PkgAcc = make(map[string]map[string]float64)
}

func GetPcapHandle(ip string) (*pcap.Handle, error) {
	devs, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}
	var device string
	for _, dev := range devs {
		for _, v := range dev.Addresses {
			if v.IP.String() == ip {
				device = dev.Name
				break
			}
		}
	}
	if device == "" {
		return nil, errors.New("find device error")
	}

	h, err := pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		return nil, err
	}

	var filter string = "tcp and (not broadcast and not multicast)"
	err = h.SetBPFFilter(filter)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func NetSniff(ctx context.Context, ipAddr string) {
	pcapHandle, err := GetPcapHandle(ipAddr)
	if err != nil {
		log.Error(err)
		return
	}

	packetSource := gopacket.NewPacketSource(pcapHandle, pcapHandle.LinkType())

	go accumulator(ctx)

	for packet := range packetSource.Packets() {
		select {
		case <-ctx.Done():
			log.Warningln("Packet sniff Stop")
			pcapHandle.Close()

			return
		default:
			parser := gopacket.NewDecodingLayerParser(
				layers.LayerTypeEthernet,
				&ethLayer,
				&ipLayer,
				&tcpLayer,
				&payload,
			)
			foundLayerTypes := []gopacket.LayerType{}

			err := parser.DecodeLayers(packet.Data(), &foundLayerTypes)
			if err != nil {
				log.Error("Trouble decoding layers: ", err)
			}

			for _, layerType := range foundLayerTypes {
				if layerType == layers.LayerTypeIPv4 {

					if ipLayer.SrcIP.String() != ipAddr {
						LocalIP = ipLayer.DstIP
						RemoteIP = ipLayer.SrcIP
						LocalPort = tcpLayer.DstPort
						RemotePort = tcpLayer.SrcPort
						pLen = ipLayer.Length

						Dir = "in"

						itemId := fmt.Sprintf("%s:%d-%s:%d", LocalIP, LocalPort, RemoteIP, RemotePort)

						if _, ok := PkgAcc[itemId]; !ok {
							PkgAcc[itemId] = map[string]float64{
								"in":          0,
								"out":         0,
								"inRate":      0,
								"outRate":     0,
								"lastAccTime": 0,
								"lastAccIn":   0,
								"lastAccOut":  0,
							}
						}

						PkgAcc[itemId][Dir] = PkgAcc[itemId][Dir] + float64(pLen)

					} else {
						LocalIP = ipLayer.SrcIP
						RemoteIP = ipLayer.DstIP
						LocalPort = tcpLayer.SrcPort
						RemotePort = tcpLayer.DstPort
						pLen = ipLayer.Length

						Dir = "out"

						itemId := fmt.Sprintf("%s:%d-%s:%d", LocalIP, LocalPort, RemoteIP, RemotePort)

						if _, ok := PkgAcc[itemId]; !ok {
							PkgAcc[itemId] = map[string]float64{
								"in":          0,
								"out":         0,
								"inRate":      0,
								"outRate":     0,
								"lastAccTime": 0,
								"lastAccIn":   0,
								"lastAccOut":  0,
							}
						}
						PkgAcc[itemId][Dir] = PkgAcc[itemId][Dir] + float64(pLen)
					}
				}
			}
		}
	}

}

func accumulator(ctx context.Context) {
	log.Infoln("net accumulator thread is starting...")
	for {
		select {
		case <-ctx.Done():
			log.Warningln("net accumulator thread is stop.")
			Start = false
			return
		default:
			for _, pkgMap := range PkgAcc {
				start := float64(time.Now().Unix())
				in := pkgMap["in"]
				out := pkgMap["out"]

				if pkgMap["lastAccTime"] == 0 {
					pkgMap["lastAccTime"] = start - float64(AccIntv)
				}

				last := pkgMap["lastAccTime"]
				pkgMap["lastAccTime"] = start

				durSec := start - last

				if in == 0 {
					pkgMap["inRate"] = 0
				} else {
					pkgMap["inRate"] = (in - pkgMap["lastAccIn"]) / durSec
				}

				if out == 0 {
					pkgMap["outRate"] = 0
				} else {
					pkgMap["outRate"] = (out - pkgMap["lastAccOut"]) / durSec

				}

				pkgMap["lastAccIn"] = in
				pkgMap["lastAccOut"] = out
			}

			time.Sleep(time.Duration(AccIntv) * time.Second)

		}

	}
}
