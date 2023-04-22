package summaryProcessor

import (
	"fmt"

	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/emailSender"
	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/logger"
	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/summaryProcessor/csvProcessor"
)

func SummaryReportGenerator(testFile string) {
	transactions, err := csvProcessor.ReadFile(testFile)
	if err != nil {
		// Handle the error
		logger := logger.NewLogger("FILE_ERROR: ")
		logger.Error(fmt.Sprintf("Error opening csv file: %s", err.Error()))
	}

	debits, credits, balance, monthlyTransactions, err := csvProcessor.ProcessData(transactions, testFile)
	if err != nil {
		// Handle the error
		logger := logger.NewLogger("FILE_ERROR: ")
		logger.Error(fmt.Sprintf("Error opening csv file: %s", err.Error()))
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

	//Sending email
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
