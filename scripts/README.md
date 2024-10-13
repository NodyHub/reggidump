# Script usage

## Dump registries and run trufflehog to extract secrets

```shell
$ scripts/dump_and_search.sh input.lst tmp
time=2024-10-13T19:38:05.661+02:00 level=DEBUG msg="processing next" target=google.com
time=2024-10-13T19:38:05.661+02:00 level=DEBUG msg="start dump" path=tmp server=google.com
time=2024-10-13T19:38:05.661+02:00 level=DEBUG msg="dump layers from server" proto="" address=google.com port=""
tee: tmp/google.com/truffle.out: No such file or directory
{"level":"info-0","ts":"2024-10-13T19:38:07+02:00","logger":"trufflehog","msg":"--directory flag is deprecated, please pass directories as arguments"}
{"level":"info-0","ts":"2024-10-13T19:38:07+02:00","logger":"trufflehog","msg":"running source","source_manager_worker_id":"6trsF","with_units":true}
{"level":"info-0","ts":"2024-10-13T19:38:07+02:00","logger":"trufflehog","msg":"finished scanning","chunks":0,"bytes":0,"verified_secrets":0,"unverified_secrets":0,"scan_duration":"455.416µs","trufflehog_version":"3.81.9"}
find: tmp/layer: No such file or directory
time=2024-10-13T19:38:07.066+02:00 level=DEBUG msg="processing next" target=localhost:5000
time=2024-10-13T19:38:07.066+02:00 level=DEBUG msg="start dump" path=tmp server=localhost
time=2024-10-13T19:38:07.066+02:00 level=DEBUG msg="dump layers from server" proto="" address=localhost port=5000
time=2024-10-13T19:38:07.073+02:00 level=DEBUG msg="fetch tags" image=hashref-srv-hashref
time=2024-10-13T19:38:07.075+02:00 level=DEBUG msg="fetch tags" image=instrumentisto/nmap
time=2024-10-13T19:38:07.078+02:00 level=DEBUG msg="fetch tags" image=postgres
time=2024-10-13T19:38:07.079+02:00 level=DEBUG msg="fetch tags" image=ubuntu
time=2024-10-13T19:38:07.085+02:00 level=DEBUG msg="manifest stored" image=hashref-srv-hashref tag=latest path=tmp/localhost/hashref-srv-hashref/latest/manifest.json
time=2024-10-13T19:38:07.085+02:00 level=DEBUG msg="download layers" image=hashref-srv-hashref tag=latest layers=5
time=2024-10-13T19:38:07.085+02:00 level=DEBUG msg="downloading layer" address=localhost image=hashref-srv-hashref tag=latest digest=a3ed95caeb02
time=2024-10-13T19:38:07.090+02:00 level=DEBUG msg="downloading layer" address=localhost image=hashref-srv-hashref tag=latest digest=ec7b167827d3
time=2024-10-13T19:38:07.113+02:00 level=DEBUG msg="downloading layer" address=localhost image=hashref-srv-hashref tag=latest digest=a9eaa45ef418
time=2024-10-13T19:38:07.137+02:00 level=DEBUG msg="manifest stored" image=instrumentisto/nmap tag=latest path=tmp/localhost/instrumentisto/nmap/latest/manifest.json
time=2024-10-13T19:38:07.137+02:00 level=DEBUG msg="download layers" image=instrumentisto/nmap tag=latest layers=7
time=2024-10-13T19:38:07.137+02:00 level=DEBUG msg="downloading layer" address=localhost image=instrumentisto/nmap tag=latest digest=e5e7f9f80bcd
time=2024-10-13T19:38:07.185+02:00 level=DEBUG msg="downloading layer" address=localhost image=instrumentisto/nmap tag=latest digest=12f3d8bd77dd
time=2024-10-13T19:38:07.195+02:00 level=DEBUG msg="downloading layer" address=localhost image=instrumentisto/nmap tag=latest digest=c30352492317
time=2024-10-13T19:38:07.218+02:00 level=DEBUG msg="manifest stored" image=postgres tag=latest path=tmp/localhost/postgres/latest/manifest.json
time=2024-10-13T19:38:07.218+02:00 level=DEBUG msg="download layers" image=postgres tag=latest layers=26
time=2024-10-13T19:38:07.218+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=72bb9abb0bff
time=2024-10-13T19:38:07.220+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=0b3a5c24523f
time=2024-10-13T19:38:07.222+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=be63ee45aebd
time=2024-10-13T19:38:07.224+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=b4704d542ca4
time=2024-10-13T19:38:07.225+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=000a2849ebb4
time=2024-10-13T19:38:07.226+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=60da88fd1a76
time=2024-10-13T19:38:07.815+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=5d91cf4cf4f7
time=2024-10-13T19:38:07.818+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=b0e138aa7cc7
time=2024-10-13T19:38:07.821+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=a98eb4aa071f
time=2024-10-13T19:38:07.831+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=d23c85d6eacd
time=2024-10-13T19:38:07.877+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=ba82aab0ebf7
time=2024-10-13T19:38:07.887+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=315e4f3b46c2
time=2024-10-13T19:38:07.917+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=7495dd9a7534
time=2024-10-13T19:38:07.919+02:00 level=DEBUG msg="downloading layer" address=localhost image=postgres tag=latest digest=aa6fbc30c84e
time=2024-10-13T19:38:08.096+02:00 level=DEBUG msg="manifest stored" image=ubuntu tag=latest path=tmp/localhost/ubuntu/latest/manifest.json
time=2024-10-13T19:38:08.096+02:00 level=DEBUG msg="download layers" image=ubuntu tag=latest layers=6
time=2024-10-13T19:38:08.096+02:00 level=DEBUG msg="downloading layer" address=localhost image=ubuntu tag=latest digest=1567e7ea90b6
{"level":"info-0","ts":"2024-10-13T19:38:09+02:00","logger":"trufflehog","msg":"--directory flag is deprecated, please pass directories as arguments"}
{"level":"info-0","ts":"2024-10-13T19:38:09+02:00","logger":"trufflehog","msg":"running source","source_manager_worker_id":"kNhc6","with_units":true}
{"level":"info-0","ts":"2024-10-13T19:38:12+02:00","logger":"trufflehog","msg":"finished scanning","chunks":25026,"bytes":149902842,"verified_secrets":0,"unverified_secrets":0,"scan_duration":"2.908269125s","trufflehog_version":"3.81.9"}
```

## Summarize results of the trufflerun

```shell
$ scripts/extract-uniq.sh tmp
Lines:        0
```