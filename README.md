# terraform-provider-opennebula

## Table of Contents
  * [About](#about)
  * [Instalation](#installation)
  * [Examples](#examples)
  * [Currently implemented](#currently-implemented)
  * [ToDo](#todo)
  * [Notes](#notes)
  * [Maintainer](#maintainer)
  * [Collaborators](#collaborators)
  * [Contributing](#contributing)
  * [License](#license)

## About
This is the official leboncoin [OpenNebula](https://opennebula.org/) [Terraform](https://www.terraform.io/) provider forked from:
  - [Runtastic](https://github.com/runtastic/terraform-provider-opennebula)
  - [ImmoweltGroup](https://github.com/ImmoweltGroup/terraform-provider-opennebula)

It leverages the [OpenNebula's XML/RPC API](https://docs.opennebula.org/5.6/integration/system_interfaces/api.html) and is tested for versions 5.x

## Installation

Download the latest release on [Github](https://github.com/leboncoin/terraform-provider-opennebula/releases)

Untar your release binary on `~/.terraform.d/plugins` and rename with pattern `terraform-provider-opennebula_vX.Y.Z`

For more information, follow the official procedure to install a [terraform plugin](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)

## Examples

See the [example](example/) folder for a quick overview

## Currently implemented

### Resources
* [X] [onevm](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onevm)
* [X] [onetemplate](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onetemplate)
* [X] [onevnet](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#onevnet)
* [X] [oneimage](https://docs.opennebula.org/5.2/integration/system_interfaces/api.html#oneimage)

### Data Sources
* [X] template_id - Get the first template id by a template name

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

* Regresion:
  - disabled disk resize
  - disabled resource NIC


## Maintainer

- [leboncoin](https://github.com/leboncoin)
  - [Andy Ladjadj](https://github.com/aladjadj)

## Collaborators

- [Runtastic](https://github.com/runtastic)
- [Lorenzo Arribas](https://github.com/larribas)
- [Jason Tevnan](https://github.com/tnosaj)


## Contributing

Bug reports and pull requests are welcome on GitHub at
https://github.com/leboncoin/terraform-provider-opennebula. This project is
intended to be a safe, welcoming space for collaboration, and contributors are
expected to adhere to the
[Contributor Covenant](http://contributor-covenant.org) code of conduct.

## License

The gem is available as open source under the terms of
the [MIT License](http://opensource.org/licenses/MIT).
