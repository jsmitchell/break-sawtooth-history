package client

import (
	"encoding/json"
	"github.com/taekion-org/break-sawtooth-history/tp/handler"
)

type BreakSawtoothClientImpl struct {}

func (self *BreakSawtoothClientImpl) GetFamilyName() string {
	return handler.FAMILY_NAME
}

func (self *BreakSawtoothClientImpl) GetFamilyVersion() string {
	return handler.FAMILY_VERSION
}

func (self *BreakSawtoothClientImpl) EncodePayload(payload interface{}) ([]byte, error) {
	return json.Marshal(payload)
}

func (self *BreakSawtoothClientImpl) DecodePayload(bytes []byte, ptr interface{}) error {
	return json.Unmarshal(bytes, ptr)
}

func (self *BreakSawtoothClientImpl) EncodeData(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (self *BreakSawtoothClientImpl) DecodeData(bytes []byte, ptr interface{}) error {
	return json.Unmarshal(bytes, ptr)
}

func (self *BreakSawtoothClientImpl) GetPayloadInputAddresses(payload interface{}) []string {
	id := payload.(*handler.BreakSawtoothPayload).Id
	return []string{handler.GetAddress(id)}
}

func (self *BreakSawtoothClientImpl) GetPayloadOutputAddresses(payload interface{}) []string {
	id := payload.(*handler.BreakSawtoothPayload).Id
	return []string{handler.GetAddress(id)}
}
