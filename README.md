# EasyOpenVPNSetup

**This is a proof of concept implementation, a secure version would require a HTTPS endpoint rather than HTTP. The setup is fully functioning, but this is something to keep in mind.**

## Overview
This project was inspired after considering how I could somewhat automate the [AWS OpenVPN repo](https://github.com/jdockerty/OpenVPNAWS), using Go and adding a REST API onto a server for creating clients was the idea that I came to. This means that you can send a new client name as JSON, such as `{ "Name" : "Jack" }` *(these names must be unique on the deployed server)* via `/api/addclient` and it adds a new client to the VPN server, responding with the relevant `.ovpn` configuration to paste into a file and load into the OpenVPN UI client. 

## Terraform
[Terraform](https://www.terraform.io/downloads.html) is used to provision the instance and an elastic IP which is associated to it, this means that the address of the instance remains static.

The command `terraform init` is required to enable the usage of the files once it is installed. 

## REST API Response
Sending a POST request to the server containing the appropriate format will mean that it responds with the relevant output required to paste into an `.ovpn` file and load into the OpenVPN UI client. 
