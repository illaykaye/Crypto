package account

import (
	"bufio"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (acs *Accounts) SignUp() {
	var (
		username       string
		password       string
		balance        float64
		accountNum     string
		isUsernameFree bool = true
		lineNum        int  = 0
	)

	fmt.Print(enterUsername)
	fmt.Scanln(&username)
	if username == "0" {
		return
	}

	f := acs.openFile()
	reader := bufio.NewReader(f)

	for lineNum < acs.numLinesInFile && isUsernameFree {
		accountLine, _, err := reader.ReadLine()
		isErr(err)
		if username == DecodeAccount(accountLine).Username {
			isUsernameFree = false
			fmt.Println("This username is taken.")
			return
		}
		lineNum++
	}
	fmt.Print(enterPassword)
	fmt.Scanln(&password)
	if len(password) < 8 {
		fmt.Println("Password need to be more than 8 charcters.")
		return
	}
	fmt.Print(enterBalance)
	fmt.Scanln(&balance)

	accountNum = acs.AccountNumGen()
	user := AccountS{Username: username, Password: password, AccountNum: accountNum, Balance: balance}

	accountJSON := EncodeAccount(user)

	f.Write(accountJSON)
	f.WriteString("\n")
	acs.numLinesInFile++
}

//Creates a unique account number////////////////////////////
func (acs *Accounts) AccountNumGen() string {
	var (
		accountNum               = strings.Builder{}
		allAccountaNums          = []string{}
		user            AccountS = AccountS{}
	)

	//scanner := bufio.NewScanner(acs.file)
	f := acs.openFile()
	reader := bufio.NewReader(f)

	for i := 0; i < acs.numLinesInFile; i++ {
		accountJSON, _, err := reader.ReadLine()
		isErr(err)
		//accountJSON := scanner.Bytes()
		user = DecodeAccount(accountJSON)
		allAccountaNums = append(allAccountaNums, user.AccountNum)
	}

	accountNum = createAccountNum()
	for findInSlice(allAccountaNums, accountNum.String()) >= 0 {
		accountNum = createAccountNum()

	}
	return accountNum.String()
}

func createAccountNum() strings.Builder {
	var accountNum = strings.Builder{}

	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < 8; i++ {
		dig := rand.Intn(10)
		accountNum.WriteString(strconv.Itoa(dig))
	}
	return accountNum
}

func findInSlice(slice []string, val string) int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == val {
			return i
		}
	}
	return -1
}

//////////////////////////////////////////////////////////////

func (acs *Accounts) LogIn() bool {

	if acs.bothFilesEmt {
		fmt.Println("Error: file is empty.")
		return false
	}
	var (
		username string
		password string
		account  AccountS
		lineNum  int = 0
	)

	fmt.Print(enterUsername)
	fmt.Scanln(&username)
	fmt.Print(enterPassword)
	fmt.Scanln(&password)

	f := acs.openFile()
	reader := bufio.NewReader(f)

	for lineNum < acs.numLinesInFile {
		accountJSON, _, err := reader.ReadLine()
		isErr(err)

		account = DecodeAccount(accountJSON)

		if username == account.Username {
			if password == account.Password {
				fmt.Println("Logged in succesfuly.")
				acs.UserLogged = account
				return true

			}
		}
		lineNum++
	}

	fmt.Println("Wrong username or password.")
	return false
}

func (acs *Accounts) BankTransfer() {
	var (
		amount         float64 = 0
		toAccountName  string
		toAccount      AccountS = AccountS{}
		toAccountFound bool     = false
		lineNum        int      = 0
		numLines       int      = acs.numLinesInFile
	)

	fmt.Print("Transfer to: ")
	fmt.Scanln(&toAccountName)

	f := acs.openFile()
	reader := bufio.NewReader(f)

	for lineNum < numLines && !toAccountFound {
		accountJSON, _, err := reader.ReadLine()
		isErr(err)
		toAccount = DecodeAccount(accountJSON)
		if toAccount.Username == toAccountName {
			toAccountFound = true
		}
		lineNum++
	}
	if !toAccountFound {
		fmt.Println(accountNotFound)
		return
	}
	fmt.Print("Enter amount to transfer: ")
	fmt.Scanln(&amount)
	if amount == 0 {
		fmt.Println("You didn't enter amount to transfer.")
		return
	} else if amount < 0 {
		fmt.Println("Not valid amount.")
		return
	} else if amount > acs.UserLogged.Balance {
		fmt.Println("You don't have enough money.")
		return
	}

	acs.UserLogged.Balance -= amount
	toAccount.Balance += amount

	acs.Bc.AddData(acs.UserLogged.Username, toAccount.Username, amount)
	acs.Bc.ConstructBlock()

	acs.UpdateFile(acs.UserLogged, toAccount)
	fmt.Println("Transfer went succesfuly.")

}
