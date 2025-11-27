# mbzr - MalwareBazaar CLI Client

A user-friendly command-line tool for interacting with the [MalwareBazaar API](https://bazaar.abuse.ch/api/).

[![Go Report Card](https://goreportcard.com/badge/github.com/andpalmier/mbzr)](https://goreportcard.com/report/github.com/andpalmier/mbzr)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

## Features

- ✅ Uses only Go standard libraries
- 📝 JSON logging and output for easy parsing
- ⚡️ Built-in rate limiting (10 req/s) to prevent API bans
- 💅🏻 Clean command-line interface

## Installation

### Using Homebrew (macOS/Linux)

```bash
brew tap andpalmier/tap
brew install mbzr
```

### Using Container Image

You can pull the pre-built image from GitHub Container Registry:

```bash
docker pull ghcr.io/andpalmier/mbzr:latest
```

### Using Go

```bash
go install github.com/andpalmier/mbzr@latest
```

### From Source

```bash
git clone https://github.com/andpalmier/mbzr.git
cd mbzr

# Build with Makefile (recommended - includes version info)
make build

# Or build manually
go build -o mbzr main.go
```

### Using Makefile

The project includes a Makefile:

```bash
make build          # Build with version information
make build-release  # Build optimized release binary
make install        # Install to $GOPATH/bin
make test           # Run tests
make test-coverage  # Run tests with coverage report
make lint           # Run linters
make clean          # Remove built binaries
make help           # Show all available commands
```

### Using Pre-built Binaries

Download the latest release from the [releases page](https://github.com/andpalmier/mbzr/releases).

### Using Containers

You can run `mbzr` using **Docker**, **Podman**, or Apple's **`container`** tool without installing Go. The commands are compatible across these runtimes.

1.  **Build the image**:

    ```bash
    # Docker
    docker build -t mbzr .

    # Podman
    podman build -t mbzr .

    # Apple container
    container build -t mbzr .
    ```

2.  **Run the container**:

    ```bash
    # Replace 'docker' with 'podman' or 'container' as needed
    docker run --rm mbzr --help
    ```

    To persist configuration or download files, mount a volume:

    ```bash
    docker run --rm -v $(pwd):/data mbzr download -sha256 <hash>
    ```

    Pass your API key as an environment variable:

    ```bash
    docker run --rm -e MALWAREBAZAAR_API_KEY="your_key" mbzr query -tag Emotet
    ```

    Or use an environment file (e.g., `.env`):

    ```bash
    # .env file content:
    # MALWAREBAZAAR_API_KEY=your_key

    docker run --rm --env-file .env mbzr query -tag Emotet
    ```

## Quick Start

1. **Set your API key** (get one from [MalwareBazaar](https://bazaar.abuse.ch/api/#account)):

```bash
export MALWAREBAZAAR_API_KEY="your_api_key_here"
```

2. **Query samples by tag**:

```bash
mbzr query -tag Emotet -limit 10
```

3. **Download a sample**:

```bash
mbzr download -sha256 ac25758feaf1ba3fe21e02e29681b2addc0246b507e4f6641a68d4baf73c9652
```

## Usage

### Global Flags

```
-v, --verbose      Enable verbose output with structured logging
-V, --version      Show version information
-h, --help         Show help message
```

### Commands

#### Query Samples

Search for malware samples using various criteria:

```bash
# Query by hash (SHA256, SHA1, or MD5)
mbzr query -hash ac25758feaf1ba3fe21e02e29681b2addc0246b507e4f6641a68d4baf73c9652

# Query by tag
mbzr query -tag Emotet -limit 50

# Query by signature
mbzr query -signature "Trojan.Generic"

# Query by file type
mbzr query -file_type exe

# Query by ClamAV signature
mbzr query -clamav "Win.Trojan.Agent"

# Query by YARA rule
mbzr query -yara rule_name

# Query by imphash
mbzr query -imphash 1234567890abcdef

# Query by telfhash
mbzr query -telfhash hash_value

# Query by dhash
mbzr query -dhash hash_value

# Query by gimphash
mbzr query -gimphash hash_value

# Query by TLSH
mbzr query -tlsh hash_value

# Query by certificate issuer
mbzr query -issuer_cn "Example CA"

# Query by certificate subject
mbzr query -subject_cn "Example Corp"

# Query by serial number
mbzr query -serial_number 1234567890
```

#### Download Samples

Download malware samples by SHA256 hash:

```bash
mbzr download -sha256 <sha256_hash>
```

**Security Note:** Downloaded files are saved as `<sha256>.zip` and are password-protected (password: `infected`).

#### Upload Samples

Upload malware samples to MalwareBazaar:

```bash
# Upload a single file
mbzr upload -file malware.exe -tags trojan,banker

# Upload all files in a directory
mbzr upload -dir /path/to/samples -tags malware

# Upload anonymously
mbzr upload -file sample.exe -anonymous
```

**Note:** Hidden files (starting with `.`) are automatically skipped during directory uploads.

#### Add Comments

Add comments to existing samples:

```bash
mbzr comment -sha256 <hash> -comment "This sample is related to campaign X"
```

#### Update Sample Metadata

Update tags or other metadata for a sample:

```bash
mbzr update -sha256 <hash> -key tags -value "newTag1,newTag2"
```

#### Get Recent Detections

Retrieve recently detected samples:

```bash
# Last 24 hours
mbzr recent_detections -hours 24

# Last week
mbzr recent_detections -hours 168
```

#### Get Latest Samples

Retrieve the latest malware samples added to MalwareBazaar:

```bash
# Last 60 minutes (default)
mbzr latest

# Last 100 samples
mbzr latest -selector 100
```

#### Query Code Signing Certificate Blocklist (CSCB)

```bash
mbzr cscb
```

## License

This project is licensed under the AGPLv3 License - see the [LICENSE](LICENSE) file for details.

## Thank you to

- [MalwareBazaar](https://bazaar.abuse.ch) for providing the API
- [abuse.ch](https://abuse.ch) for their work in fighting malware
- [@cocaman](https://github.com/cocaman) for his [MalwareBazaar scripts](https://github.com/cocaman/malware-bazaar)
