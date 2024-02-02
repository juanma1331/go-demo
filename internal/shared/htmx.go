package shared

import (
	"encoding/json"
	"net/http"
)

type HtmxTrigger struct {
	Name  string
	Value map[string]string
}

func SetHtmxTriggers(w http.ResponseWriter, triggers ...HtmxTrigger) error {
	triggerMap := make(map[string]string)

	for _, trigger := range triggers {
		valueJson, err := json.Marshal(trigger.Value)
		if err != nil {
			return err
		}
		triggerMap[trigger.Name] = string(valueJson)
	}

	triggerJson, err := json.Marshal(triggerMap)
	if err != nil {
		return err
	}

	w.Header().Set("HX-Trigger", string(triggerJson))

	return nil
}

func SetHtmxRetarget(w http.ResponseWriter, target string) {
	w.Header()["HX-Retarget"] = []string{
		target,
	}
}

func SetHtmxReswap(w http.ResponseWriter, swap string) {
	w.Header()["HX-Reswap"] = []string{
		swap,
	}
}
