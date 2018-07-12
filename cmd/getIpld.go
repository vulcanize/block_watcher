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
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/factories"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/go-cid"
	"log"
)

// getIpldCmd represents the getIpld command
var getIpldCmd = &cobra.Command{
	Use:   "getIpld",
	Short: "Get the IPLD for an Ethereum data structure",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		getIpld()
	},
}

func init() {
	rootCmd.AddCommand(getIpldCmd)
	getIpldCmd.Flags().StringVarP(&cid, "cid", "c", "", "Cid to fetch IPLD for")
}

func getIpld() {
	// init deps
	getter, err := ipfs.InitIPFSNode(ipfsPath)
	if err != nil {
		log.Fatal("Error connecting to ipfs: ", err)
	}
	cidDecoder := go_cid.NewCidDecoder()
	decodedCid, err := cidDecoder.Decode(cid)
	if err != nil {
		log.Fatal("Cannot decode cid: ", err)
	}
	resolverFactory := factories.ResolverFactory{}
	resolver, err := resolverFactory.GetResolver(decodedCid)
	if err != nil {
		log.Fatal("Cannot locate resolver for cid: ", err)
	}

	// init ipfs reader
	reader := ipfs.NewIpldReader(getter, resolver)

	// fetch and print IPLD data
	node, err := reader.Read(decodedCid)
	if err != nil {
		log.Fatal("Error reading cid: ", err)
	}
	marshalled, err := json.Marshal(node)
	log.Println(string(marshalled))
}
