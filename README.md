# Multi-Cluster tool
<!-- markdown-link-check-disable-next-line -->
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub Super-Linter](https://github.com/electrocucaracha/multicluster/workflows/Lint%20Code%20Base/badge.svg)](https://github.com/marketplace/actions/super-linter)
[![Ruby Style Guide](https://img.shields.io/badge/code_style-rubocop-brightgreen.svg)](https://github.com/rubocop/rubocop)
[![Go Report Card](https://goreportcard.com/badge/github.com/electrocucaracha/multicluster)](https://goreportcard.com/report/github.com/electrocucaracha/multicluster)
[![GoDoc](https://godoc.org/github.com/electrocucaracha/multicluster?status.svg)](https://godoc.org/github.com/electrocucaracha/multicluster)
![visitors](https://visitor-badge.glitch.me/badge?page_id=electrocucaracha.multicluster)

This tool provisions and interconnects Kubernetes clusters through the usage of
the [KinD API][1]. It uses a [configuration yaml file](scripts/config.yml) to
define the topology.

```bash
go install github.com/electrocucaracha/multicluster/...
multicluster create --config=scripts/config.yml
```

## Provisioning process

The [Vagrant tool][2] can be used for provisioning an Ubuntu Focal
Virtual Machine. It's highly recommended to use the  *setup.sh* script
of the [bootstrap-vagrant project][3] for installing Vagrant
dependencies and plugins required for this project. That script
supports two Virtualization providers (Libvirt and VirtualBox) which
are determine by the **PROVIDER** environment variable.

```bash
curl -fsSL http://bit.ly/initVagrant | PROVIDER=libvirt bash
```

Once Vagrant is installed, it's possible to provision a Virtual
Machine using the following instructions:

```bash
vagrant up
```

The provisioning process will take some time to install all
dependencies required by this project and perform a Kubernetes
deployment on it.

The following diagram shows the result after its execution.

```text
+---------------------------------+     +---------------------------------+     +---------------------------------+     +---------------------------------+
| central (k8s)                   |     | regional (k8s)                  |     | edge-1 (k8s)                    |     | edge-2 (k8s)                    |
| +-----------------------------+ |     | +-----------------------------+ |     | +-----------------------------+ |     | +-----------------------------+ |
| | central-control-plane       | |     | | regional-control-plane      | |     | | edge-1-control-plane        | |     | | edge-2-control-plane        | |
| | podSubnet: 10.196.0.0/16    | |     | | podSubnet: 10.197.0.0/16    | |     | | podSubnet: 10.198.0.0/16    | |     | | podSubnet: 10.199.0.0/16    | |
| | serviceSubnet: 10.96.0.0/16 | |     | | serviceSubnet: 10.97.0.0/16 | |     | | serviceSubnet: 10.98.0.0/16 | |     | | serviceSubnet: 10.99.0.0/16 | |
| +-----------------------------+ |     | +-----------------------------+ |     | +-----------------------------+ |     | +-----------------------------+ |
| | eth0(172.88.0.2/16)         | |     | | eth0(172.89.0.2/16)         | |     | | eth0(172.90.0.2/16)         | |     | | eth0(172.91.0.2/16)         | |
| +------------+----------------+ |     | +------------+----------------+ |     | +------------+----------------+ |     | +------------+----------------+ |
|              |                  |     |              |                  |     |              |                  |     |              |                  |
+--------------+------------------+     +--------------+------------------+     +--------------+------------------+     +--------------+------------------+
               |                                       |                                       |                                       |
     +=========+============+                +=========+============+                +=========+===========+                 +=========+===========+
     |  net-central(bridge) |                | net-regional(bridge) |                |  net-edge-1(bridge) |                 |  net-edge-2(bridge) |
     |    172.88.0.0/16     |                |    172.89.0.0/16     |                |    172.90.0.0/16    |                 |    172.91.0.0/16    |
     +=========+============+                +=========+============+                +=========+===========+                 +=========+===========+
               |                                       |                                       |                                       |
+--------------+---------------------------------------+---------------------------------------+---------------------------------------+-----------+
| wan-test (emulator)                                                                                                                              |
+--------------------------------------------------------------------------------------------------------------------------------------------------+
| eth0(172.80.0.2/24)                                                                                                                              |
| eth1(172.90.0.254/16)                                                                                                                            |
| eth2(172.91.0.254/16)                                                                                                                            |
| eth3(172.89.0.254/16)                                                                                                                            |
| eth4(172.88.0.254/16)                                                                                                                            |
+--------------------------------------------------------------------------------------------------------------------------------------------------+
```

[1]: https://kind.sigs.k8s.io/
[2]: https://www.vagrantup.com/
[3]: https://github.com/electrocucaracha/bootstrap-vagrant
