package skip

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	ibctypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/tessellated-io/mail-in-rebates/paymaster/crypto"
	"github.com/tessellated-io/pickaxe/chains"

	"github.com/cosmos/cosmos-sdk/codec"
	cdc "github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type SkipClient struct {
	registry *chains.OfflineChainRegistry
	cdc      *codec.ProtoCodec
}

func NewSkipClient(
	registry *chains.OfflineChainRegistry,
	cdc *codec.ProtoCodec,
) *SkipClient {
	return &SkipClient{
		registry: registry,
		cdc:      cdc,
	}
}

// Send funds without discretion. Sanity checks and safety meausres are not enforced in this call.
func (sc *SkipClient) GetMessages(
	senderAddress string,
	senderPublicKey cryptotypes.PubKey,
	recipientAddress string,
	amountIn string,
	denomIn string,
	sourceChainID string,
	destDenom string,
	destChainID string,
) (sdk.Msg, error) {
	// TODO: this hsould also be refactored
	if sourceChainID == destChainID && denomIn == destDenom {
		coinAmount, ok := sdk.NewIntFromString(amountIn)
		if !ok {
			return nil, fmt.Errorf("could not parse amount \"%s\" to an int", amountIn)
		}
		amount := sdk.NewCoin(denomIn, coinAmount)

		bankSendMsg := &banktypes.MsgSend{
			FromAddress: senderAddress,
			ToAddress:   recipientAddress,
			Amount:      []sdk.Coin{amount},
		}
		return bankSendMsg, nil
	}

	// Get the route
	routeResponse, err := sc.getRoute(amountIn, denomIn, sourceChainID, destDenom, destChainID)
	if err != nil {
		return nil, err
	}

	// Get messages
	messagesResponse, err := sc.getMessages(
		senderAddress,
		senderPublicKey,
		recipientAddress,
		amountIn,
		denomIn,
		sourceChainID,
		destDenom,
		destChainID,
		routeResponse.Operations,
		*routeResponse.EstimatedAmountOut,
		routeResponse.ChainIds,
	)
	if err != nil {
		return nil, err
	}

	if len(messagesResponse.Msgs) > 1 {
		return nil, fmt.Errorf("cannot sign multiple messages")
	}
	msgType := messagesResponse.Msgs[0].MsgTypeUrl
	msgJson := messagesResponse.Msgs[0].Msg

	if msgType == "/ibc.applications.transfer.v1.MsgTransfer" {
		msg := ibctypes.MsgTransfer{}
		cdc.JSONCodec.UnmarshalJSON(&cdc.ProtoCodec{}, []byte(msgJson), &msg)
		return &msg, nil
	}

	return nil, fmt.Errorf("unexpected or unsupported message type %s", msgType)
}

func (sc *SkipClient) getMessages(
	senderAddress string,
	senderPublicKey cryptotypes.PubKey,
	recipientAddress string,
	amountIn string,
	denomIn string,
	sourceChainID string,
	destDenom string,
	destChainID string,
	operations []*Operation,
	estimatedAmountOut string,
	chainIDs []string,
) (*MessagesResponse, error) {
	url := "https://api.skip.money/v1/fungible/msgs"

	// Encode operations to JSON
	jsonOperations := ""
	for opIdx, operation := range operations {
		jsonOperation, err := json.Marshal(operation)
		if err != nil {
			return nil, err
		}

		jsonOperations = fmt.Sprintf("%s%s", jsonOperations, jsonOperation)
		if opIdx != len(operations)-1 {
			jsonOperations = fmt.Sprintf("%s, ", jsonOperations)
		}
	}

	// Set up chain IDS
	chainIDMapJSON := ""
	for chainIDIdx, chainID := range chainIDs {
		prefix := sc.registry.ChainIDToData[chainID].AccountPrefix
		address, err := crypto.PubKeyToAddress(senderPublicKey, prefix)
		if err != nil {
			return nil, err
		}

		// If last address, change to the recipient address.
		if strings.HasPrefix(recipientAddress, prefix) {
			address = recipientAddress
		}

		chainIDAndAddress := fmt.Sprintf("\"%s\": \"%s\"", chainID, address)

		chainIDMapJSON = fmt.Sprintf("%s%s", chainIDMapJSON, chainIDAndAddress)
		if chainIDIdx != len(chainIDs)-1 {
			chainIDMapJSON = fmt.Sprintf("%s, ", chainIDMapJSON)
		}
	}

	slippageTolerancePercent := "1"
	payload := fmt.Sprintf(`
	{
		"amount_in": "%s",
		"source_asset_denom": "%s",
		"source_asset_chain_id": "%s",
		"dest_asset_denom": "%s",
		"dest_asset_chain_id": "%s",
		"operations": [
			%s
		],   
		"estimated_amount_out": "%s", 		
		"slippage_tolerance_percent": "%s", 
		"chain_ids_to_addresses": {
			%s
		}
	}`,
		amountIn,
		denomIn,
		sourceChainID,
		destDenom,
		destChainID,
		jsonOperations,
		estimatedAmountOut,
		slippageTolerancePercent,
		chainIDMapJSON,
	)
	fmt.Printf("Payload to messages: %s\n", payload)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))

	messagesResponse, err := ParseMessagesResponse(string(body))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf(string(body))
	}

	return messagesResponse, nil
}

func (sc *SkipClient) getRoute(
	amountIn string,
	denomIn string,
	sourceChainID string,
	destDenom string,
	destChainID string,
) (*RouteResponse, error) {
	url := "https://api.skip.money/v1/fungible/route"

	payload := fmt.Sprintf(`
		{
			"amount_in": "%s",
			"source_asset_denom": "%s",
			"source_asset_chain_id": "%s",
			"dest_asset_denom": "%s",
			"dest_asset_chain_id": "%s",
			"cumulative_affiliate_fee_bps": "0"
		}`,
		amountIn,
		denomIn,
		sourceChainID,
		destDenom,
		destChainID,
	)
	fmt.Printf("Payload for routes: %s\n", payload)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Response from routes/: %s\n", string(body))

	routeResponse, err := ParseRouteResponse(string(body))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf(string(body))
	}

	return routeResponse, nil
}
