provider "aws" {
    region = "eu-west-2"
    profile = "default"
}


resource "aws_instance" "VPNServer" {
    ami = "ami-0a0cb6c7bcb2e4c51" # RHEL instance used in AWS OpenVPN setup repo
    instance_type = "t2.micro" 

    key_name = "AWSOpenVPN" # Replace this with your AWS key pair

    security_groups = ["${aws_security_group.OpenVPNSG.name}"] 
    # Specifying .id causes an error due to no subnet being defined, we're using the default VPC.

    # Script to run on startup for installing relevant packages
    # Note: This takes a few minutes (approx 5mins) to install everything and have the Go code running.
    user_data = <<EOT
#!/bin/bash
cd /home/ec2-user/
sudo yum -y install git -y
sudo yum install iptables -y
sudo yum install wget -y
sudo yum install https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm -y
git clone https://github.com/Nyr/openvpn-install.git
cd openvpn-install/
chmod +x openvpn-install.sh
yes " " | sudo ./openvpn-install.sh
sudo yum install golang -y
cd /home/ec2-user/
git clone https://github.com/jdockerty/EasyOpenVPNSetup.git
cd EasyOpenVPNSetup/
go mod download
sudo go build main.go
sudo ./main &
EOT

    tags = {
        Name = "Terraform OpenVPN"
    }
}

# Provisioning an EIP to attach to the instance, these keeps the IP on the instance static.
resource "aws_eip" "ServerIP" {
    instance = "${aws_instance.VPNServer.id}"
}

# Associating the EIP to the instance.
resource "aws_eip_association" "EIPAssoc" {
    instance_id = "${aws_instance.VPNServer.id}"
    allocation_id = "${aws_eip.ServerIP.id}"
}

# Security group for our instance, virtual firewall.
resource "aws_security_group" "OpenVPNSG" {

    name = "OpenVPN-SG"
    ingress {
        # SSH access, this is not neccessary but provides a way to tweak/troubleshoot anything.
        from_port = 22
        to_port = 22
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    ingress { 
        # OpenVPN UDP port
        from_port = "1194"
        to_port = "1194"
        protocol = "udp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    ingress {
        # HTTP on 8080, REST API usage.
        from_port = "8080"
        to_port = "8080"
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    egress {
        from_port   = 0
        to_port     = 0
        protocol    = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }

   
}

