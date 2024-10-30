Blockless B7S Node Configuration

# Blockless B7S Node Configuration

This page lists all of the configuration options supported by the b7s daemon. It showcases the configuration structure, as accepted in a YAML config file, environment variables that can be used to set those options and, where applicable, the CLI flags and their default values.

## role [🔗](\#role)

Type: string

Path: role

Environment variable: B7S\_Role

CLI flag:

- Flag: `--role`
- Shorthand: `-r`
- Default: worker
- Description: role this node will have in the Blockless protocol (head or worker)

## concurrency [🔗](\#concurrency)

Type: uint

Path: concurrency

Environment variable: B7S\_Concurrency

CLI flag:

- Flag: `--concurrency`
- Shorthand: `-c`
- Default: 10
- Description: maximum number of requests node will process in parallel

## boot-nodes [🔗](\#boot-nodes)

Type: list (string)

Path: boot-nodes

Environment variable: B7S\_BootNodes

CLI flag:

- Flag: `--boot-nodes`
- Default: \[\]
- Description: list of addresses that this node will connect to on startup, in multiaddr format

## workspace [🔗](\#workspace)

Type: string

Path: workspace

Environment variable: B7S\_Workspace

CLI flag:

- Flag: `--workspace`
- Default: N/A
- Description: directory that the node can use for file storage

## load-attributes [🔗](\#load-attributes)

Type: bool

Path: load-attributes

Environment variable: B7S\_LoadAttributes

CLI flag:

- Flag: `--load-attributes`
- Default: false
- Description: node should try to load its attribute data from IPFS

## topics [🔗](\#topics)

Type: list (string)

Path: topics

Environment variable: B7S\_Topics

CLI flag:

- Flag: `--topics`
- Default: \[\]
- Description: topics node should subscribe to

## db [🔗](\#db)

Type: string

Path: db

Environment variable: B7S\_Db

CLI flag:

- Flag: `--db`
- Default: N/A
- Description: path to the database used for persisting peer and function data

## log [🔗](\#log)

Path: log

- ## level [🔗](\#log.level)


Type: string

Path: log.level

Environment variable: B7S\_Log\_Level
CLI flag:

- Flag: `--log-level`
- Shorthand: `-l`
- Default: info
- Description: log level to use

## connectivity [🔗](\#connectivity)

Path: connectivity

- ## address [🔗](\#connectivity.address)


Type: string

Path: connectivity.address

Environment variable: B7S\_Connectivity\_Address
CLI flag:

- Flag: `--address`
- Shorthand: `-a`
- Default: 0.0.0.0
- Description: address that the b7s host will use

- ## port [🔗](\#connectivity.port)


Type: uint

Path: connectivity.port

Environment variable: B7S\_Connectivity\_Port
CLI flag:

- Flag: `--port`
- Shorthand: `-p`
- Default: 0
- Description: port that the b7s host will use

- ## private-key [🔗](\#connectivity.private-key)


Type: string

Path: connectivity.private-key

Environment variable: B7S\_Connectivity\_PrivateKey
CLI flag:

- Flag: `--private-key`
- Default: N/A
- Description: private key that the b7s host will use

- ## dialback-address [🔗](\#connectivity.dialback-address)


Type: string

Path: connectivity.dialback-address

Environment variable: B7S\_Connectivity\_DialbackAddress
CLI flag:

- Flag: `--dialback-address`
- Default: N/A
- Description: external address that the b7s host will advertise

- ## dialback-port [🔗](\#connectivity.dialback-port)


Type: uint

Path: connectivity.dialback-port

Environment variable: B7S\_Connectivity\_DialbackPort
CLI flag:

- Flag: `--dialback-port`
- Default: 0
- Description: external port that the b7s host will advertise

- ## websocket [🔗](\#connectivity.websocket)


Type: bool

Path: connectivity.websocket

Environment variable: B7S\_Connectivity\_Websocket
CLI flag:

- Flag: `--websocket`
- Shorthand: `-w`
- Default: false
- Description: should the node use websocket protocol for communication

- ## websocket-port [🔗](\#connectivity.websocket-port)


Type: uint

Path: connectivity.websocket-port

Environment variable: B7S\_Connectivity\_WebsocketPort
CLI flag:

- Flag: `--websocket-port`
- Default: 0
- Description: port to use for websocket connections

- ## websocket-dialback-port [🔗](\#connectivity.websocket-dialback-port)


Type: uint

Path: connectivity.websocket-dialback-port

Environment variable: B7S\_Connectivity\_WebsocketDialbackPort
CLI flag:

- Flag: `--websocket-dialback-port`
- Default: 0
- Description: external port that the b7s host will advertise for websocket connections

- ## no-dialback-peers [🔗](\#connectivity.no-dialback-peers)


Type: bool

Path: connectivity.no-dialback-peers

Environment variable: B7S\_Connectivity\_NoDialbackPeers
CLI flag:

- Flag: `--no-dialback-peers`
- Default: false
- Description: start without dialing back peers from previous runs

- ## must-reach-boot-nodes [🔗](\#connectivity.must-reach-boot-nodes)


