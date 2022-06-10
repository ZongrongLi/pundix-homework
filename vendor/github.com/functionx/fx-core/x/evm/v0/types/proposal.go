package types

import (
	"fmt"
	"strings"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeInitEvmParams defines the type for a InitCrossChainParamsProposal
	ProposalTypeInitEvmParams = "InitEvmParams"
)

var (
	_ govtypes.Content = &InitEvmParamsProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeInitEvmParams)
}

// Proposal handler

func (m *InitEvmParamsProposal) GetTitle() string {
	return m.Title
}

func (m *InitEvmParamsProposal) GetDescription() string {
	return m.Description
}

func (m *InitEvmParamsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *InitEvmParamsProposal) ProposalType() string {
	return ProposalTypeInitEvmParams
}

func (m *InitEvmParamsProposal) ValidateBasic() error {

	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}
	if err := m.EvmParams.Validate(); err != nil {
		return err
	}
	if err := m.FeemarketParams.Validate(); err != nil {
		return err
	}
	return nil
}

func (m *InitEvmParamsProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Init Evm Params Proposal:
  Title:       %s
  Description: %s
  EvmParams: %v
  FeeMarketParams: %v
`, m.Title, m.Description, m.EvmParams, m.FeemarketParams))
	return b.String()
}
