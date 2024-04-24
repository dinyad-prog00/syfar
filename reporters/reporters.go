package reporters

import (
	"fmt"
	"strings"
	t "syfar/types"
)

func ConsoleReporter(result t.SyfarResult) {
	nbTests := result.NbTestsPassed + result.NbTestsFailed + result.NbTestSkipped
	passedColor := ""
	failedColor := ""
	skippedColor := ""

	if result.NbTestsPassed == 0 {
		passedColor = "\x1b[30m"
	} else {
		passedColor = "\x1b[32m"
	}

	if result.NbTestsFailed == 0 {
		failedColor = "\x1b[30m"
	} else {
		failedColor = "\x1b[31m"
	}

	if result.NbTestSkipped == 0 {
		skippedColor = "\x1b[30m"
	} else {
		skippedColor = "\x1b[33m"
	}

	fmt.Println("___________________________________________________________________")
	fmt.Println("\nTests result")
	fmt.Println("___________________________________________________________________")
	fmt.Printf("\n%s%d/%d Passed\x1b[0m\t%s%d/%d Failed\x1b[0m\t %s%d/%d Skipped\x1b[0m", passedColor, result.NbTestsPassed, nbTests, failedColor, result.NbTestsFailed, nbTests, skippedColor, result.NbTestSkipped, nbTests)
	fmt.Println("\n___________________________________________________________________")

	for i, tr := range result.TestsResult {

		switch tr.State {
		case t.StatePassed:
			fmt.Printf("\n\x1b[32m%d - PASSED\x1b[0m: %s\n", i+1, tr.Description)
		case t.StateFailed:
			fmt.Printf("\n\x1b[31m%d - FAILED\x1b[0m: %s \n", i+1, tr.Description) // Red color for failed
			for _, er := range tr.Expectations {
				for _, ck := range er.Items {
					if !er.Passed {
						msg := strings.ReplaceAll(ck.Message, "{{SEXP}}", "\x1b[33m")
						msg = strings.ReplaceAll(msg, "{{EEXP}}", "\x1b[0m")
						msg = strings.ReplaceAll(msg, "{{SGOT}}", "\x1b[31m")
						msg = strings.ReplaceAll(msg, "{{EGOT}}", "\x1b[0m")
						fmt.Printf("\t%s\n", msg)
					}
				}
			}
		case t.StateSkipped:
			fmt.Printf("\n\x1b[33m%d - SKIPPED\x1b[0m: %s\n", i+1, tr.Description)
		}
	}
}
