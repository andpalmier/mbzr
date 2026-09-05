# mbzr - MalwareBazaar CLI client

A command-line client for the [MalwareBazaar API](https://bazaar.abuse.ch/api/). It queries samples, downloads them, submits new ones, and edits entries you own.

> Part of the abuse.ch CLI toolkit, a set of clients for [abuse.ch](https://abuse.ch) services:
> - [urlhs](https://github.com/andpalmier/urlhs) for URLhaus, the malware URL database
> - [tfox](https://github.com/andpalmier/tfox) for ThreatFox, the IOC database
> - [yrfy](https://github.com/andpalmier/yrfy) for YARAify, YARA scanning
> - [mbzr](https://github.com/andpalmier/mbzr) for MalwareBazaar, malware samples

[![CI](https://github.com/andpalmier/mbzr/actions/workflows/ci.yml/badge.svg)](https://github.com/andpalmier/mbzr/actions/workflows/ci.yml)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

## Features

- Built on the Go standard library, with no third party dependencies
- Prints JSON, so you can pipe it into jq or anything else
- Rate limits itself to 10 requests per second
- Runs under Docker, Podman, and Apple container

## Installation

### Homebrew

```bash
brew install --cask andpalmier/tap/mbzr
```

Homebrew casks are macOS only. On Linux, use `go install` or a pre-built binary.

### Go

```bash
go install github.com/andpalmier/mbzr@latest
```

### Container

```bash
# Pull the pre-built image
docker pull ghcr.io/andpalmier/mbzr:latest

# Or build it yourself
docker build -t mbzr .
```

### From source

```bash
git clone https://github.com/andpalmier/mbzr.git
cd mbzr
make build
```

## Quick start

Get an API key from the [abuse.ch Authentication Portal](https://auth.abuse.ch/), export it, then query something:

```bash
export ABUSECH_API_KEY="your_api_key_here"
mbzr query -tag Emotet -limit 10
```

Every command reads the key from `ABUSECH_API_KEY`. When the API refuses a request, mbzr prints the reason it gave, so "that hash is unknown to MalwareBazaar" instead of a bare status code.

## Usage

### Global flags

These go before the command name.

| Flag | Description |
|------|-------------|
| `-v`, `--verbose` | Print what the client is doing |
| `-t`, `--timeout` | Timeout per request, as a duration such as `45s` or `2m` (default `30s`) |
| `-V`, `--version` | Print version information |
| `-h`, `--help` | Print help |

Raise the timeout for large result sets. A query with `-limit 1000` takes close to 30 seconds on its own, so it tends to need `-t 90s`.

### Commands

| Command | Description |
|---------|-------------|
| `query` | Look up samples by hash, tag, signature, file type, and more |
| `download` | Fetch a sample by SHA256 hash |
| `upload` | Submit a file or a directory of files |
| `comment` | Comment on a sample |
| `update` | Edit metadata on an entry you submitted |
| `recent_detections` | List samples detected in the last n hours |
| `latest` | List recent additions |
| `cscb` | Dump the Code Signing Certificate Blocklist |
| `version` | Print version information |

### Querying samples

```bash
# By hash. SHA256, SHA1 and MD5 all work
mbzr query -hash ac25758feaf1ba3fe21e02e29681b2addc0246b507e4f6641a68d4baf73c9652
mbzr query -hash bab94357d255c22ec55e60dc55745d58b4d7ef12

# By tag or signature
mbzr query -tag Emotet -limit 50
mbzr query -signature "Trojan.Generic"

# By file type
mbzr query -file_type elf

# By detection or rule name
mbzr query -clamav "Win.Trojan.Agent"
mbzr query -yara win_remcos_g0

# By fuzzy or structural hash
mbzr query -imphash 1234567890abcdef1234567890abcdef
mbzr query -tlsh T1A5B...
mbzr query -telfhash 1E634BC4B643D9F2ED0602B52477EF338E76F5B...
mbzr query -gimphash 3870859e16c5541b4a6d2b3ce6e8b...
mbzr query -dhash f8dcbeffbffecee8

# By code signing certificate
mbzr query -issuer_cn "Sectigo RSA Code Signing CA"
mbzr query -subject_cn "Elite Web Development Ltd."
mbzr query -serial_number 05a6cf9108f6941e492c3d2a1dc4dc9631b2
```

`-limit` defaults to 100. The API caps it at 1000, and mbzr clamps anything higher with a warning rather than sending a request the server will trim anyway.

### Downloading samples

```bash
mbzr download -sha256 ac25758feaf1ba3fe21e02e29681b2addc0246b507e4f6641a68d4baf73c9652

# Choose where it lands
mbzr download -sha256 <hash> -out /tmp/sample.zip
```

Samples arrive as a ZIP protected with the password `infected`. Without `-out` the file is written to `<sha256>.zip` in the current directory. mbzr will not overwrite a file that already exists.

Some archives use AES encryption. If your unzip tool reports `compression type 99`, it lacks AES support, so use something that has it, such as 7-Zip or Python's pyzipper.

### Uploading samples

Read the [submission policy](https://bazaar.abuse.ch/api/#policy) first. MalwareBazaar wants confirmed malware, no older than about ten days, and no adware or file infectors.

```bash
# A single file, with tags
mbzr upload -file malware.exe -tags trojan,banker

# Every file in a directory
mbzr upload -dir /path/to/samples -tags malware

# Without attributing it to your account
mbzr upload -file sample.exe -anonymous
```

You can record how the sample spread and where it came from. `-reference` and `-context` take `key=value` and can be repeated:

```bash
mbzr upload -file loader.exe \
  -delivery_method email_attachment \
  -reference any_run=https://app.any.run/tasks/1 \
  -reference twitter=https://twitter.com/abuse_ch/status/1224269018506330112 \
  -context dropped_by_malware=Gozi \
  -context comment="dropped during a Gozi campaign"
```

`-delivery_method` accepts `email_attachment`, `email_link`, `web_download`, `web_drive-by`, `multiple`, or `other`. Reference keys are `urlhaus`, `any_run`, `joe_sandbox`, `malpedia`, `twitter`, and `links`. Context keys are `dropped_by_md5`, `dropped_by_sha256`, `dropped_by_malware`, `dropping_md5`, `dropping_sha256`, `dropping_malware`, and `comment`.

A directory upload keeps going when one file fails. It reports each failure and exits non-zero if any occurred.

### Commenting and updating

```bash
mbzr comment -sha256 <hash> -comment "Seen in a phishing campaign targeting logistics firms."

mbzr update -sha256 <hash> -key add_tag -value ransomware
mbzr update -sha256 <hash> -key any_run -value https://app.any.run/tasks/1
```

You can only update entries your own account submitted. Valid keys are `add_tag`, `remove_tag`, `urlhaus`, `any_run`, `joe_sandbox`, `malpedia`, `twitter`, `links`, `dropped_by_md5`, `dropped_by_sha256`, `dropped_by_malware`, `dropping_md5`, `dropping_sha256`, `dropping_malware`, and `comment`. If the value is already set, mbzr says so and exits cleanly instead of treating it as a failure.

### Recent samples

```bash
# Added in the last 60 minutes
mbzr latest

# The 100 most recent additions
mbzr latest -selector 100

# Samples given a family label in the last n hours, up to 168
mbzr recent_detections -hours 24
```

`recent_detections` defaults to 48 hours, matching the API default.

### Code signing certificate blocklist

```bash
mbzr cscb
```

Each entry carries the serial number, thumbprint, subject and issuer, validity dates, and the reason it was blocklisted.

### Running in a container

```bash
docker run --rm -e ABUSECH_API_KEY="your_key" ghcr.io/andpalmier/mbzr query -tag Emotet

podman run --rm -e ABUSECH_API_KEY="your_key" ghcr.io/andpalmier/mbzr query -tag Emotet

container run --rm -e ABUSECH_API_KEY="your_key" ghcr.io/andpalmier/mbzr query -tag Emotet
```

Downloads need a mounted volume, otherwise the sample disappears with the container:

```bash
docker run --rm -e ABUSECH_API_KEY="your_key" -v $(pwd):/data ghcr.io/andpalmier/mbzr download -sha256 <hash> -out /data/sample.zip
```

## Environment variables

| Variable | Description |
|----------|-------------|
| `ABUSECH_API_KEY` | Your abuse.ch API key. Required. |

## License

Licensed under the AGPLv3. See [LICENSE](LICENSE) for the full text.

## Acknowledgments

- [MalwareBazaar](https://bazaar.abuse.ch) by abuse.ch
- [abuse.ch](https://abuse.ch) for their work against malware
