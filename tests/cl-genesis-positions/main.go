package main

import (
	"flag"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// operation defines the desired operation to be run by this script.
type operation int

const (
	// getData retrieves the data from the Uniswap subgraph and writes it to disk
	// under pathToFilesFromRoot + positionsFileName path.
	getData operation = iota
	// convertPositions converts the data from the Uniswap subgraph into Furya
	// genesis. It reads pathToFilesFromRoot + positionsFileName path
	// run Furya app via apptesting, creates positions and writes the genesis
	// under pathToFilesFromRoot + furyaStateFileName path.
	convertPositions
	// mergeSubgraphAndLocalFuryaGenesis merges the genesis created from the subgraph data
	// with the localfurya genesis. This command is meant to be called inside the localfurya
	// container during setup (see setup.sh). It reads the existing genesis from localfuryaHomePath,
	// updates the concentrated liquidity section to append the CL pool created from the subgraph data,
	// its positions, ticks and accumulators.
	mergeSubgraphAndLocalFuryaGenesis
)

const (
	pathToFilesFromRoot = "tests/cl-genesis-positions/"

	positionsFileName       = "subgraph_positions.json"
	furyaGenesisFileName  = "genesis.json"
	bigbangPosiionsFileName = "bigbang_positions.json"

	localFuryaHomePath = "/furya/.furyad/"

	denom0 = "uusdc"
	denom1 = "ufury"
)

var (
	// This is lo-test1 address in localfurya
	defaultCreatorAddresses = []sdk.AccAddress{sdk.MustAccAddressFromBech32("osmo1cyyzpxplxdzkeea7kwsydadg87357qnahakaks"), sdk.MustAccAddressFromBech32("osmo18s5lynnmx37hq4wlrw9gdn68sg2uxp5rgk26vv")}

	useKeyringAccounts bool

	writeGenesisToDisk bool

	writeBigBangConfigToDisk bool
)

func main() {
	var (
		desiredOperation int
		isLocalFurya   bool
	)

	flag.BoolVar(&writeBigBangConfigToDisk, "big-bang", false, fmt.Sprintf("flag indicating whether to write the big bang config to disk at path %s", bigbangPosiionsFileName))
	flag.BoolVar(&writeGenesisToDisk, "genesis", false, fmt.Sprintf("flag indicating whether to write the genesis file to disk at path %s", furyaGenesisFileName))
	flag.BoolVar(&useKeyringAccounts, "keyring", false, "flag indicating whether to use local test keyring accounts")
	flag.BoolVar(&isLocalFurya, "localfurya", false, "flag indicating whether this is being run inside the localfurya container")
	flag.IntVar(&desiredOperation, "operation", 0, fmt.Sprintf("operation to run:\nget subgraph data: %v, convert subgraph positions to fury genesis: %v\nmerge converted subgraph genesis and localfurya genesis: %v", getData, convertPositions, mergeSubgraphAndLocalFuryaGenesis))

	flag.Parse()

	fmt.Println("isLocalFurya:", isLocalFurya)

	pathToSaveFilesAt := pathToFilesFromRoot
	if isLocalFurya {
		pathToSaveFilesAt = ""
	}

	// Set this to one of the 'operation' values
	switch operation(desiredOperation) {
	// See definition for more info.
	case getData:
		fmt.Println("Getting data from Uniswap subgraph...")

		GetUniV3SubgraphData(pathToSaveFilesAt + positionsFileName)
		// See definition for more info.
	case convertPositions:
		fmt.Println("Converting positions from subgraph data to Furya genesis...")

		var creatorAddresses []sdk.AccAddress
		if useKeyringAccounts {
			fmt.Println("Using local keyring addresses as creators")
			creatorAddresses = GetLocalKeyringAccounts()
		} else {
			fmt.Println("Using default creator addresses")
			creatorAddresses = defaultCreatorAddresses
		}

		ConvertSubgraphToFuryaGenesis(creatorAddresses, pathToSaveFilesAt+positionsFileName)
		// See definition for more info.
	case mergeSubgraphAndLocalFuryaGenesis:
		fmt.Println("Merging subgraph and local Furya genesis...")
		clState, bankState := ConvertSubgraphToFuryaGenesis(defaultCreatorAddresses, pathToSaveFilesAt+positionsFileName)

		EditLocalFuryaGenesis(clState, bankState)
	default:
		panic("Invalid operation")
	}
}
