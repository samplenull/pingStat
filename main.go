package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-ini/ini"
	"github.com/sparrc/go-ping"
)

func main() {
	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, "testping.conf")
	if err != nil {
		fmt.Printf("ERROR loading config file: %s\n", err.Error())
		return
	}
	hosts := cfg.Section("hosts").KeyStrings()
	for ihost := 0; ihost < len(hosts); ihost++ {
		fmt.Println("testing ping of " + hosts[ihost])
		pingHost(hosts[ihost], true)
	}

}

func getCfgPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return ex
}

func pingHost(host string, debug bool) {
	timeout := time.Duration(5 * time.Second)
	interval := time.Duration(500 * time.Millisecond)
	count := 5
	privileged := false

	pinger, err := ping.NewPinger(host)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	pinger.Count = count
	pinger.Interval = interval
	pinger.Timeout = timeout
	pinger.SetPrivileged(privileged)

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	pinger.Run()
}
