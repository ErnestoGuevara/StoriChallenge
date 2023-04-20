package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

//Generate the struct of Transactions
type Transactions struct {
	Id    int
	Date  string
	Value float64
}

func main() {
	csvFile()
}

func csvFile() {
	file, err := os.Open("./test/test1.csv")
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
		_, err := strconv.Atoi(row[0])

		if err != nil {
			log.Fatal(err)
		}

		date := row[1]

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
	}

	// Generate summary report
	averageCredit := calculateAverage(credits)
	averageDebit := calculateAverage(debits)

	summary := fmt.Sprintf("Total balance is %.2f\nAverage debit amount: %.2f\nAverage credit amount: %.2f",
		balance, averageDebit, averageCredit)
	fmt.Println(summary)
	//Print Monthly transactions
	for monthName, numTransactions := range monthlyTransactions {
		fmt.Printf("Number of transaction in %s: %d\n", monthName, numTransactions)
	}
	// Create a new session to the SES service in the us-west-2 region
	sess, err := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1")},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Construct the email message
	message := fmt.Sprintf("Summary:\n%s\n\nMonthly Transactions:\n", summary)
	for monthName, numTransactions := range monthlyTransactions {
		message += fmt.Sprintf("Number of transaction in %s: %d\n", monthName, numTransactions)
	}

	// Construct the email parameters
	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String("neto_1208@hotmail.com"),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(message),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("Monthly Transactions Report"),
			},
		},
		Source: aws.String("neto_1208@hotmail.com"),
	}

	// Send the email
	svc := ses.New(sess)
	_, err = svc.SendEmail(params)
	if err != nil {
		log.Fatal(err)
	}
}

func calculateAverage(numbers []float64) float64 {
	var sum float64
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}
