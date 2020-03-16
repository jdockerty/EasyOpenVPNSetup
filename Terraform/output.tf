output "VPNServerIP" {
  value = "${aws_eip.ServerIP.public_ip}"
}
