package core

import (
	"encoding/json"
	"fmt"
)

func GraphQLSmartErrorMessageBytes(msg interface{}) []byte {
	messageFormat := `{"data":null,"error":%s}`
	message, err := json.Marshal(msg)
	if err == nil {
		output := fmt.Sprintf(messageFormat, message)
		return []byte(output)
	} else {
		output := []byte(`{"data":null,"error":[]}`)
		return output
	}
}
