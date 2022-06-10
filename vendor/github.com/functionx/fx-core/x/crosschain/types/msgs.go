package types

import (
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// cross chain message types
const (
	TypeMsgSetOrchestratorAddress = "set_orchestrator_address"
	TypeMsgAddOracleDeposit       = "add_oracle_deposit"
	TypeMsgOracleSetConfirm       = "valset_confirm"
	TypeMsgOracleSetUpdatedClaim  = "valset_updated_claim"

	TypeMsgBridgeTokenClaim = "bridge_token_claim"

	TypeMsgSendToFxClaim = "send_to_fx_claim"

	TypeMsgSendToExternal       = "send_to_external"
	TypeMsgCancelSendToExternal = "cancel_send_to_external"
	TypeMsgSendToExternalClaim  = "send_to_external_claim"

	TypeMsgRequestBatch = "request_batch"
	TypeMsgConfirmBatch = "confirm_batch"
)

type (
	// CrossChainMsg cross msg must implement GetChainName interface.. using in router
	CrossChainMsg interface {
		GetChainName() string
	}
)

var (
	_ sdk.Msg       = &MsgSetOrchestratorAddress{}
	_ CrossChainMsg = &MsgSetOrchestratorAddress{}
	_ sdk.Msg       = &MsgAddOracleDeposit{}
	_ CrossChainMsg = &MsgAddOracleDeposit{}
	_ sdk.Msg       = &MsgOracleSetConfirm{}
	_ CrossChainMsg = &MsgOracleSetConfirm{}
	_ sdk.Msg       = &MsgOracleSetUpdatedClaim{}
	_ CrossChainMsg = &MsgOracleSetUpdatedClaim{}

	_ sdk.Msg       = &MsgBridgeTokenClaim{}
	_ CrossChainMsg = &MsgBridgeTokenClaim{}

	_ sdk.Msg       = &MsgSendToFxClaim{}
	_ CrossChainMsg = &MsgSendToFxClaim{}

	_ sdk.Msg       = &MsgSendToExternal{}
	_ CrossChainMsg = &MsgSendToExternal{}
	_ sdk.Msg       = &MsgCancelSendToExternal{}
	_ CrossChainMsg = &MsgCancelSendToExternal{}
	_ sdk.Msg       = &MsgSendToExternalClaim{}
	_ CrossChainMsg = &MsgSendToExternalClaim{}

	_ sdk.Msg       = &MsgRequestBatch{}
	_ CrossChainMsg = &MsgRequestBatch{}
	_ sdk.Msg       = &MsgConfirmBatch{}
	_ CrossChainMsg = &MsgConfirmBatch{}
)

type MsgValidateBasic interface {
	MsgSetOrchestratorAddressValidate(m MsgSetOrchestratorAddress) (err error)
	MsgAddOracleDepositValidate(m MsgAddOracleDeposit) (err error)

	MsgOracleSetConfirmValidate(m MsgOracleSetConfirm) (err error)
	MsgOracleSetUpdatedClaimValidate(m MsgOracleSetUpdatedClaim) (err error)
	MsgBridgeTokenClaimValidate(m MsgBridgeTokenClaim) (err error)
	MsgSendToExternalClaimValidate(m MsgSendToExternalClaim) (err error)

	MsgSendToFxClaimValidate(m MsgSendToFxClaim) (err error)
	MsgSendToExternalValidate(m MsgSendToExternal) (err error)

	MsgCancelSendToExternalValidate(m MsgCancelSendToExternal) (err error)
	MsgRequestBatchValidate(m MsgRequestBatch) (err error)
	MsgConfirmBatchValidate(m MsgConfirmBatch) (err error)
}

var (
	// Denominations can be 3 ~ 128 characters long and support letters, followed by either
	// a letter, a number or a separator ('/').
	reModuleNameString = `[a-zA-Z][a-zA-Z0-9/]{1,32}`
	reModuleName       *regexp.Regexp
)

func init() {
	SetModuleNameRegex(DefaultCoinDenomRegex)
}

// DefaultCoinDenomRegex returns the default regex string
func DefaultCoinDenomRegex() string {
	return reModuleNameString
}

// coinDenomRegex returns the current regex string and can be overwritten for custom validation
var coinDenomRegex = DefaultCoinDenomRegex

// SetModuleNameRegex allows for coin's custom validation by overriding the regular
// expression string used for module name validation.
func SetModuleNameRegex(reFn func() string) {
	coinDenomRegex = reFn
	reModuleName = regexp.MustCompile(fmt.Sprintf(`^%s$`, coinDenomRegex()))
}

// ValidateModuleName is the default validation function for crosschain moduleName.
func ValidateModuleName(moduleName string) error {
	if !reModuleName.MatchString(moduleName) {
		return fmt.Errorf("invalid module name: %s", moduleName)
	}
	return nil
}

var msgValidatorBasicRouter map[string]MsgValidateBasic

func InitMsgValidatorBasicRouter() {
	msgValidatorBasicRouter = make(map[string]MsgValidateBasic)
}

func RegisterValidatorBasic(chainName string, validate MsgValidateBasic) {
	if err := ValidateModuleName(chainName); err != nil {
		panic(sdkerrors.Wrap(ErrInvalidChainName, chainName))
	}
	if _, ok := msgValidatorBasicRouter[chainName]; ok {
		panic(fmt.Sprintf("duplicate registry msg validateBasic! chainName:%s", chainName))
	}
	msgValidatorBasicRouter[chainName] = validate
}

// MsgSetOrchestratorAddress

func (m MsgSetOrchestratorAddress) Route() string {
	return RouterKey
}

func (m MsgSetOrchestratorAddress) Type() string {
	return TypeMsgSetOrchestratorAddress
}

func (m MsgSetOrchestratorAddress) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalidChainName, m.ChainName)
	}
	if router, ok := msgValidatorBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", m.ChainName))
	} else {
		return router.MsgSetOrchestratorAddressValidate(m)
	}
}

