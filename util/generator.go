package util

import (
	"encoding/json"

	"github.com/ipfs/go-datastore"
	"github.com/rs/xid"
)

// GenerateRandomID prefixes a table name and suffixes it with a random HEX string
func GenerateRandomID() string {
	return xid.New().String()
}

func CreateKey(table string, id string) datastore.Key {
	return datastore.NewKey("/" + table + "/" + id)
}

// marshal any JSON compatible object to bytes
func MarshalToBytes(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

// unmarhsal from bytes back to any object
func UnmarshalFromBytes(obj []byte) (interface{}, error) {
	var out interface{}
	err := json.Unmarshal(obj, &out)
	return out, err
}
