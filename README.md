# nsctl

![nsctl](./nsctl.png)

<p align="center">
<a href="https://github.com/Alonza0314/nsctl/releases"><img src="https://img.shields.io/github/v/release/Alonza0314/nsctl?color=orange" alt="Release"/></a>
<a href="https://github.com/Alonza0314/nsctl/blob/main/LICENSE"><img src="https://img.shields.io/github/license/Alonza0314/nsctl?color=blue" alt="License"/></a>
<a href="https://www.codefactor.io/repository/github/alonza0314/nsctl"><img src="https://www.codefactor.io/repository/github/alonza0314/nsctl/badge" alt="CodeFactor" /></a>
<a href="https://goreportcard.com/badge/github.com/Alonza0314/nsctl"><img src="https://goreportcard.com/badge/github.com/Alonza0314/nsctl" alt="goReport"/></a>
</p>

`nsctl` is a CLI tool for building your own network topology with pure Linux namespaces — no containers, no overhead, full control.

## Develop Environment

| Environment | Value                |
|-------------|----------------------|
| OS          | Ubuntu 25.04         |
| Go          | go1.25.5 linux/amd64 |
| Golint      | v2.7.2               |

## Usage

### Manual Build

```bash
git clone https://github.com/Alonza0314/nsctl.git
cd nsctl
make
```

After built, the binary file, `nsctl`, will be placed in the `build` directory.

## Command Description

- `ns` series

    | Command                  | Description            |
    | ------------------------ | ---------------------- |
    | `nsctl ns create <name>` | Create a new namespace |
    | `nsctl ns delete <name>` | Delete a namespace     |
    | `nsctl ns list`          | List all namespaces    |

- `net` series

    | Command                            | Description                    |
    | ---------------------------------- | ------------------------------ |
    | `nsctl net connect <ns1> <ns2>`    | Connect two namespaces         |
    | `nsctl net disconnect <ns1> <ns2>` | Disconnect two namespaces      |
    | `nsctl net list <ns>`              | List interfaces in a namespace |
    | `nsctl net set-ip <ns> <if> <ip>`  | Set IP address of an interface |
    | `nsctl net up <ns> <if>`           | Bring an interface up          |
    | `nsctl net down <ns> <if>`         | Bring an interface down        |

- `exec` series

    | Command                    | Description                      |
    | -------------------------- | -------------------------------- |
    | `nsctl exec <ns> -- <cmd>` | Execute a command in a namespace |

- `topo` series

    | Command                    | Description                   |
    | -------------------------- | ----------------------------- |
    | `nsctl topo apply <file>`  | Apply a topology from a file  |
    | `nsctl topo delete <file>` | Delete a topology from a file |

    The topology template file could be found at [YALM](./example/basicTopo/).
