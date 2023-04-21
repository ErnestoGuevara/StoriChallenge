package model

//Generate the struct for number of transaction per month
type Transaction struct {
	Month   string `json:"month"`
	NumTran int    `json:"numTran"`
}
