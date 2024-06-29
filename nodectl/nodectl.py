#!/usr/bin/python3
import sys
import os
import subprocess
import yaml
import ruamel.yaml
import argparse
parser = argparse.ArgumentParser()


parser.add_argument("cmd", help="Command to run")
parser.add_argument("-c", "--count", help="Number of nodes to create")
parser.add_argument("-b", "--bootnode", help="Boot node to connect to")

args = parser.parse_args()


nodes_dir = ".b7s_nodes"
pidfile = ".b7s_node_pids"
keyforge_executable = "/home/aco/code/blockless/b7s/cmd/keyforge/keyforge"
node_executable = "/home/aco/code/blockless/b7s/cmd/node/node"

default_cfg = {
    'role': 'worker',
    'log': { 'level': 'info' },
    'worker': {
        'runtime-path': '/home/aco/.local/blockless-runtime/bin',
        'runtime-cli': 'bls-runtime',
    },
    'connectivity': {
        'address': '127.0.0.1',
        'port': 0,              # set for each node
        'private-key': '',      # set for each node
    }
}

def create_keys(count):
    os.makedirs(nodes_dir, exist_ok=True)

    for i in range(count):
        print(f"Creating key for node {i}")
        create_key(i)


def create_key(i):
    dir = node_path(i)
    if os.path.exists(dir):
        return
    
    os.mkdir(dir)
    subprocess.run([keyforge_executable, "--output", dir])


def create_configs(count):
    os.makedirs(nodes_dir, exist_ok=True)

    for i in range(count):
        dir = config_path(i)
        if os.path.exists(dir):
            continue
        
        print(f"Creating config for node {i}")
        create_config(i)

def key_path(i):
    return node_path(i) + "/priv.bin"

def config_path(i):
    return node_path(i) + "/config.yaml"

def node_path(i):
    return nodes_dir + "/" + "node_" + str(i)

def keys_exist(count):
    for i in range(count):
        if not os.path.exists(key_path(i)):
            return False
    return True

def create_config(i):

    cfg = default_cfg.copy()
    cfg['connectivity']['port'] = 9000+i
    cfg['connectivity']['private-key'] = os.path.abspath(key_path(i))

    path = config_path(i)
    with open(path, 'w') as f:
        yaml = ruamel.yaml.YAML()
        yaml.indent(sequence=4, offset=2)
        yaml.dump(cfg, f)


def start_nodes(count):
    pf = os.open(pidfile, os.O_CREAT | os.O_WRONLY | os.O_TRUNC)
    for i in range(count):
        os.makedirs(node_path(i), exist_ok=True)

        if not os.path.exists(key_path(i)):
           create_key(i)
        
        if not os.path.exists(config_path(i)):
            create_config(i)

        cmd = [node_executable, "--config", config_path(i)]
        if args.bootnode:
            cmd.append("--boot-nodes")
            cmd.append(args.bootnode)

        path = node_path(i)
        stderr_path = f"{path}/stderr.log"
        stdout_path = f"{path}/stdout.log"

        print(f"starting node {i} with command {' '.join(cmd)}")
        with open(stdout_path, 'w') as stdout_file, open(stderr_path, 'w') as stderr_file:
            proc = subprocess.Popen(
                cmd,
                stdout=stdout_file,
                stderr=stderr_file
            )

            pid = proc.pid
            os.write(pf, (str(pid) + "\n").encode())
            print("started node " + str(i) + " with pid " + str(pid))

    os.close(pf)

def stop_nodes():
    pf = os.open(pidfile, os.O_RDONLY)
    pids = []
    for line in os.read(pf, os.path.getsize(pidfile)).decode().splitlines():
        pids.append(int(line))
    os.close(pf)

    print(f"Stopping {len(pids)} nodes with pids {pids}")

    cmd = ["kill", "-9"] + list(map(str, pids))
    subprocess.Popen(cmd)

## main


if args.cmd == "keys":
    create_keys(int(args.count))
elif args.cmd == "configs":
    create_configs(int(args.count))
elif args.cmd == "start":
    start_nodes(int(args.count))
elif args.cmd == "stop":
    stop_nodes()
else:
    print(parser.format_help)
