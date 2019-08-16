package client

import (
	"github.com/google/uuid"
	"github.com/taekion-org/break-sawtooth-history/tp/handler"
	"github.com/taekion-org/sawtooth-client-sdk-go"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport"
)

type BreakSawtoothClient struct {
	*sawtooth_client_sdk_go.SawtoothClient
}

// NewBreakSawtoothClient returns a new instance of BreakSawtoothClient.
func NewBreakSawtoothClient(url string, keyFile string) (*BreakSawtoothClient, error) {
	args := &sawtooth_client_sdk_go.SawtoothClientArgs{
		URL: url,
		KeyFile: keyFile,
		TransportType: transport.TRANSPORT_REST,
		Impl: &BreakSawtoothClientImpl{},
	}

	sawtoothClient, err := sawtooth_client_sdk_go.NewClient(args)
	if err != nil {
		return nil, err
	}

	client := &BreakSawtoothClient{
		SawtoothClient: sawtoothClient,
	}

	return client, nil
}

// List returns the current mapping of keys to values.
func (self *BreakSawtoothClient) GetBlock(id uuid.UUID, head string) ([]byte, string, error) {
	address := handler.GetAddress(id)

	if head == "" {
		result, err := self.Transport.GetState(address)
		if err != nil {
			return nil, "", err
		}
		return result.Data, result.Head, nil
	} else {
		result, err := self.Transport.GetStateAtHead(address, head)
		if err != nil {
			return nil, "", err
		}
		return result.Data, result.Head, nil
	}
}

// sendTransaction is a common method used to construct a payload and submit it for processing.
func (self *BreakSawtoothClient) SetBlock(id uuid.UUID, data []byte, wait uint) (string, error) {
	payload := handler.BreakSawtoothPayload {
		Id: id,
		Block: data,
		Action: handler.WRITE_BLOCK,
	}

	batchId, err := self.ExecutePayload(&payload)
	if err != nil {
		return batchId, err
	}

	if wait != 0 {
		_, err := self.WaitBatch(batchId, int(wait), 1)
		if err != nil {
			return batchId, err
		}
	}

	return batchId, nil
}
