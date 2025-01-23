package enum

import (
	"encoding/xml"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/net/ipv4"

	models "github.com/awareseven/mobilesniper/pkg/models/soap"
	"github.com/awareseven/mobilesniper/pkg/utils"
)

func DiscoverSDCDevices(interfaceName string, sdcChan chan<- models.ProbeMatch, wg *sync.WaitGroup, verbose bool) error {
	defer wg.Done()

	// Get source interface either by given name or
	// the first none-loopback device.

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("Failed to enumerate network interfaces: %v", err)
	}
	utils.LogVerbosef(verbose, "Network interfaces: %v", interfaces)

	var srcInterface *net.Interface
	var srcIP *net.IPNet
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Printf("Failed to enumerate network addresses: %v", err)
			continue
		}
		utils.LogVerbosef(verbose, "Interface %s addresses: %v", i.Name, addrs)

		// just use the first address of the network
		srcIP = addrs[0].(*net.IPNet)
		utils.LogVerbosef(verbose, "Got source IP %v on %s.", srcIP.IP, i.Name)

		if interfaces != nil {
			if i.Name == interfaceName {
				srcInterface = &i
				break
			}
		} else {
			if !srcIP.IP.IsLoopback() {
				srcInterface = &i
				break
			}
		}
	}

	if srcIP.IP == nil || srcInterface == nil {
		log.Fatalln("Failed to detect none-loopback interface to bind.")
	}
	utils.LogVerbosef(verbose, "Will bind to %s on interface %s", srcIP.IP, srcInterface.Name)

	// Create UDP addresses to send and recive discovery data
	srcAddr := &net.UDPAddr{IP: srcIP.IP, Port: 4747}
	discoveryAddr := &net.UDPAddr{IP: net.IPv4(239, 255, 255, 250), Port: 3702}

	multicastReceiver, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		log.Fatalf("Failed to create UDP connection: %v", err)
	}
	defer multicastReceiver.Close()

	// Join the multicast group
	udpConn := ipv4.NewPacketConn(multicastReceiver)
	if err := udpConn.JoinGroup(srcInterface, discoveryAddr); err != nil {
		log.Fatalln(err)
	}
	defer udpConn.LeaveGroup(srcInterface, discoveryAddr)

	// Create SOAP message - probe
	probe := models.NewSOAPMessage(
		"http://docs.oasis-open.org/ws-dd/ns/discovery/2009/01/Probe",
		"urn:docs-oasis-open-org:ws-dd:ns:discovery:2009:01",
		models.NewProbeSOAPBody(),
	)

	// Marshal the Probe message to XML
	xmlBytes, err := probe.XMLMarshal()
	if err != nil {
		log.Fatalln(err)
	}

	// Send the Probe message
	_, err = udpConn.WriteTo(xmlBytes, nil, discoveryAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Send Probe to %s", discoveryAddr)

	// Set deadline for connection at one second
	if err := udpConn.SetReadDeadline(time.Now().Add(time.Second * 5)); err != nil {
		log.Fatal(err)
	}

	// Read UDP multicast response
	buffer := make([]byte, 8192)
	n, _, _, err := udpConn.ReadFrom(buffer)
	if err != nil {
		log.Fatal(err)
	}
	utils.LogVerbosef(verbose, "Received %d bytes Anwser:\n%v", n, string(buffer))

	// Reading SOAP response from UDP response
	res := &models.ProbeMatch{}
	err = xml.Unmarshal(buffer, res)

	if err == nil {
		sdcChan <- *res
		log.Printf("Found SDC endpoint: %s", res.GetXAddrs())
	}

	return err
}