Type: bool

Path: connectivity.must-reach-boot-nodes

Environment variable: B7S\_Connectivity\_MustReachBootNodes
CLI flag:

- Flag: `--must-reach-boot-nodes`
- Default: false
- Description: halt node if we fail to reach boot nodes on start

- ## disable-connection-limits [🔗](\#connectivity.disable-connection-limits)


Type: bool

Path: connectivity.disable-connection-limits

Environment variable: B7S\_Connectivity\_DisableConnectionLimits
CLI flag:

- Flag: `--disable-connection-limits`
- Default: false
- Description: disable libp2p connection limits (experimental)

- ## connection-count [🔗](\#connectivity.connection-count)


Type: uint

Path: connectivity.connection-count

Environment variable: B7S\_Connectivity\_ConnectionCount
CLI flag:

- Flag: `--connection-count`
- Default: 0
- Description: maximum number of connections the b7s host will aim to have

## head [🔗](\#head)

Path: head

- ## rest-api [🔗](\#head.rest-api)


Type: string

Path: head.rest-api

Environment variable: B7S\_Head\_RestApi
CLI flag:

- Flag: `--rest-api`
- Default: N/A
- Description: address where the head node REST API will listen on

## worker [🔗](\#worker)

Path: worker

- ## runtime-path [🔗](\#worker.runtime-path)


Type: string

Path: worker.runtime-path

Environment variable: B7S\_Worker\_RuntimePath
CLI flag:

- Flag: `--runtime-path`
- Default: N/A
- Description: Blockless Runtime location (used by the worker node)

- ## runtime-cli [🔗](\#worker.runtime-cli)


Type: string

Path: worker.runtime-cli

Environment variable: B7S\_Worker\_RuntimeCli
CLI flag:

- Flag: `--runtime-cli`
- Default: N/A
- Description: runtime CLI name (used by the worker node)

- ## cpu-percentage-limit [🔗](\#worker.cpu-percentage-limit)


Type: float64

Path: worker.cpu-percentage-limit

Environment variable: B7S\_Worker\_CpuPercentageLimit
CLI flag:

- Flag: `--cpu-percentage-limit`
- Default: 0
- Description: amount of CPU time allowed for Blockless Functions in the 0-1 range, 1 being unlimited

- ## memory-limit [🔗](\#worker.memory-limit)


Type: int64

Path: worker.memory-limit

Environment variable: B7S\_Worker\_MemoryLimit
CLI flag:

- Flag: `--memory-limit`
- Default: 0
- Description: memory limit (kB) for Blockless Functions

## telemetry [🔗](\#telemetry)

Path: telemetry

- ## tracing [🔗](\#telemetry.tracing)


Path: telemetry.tracing

  - ## enable [🔗](\#telemetry.tracing.enable)


    Type: bool

    Path: telemetry.tracing.enable

    Environment variable: B7S\_Tracing\_Telemetry\_Enable
    CLI flag:

- Flag: `--enable-tracing`
- Default: false
- Description: emit tracing data

  - ## exporter-batch-timeout [🔗](\#telemetry.tracing.exporter-batch-timeout)


    Type: int64

    Path: telemetry.tracing.exporter-batch-timeout

    Environment variable: B7S\_Tracing\_Telemetry\_ExporterBatchTimeout

  - ## grpc [🔗](\#telemetry.tracing.grpc)


    Path: telemetry.tracing.grpc

    - ## endpoint [🔗](\#telemetry.tracing.grpc.endpoint)


      Type: string

      Path: telemetry.tracing.grpc.endpoint

      Environment variable: B7S\_Grpc\_Tracing\_Telemetry\_Endpoint
      CLI flag:

- Flag: `--tracing-grpc-endpoint`
- Default: N/A
- Description: tracing exporter GRPC endpoint
  - ## http [🔗](\#telemetry.tracing.http)


    Path: telemetry.tracing.http

    - ## endpoint [🔗](\#telemetry.tracing.http.endpoint)


      Type: string

      Path: telemetry.tracing.http.endpoint

      Environment variable: B7S\_Http\_Tracing\_Telemetry\_Endpoint
      CLI flag:

- Flag: `--tracing-http-endpoint`
- Default: N/A
- Description: tracing exporter HTTP endpoint
- ## metrics [🔗](\#telemetry.metrics)


Path: telemetry.metrics

  - ## enable [🔗](\#telemetry.metrics.enable)


    Type: bool

    Path: telemetry.metrics.enable

    Environment variable: B7S\_Metrics\_Telemetry\_Enable
    CLI flag:

- Flag: `--enable-metrics`
- Default: false
- Description: emit metrics

  - ## prometheus-address [🔗](\#telemetry.metrics.prometheus-address)


    Type: string

    Path: telemetry.metrics.prometheus-address

    Environment variable: B7S\_Metrics\_Telemetry\_PrometheusAddress
    CLI flag:

- Flag: `--prometheus-address`
- Default: N/A
- Description: address where prometheus metrics will be served

Node version: a3d41731844467fb8afea533e2a0cb9ae01fd0e8:2024-10-30T10:59:37Z