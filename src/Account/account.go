package account

import (
	"encoding/json"
	"fmt"
)

type AccountS struct {
	Username   string  `json:"Username"`
	Password   string  `json:"Password"`
	AccountNum string  `json:"AccountNum"`
	Balance    float64 `json:"Balance"`
}

func EncodeAccount(account AccountS) []byte {
	accountJSON, err := json.Marshal(account)
	if err != nil {
		panic(err)
	}
	return accountJSON
}

func DecodeAccount(accountJSON []byte) AccountS {
	var account AccountS
	err := json.Unmarshal([]byte(accountJSON), &account)
	if err != nil {
		panic(err)
	}
	return account
}

func (ac AccountS) Data() {
	fmt.Println("Username:", ac.Username)
	fmt.Println("Account number:", ac.AccountNum)
	fmt.Println("Balance:", ac.Balance)
}

func isErr(err error) {
	if err != nil {
		panic(err)
	}
}
