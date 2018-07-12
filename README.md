# terraform-provider-opennebula

[OpenNebula](https://opennebula.org/) provider for [Terraform](https://www.terraform.io/).
 
* Leverages [OpenNebula's XML/RPC API](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html) 
* Tested for versions 5.X


The provider tries to impose a lightweight level of abstraction on OpenNebula's resources. This means that only the most fundamental attributes are directly accessible (i.e. names, IDs, permissions and user/group identities). For maximum flexibility and portability, the remaining attributes can be specified using any of the formats natively accepted by OpenNebula (XML and String).

**UPDATE:**  The above-mentioned note applies now for every ressource except the *vm resource*:
The vm resource has been updated to work like other terraform providers/resources. You define VM attributes directly via terraform now instead of feeding a rendered template. This minimizes the created templates in opennebula and allows better multi-instance deployment of vms. Important note: this breaks the previous way vms are spawned via this provider...


## EXAMPLE

See [example](example/) folder for a quick overview  
TODO: Add examples for every resource and describe them better...

## NOTES ABOUT UPDATING RESOURCES

To update some vm resources the VM has to be in state poweroff.

Current flow:  
* resize disk: disk will be resized but vm won't be restarted
* resize cpu/vcpu/memory: requires new resource
* change ip address: requires new resource 

## ROADMAP

The following list represent's all of OpenNebula's resources reachable through their API. The checked items are the ones that are fully functional and tested:

* [X] [onevm](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onevm)
* [X] [onetemplate](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onetemplate)
* [ ] [onehost](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onehost)
* [ ] [onecluster](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onecluster)
* [ ] [onegroup](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onegroup)
* [ ] [onevdc](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onevdc)
* [X] [onevnet](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onevnet)
* [ ] [oneuser](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#oneuser)
* [ ] [onedatastore](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onedatastore)
* [X] [oneimage](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#oneimage)
* [ ] [onemarket](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onemarket)
* [ ] [onemarketapp](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onemarketapp)
* [ ] [onevrouter](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onevrouter)
* [ ] [onezone](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onezone)
* [ ] [onesecgroup](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onesecgroup)
* [ ] [oneacl](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#oneacl)
* [ ] [oneacct](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#oneacct)


## Collaborators

- [Lorenzo Arribas](https://github.com/larribas)
- [Jason Tevnan](https://github.com/tnosaj)
- [Immowelt Group](https://github.com/immoweltgroup)
  - [Cemal Acar](https://github.com/cacar)
  - [Dennis Kribl](https://github.com/dkribl)

## Contributing

Bug reports and pull requests are welcome on GitHub at
https://github.com/runtastic/terraform-provider-opennebula. This project is
intended to be a safe, welcoming space for collaboration, and contributors are
expected to adhere to the
[Contributor Covenant](http://contributor-covenant.org) code of conduct.

## License

The gem is available as open source under the terms of
the [MIT License](http://opensource.org/licenses/MIT).
