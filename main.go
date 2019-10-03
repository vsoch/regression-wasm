// Copyright 2019 Vanessa Sochat. All rights reserved.
// Use of this source code is governed by the Polyform Strict license
// that can be found in the LICENSE file and available at
// https://polyformproject.org/licenses/noncommercial/1.0.0

package main
 
import (
	"fmt"
	"syscall/js"
)

func main() {

	c := make(chan struct{}, 0)
	js.Global().Set("runRegression", js.FuncOf(runRegression))
	<-c
}

// runRegression is the entrypoint
func runRegression(this js.Value, val []js.Value) interface{} {
	fmt.Println("The input csv string is:", val[0])
	fmt.Println("Header is present:", val[1])
	fmt.Println("Delimiter is:", val[2])
	fmt.Println("Predictor column is:", val[3])

	runner := RegressionRunner{}

	// Browser index is 1, Golang is 0
	runner.predictCol = val[3].Int() - 1

	// read string, true/false for header, and delim
        if err := runner.readCsv(val[0].String(), val[1].Bool(), val[2].String()); err != nil {
		sendMessage("There was an error reading the csv.", "message")
		return nil
	}

	// Ensure that we have data, period
	if len(runner.records) == 0 {
		sendMessage("No records were provided in this dataset.", "message")
		return nil
	}

	// Need more than one columns
	if len(runner.records[0]) == 1 {
		sendMessage("You must provide more than one column of data.", "message")
		return nil
	}

	// Ensure the predictor column is present in the data
	if len(runner.records[0]) < runner.predictCol {
		message := fmt.Sprintf("The predictor col %s is > number cols, %s", val[3].Int(), len(runner.records[0]))
		sendMessage(message, "message")
		return nil
	} 

	fmt.Println("Header is:", runner.header)
	fmt.Println("Records are:", runner.records)

	// All goes well, hide previous messages
	hideMessage("message")

	// run, calculate residuals, and plot. This could be broken up
	// into separate functions for the user to control, if desired, but
	// we would need to maintain state of the model somewhere
	runner.runRegression()
	runner.calculateResiduals()
	runner.plotRegression()	
	return nil
}
