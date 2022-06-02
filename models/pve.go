package models

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/luthermonson/go-proxmox"
)

var (
	username, _ = os.LookupEnv("PVE_USER")
	password, _ = os.LookupEnv("PVE_PASS")
)

type PVE struct {
	Client   *proxmox.Client
	ID       int
	node     *proxmox.Node
	template *proxmox.VirtualMachine
	vm       *proxmox.VirtualMachine

	ApiUrl            string
	Username          string
	Password          string
	TwoFactorAuthCode string
	Insecure          bool
	Timeout           int
	TemplateId        int
	Node              string
	TokenID           string
	Secret            string
}

type ResultOSInfo struct {
	Result OSInfo `json:"result"`
}

type OSInfo struct {
	ID            string `json:"id"`
	KernelRelease string `json:"kernel-release"`
	KernelVersion string `json:"kernel-version"`
	Machine       string `json:"machine"`
	Name          string `json:"name"`
	PrettyName    string `json:"pretty-name"`
	Version       string `json:"version"`
	VersionID     string `json:"version-id"`
}

type CurrentState struct {
	Agent     int    `json:"agent"`
	Name      string `json:"name"`
	VMID      int    `json:"vmid"`
	Status    string `json:"status"`
	UpTime    int    `json:"uptime"`
	QMPStatus string `json:"qmpstatus"`
}

type VMState struct {
	Name      string `json:"name"`
	VMID      int    `json:"vmid"`
	Status    string `json:"status"`
	UpTime    int    `json:"uptime"`
	Agent     int    `json:"agent"`
	OSRelease string `json:"os_release"`
	OSKernel  string `json:"os_kernel"`
}

func (pve *PVE) Login() *proxmox.Client {
	insecureHTTPClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	client := proxmox.NewClient("https://node04.gonkar.internal:8006/api2/json",
		proxmox.WithClient(&insecureHTTPClient),
		proxmox.WithLogins(username, password),
	)

	return client
}

// DurationUpTime - Get Uptime as a string
func DurationUpTime(timestamp int) string {

	d := time.Duration(timestamp) * time.Second

	days := d.Truncate(24 * time.Hour) // divide by 24 for display
	d -= days

	hours := d.Truncate(time.Hour)
	d -= hours

	minutes := d.Truncate(time.Minute)
	d -= minutes

	seconds := d.Truncate(time.Second)
	d -= seconds

	return fmt.Sprintf("%d days %02d:%02d:%02d\n", int(days.Hours()/24), int(hours.Hours()), int(minutes.Minutes()), int(seconds.Seconds()))
}

// GetNodeByVMID - Fetch Cluster Resources to find VM's host node
func GetNodeByVMID(vmid int) string {
	pve := &PVE{}
	client := pve.Login()
	cluster, _ := client.Cluster()
	rs, _ := cluster.Resources()

	for k, v := range rs {
		if v.ID == fmt.Sprintf("qemu/%v", vmid) {
			fmt.Println(k)
			fmt.Println(v)
			return v.Node

		}

	}
	return ""
}

func (pve *PVE) GetVMStatus(vmId int) string {
	pve.Client = pve.Login()

	pve.ID = vmId

	var err2 error
	var err3 error

	pve.Node = GetNodeByVMID(vmId)

	pve.node, err2 = pve.Client.Node(pve.Node)
	if err2 != nil {
		fmt.Println(err2)
	}

	pve.vm, err3 = pve.node.VirtualMachine(pve.ID)
	if err3 != nil {
		fmt.Println(err3)
	}
	var vmOS *ResultOSInfo
	var vmCurrent *CurrentState

	// VM Current State - We get this first to check if QEMU Agent is active or now
	// If QEMU Agent is active then we can fetch OS info
	err := pve.Client.Get(fmt.Sprintf("/nodes/%s/qemu/%d/status/current", pve.Node, vmId), &vmCurrent)
	if err != nil {
		fmt.Printf("Current State Error: %v\n", err)
	}

	var msg string

	if vmCurrent.Status == "running" {
		err4 := pve.Client.Get(fmt.Sprintf("/nodes/%s/qemu/%d/agent/get-osinfo", pve.Node, vmId), &vmOS)
		if err4 != nil && err4.Error() != "500 QEMU guest agent is not running" {
			fmt.Printf("Get OS Info Error: %v\n", err4)
		}

		if err4 != nil {
			//Now we fill our response
			msg = fmt.Sprintf("| VM Status |  |\n| :------------ |:---------------:|\n| Name | %v |\n| Status | %v |\n| Uptime | %v |", vmCurrent.Name, vmCurrent.Status, DurationUpTime(vmCurrent.UpTime))

		} else {
			//Now we fill our response
			msg = fmt.Sprintf("| VM Status |  |\n| :------------ |:---------------:|\n| Name | %v |\n| Status | %v |\n| Uptime | %v |\n| OS | %v |\n|Kernel | %v |", vmCurrent.Name, vmCurrent.Status, DurationUpTime(vmCurrent.UpTime), vmOS.Result.PrettyName, vmOS.Result.KernelRelease)

		}
	} else {
		fmt.Println("Agent not active in VM")
		//Now we fill our response
		msg = fmt.Sprintf("| VM Status |  |\n| :------------ |:---------------:|\n| Name | %v |\n| Status | %v |", vmCurrent.Name, vmCurrent.Status)
	}

	return msg
}

func (pve *PVE) PowerState(vmId int, state string) error {
	pve.Client = pve.Login()

	pve.ID = vmId

	var err2 error
	var err3 error

	pve.Node = GetNodeByVMID(vmId)

	fmt.Printf("PVE NODE: %v\n", pve.Node)

	pve.node, err2 = pve.Client.Node(pve.Node)
	if err2 != nil {
		fmt.Printf("Error on setting pve.Client.Node: %v\n", err2)
	}

	pve.vm, err3 = pve.node.VirtualMachine(pve.ID)
	if err3 != nil {
		fmt.Println(err3)
	}

	switch state {
	case "start":
		t, err := pve.vm.Start()
		if err != nil {
			return err
		}
		return t.WaitFor(15)

	case "shutdown":
		t, err := pve.vm.Stop()
		if err != nil {
			return err
		}
		return t.WaitFor(15)
	case "reset":
		t, err := pve.vm.Reset()
		if err != nil {
			return err
		}

		return t.WaitFor(15)
	default:
		return nil

	}
}
