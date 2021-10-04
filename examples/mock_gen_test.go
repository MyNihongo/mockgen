// Code generated by my-nihongo-mockgen. DO NOT EDIT.
package mocking

import (
	"github.com/stretchr/testify/mock"
	"testing"
)

type fixtureImpl1Service struct {
	ser1 *MockService1_1
	ser2 *MockService2_1
}

// AssertExpectations asserts that everything specified with On and Return was in fact called as expected. Calls may have occurred in any order.
func (f *fixtureImpl1Service) AssertExpectations(t *testing.T) {
	f.ser1.AssertExpectations(t)
	f.ser2.AssertExpectations(t)
}

// createFixtureImpl1Service creates a new fixture with all mocks
func createFixtureImpl1Service() (Impl1Service, *fixtureImpl1Service) {
	ser1 := new(MockService1_1)
	ser2 := new(MockService2_1)
	fixture := &impl1{ser1: ser1, ser2: ser2}
	return fixture, &fixtureImpl1Service{ser1: ser1, ser2: ser2}
}

type MockService1_1 struct {
	mock.Mock
}

func (m *MockService1_1) Boo(param string) (uint64, error) {
	ret := m.Called(param)
	return ret.Get(0).(uint64), ret.Error(1)
}
func (m *MockService1_1) OnBoo(param string) *setup_MockService1_1_Boo {
	call := m.On("Boo", param)
	return &setup_MockService1_1_Boo{call: call}
}

type setup_MockService1_1_Boo struct {
	call *mock.Call
}

func (s *setup_MockService1_1_Boo) Return(param1 uint64, param2 error) {
	s.call.Return(param1, param2)
}
func (m *MockService1_1) Foo(param1 string, param2 int16) string {
	ret := m.Called(param1, param2)
	return ret.String(0)
}
func (m *MockService1_1) OnFoo(param1 string, param2 int16) *setup_MockService1_1_Foo {
	call := m.On("Foo", param1, param2)
	return &setup_MockService1_1_Foo{call: call}
}

type setup_MockService1_1_Foo struct {
	call *mock.Call
}

func (s *setup_MockService1_1_Foo) Return(param1 string) {
	s.call.Return(param1)
}

type MockService2_1 struct {
	mock.Mock
}

func (m *MockService2_1) Foo(arg1 string, arg2 string) (string, int, error) {
	ret := m.Called(arg1, arg2)
	return ret.String(0), ret.Int(1), ret.Error(2)
}
func (m *MockService2_1) OnFoo(arg1 string, arg2 string) *setup_MockService2_1_Foo {
	call := m.On("Foo", arg1, arg2)
	return &setup_MockService2_1_Foo{call: call}
}

type setup_MockService2_1_Foo struct {
	call *mock.Call
}

func (s *setup_MockService2_1_Foo) Return(param1 string, param2 int, param3 error) {
	s.call.Return(param1, param2, param3)
}

type fixtureImpl2Service struct {
	ser11 *MockService1_2
	ser3  *MockService2_1
}

// AssertExpectations asserts that everything specified with On and Return was in fact called as expected. Calls may have occurred in any order.
func (f *fixtureImpl2Service) AssertExpectations(t *testing.T) {
	f.ser11.AssertExpectations(t)
	f.ser3.AssertExpectations(t)
}

// createFixtureImpl2Service creates a new fixture with all mocks
func createFixtureImpl2Service() (Impl2Service, *fixtureImpl2Service) {
	ser11 := new(MockService1_2)
	ser3 := new(MockService2_1)
	fixture := &impl2{ser11: ser11, ser3: ser3}
	return fixture, &fixtureImpl2Service{ser11: ser11, ser3: ser3}
}

type MockService1_2 struct {
	mock.Mock
}

func (m *MockService1_2) Boo(param string) {
	m.Called(param)
	return
}
func (m *MockService1_2) OnBoo(param string) {
	m.On("Boo", param)
	return
}
func (m *MockService1_2) Foo(param1 string, param2 int16) (int, bool) {
	ret := m.Called(param1, param2)
	return ret.Int(0), ret.Bool(1)
}
func (m *MockService1_2) OnFoo(param1 string, param2 int16) *setup_MockService1_2_Foo {
	call := m.On("Foo", param1, param2)
	return &setup_MockService1_2_Foo{call: call}
}

type setup_MockService1_2_Foo struct {
	call *mock.Call
}

func (s *setup_MockService1_2_Foo) Return(param1 int, param2 bool) {
	s.call.Return(param1, param2)
}