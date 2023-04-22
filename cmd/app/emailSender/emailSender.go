package emailSender

import (
	"fmt"
	"os"
	"strings"

	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/config"
	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/logger"
	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(balance string, averageCredit string, averageDebit string, transactionsMonthly map[string]int) error {
	// Load the configuration values
	cfg, err := config.LoadConfig()
	if err != nil {
		logger := logger.NewLogger("CONFIG_ERROR: ")
		logger.Error(fmt.Sprintf("Error loading configuration values: %v", err))
	}

	from := mail.NewEmail("Stori Challenge", "neto120899@hotmail.com")
	emailAddress := os.Getenv("EMAIL_ADDRESS")
	//Validate if it is a valid email address
	if !strings.Contains(emailAddress, "@") {
		// Handle the error case
		logger := logger.NewLogger("EMAIL_ERROR: ")
		logger.Error("¡Invalid email address!")
	} else {
		to := mail.NewEmail("Client", os.Getenv("EMAIL_ADDRESS"))

		message := mail.NewV3Mail()
		message.SetFrom(from)
		message.Subject = "Summary Report"

		personalization := mail.NewPersonalization()
		personalization.AddTos(to)
		personalization.Subject = "Summary Report"

		personalization.SetDynamicTemplateData("balance", balance)
		personalization.SetDynamicTemplateData("debitavg", averageDebit)
		personalization.SetDynamicTemplateData("creditavg", averageCredit)

		var transactions []model.Transaction

		for month, numTran := range transactionsMonthly {
			transactions = append(transactions, model.Transaction{
				Month:   month,
				NumTran: numTran,
			})
		}
		personalization.SetDynamicTemplateData("transactionsMonthly", transactions)
		message.AddPersonalizations(personalization)
		message.SetTemplateID(cfg.SendGrid.Template)

		client := sendgrid.NewSendClient(cfg.SendGrid.Api)
		response, err := client.Send(message)
		if err != nil {
			logger := logger.NewLogger("EMAIL_ERROR: ")
			logger.Error(fmt.Sprintf("Error consuming email API: %v", err))

		} else if response.StatusCode == 202 {
			logger := logger.NewLogger("EMAIL_INFO: ")
			logger.Info("¡Email sended!")
		}

	}

	return nil
}
