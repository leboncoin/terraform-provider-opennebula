data "opennebula_template_id" "base" {
  template_name = "Base Template v1"
}

resource "opennebula_vm" "vm-base" {
  count                 = "1"                                      # required | count of instances to create
  name                  = "MyVM"                                   # required | custom name of the vm to create
  template_id           = "${data.opennebula_template_id.base.id}" # required | template_id to use. can be data source callback
  permissions           = "640"                                    # required | permissions for the vm (unix-style)
  cpu                   = "1"                                      # required | cpu count 
  vcpu                  = "1"                                      # required | vcpu count
  memory                = "1024"                                   # required | memory count in mb
  image                 = "Debian 9.3"                             # required | example image name stored in opennebula
  size                  = "20480"                                  # required | image size in mb
  image_driver          = "qcow2"                                  # required | image driver of the image to use
  image_uname           = "oneadmin"                               # required | owner of the image to use
  network               = "MyCustomNetwork"                        # required | network name stored in opennebula
  network_search_domain = "mycorp.com"                             # required | search Domain 
  network_uname         = "oneadmin"                               # required | owner of the network to use
  security_group_id     = "0"                                      # required | security group id to use
  ip                    = "0.0.0.0"                                # optional | ip to use in network segment
}
