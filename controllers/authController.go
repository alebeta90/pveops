package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"git.gonkar.com/gonkar/infra-cmd/models"
	u "git.gonkar.com/gonkar/infra-cmd/utils"
)

func Webhook(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		u.Respond(w, u.Message(true, "OPTIONS"))
		return
	}

	fmt.Println(r.Body)
	b, err02 := httputil.DumpRequest(r, true)
	if err02 != nil {
		log.Fatalln(err02)
	}

	fmt.Println(string(b))

	r.ParseForm()

	fmt.Printf("\nCommand Arguments: %s\n", r.FormValue("text"))
	commandRaw := r.FormValue("text")
	commandSplit := strings.Split(commandRaw, " ")

	if len(commandSplit) < 2 {
		u.Respond(w, map[string]interface{}{"status": true, "message": "test", "text": "command malformed"})
		return
	}

	fmt.Printf("\nAction: %s\n", commandSplit[0])
	fmt.Printf("\nArgument: %s\n", commandSplit[1])

	Action := commandSplit[0]
	VMID := commandSplit[1]

	pve := &models.PVE{}
	switch Action {
	case "vmstatus":
		//Do something
		vmid, err := strconv.Atoi(VMID)
		if err != nil {
			fmt.Printf("Error convertir argument to vmid, error: %v\n", err)
		}
		vmStatus := pve.GetVMStatus(vmid)
		u.Respond(w, map[string]interface{}{"status": true, "message": "test", "text": vmStatus})
		return
	case "state":
		vmid, err := strconv.Atoi(VMID)
		if err != nil {
			fmt.Printf("Error convertir argument to vmid, error: %v\n", err)
		}
		pve.PowerState(vmid, commandSplit[2])
		// Sleep few seconds till VM state changes
		u.Respond(w, map[string]interface{}{"status": true, "message": "test", "text": fmt.Sprintf("VM %v successfully %v", vmid, commandSplit[2])})
		return
	default:
		fmt.Println("No action detected")
		u.Respond(w, map[string]interface{}{"status": true, "message": "test", "text": "no action detected"})
		return
	}

}
