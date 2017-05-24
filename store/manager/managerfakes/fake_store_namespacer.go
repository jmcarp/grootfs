// This file was generated by counterfeiter
package managerfakes

import (
	"sync"

	"code.cloudfoundry.org/grootfs/groot"
	"code.cloudfoundry.org/grootfs/store/manager"
)

type FakeStoreNamespacer struct {
	ApplyMappingsStub        func(uidMappings, gidMappings []groot.IDMappingSpec) error
	applyMappingsMutex       sync.RWMutex
	applyMappingsArgsForCall []struct {
		uidMappings []groot.IDMappingSpec
		gidMappings []groot.IDMappingSpec
	}
	applyMappingsReturns struct {
		result1 error
	}
	applyMappingsReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeStoreNamespacer) ApplyMappings(uidMappings []groot.IDMappingSpec, gidMappings []groot.IDMappingSpec) error {
	var uidMappingsCopy []groot.IDMappingSpec
	if uidMappings != nil {
		uidMappingsCopy = make([]groot.IDMappingSpec, len(uidMappings))
		copy(uidMappingsCopy, uidMappings)
	}
	var gidMappingsCopy []groot.IDMappingSpec
	if gidMappings != nil {
		gidMappingsCopy = make([]groot.IDMappingSpec, len(gidMappings))
		copy(gidMappingsCopy, gidMappings)
	}
	fake.applyMappingsMutex.Lock()
	ret, specificReturn := fake.applyMappingsReturnsOnCall[len(fake.applyMappingsArgsForCall)]
	fake.applyMappingsArgsForCall = append(fake.applyMappingsArgsForCall, struct {
		uidMappings []groot.IDMappingSpec
		gidMappings []groot.IDMappingSpec
	}{uidMappingsCopy, gidMappingsCopy})
	fake.recordInvocation("ApplyMappings", []interface{}{uidMappingsCopy, gidMappingsCopy})
	fake.applyMappingsMutex.Unlock()
	if fake.ApplyMappingsStub != nil {
		return fake.ApplyMappingsStub(uidMappings, gidMappings)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.applyMappingsReturns.result1
}

func (fake *FakeStoreNamespacer) ApplyMappingsCallCount() int {
	fake.applyMappingsMutex.RLock()
	defer fake.applyMappingsMutex.RUnlock()
	return len(fake.applyMappingsArgsForCall)
}

func (fake *FakeStoreNamespacer) ApplyMappingsArgsForCall(i int) ([]groot.IDMappingSpec, []groot.IDMappingSpec) {
	fake.applyMappingsMutex.RLock()
	defer fake.applyMappingsMutex.RUnlock()
	return fake.applyMappingsArgsForCall[i].uidMappings, fake.applyMappingsArgsForCall[i].gidMappings
}

func (fake *FakeStoreNamespacer) ApplyMappingsReturns(result1 error) {
	fake.ApplyMappingsStub = nil
	fake.applyMappingsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStoreNamespacer) ApplyMappingsReturnsOnCall(i int, result1 error) {
	fake.ApplyMappingsStub = nil
	if fake.applyMappingsReturnsOnCall == nil {
		fake.applyMappingsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.applyMappingsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeStoreNamespacer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.applyMappingsMutex.RLock()
	defer fake.applyMappingsMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeStoreNamespacer) recordInvocation(key string, args []interface{}) {
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

var _ manager.StoreNamespacer = new(FakeStoreNamespacer)
