package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PrettyJSON(rawJSON []byte) error {
	var pretty bytes.Buffer
	err := json.Indent(&pretty, rawJSON, "", "\t")
	if err != nil {
		return err
	}

	fmt.Printf("\n%s\n", pretty.String())
	return nil
}
