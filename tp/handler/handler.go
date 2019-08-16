/**
 * Copyright 2017 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ------------------------------------------------------------------------------
 */

package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
)

const (
	FAMILY_NAME    = "break_sawtooth"
	FAMILY_VERSION = "0.1"
	ADDRESS_PREFIX = "BLOCK"
)

type BreakSawtoothHandler struct {
}

func NewBreakSawtoothHandler() *BreakSawtoothHandler {
	return &BreakSawtoothHandler{}
}

func (self *BreakSawtoothHandler) FamilyName() string {
	return FAMILY_NAME
}

func (self *BreakSawtoothHandler) FamilyVersions() []string {
	return []string{FAMILY_VERSION}
}

func (self *BreakSawtoothHandler) Namespaces() []string {
	return []string{Hexdigest(FAMILY_NAME)[:6]}
}

func (self *BreakSawtoothHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	// Make sure that we have payload data
	payloadBytes := request.GetPayload()
	if payloadBytes == nil {
		return &processor.InvalidTransactionError{
			Msg: "Transaction must contain a payload",
		}
	}

	// Decode the payload from json format
	var payload BreakSawtoothPayload
	err := json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return &processor.InvalidTransactionError{
			Msg: fmt.Sprint("Failed to decode payload with error: ", err),
		}
	}

	// Get the address for the block
	address := GetAddress(payload.Id)

	// Retrieve any existing blockchain state
	state, err := context.GetState([]string{address})
	if err != nil {
		return &processor.InvalidTransactionError{
			Msg: fmt.Sprintf("Failed to get state for block %s with error %s", payload.Id, err),
		}
	}

	// Save the state
	state[address] = payload.Block

	return self.saveState(state, context)
}

func (self *BreakSawtoothHandler) saveState(state map[string][]byte, context *processor.Context) error {
	resultAddresses, err := context.SetState(state)
	if err != nil {
		return &processor.InvalidTransactionError{
			Msg: fmt.Sprint("Failed to set state with error %s", err),
		}
	}

	if len(resultAddresses) != len(state) {
		return &processor.InternalError{
			Msg: fmt.Sprintf("Wrong nunber of address in state response: %d != %d", len(resultAddresses), len(state)),
		}
	}

	return nil
}

func GetAddressPrefix() string {
	namespace := Hexdigest(FAMILY_NAME)[:6]
	prefix := Hexdigest(ADDRESS_PREFIX)[:6]
	return namespace + prefix
}
func GetAddress(blockId uuid.UUID) string {
	prefix := GetAddressPrefix()
	suffix := Hexdigest(blockId.String())[:58]
	return prefix + suffix
}
