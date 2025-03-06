package corebluetooth

import (
	"testing"

	"github.com/progrium/darwinkit/internal/assert"
)
	
func TestCoreBluetoothValid(t *testing.T) {
	// Test that the classes can be accessed
	assert.NotNil(t, CentralManagerClass)
	assert.NotNil(t, PeripheralClass)
	assert.NotNil(t, ServiceClass)
	assert.NotNil(t, CharacteristicClass)
	assert.NotNil(t, DescriptorClass)
	assert.NotNil(t, PeripheralManagerClass)
	assert.NotNil(t, MutableServiceClass)
	assert.NotNil(t, MutableCharacteristicClass)
	assert.NotNil(t, CBUUIDClass)

	// Test the constants
	assert.Equal(t, StateUnknown, State(0))
	assert.Equal(t, StatePoweredOn, State(5))
	assert.Equal(t, PeripheralStateDisconnected, PeripheralState(0))
	assert.Equal(t, PeripheralStateConnected, PeripheralState(2))
	assert.Equal(t, CharacteristicPropertyRead, CharacteristicProperties(0x02))
	assert.Equal(t, CharacteristicPropertyWrite, CharacteristicProperties(0x08))
	assert.Equal(t, AttributePermissionsReadable, AttributePermissions(0x01))
	assert.Equal(t, WriteWithResponse, WriteType(0))
	assert.Equal(t, WriteWithoutResponse, WriteType(1))
}