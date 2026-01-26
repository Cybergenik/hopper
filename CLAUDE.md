# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Hopper is a distributed fuzzer inspired by AFL++. It coordinates fuzzing across multiple nodes using a master-node architecture with code coverage guidance and mutation techniques.

## Build & Test Commands

```bash
# Build binaries
go build ./cmd/hopper-master
go build ./cmd/hopper-node

# Run all tests
go test -v ./...

# Docker build
docker build -t hopper-node .
```

## Running the Fuzzer

```bash
# Start master (port 6969 default)
./hopper-master --host <hostname> --port <port> --seed <initial_seed>

# Start node(s)
./hopper-node --target <path_to_instrumented_binary> --host <master_host> --port <master_port>
```

**Target compilation** (requires clang with ASAN):
```bash
clang -fsanitize=address -fsanitize-coverage=edge,trace-pc-guard <source.c> -o <binary>
```

## Architecture

### Master-Node Model
- **Master** (`master/`): RPC server orchestrating fuzzing. Manages seed queue, mutation engine, coverage deduplication via bloom filters, and energy-based priority scheduling.
- **Node** (`node/`): Worker executing target binary. Extracts coverage via sancov, parses ASAN reports, reports results to master.

### Key Components
- `cmd/hopper-master/main.go` - Master entry point with TUI
- `cmd/hopper-node/main.go` - Node entry point
- `master/master.go` - Core master logic: RPC handlers (GetFTask, UpdateFTask), task queue management
- `master/mut.go` - Mutation engine with 7 mutation types (MUT, DEL, ADD, SWP, FLP, REV, ROT)
- `master/pq.go` - Energy-based priority queue (max-heap, size 5000)
- `master/bloomfilter.go` - Dual bloom filters for seed and coverage deduplication
- `node/node.go` - Fuzz loop, coverage extraction, target execution
- `node/utils.go` - ASAN parsing, crash persistence
- `common/common.go` - RPC message types (FTask, TaskUpdate)
- `common/murmur.go` - Zero-heap-alloc MurmurHash128 implementation
- `tui/tui.go` - BubbleTea terminal UI

### RPC Protocol
Master exposes HTTP/RPC on configurable port:
- `Master.GetFTask` - Node requests work (returns seed + mutation params)
- `Master.UpdateFTask` - Node reports results (coverage hash, crash status)

### Environment Variables
- `HOPPER_OUT` - Output directory for crashes/reports
- `HOPPER_LOG` - Enable periodic reporting
- `HOPPER_LOG_INTERVAL` - Report interval in minutes (default: 30)
- `SANCOV_BIN` - Path to sancov binary

### Energy-Based Scheduling
Seeds are prioritized by: `min(1, covEdges/maxEdges) * (2 if crash else 1)`

Havoc level parameter controls mutation intensity per seed.

## Examples

Working examples in `examples/`:
- `parse/` - Simple parser demo
- `binutils/` - GNU readelf fuzzing
- `jq/` - JSON processor fuzzing
- `unsafe/` - Unsafe code demo

Each includes Docker scripts (`master_docker.sh`, `node_docker.sh`) for containerized runs.
