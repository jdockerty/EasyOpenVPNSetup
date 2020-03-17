# EasyOpenVPNSetup

**This is a proof of concept implementation, a secure version would require a HTTPS endpoint rather than HTTP. The setup is fully functioning, but this is something to keep in mind.**

## Overview
This project was inspired after considering how I could somewhat automate the [AWS OpenVPN repo](https://github.com/jdockerty/OpenVPNAWS), using Go and adding a REST API onto a server for creating clients was the idea that I came to. This means that you can send a new client name as JSON, such as `{ "Name" : "Jack" }` *(these names must be unique on the deployed server)* via `/api/addclient` and it adds a new client to the VPN server, responding with the relevant `.ovpn` configuration to paste into a file and load into the OpenVPN UI client. 

## Terraform
[Terraform](https://www.terraform.io/downloads.html) is used to provision the instance and an elastic IP which is associated to it, this means that the address of the instance remains static.

The command `terraform init` is required to enable the usage of the files once it is installed. 

Once Terraform is installed and configured, using `terraform apply` from within the folder will provision the appropriate resources and complete the initial OpenVPN configuration steps. 

*Note: The full configuration of the server can take a couple of minutes to install everything that is required.*

## REST API Response
Sending a POST request to the server containing the appropriate JSON format will cause it to respond with the relevant output required to paste into an `.ovpn` file and load into the OpenVPN UI client. This is demonstrated using Postman below.

![POST Request](https://github.com/jdockerty/EasyOpenVPNSetup/blob/master/Images/POST%20request.png)
