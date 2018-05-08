// Copyright © 2018 Rob Mulholand <rmulholand@8thlight.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/8thlight/block_watcher/pkg/db"
	"github.com/8thlight/block_watcher/pkg/ipfs"
	"github.com/8thlight/block_watcher/pkg/ipfs/eth_block_header"
	"github.com/8thlight/block_watcher/pkg/transformers"
)

// createIpldForBlockHeadersCmd represents the createIpldForBlockHeaders command
var createIpldForBlockHeadersCmd = &cobra.Command{
	Use:   "createIpldForBlockHeaders",
	Short: "Create IPLD objects for multiple blocks.",
	Long: `Create IPLD objects for multiple blocks.

e.g. ./block_watcher createIpldForBlockHeaders -s 1234567 -e 4567890

Under the hood, the command fetches the block header RLP data from LevelDB and
puts it in IPFS, converting the data as an 'eth-block'.`,
	Run: func(cmd *cobra.Command, args []string) {
		createIpldForBlockHeaders()
	},
}

func init() {
	rootCmd.AddCommand(createIpldForBlockHeadersCmd)
	createIpldForBlockHeadersCmd.Flags().Int64VarP(&startingBlockNumber, "starting-block-number", "s", 0, "First block number to create IPLD for.")
	createIpldForBlockHeadersCmd.Flags().Int64VarP(&endingBlockNumber, "ending-block-number", "e", 5430000, "Last block number to create IPLD for.")
	createIpldForBlockHeadersCmd.Flags().BoolVarP(&useParity, "parity", "p", false, "Use Parity's Rocks DB instead of Geth's Level DB")
}

func createIpldForBlockHeaders() {
	// init eth db
	var ethDBConfig db.DatabaseConfig
	if useParity {
		ethDBConfig = db.CreateDatabaseConfig(db.Rocks, rocksDbPath)
	} else {
		ethDBConfig = db.CreateDatabaseConfig(db.Level, levelDbPath)
	}
	ethDB, err := db.CreateDatabase(ethDBConfig)
	if err != nil {
		log.Fatal("Error connecting to ethereum db: ", err)
	}

	// init ipfs publisher
	ipfsNode, err := ipfs.InitIPFSNode(ipfsPath)
	if err != nil {
		log.Fatal("Error connecting to IPFS: ", err)
	}
	decoder := db.RlpDecoder{}
	dagPutter := eth_block_header.NewBlockHeaderDagPutter(*ipfsNode, decoder)
	publisher := ipfs.NewIpfsPublisher(dagPutter)

	// execute transformer
	transformer := transformers.NewEthBlockHeaderTransformer(ethDB, publisher)
	err = transformer.Execute(startingBlockNumber, endingBlockNumber)
	if err != nil {
		log.Fatal("Error executing transformer: ", err.Error())
	}
}
