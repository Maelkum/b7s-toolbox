# Workito

Workito is small command-line tool to directly instruct a Blockless Worker node to execute an execution request.
This makes it possible to test stuff by running a headless Worker node.
It supports sending installation messages to the node, as well as request executions.

It is possible to request installation of multiple functions, specified using their CIDs.
The tool will try to wait a limited amount of time for the node to confirm installation.

## Notes

Right now the tool generates a random identity on each run. This can be a hassle because the Worker node will remember identities of past peers, and this will in fact pollute its database of known peers.

The tool does not support some more ellaborate scenarios like roll calls or clustered execution.

## Examples

### Install a Function

```console
$ workito install --address /ip4/127.0.0.1/tcp/9000/p2p/12D3KooWHUeKgXT4aj8oKvtwovVMki468igSsa5F8izZY3U5UyMD bafybeia24v4czavtpjv2co3j54o4a5ztduqcpyyinerjgncx7s2s22s7ea
2024/01/27 15:56:10 INFO node address address=/ip4/127.0.0.1/tcp/9000/p2p/12D3KooWHUeKgXT4aj8oKvtwovVMki468igSsa5F8izZY3U5UyMD
2024/01/27 15:56:10 INFO libp2p host we use id=12D3KooWFNEE9eo7rSeMPLNvFZuhVJ3HLxeuj255Exrp9tKNKzJ2
2024/01/27 15:56:10 INFO connected to node
2024/01/27 15:56:10 INFO installing function cid=bafybeia24v4czavtpjv2co3j54o4a5ztduqcpyyinerjgncx7s2s22s7ea
2024/01/27 15:56:12 INFO received install function response cid=bafybeia24v4czavtpjv2co3j54o4a5ztduqcpyyinerjgncx7s2s22s7ea status=installed
2024/01/27 15:56:12 INFO all functions installed
```

### Execute a Function

```console
$ workito --address /ip4/127.0.0.1/tcp/9000/p2p/12D3KooWHUeKgXT4aj8oKvtwovVMki468igSsa5F8izZY3U5UyMD execute --function-id bafybeia24v4czavtpjv2co3j54o4a5ztduqcpyyinerjgncx7s2s22s7ea --method hello-world.wasm --nodes 1 --parameter name=val --env-var NAME=VAL --env-var LEN=45  | jq .| jq
2024/01/27 15:57:12 INFO node address address=/ip4/127.0.0.1/tcp/9000/p2p/12D3KooWHUeKgXT4aj8oKvtwovVMki468igSsa5F8izZY3U5UyMD
2024/01/27 15:57:12 INFO libp2p host we use id=12D3KooWSdpSQwFUvKGUvg67Dze2ayA8HzQ6M4ZJJBviWaHd74ak
2024/01/27 15:57:12 INFO connected to node
2024/01/27 15:57:12 INFO received execution response
{
  "type": "MsgExecuteResponse",
  "request_id": "24e7146f-1e81-44f8-9388-3ae947a63358",
  "code": "200",
  "results": {
    "12D3KooWHUeKgXT4aj8oKvtwovVMki468igSsa5F8izZY3U5UyMD": {
      "code": "200",
      "result": {
        "stdout": "Hello, world!\n",
        "stderr": "",
        "exit_code": 0
      },
      "request_id": "24e7146f-1e81-44f8-9388-3ae947a63358",
      "usage": {
        "wall_clock_time": 49012128,
        "cpu_user_time": 124668000,
        "cpu_sys_time": 24933000,
        "memory_max_kb": 45700
      }
    }
  },
  "cluster": {},
  "pbft": {
    "view": 0,
    "request_timestamp": "0001-01-01T00:00:00Z"
  }
}
```