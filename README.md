# EasyOpenVPNSetup

**This is a proof of concept implementation, a secure version would require a HTTPS endpoint rather than HTTP. The setup is fully functioning, but this is something to keep in mind.**

## Overview
This project was inspired after considering how I could somewhat automate the [AWS OpenVPN repo](https://github.com/jdockerty/OpenVPNAWS), using Go and adding a REST API onto a server for creating clients was the idea that I came to. This means that you can send a new client name as JSON, such as `{ "Name" : "Jack" }` *(these names must be unique on the deployed server)* via `/api/AddClient` and it adds a new client to the VPN server, responding with the relevant `.ovpn` configuration to paste into a file and load into the OpenVPN UI client. 