func (m MsgSetOrchestratorAddress) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgSetOrchestratorAddress) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Oracle)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgAddOracleDeposit

func (m *MsgAddOracleDeposit) Route() string {
	return RouterKey
}

func (m *MsgAddOracleDeposit) Type() string {
	return TypeMsgAddOracleDeposit
}

func (m MsgAddOracleDeposit) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalidChainName, m.ChainName)
	}
	if router, ok := msgValidatorBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", m.ChainName))
	} else {
		return router.MsgAddOracleDepositValidate(m)
	}
}

func (m *MsgAddOracleDeposit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgAddOracleDeposit) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Oracle)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgOracleSetConfirm

// Route should return the name of the module
func (m MsgOracleSetConfirm) Route() string {
	return RouterKey
}

// Type should return the action
func (m MsgOracleSetConfirm) Type() string {
	return TypeMsgOracleSetConfirm
}

// ValidateBasic performs stateless checks
func (m MsgOracleSetConfirm) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalidChainName, m.ChainName)
	}
	if router, ok := msgValidatorBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", m.ChainName))
	} else {
		return router.MsgOracleSetConfirmValidate(m)
	}
}

// GetSignBytes encodes the message for signing
func (m MsgOracleSetConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgOracleSetConfirm) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.OrchestratorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgSendToExternal

// Route should return the name of the module
func (m MsgSendToExternal) Route() string {
	return RouterKey
}

// Type should return the action
func (m MsgSendToExternal) Type() string {
	return TypeMsgSendToExternal
}

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (m MsgSendToExternal) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalidChainName, m.ChainName)
	}
	if router, ok := msgValidatorBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", m.ChainName))
	} else {
		return router.MsgSendToExternalValidate(m)
	}
}

// GetSignBytes encodes the message for signing
func (m MsgSendToExternal) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgSendToExternal) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgRequestBatch

// Route should return the name of the module
func (m MsgRequestBatch) Route() string {
	return RouterKey
}

// Type should return the action
func (m MsgRequestBatch) Type() string {
	return TypeMsgRequestBatch
}

// ValidateBasic performs stateless checks
func (m MsgRequestBatch) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalidChainName, m.ChainName)
	}
	if router, ok := msgValidatorBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", m.ChainName))
	} else {
		return router.MsgRequestBatchValidate(m)
	}
}

