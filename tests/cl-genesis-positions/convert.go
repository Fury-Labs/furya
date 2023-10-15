package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/furya-labs/furya/osmomath"
	"github.com/furya-labs/furya/osmoutils/accum"
	"github.com/furya-labs/furya/v19/app"
	"github.com/furya-labs/furya/v19/app/apptesting"
	cl "github.com/furya-labs/furya/v19/x/concentrated-liquidity"
	"github.com/furya-labs/furya/v19/x/concentrated-liquidity/math"
	"github.com/furya-labs/furya/v19/x/concentrated-liquidity/model"
	cltypes "github.com/furya-labs/furya/v19/x/concentrated-liquidity/types"
	clgenesis "github.com/furya-labs/furya/v19/x/concentrated-liquidity/types/genesis"
)

type BigBangPositions struct {
	PositionsData []clgenesis.PositionData `json:"position_data"`
}

type FuryaApp struct {
	App         *app.FuryaApp
	Ctx         sdk.Context
	QueryHelper *baseapp.QueryServiceTestHelper
	TestAccs    []sdk.AccAddress
}

var furyaPrecision = 6

func ReadSubgraphDataFromDisk(subgraphFilePath string) []SubgraphPosition {
	// read in the data from file
	data, err := os.ReadFile(subgraphFilePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// unmarshal the data into a slice of Position structs
	var positions []SubgraphPosition
	err = json.Unmarshal(data, &positions)
	if err != nil {
		fmt.Println("Error unmarshalling data:", err)
		os.Exit(1)
	}

	return positions
}

func ConvertSubgraphToFuryaGenesis(positionCreatorAddresses []sdk.AccAddress, subgraphFilePath string) (*clgenesis.GenesisState, *banktypes.GenesisState) {
	positions := ReadSubgraphDataFromDisk(subgraphFilePath)

	furya := apptesting.KeeperTestHelper{}
	furya.Setup()

	if len(positionCreatorAddresses) == 0 {
		panic("no accounts found")
	}

	furya.TestAccs = make([]sdk.AccAddress, len(positionCreatorAddresses))
	for i := 0; i < len(positionCreatorAddresses); i++ {
		furya.TestAccs[i] = positionCreatorAddresses[i]
		fmt.Println(furya.TestAccs[i].String())
	}

	initAmounts := sdk.NewCoins(
		sdk.NewCoin(denom0, osmomath.NewInt(1000000000000000000)),
		sdk.NewCoin(denom1, osmomath.NewInt(1000000000000000000)),
	)

	// fund all accounts
	for _, acc := range furya.TestAccs {
		err := simapp.FundAccount(furya.App.BankKeeper, furya.Ctx, acc, initAmounts)
		if err != nil {
			panic(err)
		}
	}

	msgCreatePool := model.MsgCreateConcentratedPool{
		Sender:       furya.TestAccs[0].String(),
		Denom0:       denom0,
		Denom1:       denom1,
		TickSpacing:  100,
		SpreadFactor: osmomath.MustNewDecFromStr("0.0005"),
	}

	err := msgCreatePool.ValidateBasic()
	if err != nil {
		panic(err)
	}

	poolId, err := furya.App.PoolManagerKeeper.CreatePool(furya.Ctx, msgCreatePool)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created pool id of: ", poolId)

	pool, err := furya.App.ConcentratedLiquidityKeeper.GetConcentratedPoolById(furya.Ctx, poolId)
	if err != nil {
		panic(err)
	}

	// Initialize first position to be 1:1 price
	// this is because the first position must have non-zero token0 and token1 to initialize the price
	// however, our data has first position with non-zero amount.
	_, err = furya.App.ConcentratedLiquidityKeeper.CreateFullRangePosition(furya.Ctx, pool.GetId(), furya.TestAccs[0], sdk.NewCoins(sdk.NewCoin(msgCreatePool.Denom0, osmomath.NewInt(100)), sdk.NewCoin(msgCreatePool.Denom1, osmomath.NewInt(100))))
	if err != nil {
		panic(err)
	}

	clMsgServer := cl.NewMsgServerImpl(furya.App.ConcentratedLiquidityKeeper)

	numberOfSuccesfulPositions := 0

	bigBangPositions := make([]clgenesis.PositionData, 0)

	for _, uniV3Position := range positions {
		lowerPrice := parsePrice(uniV3Position.TickLower.Price0)
		upperPrice := parsePrice(uniV3Position.TickUpper.Price0)
		if err != nil {
			panic(err)
		}

		if lowerPrice.GTE(upperPrice) {
			fmt.Printf("lowerPrice (%s) >= upperPrice (%s), skipping", lowerPrice, upperPrice)
			continue
		}

		sqrtPriceLower, err := lowerPrice.ApproxRoot(2)
		if err != nil {
			panic(err)
		}
		lowerTickFurya, err := math.SqrtPriceToTickRoundDownSpacing(osmomath.BigDecFromDec(sqrtPriceLower), pool.GetTickSpacing())
		if err != nil {
			panic(err)
		}

		sqrtPriceUpper, err := upperPrice.ApproxRoot(2)
		if err != nil {
			panic(err)
		}
		upperTickFurya, err := math.SqrtPriceToTickRoundDownSpacing(osmomath.BigDecFromDec(sqrtPriceUpper), pool.GetTickSpacing())
		if err != nil {
			panic(err)
		}

		if lowerTickFurya > upperTickFurya {
			fmt.Printf("lowerTickFurya (%d) > upperTickFurya (%d), skipping", lowerTickFurya, upperTickFurya)
			continue
		}

		if lowerTickFurya == upperTickFurya {
			// bump up the upper tick by one. We don't care about having exactly the same tick range
			// Just a roughly similar breakdown
			upperTickFurya = upperTickFurya + 1
		}

		depositedAmount0, failedParsing := parseStringToInt(uniV3Position.DepositedToken0)
		if failedParsing {
			fmt.Printf("Failed parsing %s, skipping", uniV3Position.DepositedToken0)
			continue
		}

		depositedAmount1, failedParsing := parseStringToInt(uniV3Position.DepositedToken1)
		if failedParsing {
			fmt.Printf("Failed parsing %s, skipping", uniV3Position.DepositedToken0)
			continue
		}

		tokensProvided := sdk.NewCoins(sdk.NewCoin(msgCreatePool.Denom0, depositedAmount0), sdk.NewCoin(msgCreatePool.Denom1, depositedAmount1))

		randomCreator := furya.TestAccs[rand.Intn(len(furya.TestAccs))]

		position, err := clMsgServer.CreatePosition(sdk.WrapSDKContext(furya.Ctx), &cltypes.MsgCreatePosition{
			PoolId:          poolId,
			Sender:          randomCreator.String(),
			LowerTick:       lowerTickFurya,
			UpperTick:       upperTickFurya,
			TokensProvided:  tokensProvided,
			TokenMinAmount0: osmomath.ZeroInt(),
			TokenMinAmount1: osmomath.ZeroInt(),
		})
		if err != nil {
			fmt.Printf("\n\n\nWARNING: Failed to create position: %v\n\n\n", err)
			fmt.Printf("attempted creation between ticks (%d) and (%d), desired amount 0: (%s), desired amount 1 (%s)\n", lowerTickFurya, upperTickFurya, depositedAmount0, depositedAmount1)
			fmt.Println()
			continue
		}

		fmt.Printf("created position with liquidity (%s) between ticks (%d) and (%d)\n", position.LiquidityCreated, lowerTickFurya, upperTickFurya)
		numberOfSuccesfulPositions++

		bigBangPositions = append(bigBangPositions, clgenesis.PositionData{
			Position: &model.Position{
				Address:    randomCreator.String(),
				PoolId:     poolId,
				JoinTime:   furya.Ctx.BlockTime(),
				Liquidity:  position.LiquidityCreated,
				PositionId: position.PositionId,
				LowerTick:  lowerTickFurya,
				UpperTick:  upperTickFurya,
			},
			SpreadRewardAccumRecord: accum.Record{
				NumShares: position.LiquidityCreated,
			},
		})
	}

	fmt.Printf("\nout of %d uniswap positions, %d were successfully created\n", len(positions), numberOfSuccesfulPositions)

	if writeBigBangConfigToDisk {
		writeBigBangPositionsToState(bigBangPositions)
	}

	if writeGenesisToDisk {
		state := furya.App.ExportState(furya.Ctx)
		writeStateToDisk(state)
	}

	clGenesis := furya.App.ConcentratedLiquidityKeeper.ExportGenesis(furya.Ctx)
	bankGenesis := furya.App.BankKeeper.ExportGenesis(furya.Ctx)
	return clGenesis, bankGenesis
}

func parsePrice(strPrice string) (result osmomath.Dec) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Printf("Recovered in price parsing %s, %s\n", r, strPrice)
			if strPrice[0] == '0' {
				result = cltypes.MinSpotPrice
			} else {
				result = cltypes.MaxSpotPrice
			}
		}

		if result.GT(cltypes.MaxSpotPrice) {
			result = cltypes.MaxSpotPrice
		}

		if result.LT(cltypes.MinSpotPrice) {
			result = cltypes.MinSpotPrice
		}
	}()
	result = osmomath.MustNewBigDecFromStr(strPrice).Dec()
	return result
}

func parseStringToInt(strInt string) (result osmomath.Int, failedParsing bool) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Printf("Recovered in int parsing %s, %s\n", r, strInt)
			failedParsing = true
		}
	}()
	result = osmomath.MustNewBigDecFromStr(strInt).Dec().MulInt64(int64(furyaPrecision)).TruncateInt()
	return result, failedParsing
}

func writeStateToDisk(state map[string]json.RawMessage) {
	stateBz, err := json.MarshalIndent(state, "", "    ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(pathToFilesFromRoot+furyaGenesisFileName, stateBz, 0o644)
	if err != nil {
		panic(err)
	}
}

func writeBigBangPositionsToState(positions []clgenesis.PositionData) {
	fmt.Println("writing big bang positions to disk")
	positionsBytes, err := json.MarshalIndent(BigBangPositions{
		PositionsData: positions,
	}, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(pathToFilesFromRoot+bigbangPosiionsFileName, positionsBytes, 0o644)
	if err != nil {
		panic(err)
	}
}
