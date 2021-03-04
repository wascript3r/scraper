package main

import (
	"encoding/json"
	"time"
)

type Duration struct {
	time.Duration
}

func (d Duration) MarhsalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	pd, err := time.ParseDuration(str)
	if err != nil {
		return err
	}

	d.Duration = pd
	return nil
}