// GetSignBytes encodes the message for signing
func (m MsgRequestBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgRequestBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgConfirmBatch

// Route should return the name of the module
func (m MsgConfirmBatch) Route() string { return RouterKey }

// Type should return the action
func (m MsgConfirmBatch) Type() string { return TypeMsgConfirmBatch }

// ValidateBasic performs stateless checks
func (m MsgConfirmBatch) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalidChainName, m.ChainName)
	}
	if router, ok := msgValidatorBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", m.ChainName))
	} else {
		return router.MsgConfirmBatchValidate(m)
	}
}

// GetSignBytes encodes the message for signing
func (m MsgConfirmBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgConfirmBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.OrchestratorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgCancelSendToExternal

// Route should return the name of the module
func (m MsgCancelSendToExternal) Route() string {
	return RouterKey
}

// Type should return the action
func (m MsgCancelSendToExternal) Type() string {
	return TypeMsgCancelSendToExternal
}

// ValidateBasic performs stateless checks
func (m MsgCancelSendToExternal) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalidChainName, m.ChainName)
	}
	if router, ok := msgValidatorBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", m.ChainName))
	} else {
		return router.MsgCancelSendToExternalValidate(m)
	}
}

// GetSignBytes encodes the message for signing
func (m MsgCancelSendToExternal) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgCancelSendToExternal) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// ExternalClaim represents a claim on ethereum state
type ExternalClaim interface {
	// GetEventNonce All Ethereum claims that we relay from the Gravity contract and into the module
	// have a nonce that is monotonically increasing and unique, since this nonce is
	// issued by the Ethereum contract it is immutable and must be agreed on by all validators
	// any disagreement on what claim goes to what nonce means someone is lying.
	GetEventNonce() uint64
	// GetBlockHeight The block height that the claimed event occurred on. This EventNonce provides sufficient
	// ordering for the execution of all claims. The block height is used only for batchTimeouts + logicTimeouts
	// when we go to create a new batch we set the timeout some number of batches out from the last
	// known height plus projected block progress since then.
	GetBlockHeight() uint64
	// GetClaimer the delegate address of the claimer, for MsgSendToExternalClaim and MsgSendToFxClaim
	// this is sent in as the sdk.AccAddress of the delegated key. it is up to the user
	// to disambiguate this into a sdk.ValAddress
	GetClaimer() sdk.AccAddress
	// GetType Which type of claim this is
	GetType() ClaimType
	ValidateBasic() error
	ClaimHash() []byte
}

var (
	_ ExternalClaim = &MsgSendToFxClaim{}
	_ ExternalClaim = &MsgBridgeTokenClaim{}
	_ ExternalClaim = &MsgSendToExternalClaim{}
	_ ExternalClaim = &MsgOracleSetUpdatedClaim{}
)

// MsgSendToFxClaim

// GetType returns the type of the claim
func (m MsgSendToFxClaim) GetType() ClaimType {
	return CLAIM_TYPE_SEND_TO_FX
}

// ValidateBasic performs stateless checks
func (m MsgSendToFxClaim) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalidChainName, m.ChainName)
	}
	if router, ok := msgValidatorBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", m.ChainName))
	} else {
		return router.MsgSendToFxClaimValidate(m)
	}
}

// GetSignBytes encodes the message for signing
func (m MsgSendToFxClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgSendToFxClaim) GetClaimer() sdk.AccAddress {
	err := m.ValidateBasic()
	if err != nil {
		panic("MsgSendToFxClaim failed ValidateBasic! Should have been handled earlier")
	}
	val, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}
	return val
}

// GetSigners defines whose signature is required
func (m MsgSendToFxClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// Type should return the action
func (m MsgSendToFxClaim) Type() string {
	return TypeMsgSendToFxClaim
}

// Route should return the name of the module
func (m MsgSendToFxClaim) Route() string {
	return RouterKey
}

// ClaimHash Hash implements BridgeSendToExternal.Hash
func (m MsgSendToFxClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d%s/%s/%s/%s/%s", m.BlockHeight, m.EventNonce, m.TokenContract, m.Sender, m.Amount.String(), m.Receiver, m.TargetIbc)
	return tmhash.Sum([]byte(path))
}

// MsgSendToExternalClaim

// GetType returns the claim type
func (m MsgSendToExternalClaim) GetType() ClaimType {
	return CLAIM_TYPE_SEND_TO_EXTERNAL
}

// ValidateBasic performs stateless checks
func (m MsgSendToExternalClaim) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalidChainName, m.ChainName)
	}
	if router, ok := msgValidatorBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", m.ChainName))
	} else {
		return router.MsgSendToExternalClaimValidate(m)
	}
}

