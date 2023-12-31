package cli_test

import (
	gocontext "context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/fury-labs/furya/v20/app/apptesting"
	"github.com/fury-labs/furya/v20/x/gamm/types"
)

type QueryTestSuite struct {
	apptesting.KeeperTestHelper
	queryClient types.QueryClient
}

func (s *QueryTestSuite) SetupSuite() {
	s.Setup()
	s.queryClient = types.NewQueryClient(s.QueryHelper)
	// create a new pool
	s.PrepareBalancerPool()
	s.Commit()
}

func (s *QueryTestSuite) TestQueriesNeverAlterState() {
	var (
		fooDenom   = apptesting.DefaultPoolAssets[0].Token.Denom
		barDenom   = apptesting.DefaultPoolAssets[1].Token.Denom
		bazDenom   = apptesting.DefaultPoolAssets[2].Token.Denom
		ufuryDenom = apptesting.DefaultPoolAssets[3].Token.Denom

		basicValidTokensIn = sdk.NewCoins(
			sdk.NewCoin(fooDenom, osmomath.OneInt()),
			sdk.NewCoin(barDenom, osmomath.OneInt()),
			sdk.NewCoin(bazDenom, osmomath.OneInt()),
			sdk.NewCoin(ufuryDenom, osmomath.OneInt()))
	)

	testCases := []struct {
		name   string
		query  string
		input  interface{}
		output interface{}
	}{
		{
			"Query pools",
			"/furya.gamm.v1beta1.Query/Pools",
			&types.QueryPoolsRequest{},
			&types.QueryPoolsResponse{},
		},
		{
			"Query single pool",
			"/furya.gamm.v1beta1.Query/Pool",
			&types.QueryPoolRequest{PoolId: 1},
			&types.QueryPoolsResponse{},
		},
		{
			"Query num pools",
			"/furya.gamm.v1beta1.Query/NumPools",
			&types.QueryNumPoolsRequest{},
			&types.QueryNumPoolsResponse{},
		},
		{
			"Query pool params",
			"/furya.gamm.v1beta1.Query/PoolParams",
			&types.QueryPoolParamsRequest{PoolId: 1},
			&types.QueryPoolParamsResponse{},
		},
		{
			"Query pool type",
			"/furya.gamm.v1beta1.Query/PoolType",
			&types.QueryPoolTypeRequest{PoolId: 1},
			&types.QueryPoolTypeResponse{},
		},
		{
			"Query spot price",
			"/furya.gamm.v1beta1.Query/SpotPrice",
			&types.QuerySpotPriceRequest{PoolId: 1, BaseAssetDenom: fooDenom, QuoteAssetDenom: barDenom},
			&types.QuerySpotPriceResponse{},
		},
		{
			"Query total liquidity",
			"/furya.gamm.v1beta1.Query/TotalLiquidity",
			&types.QueryTotalLiquidityRequest{},
			&types.QueryTotalLiquidityResponse{},
		},
		{
			"Query pool total liquidity",
			"/furya.gamm.v1beta1.Query/TotalPoolLiquidity",
			&types.QueryTotalPoolLiquidityRequest{PoolId: 1},
			&types.QueryTotalPoolLiquidityResponse{},
		},
		{
			"Query total shares",
			"/furya.gamm.v1beta1.Query/TotalShares",
			&types.QueryTotalSharesRequest{PoolId: 1},
			&types.QueryTotalSharesResponse{},
		},
		{
			"Query estimate for join pool shares with no swap",
			"/furya.gamm.v1beta1.Query/CalcJoinPoolNoSwapShares",
			&types.QueryCalcJoinPoolNoSwapSharesRequest{PoolId: 1, TokensIn: basicValidTokensIn},
			&types.QueryCalcJoinPoolNoSwapSharesResponse{},
		},
		{
			"Query estimate for join pool shares with no swap",
			"/furya.gamm.v1beta1.Query/CalcJoinPoolShares",
			&types.QueryCalcJoinPoolSharesRequest{PoolId: 1, TokensIn: basicValidTokensIn},
			&types.QueryCalcJoinPoolSharesResponse{},
		},
		{
			"Query exit pool coins from shares",
			"/furya.gamm.v1beta1.Query/CalcExitPoolCoinsFromShares",
			&types.QueryCalcExitPoolCoinsFromSharesRequest{PoolId: 1, ShareInAmount: osmomath.OneInt()},
			&types.QueryCalcExitPoolCoinsFromSharesResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			s.SetupSuite()
			err := s.QueryHelper.Invoke(gocontext.Background(), tc.query, tc.input, tc.output)
			s.Require().NoError(err)
			s.StateNotAltered()
		})
	}
}

func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}
