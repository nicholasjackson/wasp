package engine

import "github.com/stretchr/testify/mock"

type mockInstance struct {
	mock.Mock
}

func (m *mockInstance) CallFunction(name string, outParam interface{}, inParam ...interface{}) error {
	return m.Called(name, outParam, inParam).Error(0)
}

func (m *mockInstance) Remove() error {
	return m.Called().Error(0)
}

func (m *mockInstance) setError(err string) {
	m.Called(err)
}

func (m *mockInstance) getImportObject() importObject {
	return m.Called().Get(0).(importObject)
}

func (m *mockInstance) getError() error {
	return m.Called().Error(0)
}

func (m *mockInstance) setStringInMemory(str string) (int32, error) {
	args := m.Called(str)

	return args.Get(0).(int32), args.Error(1)
}

func (m *mockInstance) getStringFromMemory(ptr int32) (string, error) {
	args := m.Called(ptr)

	return args.Get(0).(string), args.Error(1)
}

func (m *mockInstance) setBytesInMemory(data []byte) (int32, error) {
	args := m.Called(data)

	return args.Get(0).(int32), args.Error(1)
}

func (m *mockInstance) getBytesFromMemory(ptr int32) ([]byte, error) {
	args := m.Called(ptr)

	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockInstance) freeAllocatedMemory() {
	m.Called()
}

func (m *mockInstance) getStringSize(addr int32) (int32, error) {
	args := m.Called(addr)

	return args.Get(0).(int32), args.Error(1)
}

func (m *mockInstance) allocate(size int32) (int32, error) {
	args := m.Called(size)

	return args.Get(0).(int32), args.Error(1)
}

func (m *mockInstance) deallocate(addr int32, size int32) error {
	args := m.Called(addr, size)

	return args.Error(0)
}
