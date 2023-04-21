package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

//Generate the struct of Transactions
type Transactions struct {
	Id    int
	Date  string
	Value float64
}

func main() {
	connectDB()
	csvFile("client1.csv")

}
func connectDB() {
	//Use sql.Open to initialize a new sql.DB object
	//Pass the driver name and the connection string
	db, err := sql.Open("mysql", "admin:neto120899@tcp(database-storichallenge.ccmv5yteur66.us-east-1.rds.amazonaws.com:3306)/stori_challenge")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//Call db.Ping() to check the connection
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS stori_transactions (id BIGINT PRIMARY KEY AUTO_INCREMENT,file VARCHAR(100), idFile INT, transaction VARCHAR(100), date VARCHAR(100));")
	if err != nil {
		log.Fatal(err)
	}

}

func csvFile(testFile string) {

	file, err := os.Open("./test/" + testFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// Parse the CSV file
	reader := csv.NewReader(file)
	transactions, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	//Open MySQL connection to insert data
	db, err := sql.Open("mysql", "admin:neto120899@tcp(database-storichallenge.ccmv5yteur66.us-east-1.rds.amazonaws.com:3306)/stori_challenge")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
	var transactionsMonthly string
	for _, row := range transactions[1:] {
		//Set idFile variable with row[0] value also with Atoi convert String to int
		idFile, err := strconv.Atoi(row[0])

		if err != nil {
			log.Fatal(err)
		}

		//Set date variable with row[1] value
		date := row[1]
		//Set amount variable with row[2] value setting as float and without "+" symbol
		amount, err := strconv.ParseFloat(strings.Trim(row[2], " "), 64)

		if err != nil {
			log.Fatal(err)
		}
		month, err := strconv.Atoi(strings.Split(date, "/")[0])
		if err != nil {
			log.Fatal(err)
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
		// Check if the row already exists in the table
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM stori_transactions WHERE file = ? AND idFile = ?", testFile, idFile).Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
		// Insert the row if it doesn't already exist in the table
		if count == 0 {
			_, err = db.Exec("INSERT INTO stori_transactions(file,idFile,transaction,date) VALUES (?,?, ?, ?)", testFile, idFile, amount, date)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Generate summary report
	averageCredit := calculateAverage(credits)
	averageDebit := calculateAverage(debits)

	summary := fmt.Sprintf("\nTotal balance is %.2f\nAverage debit amount: %.2f\nAverage credit amount: %.2f\n",
		balance, averageDebit, averageCredit)

	//Print Monthly transactions
	for monthName, numTransactions := range monthlyTransactions {
		transactionsMonthly += fmt.Sprintf("Number of transaction in %s: %d \n", monthName, numTransactions)
	}
	balanceStr := fmt.Sprintf("%.2f", balance)
	averageCreditStr := fmt.Sprintf("%.2f", averageCredit)
	averageDebitStr := fmt.Sprintf("%.2f", averageDebit)
	fmt.Println(summary, transactionsMonthly)

	sendEmail(balanceStr, averageCreditStr, averageDebitStr, transactionsMonthly)

}
func sendEmail(balance string, avergaCredit string, averageDebit string, transactionsMonthly string) error {
	from := mail.NewEmail("Stori Challenge", "neto_1208@hotmail.com")
	subject := "Summary Report"
	to := mail.NewEmail("Client", "neto120899@hotmail.com")

	message := mail.NewV3MailInit(from, subject, to)
	message.SetTemplateID("126b3b3d-eb7e-4b20-8d2b-bd1b2f9f712e")
	message.Personalizations[0].SetSubstitution("-balance-", balance)
	message.Personalizations[0].SetSubstitution("-debitavg-", averageDebit)
	message.Personalizations[0].SetSubstitution("-creditavg-", avergaCredit)
	message.Personalizations[0].SetSubstitution("-monthly-", transactionsMonthly)

	client := sendgrid.NewSendClient("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")
	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}

func calculateAverage(numbers []float64) float64 {
	var sum float64
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}
