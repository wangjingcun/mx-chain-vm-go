package contracts

import (
	"math/big"

	"github.com/multiversx/mx-chain-vm-go/vmhost"
	mock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/stretchr/testify/require"
)

// ForwardAsyncCallParentBuiltinMock is an exposed mock contract method
func ForwardAsyncCallParentBuiltinMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("forwardAsyncCall", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		err := host.Metering().UseGasBounded(testConfig.GasUsedByParent)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
			return instance
		}

		arguments := host.Runtime().Arguments()
		destination := arguments[0]
		function := arguments[1]
		value := big.NewInt(testConfig.TransferFromParentToChild).Bytes()

		if testConfig.IsLegacyAsync {
			err = host.Async().RegisterLegacyAsyncCall(destination, function, value)
		} else {
			err = host.Async().RegisterAsyncCall("testGroup", &vmhost.AsyncCall{
				Status:          vmhost.AsyncCallPending,
				Destination:     destination,
				Data:            function,
				ValueBytes:      value,
				SuccessCallback: "callBack",
				ErrorCallback:   "callBack",
				GasLimit:        testConfig.GasProvidedToChild,
				GasLocked:       150,
			})
		}
		require.Nil(instance.T, err)

		return instance
	})
}

// CallBackParentBuiltinMock is an exposed mock contract method
func CallBackParentBuiltinMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*test.TestConfig)
	instanceMock.AddMockMethod("callBack", test.SimpleWasteGasMockMethod(instanceMock, testConfig.GasUsedByCallback))
}