// ClaimHash Hash implements SendToFxBatch.Hash
func (m MsgSendToExternalClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%s/%d/", m.BlockHeight, m.EventNonce, m.TokenContract, m.BatchNonce)
	return tmhash.Sum([]byte(path))
}

// GetSignBytes encodes the message for signing
func (m MsgSendToExternalClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgSendToExternalClaim) GetClaimer() sdk.AccAddress {
	err := m.ValidateBasic()
	if err != nil {
		panic("MsgSendToExternalClaim failed ValidateBasic! Should have been handled earlier")
	}
	val, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(fmt.Sprintf("invalid address %s", m.Orchestrator))
	}
	return val
}

// GetSigners defines whose signature is required
func (m MsgSendToExternalClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (m MsgSendToExternalClaim) Route() string {
	return RouterKey
}

// Type should return the action
func (m MsgSendToExternalClaim) Type() string {
	return TypeMsgSendToExternalClaim
}

// MsgBridgeTokenClaim

func (m MsgBridgeTokenClaim) Route() string {
	return RouterKey
}

func (m MsgBridgeTokenClaim) Type() string {
	return TypeMsgBridgeTokenClaim
}

func (m MsgBridgeTokenClaim) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalidChainName, m.ChainName)
	}
	if router, ok := msgValidatorBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", m.ChainName))
	} else {
		return router.MsgBridgeTokenClaimValidate(m)
	}
}

func (m MsgBridgeTokenClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgBridgeTokenClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (m MsgBridgeTokenClaim) GetClaimer() sdk.AccAddress {
	err := m.ValidateBasic()
	if err != nil {
		panic("MsgBridgeTokenClaim failed ValidateBasic! Should have been handled earlier")
	}
	val, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(fmt.Sprintf("invalid address %s", m.Orchestrator))
	}
	return val
}

func (m MsgBridgeTokenClaim) GetType() ClaimType {
	return CLAIM_TYPE_BRIDGE_TOKEN
}

func (m MsgBridgeTokenClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d%s/%s/%s/%d/%s/", m.BlockHeight, m.EventNonce, m.TokenContract, m.Name, m.Symbol, m.Decimals, m.ChannelIbc)
	return tmhash.Sum([]byte(path))
}

// MsgOracleSetUpdatedClaim

// GetType returns the type of the claim
func (m MsgOracleSetUpdatedClaim) GetType() ClaimType {
	return CLAIM_TYPE_ORACLE_SET_UPDATED
}

// ValidateBasic performs stateless checks
func (m MsgOracleSetUpdatedClaim) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalidChainName, m.ChainName)
	}
	if router, ok := msgValidatorBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", m.ChainName))
	} else {
		return router.MsgOracleSetUpdatedClaimValidate(m)
	}
}

// GetSignBytes encodes the message for signing
func (m MsgOracleSetUpdatedClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgOracleSetUpdatedClaim) GetClaimer() sdk.AccAddress {
	err := m.ValidateBasic()
	if err != nil {
		panic("MsgOracleSetUpdatedClaim failed ValidateBasic! Should have been handled earlier")
	}
	val, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(fmt.Sprintf("invalid address %s", m.Orchestrator))
	}
	return val
}

// GetSigners defines whose signature is required
func (m MsgOracleSetUpdatedClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// Type should return the action
func (m MsgOracleSetUpdatedClaim) Type() string {
	return TypeMsgOracleSetUpdatedClaim
}

// Route should return the name of the module
func (m MsgOracleSetUpdatedClaim) Route() string {
	return RouterKey
}

// ClaimHash Hash implements BridgeSendToExternal.Hash
func (m MsgOracleSetUpdatedClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%d/%s/", m.BlockHeight, m.OracleSetNonce, m.EventNonce, m.Members)
	return tmhash.Sum([]byte(path))
}
