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
	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

//Generate the struct of Transactions
type Transactions struct {
	Id    int
	Date  string
	Value float64
}

//Generate the struct for number of transaction per month
type transaction struct {
	Month   string `json:"month"`
	NumTran int    `json:"numTran"`
}

var dbUser string
var dbPassword string
var dbHost string
var dbPort string
var dbName string
var sgApi string
var sgTemplate string
var db *sql.DB
var err error

func main() {

	//Open .env File
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbName = os.Getenv("DB_NAME")
	sgApi = os.Getenv("SG_APIKEY")
	sgTemplate = os.Getenv("SG_TEMPLATEID")

	csvFile("client1.csv")

}
func connectDB() {

	//Use sql.Open to initialize a new sql.DB object
	dataBaseString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err = sql.Open("mysql", dataBaseString)
	if err != nil {
		log.Println(err)
	}
	//Call db.Ping() to check the connection
	pingErr := db.Ping()
	if pingErr != nil {
		log.Println(pingErr)
	}
	fmt.Println("Connected!")

}

func csvFile(testFile string) {

	file, err := os.Open("./test/" + testFile)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	// Parse the CSV file
	reader := csv.NewReader(file)
	transactions, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
	}
	connectDB()

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
			log.Println(err)
		}

		//Set date variable with row[1] value
		date := row[1]
		//Set amount variable with row[2] value setting as float and without "+" symbol
		amount, err := strconv.ParseFloat(strings.Trim(row[2], " "), 64)

		if err != nil {
			log.Println(err)
		}
		month, err := strconv.Atoi(strings.Split(date, "/")[0])
		if err != nil {
			log.Println(err)
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
			log.Println(err)
		}
		// Insert the row if it doesn't already exist in the table
		if count == 0 {
			_, err = db.Exec("INSERT INTO stori_transactions(file,idFile,transaction,date) VALUES (?,?, ?, ?)", testFile, idFile, amount, date)
			if err != nil {
				log.Println(err)
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
		summary += fmt.Sprintf("Number of transaction in %s: %d \n", monthName, numTransactions)
	}
	balanceStr := fmt.Sprintf("%.2f", balance)
	averageCreditStr := fmt.Sprintf("%.2f", averageCredit)
	averageDebitStr := fmt.Sprintf("%.2f", averageDebit)
	fmt.Println(summary)

	sendEmail(balanceStr, averageCreditStr, averageDebitStr, monthlyTransactions)

}

func sendEmail(balance string, averageCredit string, averageDebit string, transactionsMonthly map[string]int) error {
	from := mail.NewEmail("Stori Challenge", "neto_1208@hotmail.com")
	subject := "Summary Report"
	to := mail.NewEmail("Client", "neto120899@hotmail.com")

	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.Subject = subject

	personalization := mail.NewPersonalization()
	personalization.AddTos(to)

	personalization.SetDynamicTemplateData("balance", balance)
	personalization.SetDynamicTemplateData("debitavg", averageDebit)
	personalization.SetDynamicTemplateData("creditavg", averageCredit)

	var transactions []transaction

	for month, numTran := range transactionsMonthly {
		transactions = append(transactions, transaction{
			Month:   month,
			NumTran: numTran,
		})
	}
	personalization.SetDynamicTemplateData("transactionsMonthly", transactions)
	message.AddPersonalizations(personalization)
	message.SetTemplateID(sgTemplate)

	client := sendgrid.NewSendClient(sgApi)
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
