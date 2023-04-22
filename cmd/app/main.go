package main

import (
	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/summaryProcessor"
)

func main() {
	summaryProcessor.SummaryReportGenerator("/app/client1.csv")
	summaryProcessor.SummaryReportGenerator("/app/client2.csv")

}
