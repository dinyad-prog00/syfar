package reporters

import (
	"fmt"
	"strings"
	t "syfar/types"
)

func ConsoleReporter(result []t.TestResult) {
	for i, tr := range result {
		if tr.Passed {
			fmt.Printf("\n\x1b[32m%d - PASSED\x1b[0m: %s\n", i, tr.Description) // Green color for passed
		} else {
			fmt.Printf("\n\x1b[31m%d - FAILED\x1b[0m: %s \n", i, tr.Description) // Red color for failed
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
		}
	}
}
