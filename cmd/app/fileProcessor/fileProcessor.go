package fieProcessor

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/database"
	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/emailSender"
	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/logger"
)

func CsvFile(testFile string) {
	// create a new database object
	db, err := database.NewDB()
	if err != nil {
		logger := logger.NewLogger("DB_ERROR: ")
		logger.Error(fmt.Sprintf("Error initialazing database: %s", err.Error()))
	}

	file, err := os.Open("./test/" + testFile)
	if err != nil {
		logger := logger.NewLogger("FILE_ERROR: ")
		logger.Error(fmt.Sprintf("Error opening csv file: %v", err))
	}
	defer file.Close()

	// Parse the CSV file
	reader := csv.NewReader(file)
	transactions, err := reader.ReadAll()
	if err != nil {
		logger := logger.NewLogger("FILE_ERROR: ")
		logger.Error(fmt.Sprintf("Error parsing csv file: %v", err))
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

	// Generate summary report
	averageCredit := calculateAverage(credits)
	averageDebit := calculateAverage(debits)

	summary := fmt.Sprintf("\nTotal balance is %.2f\nAverage debit amount: %.2f\nAverage credit amount: %.2f\n",
		balance, averageDebit, averageCredit)

	//Print Monthly transactions
	for monthName, numTransactions := range monthlyTransactions {
		summary += fmt.Sprintf("Number of transaction in %s: %d \n", monthName, numTransactions)
	}
	balanceStr := fmt.Sprintf("%.2f", balance)
	averageCreditStr := fmt.Sprintf("%.2f", averageCredit)
	averageDebitStr := fmt.Sprintf("%.2f", averageDebit)
	fmt.Println(summary)

	err = emailSender.SendEmail(balanceStr, averageCreditStr, averageDebitStr, monthlyTransactions)
	if err != nil {
		logger := logger.NewLogger("EMAIL_ERROR: ")
		logger.Error(fmt.Sprintf("Error sending email: %v", err))
	}

}

func calculateAverage(numbers []float64) float64 {
	var sum float64
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}
