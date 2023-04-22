package csvProcessor

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/logger"
)

func ReadFile(testFile string) ([][]string, error) {

	file, err := os.Open(testFile)
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

	return transactions, nil
}
