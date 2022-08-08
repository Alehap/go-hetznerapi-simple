package main

import (
	"fmt"
	"os"
	"hdapi/hetznerapi"
	"time"
)

func main() {
	hetzner := hetznerapi.New(os.Getenv("HETZNER_API_KEY"))
	// fmt.Println(hetzner)
	s := hetzner.ListAllServers()
	for _, e := range s {
		fmt.Println(e.ID, e.Name, e.Status, e.PublicNet.IPv4.IP, e.PublicNet.IPv6.IP, e.IncludedTraffic, e.OutgoingTraffic, e.IngoingTraffic )
	}
	// create new server with ssd key
	sshid := hetzner.SSHKeyIdGetOrCreate(os.Getenv("SSH_PUB_KEY"))
	serverId := hetzner.CreateServer(sshid, "cx11", "centos-stream-8", "hel1")
	// fmt.Println(serverId)
	
	for {
		serverObj := hetzner.GetServerById(serverId)
		fmt.Println(serverObj.ID, serverObj.Name, serverObj.Status, serverObj.PublicNet.IPv4.IP, serverObj.PublicNet.IPv6.IP, serverObj.IncludedTraffic, serverObj.OutgoingTraffic, serverObj.IngoingTraffic )
		time.Sleep(5 * time.Second)
		if (serverObj.Status == "running") {break}
	} 
	// delete server by ID
	hetzner.DeleteServer(serverId)
	s2 := hetzner.ListAllServers()
	for _, e := range s2 {
		fmt.Println(e.ID, e.Name, e.Status, e.PublicNet.IPv4.IP, e.PublicNet.IPv6.IP, e.IncludedTraffic, e.OutgoingTraffic, e.IngoingTraffic )
	}
}