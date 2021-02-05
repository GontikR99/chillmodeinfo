// +build wasm

package ipc

import (
	"encoding/hex"
)

func Encode(data []byte) string {
	//buffer := new(bytes.Buffer)
	//base64.NewEncoder(base64.StdEncoding, buffer).Write(data)
	//return string(buffer.Bytes())
	return hex.EncodeToString(data)
}

func Decode(encoded string) ([]byte,error) {
	//decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(encoded))
	//data, err := ioutil.ReadAll(decoder)
	//if err!=nil {
	//	console.Log("Failed to decode", err)
	//}
	//return data, err
	return hex.DecodeString(encoded)
}