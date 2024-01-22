package app

import (
	"encoding/json"
	"net/http"
)

type HtmxTrigger struct {
	Name  string
	Value map[string]string
}

func SetHtmxTrigger(w http.ResponseWriter, trigger HtmxTrigger) error {
	valueJson, err := json.Marshal(trigger.Value)
	if err != nil {
		return err
	}

	triggerJson, err := json.Marshal(map[string]string{
		trigger.Name: string(valueJson),
	})
	if err != nil {
		return err
	}

	w.Header()["HX-Trigger"] = []string{
		string(triggerJson),
	}

	return nil
}
