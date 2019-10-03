// Copyright 2019 Vanessa Sochat. All rights reserved.
// Use of this source code is governed by the Polyform Strict license
// that can be found in the LICENSE file and available at
// https://polyformproject.org/licenses/noncommercial/1.0.0

package main
 
import (
	"fmt"
	"strings"
	"syscall/js"
)

// convert an array of floats to a string to send back to browser
func floatToString(floats []float64) string {
	var values []string
	var stringValues string
	for i := range floats {
		text := fmt.Sprintf("%f", floats[i])
	        values = append(values, text)
	}
	stringValues = strings.Join(values, ",")
	return stringValues
}

// floatArrayToString converts an array of floats to a single string
// we expect there only to be one value per array entry
func floatArrayToString(floats [][]float64) string {
	var values []string
	var stringValues string
	for i := range floats {
		text := fmt.Sprintf("%f", floats[i][0])
	        values = append(values, text)
	}
	stringValues = strings.Join(values, ",")
	return stringValues
}


// returnResult back to the browser, in the innerHTML of the result element
func returnResult(output string, divid string) {
	js.Global().Get("document").
		Call("getElementById", divid).
		Set("innerHTML", output)
}

func sendMessage (message string, div_id string) {
	messageFunc := js.Global().Get("showMessage")
       	messageFunc.Invoke(message, div_id)
}

func hideMessage (div_id string) {
	messageFunc := js.Global().Get("hideMessage")
       	messageFunc.Invoke(div_id)
}
