# Block Watcher

## Description
A [VulcanizeDB](https://github.com/vulcanize/VulcanizeDB) transformer for creating IPLDs for Ethereum block data.

## Dependencies
 - Go 1.9+
 - Postgres 10
 - Ethereum Node
   - [Go Ethereum](https://ethereum.github.io/go-ethereum/downloads/) (1.8+)
   - [Parity 1.8.11+](https://github.com/paritytech/parity/releases): *UNSUPPORTED*
 - [IPFS](https://github.com/ipfs/go-ipfs#build-from-source)
 - [go-ipld-eth](https://github.com/ipfs/go-ipld-eth) (Plugin enabling conversion of block headers to IPLDs in IPFS)

## Installation
1. Setup Postgres and an Ethereum node - see [VulcanizeDB README](https://github.com/vulcanize/VulcanizeDB/blob/master/README.md).
1. Sync VulcanizeDB to populate core block data (commands will read block data from VulcanizeDB to fetch and persist block RLP data).
1. `git clone git@github.com:8thlight/block_watcher.git`

  _note: `go get` does not work for this project because need to run the (fixlibcrypto)[https://github.com/8thlight/sai_watcher/blob/master/Makefile] command along with `go build`._
1. Build:
    ```
    make build
    ```

## Configuration
- To use a local Ethereum node, copy `environments/public.toml.example` to
  `environments/public.toml` and update the `levelDbPath` to the local node's levelDB filepath:
  - when using geth:
    - The LevelDB file is called `chaindata`.
    - The LevelDB file path is printed to the console when you start geth.
    - The default location is:
      - Mac: `$HOME/Library/Ethereum`
      - Linux: `$HOME/.ethereum`
  - when using parity:
    - The RocksDB file is called `db`.
    - The default location on Mac is: `/Users/$USER/Library/Application Support/io.parity.ethereum/chains/ethereum/db/906a34e69aec8c0d/overlayrecent/db`

## Running the createIpldForBlockHeader command
- This command creates an IPLD for the header of a single Ethereum block.
- `./block_watcher createIpldForBlockHeader --config <config.toml> --block-number <block-number>`
- Optionally, point at Parity's RocksDB instead of Geth's LevelDB with flag `-p`.
- Note: block number argument is required.

## Running the createIpldForBlockHeaders command
- This command creates IPLDs for headers in a range of Ethereum blocks.
- `./block_watcher createIpldForBlockHeaders --config <config.toml> --starting-block-number <block-number> --ending-block-number <block-number>`
- Optionally, point at Parity's RocksDB instead of Geth's LevelDB with flag `-p`.
- Note: starting and ending block number arguments are required, and ending block number must be greater than starting block number.

## Running the createIpldsForBlockTransactions command
- This command creates IPLDs for transactions on an Ethereum block.
- `./block_watcher createIpldsForBlockTransactions --config <config.toml> --block-number <block-number>`
- No Parity support at this time.
- Note: block number argument is required.

## Running the createIpldsForBlocksTransactions command
- This command creates IPLDs for transactions on a range of Ethereum blocks.
- `./block_watcher createIpldsForBlocksTransactions --config <config.toml> --starting-block-number <block-number> --ending-block-number <block-number>`
- No Parity support at this time.
- Note: starting and ending block number arguments are required, and ending block number must be greater than starting block number.

## Running the tests
```
ginkgo -r
```
