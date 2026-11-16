# mbzr

CLI for the [MalwareBazaar API](https://bazaar.abuse.ch/api/). Query, download, upload, comment on, update, and export malware samples.

[![CI](https://github.com/andpalmier/mbzr/actions/workflows/ci.yml/badge.svg)](https://github.com/andpalmier/mbzr/actions/workflows/ci.yml)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

## Install

```bash
brew install --cask andpalmier/tap/mbzr
go install github.com/andpalmier/mbzr@latest
docker pull ghcr.io/andpalmier/mbzr:latest
```

Homebrew works on macOS only. On Linux, use Go or a release binary.

## Quick start

```bash
export ABUSECH_API_KEY="your_api_key_here"
mbzr query -tag RemcosRAT -limit 10
```

Every command needs `ABUSECH_API_KEY`. mbzr prints API errors as received.

## Commands

- `query`: look up samples by hash, tag, signature, file type, and more.
- `download`: fetch a sample by SHA256 hash.
- `upload`: submit a file or directory.
- `comment`: comment on a sample.
- `update`: edit metadata on an entry you submitted.
- `recent_detections`: list samples detected in the last n hours.
- `latest`: list recent additions.
- `cscb`: dump the Code Signing Certificate Blocklist.
- `export`: download a bulk export.
- `version`: print version information.

## Examples

```bash
mbzr query -hash ac25758feaf1ba3fe21e02e29681b2addc0246b507e4f6641a68d4baf73c9652
mbzr query -tag Emotet -limit 50
mbzr download -sha256 <hash> -out /tmp/sample.zip
mbzr upload -file malware.exe -tags trojan,banker
mbzr upload -dir /path/to/samples -tags malware
mbzr latest -selector 100
mbzr recent_detections -hours 24
mbzr -t 15m export -file full.csv.zip -out /tmp/bazaar.csv.zip
```

## Notes

- `-limit` defaults to 100 and is capped at 1000.
- Downloads are ZIP files with the password `infected`. AES archives need 7-Zip or `pyzipper`.
- Directory uploads skip hidden files and subdirectories. They continue after failures and exit non-zero if any file fails.
- Read the [submission policy](https://bazaar.abuse.ch/api/#policy) before uploading.
- Large exports need a longer `-t` timeout. Full exports rebuild once an hour.

## Container

```bash
docker run --rm -e ABUSECH_API_KEY="your_key" ghcr.io/andpalmier/mbzr query -tag RemcosRAT
docker run --rm -e ABUSECH_API_KEY="your_key" -v "$PWD":/data ghcr.io/andpalmier/mbzr download -sha256 <hash> -out /data/sample.zip
```

## License

AGPLv3. See [LICENSE](LICENSE).
