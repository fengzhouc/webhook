package util

import (
	"encoding/json"
	"log"
)

// #################将jsonString转换成任意Model对象############################################
func JsonToAnyModel(model interface{}, jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), model)
	if err != nil {
		log.Println("[JsonToAnyModel err] ", err, "/n", jsonString)
		return err
	}
	return nil
}
