package account

import (
	"bufio"
	"fmt"
	"os"

	b "github.com/illaykaye/Crypto/Block"
)

var (
	enterUsername   string = "Enter username: "
	enterPassword   string = "Enter password: "
	enterBalance    string = "Enter balance: "
	accountNotFound string = "Account not found."
)

type Accounts struct {
	Bc             b.BlockChain
	file           string
	UserLogged     AccountS
	bothFilesEmt   bool
	numLinesInFile int
}

func NewAccounts(acf string, bcf string) Accounts {
	return Accounts{
		b.NewBC(bcf),
		acf,
		AccountS{},
		false,
		NumLinesInFile(acf),
	}
}

func (acs *Accounts) GetUser() AccountS {
	return acs.UserLogged
}

func (acs *Accounts) openFile() *os.File {
	f, err := os.OpenFile(acs.file, os.O_RDWR, 0644)
	isErr(err)

	return f
}

func (acs *Accounts) UpdateFile(userToUpdate ...AccountS) {
	var (
		user          AccountS = AccountS{}
		lineNum       int      = 0
		usersUsername []string
	)

	for i := 0; i < len(userToUpdate); i++ {
		usersUsername = append(usersUsername, userToUpdate[i].Username)
	}

	file, err := os.OpenFile(acs.file, os.O_RDWR, 0644)
	isErr(err)

	reader := bufio.NewReader(file)

	fileToUpdate, err := os.OpenFile(fmt.Sprintf("%s1", acs.file), os.O_RDWR|os.O_CREATE, 0644)
	isErr(err)

	for lineNum < acs.numLinesInFile {
		accountLine, _, err := reader.ReadLine()
		isErr(err)
		user = DecodeAccount(accountLine)

		indexInslice := findInSlice(usersUsername, user.Username)
		if indexInslice >= 0 {
			fileToUpdate.Write(EncodeAccount(userToUpdate[indexInslice]))
			fileToUpdate.WriteString("\n")
		} else {
			fileToUpdate.Write(accountLine)
			fileToUpdate.WriteString("\n")
		}
		lineNum++
	}

	os.Remove(file.Name())
	os.Rename(fileToUpdate.Name(), file.Name())

	file.Close()
	fileToUpdate.Close()
}

func NumLinesInFile(file string) int {
	f, err := os.OpenFile(file, os.O_RDONLY, 0644)
	isErr(err)

	fileScanner := bufio.NewScanner(f)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	f.Close()
	return lineCount
}
