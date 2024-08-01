# ebpf-pinnedmap-metrics

## Development Status

This project is currently under development.

**Note:** The tool only exposes metrics for pinned maps.


So far it supports metrics for:

- Hash eBPF Map
- Array eBPF Map
- Hash LRU eBPF Map

## How to Run

To run the program, follow these steps:

```
go generate
go build
sudo ./ebpf-pinned-map-metrics
```

On each host you can trigger actions on eBPF map using:

```
sudo bpftool map
sudo bpftool map update id <MAP-ID> key 0 0 0 0 value 1 0 0 0
sudo bpftool map delete id <MAP-ID> key 0 0 0 0
sudo bpftool map lookup id <MAP-ID> key 0 0 0 0
```
