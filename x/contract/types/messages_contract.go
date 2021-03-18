package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateContract{}

func NewMsgCreateContract(creator string, contractNo string, partyA string, partyB string, action string, data string) *MsgCreateContract {
	return &MsgCreateContract{
		Creator:    creator,
		ContractNo: contractNo,
		PartyA:     partyA,
		PartyB:     partyB,
		Action:     action,
		Data:       data,
	}
}

func (msg *MsgCreateContract) Route() string {
	return RouterKey
}

func (msg *MsgCreateContract) Type() string {
	return "CreateContract"
}

func (msg *MsgCreateContract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateContract{}

func NewMsgUpdateContract(creator string, id string, contractNo string, partyA string, partyB string, action string, data string) *MsgUpdateContract {
	return &MsgUpdateContract{
		Id:         id,
		Creator:    creator,
		ContractNo: contractNo,
		PartyA:     partyA,
		PartyB:     partyB,
		Action:     action,
		Data:       data,
	}
}

func (msg *MsgUpdateContract) Route() string {
	return RouterKey
}

func (msg *MsgUpdateContract) Type() string {
	return "UpdateContract"
}

func (msg *MsgUpdateContract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateContract{}

func NewMsgDeleteContract(creator string, id string) *MsgDeleteContract {
	return &MsgDeleteContract{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeleteContract) Route() string {
	return RouterKey
}

func (msg *MsgDeleteContract) Type() string {
	return "DeleteContract"
}

func (msg *MsgDeleteContract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
