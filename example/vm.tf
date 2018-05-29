resource "opennebula_vm" "vm-base" {
  # all properties listed here are mandatory for creating a vm with a basic template

  # basic informations:
  count       = "1"    # count of instances to create
  name        = "MyVM" # custom name of the vm to create
  template_id = "1"    # template_id to use. Can also be output id of resource_template 
  permissions = "640"

  # vm properties section:
  cpu    = "1"    # cpu count 
  vcpu   = "1"    # vcpu count
  memory = "1024" # memory count in mb

  # disk section: 
  image        = "Debian 9.3" # example image name stored in opennebula
  size         = "20480"      # image size in mb
  image_driver = "qcow2"      # image driver of the image to use
  image_uname  = "oneadmin"   # owner of the image to use

  #network section: 
  network               = "MyCustomNetwork" # network name stored in opennebula
  network_search_domain = "mycorp.com"      # search Domain 
  network_uname         = "oneadmin"        # owner of the network to use
  security_group_id     = "0"               # security group id to use
}
