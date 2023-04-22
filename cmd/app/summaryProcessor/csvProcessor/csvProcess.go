package csvProcessor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/database"
	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/model"

	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/logger"
)

func ProcessData(transactions [][]string, testFile string) ([]float64, []float64, float64, map[string]int, error) {

	// create a new database object
	db, err := database.NewDB()
	if err != nil {
		logger := logger.NewLogger("DB_ERROR: ")
		logger.Error(fmt.Sprintf("Error initialazing database: %s", err.Error()))
	}

	// Process the transactions
	var balance float64
	var credits, debits []float64
	monthlyCredits := make(map[string][]float64)
	monthlyDebits := make(map[string][]float64)
	monthlyTransactions := make(map[string]int)
	monthMap := map[int]string{
		1:  "January",
		2:  "February",
		3:  "March",
		4:  "April",
		5:  "May",
		6:  "June",
		7:  "July",
		8:  "August",
		9:  "September",
		10: "October",
		11: "November",
		12: "December",
	}
	var transactionList []model.Transactions
	for _, row := range transactions[1:] {
		//Set idFile variable with row[0] value also with Atoi convert String to int
		idFile, err := strconv.Atoi(row[0])

		if err != nil {
			logger := logger.NewLogger("VALUE_ERROR: ")
			logger.Error(fmt.Sprintf("Error converting from string to int: %v", err))
		}

		//Set date variable with row[1] value
		date := row[1]
		//Set amount variable with row[2] value setting as float and without "+" symbol
		amount, err := strconv.ParseFloat(strings.Trim(row[2], " "), 64)
		if err != nil {
			logger := logger.NewLogger("VALUE_ERROR: ")
			logger.Error(fmt.Sprintf("Error converting from string to float: %v", err))
		}

		// Create a new Transactions object and append it to the transactionList
		transactionList = append(transactionList, model.Transactions{
			Id:    idFile,
			Date:  date,
			Value: amount,
		})

		month, err := strconv.Atoi(strings.Split(date, "/")[0])
		if err != nil {
			logger := logger.NewLogger("VALUE_ERROR: ")
			logger.Error(fmt.Sprintf("Error converting from string to int: %v", err))
		}
		monthName := monthMap[month]

		if amount > 0 {

			credits = append(credits, amount)
			monthlyCredits[monthName] = append(monthlyCredits[monthName], amount)
		} else {

			debits = append(debits, amount)
			monthlyDebits[monthName] = append(monthlyDebits[monthName], amount)
		}

		balance += amount
		monthlyTransactions[monthName]++

		// insert a transaction into the database
		err = db.InsertTransaction(testFile, idFile, amount, date)
		if err != nil {
			logger := logger.NewLogger("DB_ERROR: ")
			logger.Error(fmt.Sprintf("Error inserting the row: %s", err.Error()))
		}

	}
	return debits, credits, balance, monthlyTransactions, nil
}
