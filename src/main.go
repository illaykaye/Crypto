package main

import (
	"fmt"

	a "github.com/illaykaye/Crypto/Account"
)

var (
	//path string = "accounts"
	acs a.Accounts
)

var (
	commandNotFound string = "Command not found."
)

func main() {
	var (
		command1 string
		running  bool   = true
		bcFile   string = "../files/blockchain"
		acsFile  string = "../files/accounts"
	)

	acs = a.NewAccounts(acsFile, bcFile)

	if acs.Bc.GetNumBlocks() == 0 {
		acs.Bc.Genesis()
	}

	for running {
		fmt.Print(">>")
		fmt.Scanln(&command1)

		switch command1 {
		case "sign":
			acs.SignUp()

		case "login":
			if acs.LogIn() {
				loggedIn()
			}

		case "quit":
			running = false

		default:
			fmt.Println(commandNotFound)

		}
		command1 = ""
	}
}

func loggedIn() {
	var (
		command string
		running bool   = true
		v1      string = "("
		v2      string = ")"
	)

	for running {
		fmt.Printf("%v%v%v %v", v1, acs.GetUser().Username, v2, ">>")
		fmt.Scanln(&command)

		switch command {
		case "data":
			acs.GetUser().Data()

		case "transfer":
			acs.BankTransfer()

		case "logout":
			fmt.Println("Logged out.")
			acs.UserLogged = a.AccountS{}
			running = false

		default:
			fmt.Println(commandNotFound)

		}
		command = ""
	}
}

/*
func isErr(err error) {
	if err != nil {
		panic(err)
	}
}
*/
