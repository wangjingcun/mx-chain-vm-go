package scenarioexec

import (
	"fmt"

	logger "github.com/multiversx/mx-chain-logger-go"
	vmi "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/hostCore"
	"github.com/multiversx/mx-chain-vm-go/vmhost/mock"
	gasSchedules "github.com/multiversx/mx-chain-vm-go/scenarioexec/gasSchedules"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/executor"
	mc "github.com/multiversx/mx-chain-vm-go/scenarios/controller"
	er "github.com/multiversx/mx-chain-vm-go/scenarios/expression/reconstructor"
	fr "github.com/multiversx/mx-chain-vm-go/scenarios/fileresolver"
	mj "github.com/multiversx/mx-chain-vm-go/scenarios/model"
	worldhook "github.com/multiversx/mx-chain-vm-go/mock/world"
)

var log = logger.GetOrCreate("arwen/scenarios")

// TestVMType is the VM type argument we use in tests.
var TestVMType = []byte{0, 0}

// ArwenTestExecutor parses, interprets and executes both .test.json tests and .scen.json scenarios with Arwen.
type ArwenTestExecutor struct {
	World              *worldhook.MockWorld
	vm                 vmi.VMExecutionHandler
	OverrideVMExecutor executor.ExecutorAbstractFactory
	vmHost             vmhost.VMHost
	checkGas           bool
	scenarioTraceGas   []bool
	fileResolver       fr.FileResolver
	exprReconstructor  er.ExprReconstructor
}

var _ mc.TestExecutor = (*ArwenTestExecutor)(nil)
var _ mc.ScenarioExecutor = (*ArwenTestExecutor)(nil)

// NewArwenTestExecutor prepares a new ArwenTestExecutor instance.
func NewArwenTestExecutor() (*ArwenTestExecutor, error) {
	world := worldhook.NewMockWorld()

	return &ArwenTestExecutor{
		World:             world,
		vm:                nil,
		checkGas:          true,
		scenarioTraceGas:  make([]bool, 0),
		fileResolver:      nil,
		exprReconstructor: er.ExprReconstructor{},
	}, nil
}

// InitVM will initialize the VM and the builtin function container.
// Does nothing if the VM is already initialized.
func (ae *ArwenTestExecutor) InitVM(mandosGasSchedule mj.GasSchedule) error {
	if ae.vm != nil {
		return nil
	}

	gasSchedule, err := ae.gasScheduleMapFromMandos(mandosGasSchedule)
	if err != nil {
		return err
	}

	err = ae.World.InitBuiltinFunctions(gasSchedule)
	if err != nil {
		return err
	}

	blockGasLimit := uint64(10000000)
	esdtTransferParser, _ := parsers.NewESDTTransferParser(worldhook.WorldMarshalizer)

	vm, err := hostCore.NewArwenVM(
		ae.World,
		&vmhost.VMHostParameters{
			VMType:                   TestVMType,
			OverrideVMExecutor:       ae.OverrideVMExecutor,
			BlockGasLimit:            blockGasLimit,
			GasSchedule:              gasSchedule,
			BuiltInFuncContainer:     ae.World.BuiltinFuncs.Container,
			ProtectedKeyPrefix: []byte(ProtectedKeyPrefix),
			ESDTTransferParser:       esdtTransferParser,
			EpochNotifier:            &mock.EpochNotifierStub{},
			EnableEpochsHandler:      worldhook.EnableEpochsHandlerStubAllFlags(),
			WasmerSIGSEGVPassthrough: false,
		})
	if err != nil {
		return err
	}

	ae.vm = vm
	ae.vmHost = vm
	return nil
}

// GetVM yields a reference to the VMExecutionHandler used.
func (ae *ArwenTestExecutor) GetVM() vmi.VMExecutionHandler {
	return ae.vm
}

func (ae *ArwenTestExecutor) getVMHost() vmhost.VMHost {
	return ae.vmHost
}

func (ae *ArwenTestExecutor) gasScheduleMapFromMandos(mandosGasSchedule mj.GasSchedule) (config.GasScheduleMap, error) {
	switch mandosGasSchedule {
	case mj.GasScheduleDefault:
		return gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV4())
	case mj.GasScheduleDummy:
		return config.MakeGasMapForTests(), nil
	case mj.GasScheduleV3:
		return gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV3())
	case mj.GasScheduleV4:
		return gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV4())
	default:
		return nil, fmt.Errorf("unknown mandos GasSchedule: %d", mandosGasSchedule)
	}
}

// PeekTraceGas returns the last position from the scenarioTraceGas, if existing
func (ae *ArwenTestExecutor) PeekTraceGas() bool {
	length := len(ae.scenarioTraceGas)
	if length != 0 {
		return ae.scenarioTraceGas[length-1]
	}
	return false
}
