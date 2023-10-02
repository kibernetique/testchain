package models

type TxInfo struct {
	Hash        string `json:"hash"`
	IsPending   bool   `json:"isPending"`
	ChainId     string `json:"chainId"`
	Cost        string `json:"cost"`
	To          string `json:"to"`
	Sender      string `json:"sender"`
	Value       string `json:"value"`
	Data        string `json:"data"`
	BlockNumber string `json:"blockNuber"`
	GasUsed     uint64 `json:"gasUsed"`
}
