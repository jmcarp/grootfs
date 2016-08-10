// This file was generated by counterfeiter
package clonerfakes

import (
	"sync"

	"code.cloudfoundry.org/grootfs/cloner"
	"code.cloudfoundry.org/lager"
)

type FakeUnpacker struct {
	UnpackStub        func(logger lager.Logger, spec cloner.UnpackSpec) error
	unpackMutex       sync.RWMutex
	unpackArgsForCall []struct {
		logger lager.Logger
		spec   cloner.UnpackSpec
	}
	unpackReturns struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeUnpacker) Unpack(logger lager.Logger, spec cloner.UnpackSpec) error {
	fake.unpackMutex.Lock()
	fake.unpackArgsForCall = append(fake.unpackArgsForCall, struct {
		logger lager.Logger
		spec   cloner.UnpackSpec
	}{logger, spec})
	fake.recordInvocation("Unpack", []interface{}{logger, spec})
	fake.unpackMutex.Unlock()
	if fake.UnpackStub != nil {
		return fake.UnpackStub(logger, spec)
	} else {
		return fake.unpackReturns.result1
	}
}

func (fake *FakeUnpacker) UnpackCallCount() int {
	fake.unpackMutex.RLock()
	defer fake.unpackMutex.RUnlock()
	return len(fake.unpackArgsForCall)
}

func (fake *FakeUnpacker) UnpackArgsForCall(i int) (lager.Logger, cloner.UnpackSpec) {
	fake.unpackMutex.RLock()
	defer fake.unpackMutex.RUnlock()
	return fake.unpackArgsForCall[i].logger, fake.unpackArgsForCall[i].spec
}

func (fake *FakeUnpacker) UnpackReturns(result1 error) {
	fake.UnpackStub = nil
	fake.unpackReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeUnpacker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.unpackMutex.RLock()
	defer fake.unpackMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeUnpacker) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ cloner.Unpacker = new(FakeUnpacker)