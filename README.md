# mbzr - MalwareBazaar CLI Client

A command-line tool for interacting with the [MalwareBazaar API](https://bazaar.abuse.ch/api/).

> **Part of the abuse.ch CLI toolkit** - This project is part of a collection of CLI tools for interacting with [abuse.ch](https://abuse.ch) services:
> - [urlhs](https://github.com/andpalmier/urlhs) - URLhaus (malware URL database)
> - [tfox](https://github.com/andpalmier/tfox) - ThreatFox (IOC database)
> - [yrfy](https://github.com/andpalmier/yrfy) - YARAify (YARA scanning)
> - [mbzr](https://github.com/andpalmier/mbzr) - MalwareBazaar (malware samples)

[![Go Report Card](https://goreportcard.com/badge/github.com/andpalmier/mbzr)](https://goreportcard.com/report/github.com/andpalmier/mbzr)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

## Features

- ✅ Uses only Go standard libraries
- 📝 JSON output for easy parsing
- ⚡️ Built-in rate limiting (10 req/s)
- 🐳 Docker, Podman, and Apple container support

## Installation

### Using Go

```bash
go install github.com/andpalmier/mbzr@latest
```

### Using Container (Docker/Podman)

```bash
# Pull pre-built image
docker pull ghcr.io/andpalmier/mbzr:latest

# Or build locally
docker build -t mbzr .
```

### From Source

```bash
git clone https://github.com/andpalmier/mbzr.git
cd mbzr
make build
```

## Quick Start

1. **Get your API key** from [abuse.ch Authentication Portal](https://auth.abuse.ch/)

2. **Set your API key**:

```bash
export ABUSECH_API_KEY="your_api_key_here"
```

3. **Query samples by tag**:

```bash
mbzr query -tag Emotet -limit 10
```

## Usage

### Commands

| Command | Description |
|---------|-------------|
| `query` | Query samples by hash, tag, signature, file type, etc. |
| `download` | Download a malware sample by SHA256 hash |
| `upload` | Upload a file or directory to MalwareBazaar |
| `comment` | Add a comment to a malware sample |
| `latest` | Get latest malware samples |
| `cscb` | Query the Code Signing Certificate Blocklist |
| `version` | Show version information |

### Query Samples

```bash
# By hash (SHA256, SHA1, or MD5)
mbzr query -hash ac25758feaf1ba3fe21e02e29681b2addc0246b507e4f6641a68d4baf73c9652

# By tag
mbzr query -tag Emotet -limit 50

# By signature
mbzr query -signature "Trojan.Generic"

# By file type
mbzr query -file_type exe

# By ClamAV signature
mbzr query -clamav "Win.Trojan.Agent"

# By YARA rule
mbzr query -yara rule_name

# By imphash
mbzr query -imphash 1234567890abcdef1234567890abcdef

# By TLSH
mbzr query -tlsh T1A5B...
```

### Download Samples

```bash
mbzr download -sha256 ac25758feaf1ba3fe21e02e29681b2addc0246b507e4f6641a68d4baf73c9652
```

> **Note**: Downloaded files are saved as `<sha256>.zip` (password: `infected`)

### Upload Samples

```bash
# Single file
mbzr upload -file malware.exe -tags trojan,banker

# Directory
mbzr upload -dir /path/to/samples -tags malware

# Anonymous upload
mbzr upload -file sample.exe -anonymous
```

### Get Latest Samples

```bash
# Last 60 minutes
mbzr latest

# Last 100 samples
mbzr latest -selector 100
```

### Container Usage

```bash
# Run with Docker
docker run --rm -e ABUSECH_API_KEY="your_key" ghcr.io/andpalmier/mbzr query -tag Emotet

# Run with Podman
podman run --rm -e ABUSECH_API_KEY="your_key" ghcr.io/andpalmier/mbzr query -tag Emotet

# Run with Apple container
container run --rm -e ABUSECH_API_KEY="your_key" ghcr.io/andpalmier/mbzr query -tag Emotet

# Mount volume for downloads
docker run --rm -e ABUSECH_API_KEY="your_key" -v $(pwd):/data ghcr.io/andpalmier/mbzr download -sha256 <hash>
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ABUSECH_API_KEY` | Your abuse.ch API key (required) |

## License

This project is licensed under the AGPLv3 License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [MalwareBazaar](https://bazaar.abuse.ch) by abuse.ch
- [abuse.ch](https://abuse.ch) for their work in fighting malware
