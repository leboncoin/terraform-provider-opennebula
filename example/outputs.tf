output "IP Address" {
  value = "${opennebula_vm.vm-base.*.ip}"
}
