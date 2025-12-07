-include .env

.PHONY: all test clean deploy fund help install format anvil

DEFAULT_ANVIL_KEY := 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

all: clean remove install update build

# Clean the repo
clean  :; forge clean

# Remove modules
remove :; rm -rf .gitmodules && rm -rf .git/modules/* && rm -rf lib && touch .gitmodules && git add . && git commit -m "modules"

install :; forge install cyfrin/foundry-devops@0.2.2 --no-commit && forge install foundry-rs/forge-std@v1.8.2 --no-commit && forge install openzeppelin/openzeppelin-contracts@v5.0.2 --no-commit

# Update Dependencies
update:; forge update

build:; forge build

test :; forge test 

format :; forge fmt

anvil :; anvil -m 'test test test test test test test test test test test junk' --steps-tracing --block-time 1

NETWORK_ARGS := --rpc-url http://localhost:8545 --private-key $(DEFAULT_ANVIL_KEY) --broadcast

ifeq ($(findstring --network sepolia,$(ARGS)),--network sepolia)
	NETWORK_ARGS := --rpc-url $(SEPOLIA_RPC_URL) --private-key $(SEPOLIA_PRIVATE_KEY) --broadcast --verify --etherscan-api-key $(ETHERSCAN_API_KEY) -vvvv
endif

ifeq ($(findstring --network polygon,$(ARGS)),--network polygon)
	NETWORK_ARGS := --rpc-url $(POLYGON_RPC_URL) --private-key $(POLYGON_PRIVATE_KEY) --broadcast --verify --etherscan-api-key $(ETHERSCAN_API_KEY) -vvvv
endif

deployDecmAccessManager:
	@forge script script/DeployDecmAccessManager.s.sol:DeployDecmAccessManager $(NETWORK_ARGS)

deployEvent:
	@forge script script/DeployEvent.s.sol:DeployEvent $(NETWORK_ARGS)

deployEventTicket:
	@forge script script/DeployEventTicket.s.sol:DeployEventTicket $(NETWORK_ARGS)

deployEventCertificate:
	@forge script script/DeployEventCertificate.s.sol:DeployEventCertificate $(NETWORK_ARGS)

# Example
# make deployDecmAccessManager ARGS="--network sepolia"
# make deployDecmAccessManager ARGS="--network polygon"
# make deployEvent ARGS="--network sepolia"
# make deployEventTicket ARGS="--network sepolia"
# make deployEventCertificate ARGS="--network sepolia"