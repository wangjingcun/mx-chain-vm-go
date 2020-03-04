package nodepart

import (
	"context"
	"os"
	"time"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
)

// NodeMessenger is
type NodeMessenger struct {
	common.Messenger
}

// NewNodeMessenger creates
func NewNodeMessenger(reader *os.File, writer *os.File) *NodeMessenger {
	return &NodeMessenger{
		Messenger: *common.NewMessenger("Node", reader, writer),
	}
}

// SendContractRequest sends
func (messenger *NodeMessenger) SendContractRequest(request *common.ContractRequest) error {
	err := messenger.Send(request)
	if err != nil {
		return common.ErrCannotSendContractRequest
	}

	common.LogDebug("Node: sent contract request %s", request)
	return nil
}

// SendHookCallResponse sends
func (messenger *NodeMessenger) SendHookCallResponse(response *common.HookCallResponse) error {
	err := messenger.Send(response)
	if err != nil {
		return common.ErrCannotSendHookCallResponse
	}

	common.LogDebug("Node: sent hook call response %s", response)
	return nil
}

// ReceiveHookCallRequestOrContractResponse waits
func (messenger *NodeMessenger) ReceiveHookCallRequestOrContractResponse() (*common.HookCallRequestOrContractResponse, error) {
	message := &common.HookCallRequestOrContractResponse{}

	timeout := 1 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errorChannel := make(chan error, 1)

	select {
	case <-ctx.Done():
		common.LogError("TIMEOUT ReceiveHookCallRequestOrContractResponse")
		messenger.Shutdown()
	case errorChannel <- ctx.Err():
		common.LogError("READ ERROR ReceiveHookCallRequestOrContractResponse")
	}

	err := messenger.Receive(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}
