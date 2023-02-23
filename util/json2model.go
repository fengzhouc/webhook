package util

import (
	"encoding/json"
	"fmt"
)

// #################将jsonString转换成任意Model对象############################################
func JsonToAnyModel(model interface{}, jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), model)
	if err != nil {
		fmt.Println("[JsonToAnyModel err] ", err, "/n", jsonString)
		return err
	}
	return nil
}
