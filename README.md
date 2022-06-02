# PVE Ops - Slash Command for Proxmox PVE Control

**PVE Ops** aim is to offer the possibility to control your proxmox infrastructure through **ChatOps** (Mattermost, Slack, Google Chat, etc).  

We use the Golang module [go-proxmox](https://github.com/luthermonson/go-proxmox) by [luthermonson](https://github.com/luthermonson/go-proxmox) and [Gorilla/mux](https://github.com/gorilla/mux) for HTTP router.  

## Supported ChatOps platforms

* Mattermost _WIP_
* Slack _TODO_
* Google Chat _TODO_

## Installation

### Environment Variables

`PVE_USER` - Proxmox PVE User  
`PVE_PASS` - Proxmox PVE Pass  
`MM_TOKEN` - Mattermost Slash Command token  

### Manual Deployment
To use PVEOps only compile the code:  

`go build`

Then run the server:  

`./pveops`

This will start the server in the port `8000`.  

### Docker Compose Deployment

Attach you will find a docker-compose file, which you could use as template to deploy the service.
The docker compose file is configured to work out of the box with *Gitlab CI/CD* and *Traefik*.  

Gitlab CI/CD - Does the build of the service and its docker image and then deploy it to docker swarm cluster.  
Traefik - Is the load balancer/reverse proxy that server the application.  

### Mattermost - Slash Command Configuration 

In the mattermost you just need to create the slash command config in the integration section for your team. You can use this information as reference:  

* Title: `PVE Ops`
* Description: `Control Proxmox PVE with Mattermost`
* Command Trigger Word: `pve`
* Request URL:  `https://example.com/webhook/mattermost` - exchange example.com by your own domain name
* Request Method: `POST`
* Response Username: `PVEOps`
* Autocomplete: `checked`
* Autocomplete Hint: `[vmstatus VMID] [state VMID start|shutdown|reset]`

Please feel free to change any of this. Also you can add an icon if wanted.

After creating the integration/slash command, you should receive the Token, this token has to be set for the environment variable `MM_TOKEN`.  

## Using the slash command

* `/pve vmstatus VMID` - This will fetch virtual machine status. Replace VMID by the id of your vm
* `/pve state VMID shutdown` - Control virtual machine state, you can use shutdown, start or reset. Replace VMID by the id of your vm  
