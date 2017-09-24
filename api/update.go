package api

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func UpdateRecord(domain string, ipaddr string, addrType string) string {
	log.Println(fmt.Sprintf("%s record update request: %s -> %s", addrType, domain, ipaddr))

	f, err := ioutil.TempFile(os.TempDir(), "dyndns")
	if err != nil {
		return err.Error()
	}

	defer os.Remove(f.Name())
	w := bufio.NewWriter(f)

	w.WriteString(fmt.Sprintf("server %s\n", dnsConf.Server))
	w.WriteString(fmt.Sprintf("zone %s\n", dnsConf.Zone))
	w.WriteString(fmt.Sprintf("update delete %s.%s A\n", domain, dnsConf.Domain))
	w.WriteString(fmt.Sprintf("update delete %s.%s AAAA\n", domain, dnsConf.Domain))
	w.WriteString(fmt.Sprintf("update add %s.%s %v %s %s\n", domain, dnsConf.Domain, dnsConf.RecordTTL, addrType, ipaddr))
	w.WriteString("send\n")

	w.Flush()
	f.Close()

	cmd := exec.Command(dnsConf.NsupdateCmd, f.Name())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return err.Error() + ": " + stderr.String()
	}

	return out.String()
}
