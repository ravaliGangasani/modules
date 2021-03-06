package nfts

import (
	"strings"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (result *sdk.Result, err error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		
		switch msg := msg.(type) {
		case MsgMintTweetNFT:
			return handleMsgMintTweetNFT(ctx, keeper, msg)
		
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized NFT message type: %T", msg)
		}
	}
}

func handleMsgMintTweetNFT(ctx sdk.Context, keeper Keeper, msg MsgMintTweetNFT) (*sdk.Result, error) {
	
	nfts := keeper.GetTweetsOfAccount(ctx, msg.Sender)
	
	for _, nft := range nfts {
		if strings.EqualFold(nft.AssetID, msg.AssetID) {
			return nil, sdkerrors.Wrap(ErrAssetIDAlreadyExist, "")
		}
	}
	
	count := keeper.GetGlobalTweetCount(ctx)
	id := GetPrimaryNFTID(count)
	tweetNFT := BaseTweetNFT{
		PrimaryNFTID:   id,
		PrimaryOwner:   msg.Sender.String(),
		SecondaryNFTID: "",
		SecondaryOwner: "",
		License:        msg.License,
		AssetID:        msg.AssetID,
		LicensingFee:   msg.LicensingFee,
		RevenueShare:   msg.RevenueShare,
		TwitterHandle:  msg.TwitterHandle,
	}
	
	keeper.MintTweetNFT(ctx, tweetNFT)
	keeper.SetTweetIDToAccount(ctx, msg.Sender, tweetNFT.PrimaryNFTID)
	keeper.SetGlobalTweetCount(ctx, count+1)
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeMsgMintTweetNFT,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(AttributePrimaryNFTID, tweetNFT.PrimaryNFTID),
			sdk.NewAttribute(AttributeAssetID, tweetNFT.AssetID),
			sdk.NewAttribute(AttributeTwitterHandle, tweetNFT.TwitterHandle),
		),
	)
	
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
	
}
