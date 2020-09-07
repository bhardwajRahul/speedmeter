package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/prometheus/procfs"
)

func getInterface() string {
	ipExec, err := exec.LookPath("ip")
	if err != nil {
		log.Fatalf("ip not found in $PATH: %s", err)
	}
	cmd := exec.Command(ipExec, "-o", "route", "show", "to", "default")
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	ipInterface := strings.Split(string(output), " ")[4]
	return ipInterface
}

func getBytes(p procfs.Proc, ipInterface string) (uint64, uint64) {
	stat, err := p.NetDev()
	if err != nil {
		log.Fatalf("could not fetch net stats: %s", err)
	}
	lo := stat[ipInterface]
	return lo.RxBytes, lo.TxBytes
}

func main() {
	p, err := procfs.Self()
	if err != nil {
		log.Fatalf("could not get process: %s", err)
	}
	ipInterface := getInterface()
	rInit, tInit := getBytes(p, ipInterface)
	for {
		r, t := getBytes(p, ipInterface)
		rDiff, tDiff := r-rInit, t-tInit
		fmt.Printf("Download speed: %d KBps, Upload speed: %d KBps\n", rDiff/1024, tDiff/1024)
		rInit, tInit = r, t
		time.Sleep(1 * time.Second)
	}
}
