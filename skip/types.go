package skip

import "encoding/json"

type RouteResponse struct {
	SourceAssetDenom   string       `json:"source_asset_denom"`
	SourceAssetChainId string       `json:"source_asset_chain_id"`
	DestAssetDenom     string       `json:"dest_asset_denom"`
	DestAssetChainId   string       `json:"dest_asset_chain_id"`
	AmountIn           string       `json:"amount_in"`
	Operations         []*Operation `json:"operations"`
	ChainIds           []string     `json:"chain_ids"`
	DoesSwap           bool         `json:"does_swap"`
	EstimatedAmountOut *string      `json:"estimated_amount_out"`
}

type Operation struct {
	Swap     *SwapOperation     `json:"swap,omitempty"`
	Transfer *TransferOperation `json:"transfer,omitempty"`
}

type SwapOperation struct {
	SwapIn           *SwapInOperation  `json:"swap_in"`
	SwapExactCoinOut *SwapOutOperation `json:"swap_out"`
}

type SwapInOperation struct {
	SwapVenue      *SwapVenue           `json:"swap_venue"`
	SwapOperations []*PoolSwapOperation `json:"swap_operations"`
	SwapAmountIn   string               `json:"swap_amount_in"`
}

type SwapOutOperation struct {
	SwapVenue      *SwapVenue           `json:"swap_venue"`
	SwapOperations []*PoolSwapOperation `json:"swap_operations"`
	SwapAmountIn   string               `json:"swap_amount_in"`
}

type PoolSwapOperation struct {
	Pool     string `json:"pool"`
	DenomIn  string `json:"denom_in"`
	DenomOut string `json:"denom_out"`
}

type SwapVenue struct {
	Name    string `json:"name"`
	ChainId string `json:"chain_id"`
}

type TransferOperation struct {
	Port         string `json:"port"`
	Channel      string `json:"channel"`
	ChainId      string `json:"chain_id"`
	PfmEnabled   bool   `json:"pfm_enabled"`
	DestDenom    string `json:"dest_denom"`
	SupportsMemo bool   `json:"supports_memo"`
}

func ParseRouteResponse(jsonString string) (*RouteResponse, error) {
	data := &RouteResponse{}

	// Unmarshal the JSON string into the data struct
	err := json.Unmarshal([]byte(jsonString), data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type MessagesResponse struct {
	Msgs []*Transaction `json:"msgs"`
}

type Transaction struct {
	ChainId    string   `json:"chain_id"`
	Path       []string `json:"path"`
	Msg        string   `json:"msg"`
	MsgTypeUrl string   `json:"msg_type_url"`
}

func ParseMessagesResponse(jsonString string) (*MessagesResponse, error) {
	data := &MessagesResponse{}

	// Unmarshal the JSON string into the data struct
	err := json.Unmarshal([]byte(jsonString), data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
