package main

import (
	"encoding/json"
	"os"
	"testing"

	nutest "github.com/nuclio/nuclio-test-go"
	"github.com/stretchr/testify/assert"
)

func TestBase64Image(t *testing.T) {
	var re request

	data, err := os.ReadFile("helloworld.txt")
	assert.Nil(t, err)

	re.ImageName = "helloworld.png"
	re.Base64Image = string(data)

	body, err := json.Marshal(re)
	assert.Nil(t, err)

	// Create TestContext and specify the function name, verbose, data
	tc, err := nutest.NewTestContext(Handler, false, nil)
	assert.Nil(t, err)

	// Optional, initialize context must have a function in the form:
	//    InitContext(context *nuclio.Context) error
	err = tc.InitContext(InitContext)
	assert.Nil(t, err)

	// Create a new test event
	testEvent := nutest.TestEvent{
		Path: "/ocr/0/0/1.png",
		Body: body,
	}

	// invoke the tested function with the new event and log it's output
	resp, err := tc.Invoke(&testEvent)
	tc.Logger.InfoWith("Run complete", "resp", resp, "err", err)
}

func TestUrlImage(t *testing.T) {
	var re request

	re.ImageName = "helloworld.png"
	re.ImageUrl = "https://www.diggernaut.com/sandbox/hello_world.png"

	body, err := json.Marshal(re)
	assert.Nil(t, err)

	// Create TestContext and specify the function name, verbose, data
	tc, err := nutest.NewTestContext(Handler, false, nil)
	assert.Nil(t, err)

	// Optional, initialize context must have a function in the form:
	//    InitContext(context *nuclio.Context) error
	err = tc.InitContext(InitContext)
	assert.Nil(t, err)

	// Create a new test event
	testEvent := nutest.TestEvent{
		Path: "/ocr",
		Body: body,
	}

	// invoke the tested function with the new event and log it's output
	resp, err := tc.Invoke(&testEvent)
	tc.Logger.InfoWith("Run complete", "resp", resp, "err", err)
}
