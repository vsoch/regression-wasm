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
)

// RegressionRunner stores input data for the file
type RegressionRunner struct {
	records [][]string		// array of records
	y []float64			// array of parsed data (y only)
	x [][]float64			// array of parsed data (x only)
	header []string			// first row of headers
	predictCol int			// prediction column
	residuals []float64		// final array of residuals
	predictions []float64		// final array of predictions
	model *regression.Regression	// the regression model
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

// runRegression and generate model to save to runner
func (runner *RegressionRunner) runRegression() {

	r := new(regression.Regression)

	// Iterate through headers to generate variables
	for index, element := range runner.header {

		// Add as regression or observed variable
		if (index != runner.predictCol) {
			r.SetVar(index, element)
		} else {
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

	// Show and save the regression model
	fmt.Printf("Regression formula:\n%v\n", r.Formula)
	//fmt.Printf("Regression:\n%s\n", r)
	runner.model = r
}


// Calculate residuals using the model. If we have more than two covariates,
// then we will plot residuals in a histogram.
func (runner *RegressionRunner) calculateResiduals() {

	// Calculate residuals and predictions
	var predictions []float64
	var residuals []float64

	for i, row := range runner.x {
		if prediction, err := runner.model.Predict(row); err == nil {
			predictions = append(predictions, prediction)
			residuals = append(residuals, runner.y[i] - prediction)
		}
	}

	fmt.Println("Residuals:", residuals)
	fmt.Println("Predictions:", predictions)

	// Save to the runner!
	runner.residuals = residuals
	runner.predictions = predictions
}

// plotResult will generate a histogram for multivariate regression (of 
// residuals) or a line plot given only two variables.
func (runner *RegressionRunner) plotRegression() {

	// Greater than two variables -> regression line
	if len(runner.header) != 2 {
		runner.plotResiduals()
	} else {
		runner.plotLinear()
	}
}

// plotLinear is called given 2 variables, and we create a line plot
func (runner *RegressionRunner) plotLinear() {

	plotFunc := js.Global().Get("plotLinear")
	describeFunc := js.Global().Get("describePlot")

	// Convert data to string to send back to browser
	X := floatArrayToString(runner.x) // only uses first entry
	Y := floatToString(runner.y)
        predictions := floatToString(runner.predictions)
	headers := strings.Join(runner.header, ",")

	result := fmt.Sprintf("Dinosaur Regression Wasm\n%s\n%s", runner.model.Formula, runner.model)

	// provide the title, X, Y, headers, and result
      	plotFunc.Invoke(runner.header[runner.predictCol], X, Y, predictions, headers, result)
       	describeFunc.Invoke(runner.model.Formula)
}

// plotResiduals calls plotResiduals on the front end and passes residuals
func (runner *RegressionRunner) plotResiduals() {

	// Get the functions to do and describe the plot
	plotFunc := js.Global().Get("plotResiduals")
	describeFunc := js.Global().Get("describePlot")

	// Comma separated list of residuals
	resultString := floatToString(runner.residuals)


	// provide the title (predictor) and string data
	result := fmt.Sprintf("Dinosaur Regression Wasm\n%s\n%s", runner.model.Formula, runner.model)
       	plotFunc.Invoke(runner.header[runner.predictCol], resultString, result)
       	describeFunc.Invoke(runner.model.Formula)
}
