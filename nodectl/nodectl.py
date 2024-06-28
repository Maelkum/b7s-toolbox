#!/usr/bin/python3
import sys
import os
import subprocess
import yaml
import ruamel.yaml

usage = """
nodectl.py <keys|start> [count] - create keys or start nodes
nodectl.py <stop|help> - stop nodes or show help
"""

nodes_dir = ".b7s_nodes"
pidfile = ".b7s_node_pids"
keyforge_executable = "/home/aco/code/blockless/b7s/cmd/keyforge/keyforge"

default_cfg = {
    'role': 'worker',
    'log': { 'level': 'info' },
    'worker': {
        'runtime-path:': '/home/aco/.local/blockless-runtime/bin',
        'runtime-cli': 'bls-runtime',
    },
    'connectivity': {
        'address': '127.0.0.1',
        'port': 0,              # set for each node
        'private-key': '',      # set for each node
    }
}


def create_keys(count):
    if not os.path.exists(nodes_dir):
        os.mkdir(nodes_dir)

    for i in range(count):
        dir = node_path(i)
        if os.path.exists(dir):
            continue
        
        os.mkdir(dir)
        subprocess.run([keyforge_executable, "--output", dir])

def create_configs(count):
    if not os.path.exists(nodes_dir):
        os.mkdir(nodes_dir)

    for i in range(count):
        dir = config_path(i)
        if os.path.exists(dir):
            continue

        create_config(i)

def key_path(i):
    return node_path(i) + "/key"

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


## main

cmd = sys.argv[1]

if cmd == "keys":
    count = int(sys.argv[2])
    create_keys(count)
elif cmd == "configs":
    count = int(sys.argv[2])
    create_configs(count)
elif cmd == "start":
    count = int(sys.argv[2])
    print("start")
elif cmd == "stop":
    print("stop")
elif cmd == "help":
    print(usage)
else:
    print("unknown command")


print ("done")

