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

Dump all Docker images from a registry.

Arguments:
  <target> ...    Dump targets. Can be a file, a registry address or - for stdin.

Flags:
  -h, --help                      Show context-sensitive help.
  -a, --auth=STRING               Registry authentication header.
  -d, --dump=STRING               Dump image layers to specified directory.
  -f, --fail-count=5              Number of failed downloads before giving up a server.
  -l, --list                      List available images+tags.
  -m, --manifest-only             Dump only image manifest.
  -o, --output="-"                Output file. Default is stdout.
  -P, --parallel=5                Number of parallel downloads.
  -p, --ping                      Check if target is a registry
  -r, --retry=5                   Number of retries a download.
  -t, --timeout=5                 Timeout in seconds for registry operations.
  -u, --user-agent="reggidump"    User agent string.
  -v, --verbose                   Enable verbose output.
      --version
```

## Execution & Result

### Run

```shell
$ cat input.lst
google.com
localhost:5000

$ reggidump input.lst -v -l -d tmp
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

### Result

```shell
$ tree tmp/
tmp
├── layer
│   ├── sha256:000a2849ebb4264e0cc7a04ca19818db1e819b09706b23136e9936cd5348001b
│   ├── sha256:0b3a5c24523f60f013679fa5ba5644904e21fa7fa7d1c3ecc090a38b7ca6441f
│   ├── sha256:12f3d8bd77ddf56505e18e92bbe9c49def2c86e484e03b1c887ce1f8e26a8852
│   ├── sha256:1567e7ea90b67fc95ccdeeec39bdc3045098dee7e0c604975b957a9f8c0e9616
│   ├── sha256:315e4f3b46c2cf50f4850869012863150f85f947526b816238e127d72e7c2543
│   ├── sha256:5d91cf4cf4f79f8fabf0a7711f3d23c3df3ee58b387038819671894ceb68dc95
│   ├── sha256:60da88fd1a769a57b1a3de2bd799334bd14be05947e2edb3fb585a0e737bf742
│   ├── sha256:72bb9abb0bff4f43266a70ea77845d20da754a53e2d09553f1d761f79b10de60
│   ├── sha256:7495dd9a7534dc031e9ab69954fb433411b3a5e2d13b6b7aa3ef1d0ea900e3c6
│   ├── sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4
│   ├── sha256:a98eb4aa071fb95db04ba96eb0269f689ee05f25335128845cfcceaa3359494f
│   ├── sha256:a9eaa45ef418e883481a13c7d84fa9904f2ec56789c52a87ba5a9e6483f2b74f
│   ├── sha256:aa6fbc30c84e14e64571d3d7b547ea801dfca8a7bd74bd930b5ea5de3eb2f442
│   ├── sha256:b0e138aa7cc779cd901cba0297fb6a61a13cd65bfd4edf9905e82bf694fbcfbb
│   ├── sha256:b4704d542ca4977c5aa1ca95435db244239e35a8fd172821bd9b53534ec91d79
│   ├── sha256:ba82aab0ebf7e429fdfca16774a458f150afc39d97d8fdddbd5a2eef016d56fb
│   ├── sha256:be63ee45aebd7d4790879c716227b84373f6871eb17a5df86cf81d7b6e8e2063
│   ├── sha256:c303524923177661067f7eb378c3dd5277088c2676ebd1cd78e68397bb80fdbf
│   ├── sha256:d23c85d6eacd685eaa7458227091cfa23c0e25707befe09f2c6717f0dd862484
│   ├── sha256:e5e7f9f80bcdd4e3041465c79c109497403abde4cdbd39d165a6daf6620bfb29
│   └── sha256:ec7b167827d344df3be0dee00c3cea73569106d52d0fb8551adc810d3c8cc6e1
└── localhost
    ├── hashref-srv-hashref
    │   └── latest
    │       ├── manifest.json
    │       ├── sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4 -> ../../../layer/sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4
    │       ├── sha256:a9eaa45ef418e883481a13c7d84fa9904f2ec56789c52a87ba5a9e6483f2b74f -> ../../../layer/sha256:a9eaa45ef418e883481a13c7d84fa9904f2ec56789c52a87ba5a9e6483f2b74f
    │       └── sha256:ec7b167827d344df3be0dee00c3cea73569106d52d0fb8551adc810d3c8cc6e1 -> ../../../layer/sha256:ec7b167827d344df3be0dee00c3cea73569106d52d0fb8551adc810d3c8cc6e1
    ├── instrumentisto
    │   └── nmap
    │       └── latest
    │           ├── manifest.json
    │           ├── sha256:12f3d8bd77ddf56505e18e92bbe9c49def2c86e484e03b1c887ce1f8e26a8852 -> ../../../../layer/sha256:12f3d8bd77ddf56505e18e92bbe9c49def2c86e484e03b1c887ce1f8e26a8852
    │           ├── sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4 -> ../../../../layer/sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4
    │           ├── sha256:c303524923177661067f7eb378c3dd5277088c2676ebd1cd78e68397bb80fdbf -> ../../../../layer/sha256:c303524923177661067f7eb378c3dd5277088c2676ebd1cd78e68397bb80fdbf
    │           └── sha256:e5e7f9f80bcdd4e3041465c79c109497403abde4cdbd39d165a6daf6620bfb29 -> ../../../../layer/sha256:e5e7f9f80bcdd4e3041465c79c109497403abde4cdbd39d165a6daf6620bfb29
    ├── postgres
    │   └── latest
    │       ├── manifest.json
    │       ├── sha256:000a2849ebb4264e0cc7a04ca19818db1e819b09706b23136e9936cd5348001b -> ../../../layer/sha256:000a2849ebb4264e0cc7a04ca19818db1e819b09706b23136e9936cd5348001b
    │       ├── sha256:0b3a5c24523f60f013679fa5ba5644904e21fa7fa7d1c3ecc090a38b7ca6441f -> ../../../layer/sha256:0b3a5c24523f60f013679fa5ba5644904e21fa7fa7d1c3ecc090a38b7ca6441f
    │       ├── sha256:315e4f3b46c2cf50f4850869012863150f85f947526b816238e127d72e7c2543 -> ../../../layer/sha256:315e4f3b46c2cf50f4850869012863150f85f947526b816238e127d72e7c2543
    │       ├── sha256:5d91cf4cf4f79f8fabf0a7711f3d23c3df3ee58b387038819671894ceb68dc95 -> ../../../layer/sha256:5d91cf4cf4f79f8fabf0a7711f3d23c3df3ee58b387038819671894ceb68dc95
    │       ├── sha256:60da88fd1a769a57b1a3de2bd799334bd14be05947e2edb3fb585a0e737bf742 -> ../../../layer/sha256:60da88fd1a769a57b1a3de2bd799334bd14be05947e2edb3fb585a0e737bf742
    │       ├── sha256:72bb9abb0bff4f43266a70ea77845d20da754a53e2d09553f1d761f79b10de60 -> ../../../layer/sha256:72bb9abb0bff4f43266a70ea77845d20da754a53e2d09553f1d761f79b10de60
    │       ├── sha256:7495dd9a7534dc031e9ab69954fb433411b3a5e2d13b6b7aa3ef1d0ea900e3c6 -> ../../../layer/sha256:7495dd9a7534dc031e9ab69954fb433411b3a5e2d13b6b7aa3ef1d0ea900e3c6
    │       ├── sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4 -> ../../../layer/sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4
    │       ├── sha256:a98eb4aa071fb95db04ba96eb0269f689ee05f25335128845cfcceaa3359494f -> ../../../layer/sha256:a98eb4aa071fb95db04ba96eb0269f689ee05f25335128845cfcceaa3359494f
    │       ├── sha256:aa6fbc30c84e14e64571d3d7b547ea801dfca8a7bd74bd930b5ea5de3eb2f442 -> ../../../layer/sha256:aa6fbc30c84e14e64571d3d7b547ea801dfca8a7bd74bd930b5ea5de3eb2f442
    │       ├── sha256:b0e138aa7cc779cd901cba0297fb6a61a13cd65bfd4edf9905e82bf694fbcfbb -> ../../../layer/sha256:b0e138aa7cc779cd901cba0297fb6a61a13cd65bfd4edf9905e82bf694fbcfbb
    │       ├── sha256:b4704d542ca4977c5aa1ca95435db244239e35a8fd172821bd9b53534ec91d79 -> ../../../layer/sha256:b4704d542ca4977c5aa1ca95435db244239e35a8fd172821bd9b53534ec91d79
    │       ├── sha256:ba82aab0ebf7e429fdfca16774a458f150afc39d97d8fdddbd5a2eef016d56fb -> ../../../layer/sha256:ba82aab0ebf7e429fdfca16774a458f150afc39d97d8fdddbd5a2eef016d56fb
    │       ├── sha256:be63ee45aebd7d4790879c716227b84373f6871eb17a5df86cf81d7b6e8e2063 -> ../../../layer/sha256:be63ee45aebd7d4790879c716227b84373f6871eb17a5df86cf81d7b6e8e2063
    │       └── sha256:d23c85d6eacd685eaa7458227091cfa23c0e25707befe09f2c6717f0dd862484 -> ../../../layer/sha256:d23c85d6eacd685eaa7458227091cfa23c0e25707befe09f2c6717f0dd862484
    └── ubuntu
        └── latest
            ├── manifest.json
            ├── sha256:1567e7ea90b67fc95ccdeeec39bdc3045098dee7e0c604975b957a9f8c0e9616 -> ../../../layer/sha256:1567e7ea90b67fc95ccdeeec39bdc3045098dee7e0c604975b957a9f8c0e9616
            └── sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4 -> ../../../layer/sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4
```


## Analysis

Use the scripts from [`scripts`](scripts/) directory. 