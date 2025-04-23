# TeSmart Control CLI

`tesmartctl` is a command-line utility for controlling TeSmart KVM switches over the network. It allows you to switch inputs, control the buzzer, manage LED timeouts, and more through a simple CLI interface.

## Installation

### Homebrew (macOS & Linux)

```bash
brew install mirceanton/taps/tesmartctl
```

### Manual Installation

Download the latest release binary for your platform from the [releases page](https://github.com/mirceanton/tesmartctl/releases).

### Building from Source

```bash
# Clone the repository
git clone https://github.com/mirceanton/tesmartctl.git
cd tesmartctl

# Build the binary
go build -o tesmartctl

# Move to a directory in your PATH (optional)
sudo mv tesmartctl /usr/local/bin/
```

## Configuration

On first run, `tesmartctl` will create a configuration file at `~/.config/tesmartctl.yaml` with default settings.

You can edit the configuration using the `config` command:

```bash
# View current configuration
tesmartctl config get

# Set KVM IP address
tesmartctl config set ip 192.168.1.10

# Set KVM port
tesmartctl config set port 5000
```

## Usage

### Basic Commands

```bash
# Test connectivity to the KVM
tesmartctl ping

# Get current active input
tesmartctl input get

# Switch to input 3 (PC3)
tesmartctl input set 3

# Mute the buzzer
tesmartctl buzzer mute
# or
tesmartctl buzzer 0

# Unmute the buzzer
tesmartctl buzzer unmute
# or
tesmartctl buzzer 1

# Set LED timeout to 10 seconds
tesmartctl timeout 10
# or
tesmartctl timeout short

# Set LED timeout to 30 seconds
tesmartctl timeout 30
# or
tesmartctl timeout long

# Disable LED timeout (LEDs always on)
tesmartctl timeout never
# or
tesmartctl timeout 0
# or
tesmartctl timeout off
```

### Advanced Usage

```bash
# Send raw hexadecimal commands
tesmartctl raw aabb031000ee  # Get current input

# Enable debug output
tesmartctl --debug input get

# Use a custom config file
tesmartctl --config /path/to/config.yaml input get
```

## Protocol Documentation

TeSmart KVM switches use a simple TCP-based protocol for control. The commands are sent as hexadecimal values.

### Common commands

- Switch to input: `0xAA 0xBB 0x03 0x11 0xXX 0xEE` (where XX is the port number in hex, starting from 00)
- Get current input: `0xAA 0xBB 0x03 0x10 0x00 0xEE`
- Mute buzzer: `0xAA 0xBB 0x03 0x02 0x00 0xEE`
- Unmute buzzer: `0xAA 0xBB 0x03 0x02 0x01 0xEE`
- Set LED timeout to 10 seconds: `0xAA 0xBB 0x03 0x03 0x0A 0xEE`
- Set LED timeout to 30 seconds: `0xAA 0xBB 0x03 0x03 0x1E 0xEE`
- Disable LED timeout: `0xAA 0xBB 0x03 0x03 0x00 0xEE`

For more details, see the [official TeSmart protocol documentation](https://support.tesmart.com/hc/en-us/article_attachments/27712362494361).

## Credits

This project was (heavily) inspired by and builds upon the work done by:

- [bbeaudoin's bash scripts for TeSmart KVM control](https://github.com/bbeaudoin/bash/tree/master/tesmart)
- [TeSmart's official protocol documentation](https://support.tesmart.com/hc/en-us/article_attachments/27712362494361)
