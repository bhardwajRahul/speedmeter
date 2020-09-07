package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/getlantern/systray"
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

func getBytes() (uint64, uint64) {
	p, err := procfs.Self()
	if err != nil {
		log.Fatalf("could not get process: %s", err)
	}
	ipInterface := getInterface()
	stat, err := p.NetDev()
	if err != nil {
		log.Fatalf("could not fetch net stats: %s", err)
	}
	lo := stat[ipInterface]
	return lo.RxBytes, lo.TxBytes
}

func onReady() {
	go func() {
		rInit, tInit := getBytes()
		for {
			r, t := getBytes()
			rDiff, tDiff := r-rInit, t-tInit
			rInit, tInit = r, t
			systray.SetTitle(fmt.Sprintf("Download: %d KBps, Upload: %d KBps\n", rDiff/1024, tDiff/1024))
			time.Sleep(1 * time.Second)
		}
	}()
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quits this app")
	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func main() {
	systray.Run(onReady, nil)
}
