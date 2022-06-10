package clients

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_QueryParams(t *testing.T) {
	paramsRes, err := DistrQueryClientInstance.QueryParams()
	require.NoError(t, err)
	require.NotNil(t, paramsRes)
	t.Log("===>> QueryParams resp info", paramsRes)
}

func Test_ValidatorCommission(t *testing.T) {
	validatorCommissionRes, err := DistrQueryClientInstance.ValidatorCommission(singaporeValidator)
	require.NoError(t, err)
	require.NotNil(t, validatorCommissionRes)
	t.Log("===>>validatorCommissionRes resp info", validatorCommissionRes)
}

func Test_ValidatorOutstandingRewards(t *testing.T) {
	validatorOutstandingRes, err := DistrQueryClientInstance.ValidatorOutstandingRewards(singaporeValidator)
	require.NoError(t, err)
	require.NotNil(t, validatorOutstandingRes)
	t.Log("===>> validatorOutstandingRes resp info", validatorOutstandingRes)
}

func Test_CommunityPool(t *testing.T) {
	communityPoolRes, err := DistrQueryClientInstance.CommunityPool()
	require.NoError(t, err)
	require.NotNil(t, communityPoolRes)
	t.Log("===>> CommunityPool resp info", communityPoolRes)
}

func Test_Balance(t *testing.T) {
	bankBalanceRes, err := BankQueryClientInstance.Balance(userAccount1)
	require.NoError(t, err)
	require.NotNil(t, bankBalanceRes)
	t.Log("===>> Balance resp info", bankBalanceRes)
}

func Test_TotalSupply(t *testing.T) {
	bankTotalRes, err := BankQueryClientInstance.TotalSupply()
	require.NoError(t, err)
	require.NotNil(t, bankTotalRes)
	t.Log("===>> TotalSupply resp info", bankTotalRes)
}
