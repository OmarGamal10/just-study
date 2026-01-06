# Just study

- Operates at the **Kernel Network Layer**, making it impossible to bypass without a reboot or explicit command.

## Requirements

- **OS:** Linux
- **Dependencies:** `go`, `iproute2` (provides the `ss` command), `systemd-resolved`.

## Installation

**Clone & install:**

```bash
sudo ./install.sh
```

## Why Sudo?

- Modifying `/etc/hosts` requires root privileges.
- Killing TCP connections at the kernel level (`ss -K dst [ip]`) also requires elevated permissions.
- Installing as a system command requires root access. (moving to `/usr/local/bin`)

## Usage

```bash
# Block all distractions & kill active sessions
  sudo focus on

# Restore access
  sudo focus off

# Check current status
  sudo focus status
```
