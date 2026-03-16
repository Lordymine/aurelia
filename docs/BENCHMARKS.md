# BENCHMARKS

## Objective

This document records measured baseline performance numbers for `Aurelia`.

The goal is to support README claims about lightweight runtime behavior with repeatable data instead of guesses.

## Benchmark Scope

Current baseline covers:

- release binary size
- process startup latency
- idle working set
- idle private memory
- idle CPU average

## Environment

Measurement date:

- `2026-03-16`

Environment:

- Windows
- Go `1.25`
- release build with `-trimpath -ldflags "-s -w"`
- local `.env` loaded privately for runtime-only dependencies
- benchmark runtime isolated in a temporary `AURELIA_HOME`
- benchmark DB and MCP config paths redirected to temporary paths

Notes:

- no secrets are recorded here
- numbers below are from local measurement on the current machine and should be treated as the current baseline, not a universal promise across all hardware

## Build Artifact

Command used:

```powershell
go build -trimpath -ldflags "-s -w" -o .\build\aurelia.exe .\cmd\aurelia
```

Result:

- binary size: `24,347,136` bytes
- binary size: `23.22 MB`

## Idle Runtime Baseline

Method:

1. start the release binary
2. wait `8s` for runtime stabilization
3. sample process memory
4. wait `5s`
5. sample CPU delta and memory again
6. stop the process

Three-run sample:

| Run | Startup (ms) | Working Set (MB) | Private Memory (MB) | Idle CPU Avg (%) |
| --- | ---: | ---: | ---: | ---: |
| 1 | 37.16 | 25.63 | 53.70 | 0.00 |
| 2 | 5.07 | 25.61 | 53.39 | 0.00 |
| 3 | 5.02 | 25.73 | 53.07 | 0.00 |

Averages:

- startup: `15.75 ms`
- idle working set: `25.66 MB`
- idle private memory: `53.39 MB`
- idle CPU average: `0.00%`

Interpretation:

- the first start is slower due to cold initialization
- after that, startup stabilizes near `5 ms` in this environment
- the runtime stays near `25-26 MB` working set while idle
- private memory stays near `53 MB`
- idle CPU usage is effectively zero in this sample window

## Reproduction Notes

For reproducibility:

- use the release build command above
- isolate `AURELIA_HOME` to a temporary directory
- redirect `DB_PATH` and `MCP_SERVERS_CONFIG_PATH` to temporary paths
- do not publish local credentials, tokens, or private config

## Pending Expansion

Still useful in future benchmark passes:

- cold versus warm startup split reported separately
- focused `go test -bench` numbers for selected internal packages
- runtime behavior under active tool execution
- memory impact of Agent Teams under concurrent tasks
