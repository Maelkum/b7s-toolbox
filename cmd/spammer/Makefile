# Variables
NODE_MULTIADDRESS = /ip4/127.0.0.1/tcp/9910/p2p/12D3KooWDqp65e5vmHm5aRBxrRh3H8vvcMjHzsJ8LSEqumRwoD7f
KEY_PATH = priv.bin

# Targets
all: build

build:
	go build

# clean does not remove the built binary
clean:
	rm -f ./spammer_*.log
	rm -f ./stats_*.txt

run:
	./spammer --node $(NODE_MULTIADDRESS) --key $(KEY_PATH)

key:
	b7s-keyforge -o .