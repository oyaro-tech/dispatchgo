# DispatchGo 🚀 - Fast & Simple SSH Task Runner

![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![GitHub stars](https://img.shields.io/github/stars/oyaro-tech/dispatchgo?style=social)

**DispatchGo** is a lightweight, concurrent command-line tool written in Go for *dispatching* tasks to multiple remote hosts via SSH. Define your infrastructure and tasks in a simple YAML playbook and let DispatchGo handle the execution, quickly and efficiently.

It's designed for developers, sysadmins, and DevOps engineers who need a straightforward way to automate tasks across multiple servers without the overhead of larger configuration management systems.

## ✨ Features

* **Simple YAML Playbooks:** Define hosts and tasks using an intuitive and easy-to-read YAML format.
* **Secure SSH Execution:** Runs commands on remote hosts securely over SSH.
* **Flexible Authentication:** Supports both **password** and **private key** (with optional passphrase) based SSH authentication.
* **Concurrent Execution:** Leverages Go's concurrency model (goroutines) to run tasks on multiple hosts simultaneously, making it fast.
* **Configurable Concurrency:** Control the number of concurrent routines using a command-line flag.
* **Debug Mode:** Get detailed output for troubleshooting with a simple debug flag.
* **Tilde Expansion:** Conveniently use `~` in your playbook for home directory paths (e.g., for private keys).
* **Known Hosts Handling:** Uses `~/.ssh/known_hosts` for host key verification.

## 🚀 Getting Started

### Prerequisites

* Go (Version 1.24 or higher) installed on your local machine.
* SSH access to your target hosts.

### Installation

1.  **Clone the repository:**

```bash
git clone https://github.com/oyaro-tech/dispatchgo.git
cd dispatchgo
```
2.  **Build the `dispatchgo` binary:**

```bash
go build ./cmd/dispatchgo/
```

This will create an executable file named `dispatchgo` in the current directory. You can move this to a directory in your `PATH` for system-wide access.

## 🛠️ Usage

Run `dispatchgo` from your terminal, pointing it to your playbook file.

```bash
./dispatchgo [flags]

Command-Line Flags
-playbook <path>: Specifies the path to your YAML playbook file. (Default: ./playbook.yaml)
-routines <number>: Sets the maximum number of concurrent goroutines (SSH connections/tasks) to run. (Default: 10)
-debug: Enables debug mode, which prints detailed output, including command outputs. (Default: false)
```

### Example

```bash
# Run tasks using the default playbook.yaml with 15 concurrent routines
./dispatchgo -routines 15

# Run tasks using a specific playbook and enable debug mode
./dispatchgo -playbook my_setup.yaml -debug
```

## 📖 Playbook Structure
The heart of DispatchGo is the playbook.yaml file. It defines the hosts you want to connect to and the jobs (collections of tasks) you want to run.

### Top-Level Structure

- `playbook_name`: (String) A descriptive name for your playbook.
- `hosts`: (List) A list defining your target hosts and their connection details.
- `jobs`: (Map) A map where each key is a job name, and the value contains the hosts to run on and the tasks to execute.

### `hosts` Section

The hosts section is a list where each item is a map. The key of the map is a unique name you give to the host (used to reference it in jobs), and the value contains its connection configuration.

- `host`: (String, Required) The IP address or hostname of the target machine.
- `port`: (Integer, Required) The SSH port number.
- `user`: (String, Required) The username for SSH login.
- `password`: (String, Optional) The password for SSH login.
- `private_key`: (String, Optional) The path to the SSH private key (supports ~ expansion).
- `passphrase`: (String, Optional) The passphrase for the private key, if it's encrypted.

### `jobs` Section

The jobs section defines what needs to be done. Each job has:

- `hosts`: (List, Required) A list of host names (as defined in the hosts section) where this job should run.
- `tasks`: (List, Required) A list of tasks to execute sequentially on each host.


### `tasks` Section

Each task within a job has:

- `name`: (String, Required) A descriptive name for the task (used for logging).
- `script`: (String, Required) The shell command or script to execute on the remote host. You can use YAML's multi-line string syntax (|) for longer scripts.

## 📝 Example Playbook (playbook.yaml.example)

Here's an example demonstrating how to set up fresh Debian-based VMs, including SSH configuration, updates, and Docker installation.

```yaml
playbook_name: "Setup fresh VM"

hosts:
  - host1:
      host: 192.168.10.1
      port: 22
      user: user
      password: "plain" # Consider using private keys for production!
      private_key: "~/.ssh/id_rsa"
  - host2:
      host: 192.168.10.2
      port: 22
      user: user
      private_key: "~/.ssh/id_rsa"
  - host3:
      host: 192.168.10.3
      port: 22
      user: user
      private_key: "~/.ssh/id_rsa"
  - host4:
      host: 192.168.10.4
      port: 22
      user: user
      private_key: "~/.ssh/id_rsa"
  - host5:
      host: 192.168.10.5
      port: 22
      user: user
      private_key: "~/.ssh/id_rsa"

jobs:
  setup:
    hosts:
      - host1
      - host2
      - host3
      - host4
      - host5
    tasks:
      - name: "Create a .ssh directory"
        script: "mkdir -p ~/.ssh"
      - name: "Add public key"
        script: "echo \"YOUR_PUBLIC_SSH_KEY_HERE\" >> ~/.ssh/authorized_keys" # Replace with your key
      - name: "Reconfigure sshd for public key auth"
        script: |
          sed -i 's/PasswordAuthentication yes/PasswordAuthentication no/g' /etc/ssh/sshd_config
          sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/g' /etc/ssh/sshd_config
          systemctl restart sshd
      - name: "Update apt package lists"
        script: "apt update -y"
      - name: "Upgrade packages"
        script: "apt upgrade -y"
      - name: "Set locale"
        script: |
          echo "LC_ALL=en_US.UTF-8" >> /etc/environment
          echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen
          echo "LANG=en_US.UTF-8" > /etc/locale.conf
          locale-gen en_US.UTF-8
      - name: "Install utils"
        script: "apt install ca-certificates curl -y"
      - name: "Add Docker's APT repository"
        script: |
          # Add Docker's official GPG key:
          install -m 0755 -d /etc/apt/keyrings
          curl -fsSL [https://download.docker.com/linux/debian/gpg](https://download.docker.com/linux/debian/gpg) -o /etc/apt/keyrings/docker.asc
          chmod a+r /etc/apt/keyrings/docker.asc

          # Add the repository to Apt sources:
          echo \
            "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] [https://download.docker.com/linux/debian](https://download.docker.com/linux/debian) \
            $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
            tee /etc/apt/sources.list.d/docker.list > /dev/null
      - name: "Update apt package lists (for Docker)"
        script: "apt update -y"
      - name: "Install the latest Docker"
        script: "apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y"
```

## 🤝 Contributing

Contributions are welcome! If you find a bug, have a feature request, or want to contribute code, please:

1.  Fork the repository (`https://github.com/oyaro-tech/dispatchgo.git`).
2.  Create a new branch (`git checkout -b feature/your-feature`).
3.  Make your changes.
4.  Commit your changes (`git commit -am 'Add some feature'`).
5.  Push to the branch (`git push origin feature/your-feature`).
6.  Create a new Pull Request.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file (you'll need to add one!) for details.

---

Made with ❤️ and Go.
