// Copyright 2019 Vanessa Sochat. All rights reserved.
// Use of this source code is governed by the Polyform Strict license
// that can be found in the LICENSE file and available at
// https://polyformproject.org/licenses/noncommercial/1.0.0

package main
 
import (
	"encoding/csv"
	"fmt"
	"strings"
	"strconv"
	"syscall/js"

	"github.com/sajari/regression"
	"github.com/aybabtme/uniplot/histogram"
)

func main() {

	c := make(chan struct{}, 0)
	js.Global().Set("runRegression", js.FuncOf(runRegression))
	<-c
}

// RegressionRunner stores input data for the file
type RegressionRunner struct {
	records [][]string	// array of records
	y []float64		// array of parsed data (y only)
	x [][]float64		// array of parsed data (x only)
	header []string		// first row of headers
	predictCol int		// prediction column
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
		returnResult("There was an error reading the csv.", "message")
		return nil
	}

	// Ensure that we have data, period
	if len(runner.records) == 0 {
		returnResult("No records were provided in this dataset.", "message")
		return nil
	}

	// Ensure the predictor column is present in the data
	if len(runner.records[0]) < runner.predictCol {
		message := fmt.Sprintf("The predictor col %s is > number cols, %s", val[3].Int(), len(runner.records[0]))
		returnResult(message, "message")
		return nil
	} 

	fmt.Println("Header is:", runner.header)
	fmt.Println("Records are:", runner.records)

	runner.runRegression()

	return nil
}


// returnResult back to the browser, in the innerHTML of the result element
func returnResult(output string, divid string) {
	js.Global().Get("document").
		Call("getElementById", divid).
		Set("innerHTML", output)
}

// readCsv file and set the records on the runner
func (runner *RegressionRunner) readCsv(csvString string, hasHeader bool, delim string) error {

	reader := csv.NewReader(strings.NewReader(csvString))
	records, err := reader.ReadAll()
	var header []string

	if err != nil {
		return err
	}

	// Remove header row, if we have it
	if hasHeader {
		header, records = records[0], records[1:]
		runner.header = header
	} else {
		// Add dummy names
		for i := range records[0] {
			header = append(header, fmt.Sprintf("Element %d",  i))
		}
	}

	runner.records = records
	return nil
}


func (runner *RegressionRunner) runRegression() {

	r := new(regression.Regression)

	// Iterate through headers to generate variables
	for index, element := range runner.header {

		fmt.Println("Index:", index)
		fmt.Println("Element:", element)

		// Add as regression or observed variable
		if (index != runner.predictCol) {
	 		fmt.Println("Adding regression variable:", element)
			r.SetVar(index, element)
		} else {
	 		fmt.Println("Adding observed variable:", element)
			r.SetObserved(element)
		}
	}

	// We will unwrap an array of data points
	var dataPoints regression.DataPoints
	var predictor float64
	var regressors []float64	

	// Iterate through records to generate dataPoints
	for _, row := range runner.records {

		regressors = nil
		for index, record := range row {

			if (index != runner.predictCol) {
				if n, err := strconv.ParseFloat(record, 64); err == nil {
			        	regressors = append(regressors, n)
				}
			} else {
				if n, err := strconv.ParseFloat(record, 64); err == nil {
			        	predictor = n
				}
			}
		}

		dataPoints = append(dataPoints, regression.DataPoint(predictor, regressors))
		runner.x = append(runner.x, regressors)
		runner.y = append(runner.y, predictor)
	}

	// Unwrap data points into function
	r.Train(dataPoints...)
	r.Run()

	// Calculate residuals and predictions
	var predictions []float64
	var residuals []float64

	for i, row := range runner.x {
		if prediction, err := r.Predict(row); err == nil {
			predictions = append(predictions, prediction)
			residuals = append(residuals, runner.y[i] - prediction)
		}
	}

	fmt.Println("Residuals:", residuals)
	fmt.Println("Predictions:", predictions)
	fmt.Printf("Regression formula:\n%v\n", r.Formula)
	fmt.Printf("Regression:\n%s\n", r)

	// Generate histogram
	hist := histogram.Hist(10, residuals)

	for _, bkt := range hist.Buckets {
		fmt.Println(bkt.Min)
		fmt.Println(bkt.Max)
		fmt.Println(bkt.Count)
	}

	// Get the function to do the plot
	plotFunc := js.Global().Get("drawPlot")
        plotFunc.Invoke("hello")
}
