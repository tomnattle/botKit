/*
	achieve deep copy by marshal/unmarshal

	Warning:
		if the type of value is struct,
		only exported filedes will be deep copy
*/
package util

import (
	"bytes"
	"encoding/gob"
)

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
