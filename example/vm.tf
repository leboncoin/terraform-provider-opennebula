data "opennebula_template_id" "base" {
  template_name = "Base Template v1"
}

resource "opennebula_vm" "vm-base" {
  count                 = "1"                                       # required | count of instances to create
  name                  = "MyVM"                                    # required | custom name of the vm to create
  template_id           = "${data.opennebula_template_id.base.id}"  # required | template_id to use. can be data source callback
  permissions           = "640"                                     # required | permissions for the vm (unix-style)
  cpu                   = "1"                                       # required | cpu count
  vcpu                  = "1"                                       # required | vcpu count
  memory                = "1024"                                    # required | memory count in mb

# https://docs.opennebula.org/5.6/operation/references/template.html#disks-section

  root_disk {
    dev_prefix          = "vd"                                      # optional | prefix device to map disk
    disk_type           = "FILE"                                    # optional | support media for the image
    image_id            = "1"                                       # optional | ID of the Image to use
    image               = "Debian 9.3"                              # optional | example image name stored in opennebula
    size                = 20480                                   # optional | image size in mb
    driver              = "qcow2"                                   # optional | specific image mapping driver
    image_uid           = 7                                        # optional | image of a given user by her ID
    image_uname         = "oneadmin"                                # required | owner of the image to use
  }

  disk {
    dev_prefix          = "vd"                                      # optional | prefix device to map disk
    disk_type           = "FILE"                                    # optional | support media for the image
    driver              = "qcow2"                                   # optional | specific image mapping driver
    format              = "qcow2"                                   # optional | format of the Image: raw or qcow2
    size                = 1024                                    # required | image size in mb
    target              = "vdb"                                     # required | device to map disk
    type                = "fs"                                      # optional | type of the disk: swap or fs
  }
}
