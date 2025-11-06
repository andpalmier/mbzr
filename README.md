# mbzr

[![License](https://img.shields.io/badge/license-AGPL--v3-blue)](https://www.gnu.org/licenses/agpl-3.0.en.html)
[![GoDoc Card](https://godoc.org/github.com/andpalmier/mbzr?status.svg)](https://godoc.org/github.com/andpalmier/mbzr)
[![Go Report Card](https://goreportcard.com/badge/github.com/andpalmier/mbzr)](https://goreportcard.com/report/github.com/andpalmier/mbzr)
[![follow on X](https://img.shields.io/twitter/follow/andpalmier?style=social&logo=x)](https://x.com/intent/follow?screen_name=andpalmier)

mbzr is a cli tool to interact with [MalwareBazaar API](https://bazaar.abuse.ch/api). [MalwareBazaar](https://bazaar.abuse.ch) is a project by [abuse.ch](https://abuse.ch) that collects and shares malware samples.
This tool allows users to query, download, upload and update malware samples and malware sample data.

This project was inspired by [@cocaman's malwarebazaar scripts](https://github.com/cocaman/malware-bazaar/) and the [official MalwareBazaar scripts](https://github.com/abusech/MalwareBazaar).

## Features

- **Query samples**: retrieve information about malware samples using various filters: tag, signature, filetype, [ClamAV](https://www.clamav.net) signature, imphash, [tlsh](https://github.com/trendmicro/tlsh), [gimphash](https://github.com/NextronSystems/gimphash), [telfhash](https://github.com/trendmicro/telfhash), icon dhash, YARA rule and Code signing certificate info (Issuer CN, Subject CN and Serial number).
- **Upload samples**: upload a malware sample or all malware samples in a directory.
- **Download samples**: download a malware sample by its SHA256 sha256.
- **Update sample metadata**: update the metadata of a malware sample, such as tags and comments.
- **Query recent samples**: retrieve the most recent malware samples.
- **Query Code Signing Certificate Blocklist (cscb)**: retrieve the content of the MalwareBazaar Code Signing Certificate Blocklist.

## Installation

You can install mbzr using `go`:

```bash
go install github.com/andpalmier/mbzr@latest
mbzr help
```

Or cloning the repo, and manually building it:

```bash
git clone https://github.com/andpalmier/mbzr.git
cd mbzr
go build -o mbzr main.go
```

In order to use mbzr, you need an API key from MalwareBazaar, you can get one for free [here](https://auth.abuse.ch/user/me).
Remember to set your MalwareBazaar API key as an environment variable:

```bash
export MBZR_API_KEY=<your_API_key>
```

## Usage

General syntax is: ```mbzr <command> [options]```

Available Commands:
- `query`: Query MalwareBazaar for information about a sample
- `download`: Download a sample by its sha256 sha256
- `upload`: Upload a file or directory to MalwareBazaar
- `comment`: Add a comment to a malware sample
- `recent_detections`: Get recent malware detections within a specified timeframe
- `cscb`: Query the Code Signing Certificate Blocklist (CSCB)
- `update`: Update metadata of a malware sample

Use `mbzr <command> -h` or `--help` for more information about a command.

### query

Query MalwareBazaar for information about malware samples. Usage:

```bash
mbzr query [flags]
```

Query options:

- `sha256`: Retrieve info about a malware sample by its sha256 (sha1, sha256 or md5).
- `tag`: Query malware sample associated with a tag.
- `signature`: Query malware samples associated with a signature.
- `filetype`: Query malware samples by filetype.
- `clamav`: Query malware samples by [ClamAV](https://www.clamav.net) signature.
- `imphash`: Query malware samples by import sha256 (only for PE files).
- `tlsh`: Query malware samples by [tlsh](https://github.com/trendmicro/tlsh) sha256.
- `telfhash`: Query malware samples by [telfhash](https://github.com/trendmicro/telfhash) sha256.
- `dhash`: Query malware samples by dhash Icon.
- `gimphash`: Query malware samples by [gimphash](https://github.com/NextronSystems/gimphash).
- `yara`: Query malware samples by YARA rule name.
- `cert_issuer`: Query code signing certificates by Issuer Common Name.
- `cert_subject`: Query code signing certificates by Subject Common Name.
- `cert_serial`: Query code signing certificates by Serial Number.

Optional flags:
- `limit`: Number of results to return (default: 100, max: 1000)

Examples:

```bash
mbzr query -sha256 <file_hash>
mbzr query -tag TrickBot -limit 10
mbzr query -signature 'Emotet' -limit 10
mbzr query -filetype 'exe' -limit 10
mbzr query -clamav 'Win.Trojan.Emotet-1234567' -limit 10
mbzr query -imphash <imphash_value> -limit 10
mbzr query -tlsh <tlsh_value> -limit 10
mbzr query -telfhash <telfhash_value> -limit 10
mbzr query -dhash <dhash_value> -limit 10
mbzr query -gimphash <gimphash_value> -limit 10
mbzr query -yara 'NETexecutableMicrosoft' -limit 10
mbzr query -cert_issuer 'Sectigo RSA Code Signing CA' -limit 10
mbzr query -cert_subject 'Microsoft Corporation' -limit 10
mbzr query -cert_serial <cert_serial> -limit 10
```

### upload

Uploads a file or all files in a directory to MalwareBazaar. Usage:

```bash
mbzr upload [flags]
```

Flags:

- `file <file_path>`: File to upload
- `dir <directory_path>`: Directory containing files to upload
- `tags <tag1,tag2,...>`: Comma separated list of tags associated with the files to upload
- `anonymous`: Upload files anonymously (no user association)

Examples:

```bash
mbzr upload -file sample.exe -tags trojan,banker
mbzr upload -dir /path/to/malware_samples -anonymous
```

### update

Updates metadata for a malware sample in MalwareBazaar. Usage:

```bash
mbzr update [flags]
```

Flags:

- `sha256 <sha256_hash>`: SHA256 sha256 of the file to update
- `key <key>`: Key to update (`add_tag`, `remove_tag`, `urlhaus`, `any_run`, `joe_sandbox`, `malpedia`, `twitter`, `links`, `dropped_by_md5`, `dropped_by_sha256`, `dropped_by_malware`, `dropping_md5`, `dropping_sha256`, `dropping_malware`, `comment`)
- `value <new_value>`: New value for the specified key

Example:

```bash
mbzr update -sha256 <sha256> -key add_tag -value ransomware
```

### comment 

Adds a comment to a malware sample in MalwareBazaar. Usage:

```bash
mbzr comment [flags]
```

Flags:

- `sha256 <sha256_hash>`: sha256 sha256 of the file to comment on
- `comment <comment>`: Comment to add to the malware sample

Examples:

```bash
mbzr comment -sha256 <sha256> -comment "This sample is part of a new campaign."
```

### download

Downloads a malware sample by its sha256 hash from MalwareBazaar. Usage:

```bash
mbzr download [flags]
```

Flag:

- `sha256 <file_hash>`: sha256 hash of the file to download

### cscb

Query the Code Signing Certificate Blocklist from MalwareBazaar. Usage:

```bash
bzr cscb [flags]
```

Example:
```
mbzr cscb
```

## Other examples

- Get a list of samples tagged with "Emotet", limited to 50 results:
```mbzr query -tag Emotet -limit 50```
- Download a sample by its SHA256 hash:
```mbzr download -sha256 ac25758feaf1ba3fe21e02e29681b2addc0246b507e4f6641a68d4baf73c9652```
- Upload a sample file:
```mbzr upload sample.exe```
- Add a comment to a sample:
```mbzr comment -sha256 ac25758feaf1ba3fe21e02e29681b2addc0246b507e4f6641a68d4baf73c9652 -comment 'This is a test comment'```
- Get recent detections from the last 12 hours:
```mbzr recent_detections -hours 12```
- Get the Code Signing Certificate Blocklist:
```mbzr cscb```
- Save the output in a JSON file by using the `tail` command to skip the first line:
```mbzr cscb | tail -n +2 > cscb.json```
- Get colored JSON output by piping the output to `jq`, for example:
```mbzr cscb | tail -n +2 | jq```

## Contributing

Contributions are welcome! If you have ideas for new features or improvements, feel free to open an issue or submit a pull request.

## License

This project is licensed under the [GNU AGPL V3 license](https://www.gnu.org/licenses/agpl-3.0.en.html).