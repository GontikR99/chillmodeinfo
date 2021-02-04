// +build wasm

package binding

import "encoding/json"

func trimOptions(prefMap *map[string]interface{}) {
	var trimKeys []string
	for k,v := range *prefMap {
		if v==nil {
			trimKeys = append(trimKeys, k)
		}
		if submap, ok :=v.(map[string]interface{}); ok {
			trimOptions(&submap)
			(*prefMap)[k]=submap
		}
		if subarr, ok:=v.([]interface{}); ok {
			var newArr []interface{}
			for _, v2:=range subarr {
				if submap, ok := v2.(map[string]interface{}); ok {
					trimOptions(&submap)
					newArr = append(newArr, submap)
				} else {
					newArr = append(newArr, v2)
				}
			}
			(*prefMap)[k]=newArr
		}
	}

	for _, k := range trimKeys {
		delete(*prefMap, k)
	}
}

func JsonifyOptions(options interface{}) map[string]interface{} {
	data, err := json.Marshal(options)
	if err != nil {
		panic(err)
	}
	parsed:=make(map[string]interface{})
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		panic(err)
	}
	trimOptions(&parsed)
	return parsed
}