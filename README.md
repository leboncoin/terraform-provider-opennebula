# terraform-provider-opennebula
[![Powered by Immowelt](https://img.shields.io/badge/powered%20by-immowelt-yellow.svg?colorB=ffb200)](https://stackshare.io/immowelt-group/)  

### This is the official Immowelt Group [OpenNebula](https://opennebula.org/) [Terraform](https://www.terraform.io/) Provider

Forked from [Runtastics Provider](https://github.com/runtastic/terraform-provider-opennebula.)
* Leverages [OpenNebula's XML/RPC API](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html) 
* Tested for versions 5.X  
  
    


## Table of Contents

* [terraform-provider-opennebula](#terraform-provider-opennebula)
  * [Table of Contents](#table-of-contents)
    * [Examples](#examples)
    * [Currently implemented](#currently-implemented)
    * [ToDo](#todo)
    * [Notes](#notes)
    * [Maintainer](#maintainer)
    * [Collaborators](#collaborators)
    * [Contributing](#contributing)
    * [License](#license)

## Examples

See [example](example/) folder for a quick overview 

## Currently implemented

The following list represent's all of OpenNebula's resources reachable through their API. The checked items are the ones that are fully functional and tested:

* [X] [onevm](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onevm)
* [X] [onetemplate](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onetemplate)
* [X] [onevnet](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onevnet)
* [X] [oneimage](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#oneimage)

## ToDo
* [ ]  Better examples of all modules
* [ ] [onemarket](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onemarket)
* [ ] [onemarketapp](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onemarketapp)
* [ ] [onevrouter](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onevrouter)
* [ ] [onezone](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onezone)
* [ ] [onesecgroup](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onesecgroup)
* [ ] [oneacl](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#oneacl)
* [ ] [oneacct](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#oneacct)
* [ ] [onehost](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onehost)
* [ ] [onecluster](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onecluster)
* [ ] [onegroup](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onegroup)
* [ ] [onevdc](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onevdc)
* [ ] [oneuser](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#oneuser)
* [ ] [onedatastore](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onedatastore)

## Notes

To update some vm resources the VM has to be in state poweroff.

Current flow:  
* resize disk: disk will be resized but vm won't be restarted
* resize cpu/vcpu/memory: requires new resource
* change ip address: requires new resource 


## Maintainer

- [Immowelt Group](https://github.com/immoweltgroup)
  - [Cemal Acar](https://github.com/cacar)
  - [Dennis Kribl](https://github.com/dkribl)
  - [Leroy Förster](https://github.com/gersilex)
  
## Collaborators

- [Lorenzo Arribas](https://github.com/larribas)
- [Jason Tevnan](https://github.com/tnosaj)


## Contributing

Bug reports and pull requests are welcome on GitHub at
https://github.com/immoweltgroup/terraform-provider-opennebula. This project is
intended to be a safe, welcoming space for collaboration, and contributors are
expected to adhere to the
[Contributor Covenant](http://contributor-covenant.org) code of conduct.

## License

The gem is available as open source under the terms of
the [MIT License](http://opensource.org/licenses/MIT).
