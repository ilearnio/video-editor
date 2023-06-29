package helpers

import "encoding/json"

func EncodeToMap(obj interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
