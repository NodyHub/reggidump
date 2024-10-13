# reggidump
Dump Docker registries.

## Install

### Build from source

```shell
make install
```

### Fetch latest from GitHub

```shell
go install github.com/NodyHub/reggidump@latest
```

## Usage

```shell
reggidump -h
Usage: reggidump <target> ... [flags]

Dump Docker images from a registry.

Arguments:
  <target> ...    Scan target. Can be a file, a registry address or - for stdin.

Flags:
  -h, --help                      Show context-sensitive help.
  -a, --auth=STRING               Registry authentication header.
  -d, --dump=STRING               Dump image layers to specified directory.
  -l, --list                      List available images+tags.
  -m, --manifest-only             Dump only image manifest.
  -o, --output="-"                Output file. Default is stdout.
  -p, --ping                      Check if target is a registry
  -r, --retry=5                   Number of retries a download.
  -t, --timeout=5                 Timeout in seconds for registry operations.
  -u, --user-agent="reggidump"    User agent string.
  -v, --verbose                   Enable verbose output.
      --version

cat input.lst
google.com
localhost:5000

reggidump input.lst -v -l -d tmp
time=2024-10-13T19:27:11.775+02:00 level=DEBUG msg="processing next" target=google.com
time=2024-10-13T19:27:11.776+02:00 level=DEBUG msg="list images and tags" server=google.com
time=2024-10-13T19:27:12.127+02:00 level=ERROR msg="failed to fetch images" err="cannot fetch images and tags: server is not a docker registry"
time=2024-10-13T19:27:12.127+02:00 level=DEBUG msg="processing next" target=localhost:5000
time=2024-10-13T19:27:12.127+02:00 level=DEBUG msg="list images and tags" server=localhost
time=2024-10-13T19:27:12.144+02:00 level=DEBUG msg="fetch tags" image=hashref-srv-hashref
time=2024-10-13T19:27:12.147+02:00 level=DEBUG msg="fetch tags" image=instrumentisto/nmap
time=2024-10-13T19:27:12.155+02:00 level=DEBUG msg="fetch tags" image=postgres
time=2024-10-13T19:27:12.158+02:00 level=DEBUG msg="fetch tags" image=ubuntu
localhost:5000/hashref-srv-hashref:latest
localhost:5000/instrumentisto/nmap:latest
localhost:5000/postgres:latest
localhost:5000/ubuntu:latest
time=2024-10-13T19:27:12.159+02:00 level=DEBUG msg="start dump" path=tmp server=localhost
time=2024-10-13T19:27:12.159+02:00 level=DEBUG msg="dump layers from server" proto=http address=localhost port=5000
time=2024-10-13T19:27:12.166+02:00 level=DEBUG msg="manifest stored" image=hashref-srv-hashref tag=latest path=tmp/localhost/hashref-srv-hashref/latest/manifest.json
time=2024-10-13T19:27:12.166+02:00 level=DEBUG msg="download layers" image=hashref-srv-hashref tag=latest layers=5
time=2024-10-13T19:27:12.166+02:00 level=DEBUG msg="downloading layer" address=localhost image=hashref-srv-hashref tag=latest digest=a3ed95caeb02
time=2024-10-13T19:27:12.178+02:00 level=DEBUG msg="downloading layer" address=localhost image=hashref-srv-hashref tag=latest digest=ec7b167827d3
time=2024-10-13T19:27:12.217+02:00 level=DEBUG msg="downloading layer" address=localhost image=hashref-srv-hashref tag=latest digest=a9eaa45ef418
time=2024-10-13T19:27:12.250+02:00 level=DEBUG msg="manifest stored" image=instrumentisto/nmap tag=latest path=tmp/localhost/instrumentisto/nmap/latest/manifest.json
time=2024-10-13T19:27:12.250+02:00 level=DEBUG msg="download layers" image=instrumentisto/nmap tag=latest layers=7
time=2024-10-13T19:27:12.250+02:00 level=DEBUG msg="downloading layer" address=localhost image=instrumentisto/nmap tag=latest digest=e5e7f9f80bcd
time=2024-10-13T19:27:12.305+02:00 level=DEBUG msg="downloading layer" address=localhost image=instrumentisto/nmap tag=latest digest=12f3d8bd77dd
time=2024-10-13T19:27:12.323+02:00 level=DEBUG msg="downloading layer" address=localhost image=instrumentisto/nmap tag=latest digest=c30352492317
time=2024-10-13T19:27:12.353+02:00 level=DEBUG msg="manifest stored" image=postgres tag=latest path=tmp/localhost/postgres/latest/manifest.json
time=2024-10-13T19:27:12.353+02:00 level=DEBUG msg="download layers" image=postgres tag=latest layers=26
time=2024-10-13T19:27:12.353+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=72bb9abb0bff
time=2024-10-13T19:27:12.358+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=0b3a5c24523f
time=2024-10-13T19:27:12.360+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=be63ee45aebd
time=2024-10-13T19:27:12.362+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=b4704d542ca4
time=2024-10-13T19:27:12.366+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=000a2849ebb4
time=2024-10-13T19:27:12.370+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=60da88fd1a76
time=2024-10-13T19:27:12.995+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=5d91cf4cf4f7
time=2024-10-13T19:27:13.000+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=b0e138aa7cc7
time=2024-10-13T19:27:13.004+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=a98eb4aa071f
time=2024-10-13T19:27:13.015+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=d23c85d6eacd
time=2024-10-13T19:27:13.065+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=ba82aab0ebf7
time=2024-10-13T19:27:13.084+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=315e4f3b46c2
time=2024-10-13T19:27:13.115+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=7495dd9a7534
time=2024-10-13T19:27:13.123+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=aa6fbc30c84e
time=2024-10-13T19:27:13.308+02:00 level=DEBUG msg="manifest stored" image=ubuntu tag=latest path=tmp/localhost/ubuntu/latest/manifest.json
time=2024-10-13T19:27:13.308+02:00 level=DEBUG msg="download layers" image=ubuntu tag=latest layers=6
time=2024-10-13T19:27:13.308+02:00 level=DEBUG msg="downloading layer" address=localhost image=ubuntu tag=latest digest=1567e7ea90b6
```

## Analysis

Use the scripts from [`scripts`](scripts/) directory. 