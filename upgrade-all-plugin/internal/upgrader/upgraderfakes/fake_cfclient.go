// Code generated by counterfeiter. DO NOT EDIT.
package upgraderfakes

import (
	"sync"

	"github.com/cloudfoundry/cloud-service-broker/upgrade-all-plugin/internal/ccapi"
	"github.com/cloudfoundry/cloud-service-broker/upgrade-all-plugin/internal/upgrader"
)

type FakeCFClient struct {
	GetServiceInstancesStub        func([]string) ([]ccapi.ServiceInstance, error)
	getServiceInstancesMutex       sync.RWMutex
	getServiceInstancesArgsForCall []struct {
		arg1 []string
	}
	getServiceInstancesReturns struct {
		result1 []ccapi.ServiceInstance
		result2 error
	}
	getServiceInstancesReturnsOnCall map[int]struct {
		result1 []ccapi.ServiceInstance
		result2 error
	}
	GetServicePlansStub        func(string) ([]ccapi.Plan, error)
	getServicePlansMutex       sync.RWMutex
	getServicePlansArgsForCall []struct {
		arg1 string
	}
	getServicePlansReturns struct {
		result1 []ccapi.Plan
		result2 error
	}
	getServicePlansReturnsOnCall map[int]struct {
		result1 []ccapi.Plan
		result2 error
	}
	UpgradeServiceInstanceStub        func(string, string) error
	upgradeServiceInstanceMutex       sync.RWMutex
	upgradeServiceInstanceArgsForCall []struct {
		arg1 string
		arg2 string
	}
	upgradeServiceInstanceReturns struct {
		result1 error
	}
	upgradeServiceInstanceReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCFClient) GetServiceInstances(arg1 []string) ([]ccapi.ServiceInstance, error) {
	var arg1Copy []string
	if arg1 != nil {
		arg1Copy = make([]string, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.getServiceInstancesMutex.Lock()
	ret, specificReturn := fake.getServiceInstancesReturnsOnCall[len(fake.getServiceInstancesArgsForCall)]
	fake.getServiceInstancesArgsForCall = append(fake.getServiceInstancesArgsForCall, struct {
		arg1 []string
	}{arg1Copy})
	stub := fake.GetServiceInstancesStub
	fakeReturns := fake.getServiceInstancesReturns
	fake.recordInvocation("GetServiceInstances", []interface{}{arg1Copy})
	fake.getServiceInstancesMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFClient) GetServiceInstancesCallCount() int {
	fake.getServiceInstancesMutex.RLock()
	defer fake.getServiceInstancesMutex.RUnlock()
	return len(fake.getServiceInstancesArgsForCall)
}

func (fake *FakeCFClient) GetServiceInstancesCalls(stub func([]string) ([]ccapi.ServiceInstance, error)) {
	fake.getServiceInstancesMutex.Lock()
	defer fake.getServiceInstancesMutex.Unlock()
	fake.GetServiceInstancesStub = stub
}

func (fake *FakeCFClient) GetServiceInstancesArgsForCall(i int) []string {
	fake.getServiceInstancesMutex.RLock()
	defer fake.getServiceInstancesMutex.RUnlock()
	argsForCall := fake.getServiceInstancesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCFClient) GetServiceInstancesReturns(result1 []ccapi.ServiceInstance, result2 error) {
	fake.getServiceInstancesMutex.Lock()
	defer fake.getServiceInstancesMutex.Unlock()
	fake.GetServiceInstancesStub = nil
	fake.getServiceInstancesReturns = struct {
		result1 []ccapi.ServiceInstance
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) GetServiceInstancesReturnsOnCall(i int, result1 []ccapi.ServiceInstance, result2 error) {
	fake.getServiceInstancesMutex.Lock()
	defer fake.getServiceInstancesMutex.Unlock()
	fake.GetServiceInstancesStub = nil
	if fake.getServiceInstancesReturnsOnCall == nil {
		fake.getServiceInstancesReturnsOnCall = make(map[int]struct {
			result1 []ccapi.ServiceInstance
			result2 error
		})
	}
	fake.getServiceInstancesReturnsOnCall[i] = struct {
		result1 []ccapi.ServiceInstance
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) GetServicePlans(arg1 string) ([]ccapi.Plan, error) {
	fake.getServicePlansMutex.Lock()
	ret, specificReturn := fake.getServicePlansReturnsOnCall[len(fake.getServicePlansArgsForCall)]
	fake.getServicePlansArgsForCall = append(fake.getServicePlansArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetServicePlansStub
	fakeReturns := fake.getServicePlansReturns
	fake.recordInvocation("GetServicePlans", []interface{}{arg1})
	fake.getServicePlansMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCFClient) GetServicePlansCallCount() int {
	fake.getServicePlansMutex.RLock()
	defer fake.getServicePlansMutex.RUnlock()
	return len(fake.getServicePlansArgsForCall)
}

func (fake *FakeCFClient) GetServicePlansCalls(stub func(string) ([]ccapi.Plan, error)) {
	fake.getServicePlansMutex.Lock()
	defer fake.getServicePlansMutex.Unlock()
	fake.GetServicePlansStub = stub
}

func (fake *FakeCFClient) GetServicePlansArgsForCall(i int) string {
	fake.getServicePlansMutex.RLock()
	defer fake.getServicePlansMutex.RUnlock()
	argsForCall := fake.getServicePlansArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCFClient) GetServicePlansReturns(result1 []ccapi.Plan, result2 error) {
	fake.getServicePlansMutex.Lock()
	defer fake.getServicePlansMutex.Unlock()
	fake.GetServicePlansStub = nil
	fake.getServicePlansReturns = struct {
		result1 []ccapi.Plan
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) GetServicePlansReturnsOnCall(i int, result1 []ccapi.Plan, result2 error) {
	fake.getServicePlansMutex.Lock()
	defer fake.getServicePlansMutex.Unlock()
	fake.GetServicePlansStub = nil
	if fake.getServicePlansReturnsOnCall == nil {
		fake.getServicePlansReturnsOnCall = make(map[int]struct {
			result1 []ccapi.Plan
			result2 error
		})
	}
	fake.getServicePlansReturnsOnCall[i] = struct {
		result1 []ccapi.Plan
		result2 error
	}{result1, result2}
}

func (fake *FakeCFClient) UpgradeServiceInstance(arg1 string, arg2 string) error {
	fake.upgradeServiceInstanceMutex.Lock()
	ret, specificReturn := fake.upgradeServiceInstanceReturnsOnCall[len(fake.upgradeServiceInstanceArgsForCall)]
	fake.upgradeServiceInstanceArgsForCall = append(fake.upgradeServiceInstanceArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.UpgradeServiceInstanceStub
	fakeReturns := fake.upgradeServiceInstanceReturns
	fake.recordInvocation("UpgradeServiceInstance", []interface{}{arg1, arg2})
	fake.upgradeServiceInstanceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeCFClient) UpgradeServiceInstanceCallCount() int {
	fake.upgradeServiceInstanceMutex.RLock()
	defer fake.upgradeServiceInstanceMutex.RUnlock()
	return len(fake.upgradeServiceInstanceArgsForCall)
}

func (fake *FakeCFClient) UpgradeServiceInstanceCalls(stub func(string, string) error) {
	fake.upgradeServiceInstanceMutex.Lock()
	defer fake.upgradeServiceInstanceMutex.Unlock()
	fake.UpgradeServiceInstanceStub = stub
}

func (fake *FakeCFClient) UpgradeServiceInstanceArgsForCall(i int) (string, string) {
	fake.upgradeServiceInstanceMutex.RLock()
	defer fake.upgradeServiceInstanceMutex.RUnlock()
	argsForCall := fake.upgradeServiceInstanceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCFClient) UpgradeServiceInstanceReturns(result1 error) {
	fake.upgradeServiceInstanceMutex.Lock()
	defer fake.upgradeServiceInstanceMutex.Unlock()
	fake.UpgradeServiceInstanceStub = nil
	fake.upgradeServiceInstanceReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCFClient) UpgradeServiceInstanceReturnsOnCall(i int, result1 error) {
	fake.upgradeServiceInstanceMutex.Lock()
	defer fake.upgradeServiceInstanceMutex.Unlock()
	fake.UpgradeServiceInstanceStub = nil
	if fake.upgradeServiceInstanceReturnsOnCall == nil {
		fake.upgradeServiceInstanceReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.upgradeServiceInstanceReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeCFClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getServiceInstancesMutex.RLock()
	defer fake.getServiceInstancesMutex.RUnlock()
	fake.getServicePlansMutex.RLock()
	defer fake.getServicePlansMutex.RUnlock()
	fake.upgradeServiceInstanceMutex.RLock()
	defer fake.upgradeServiceInstanceMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCFClient) recordInvocation(key string, args []interface{}) {
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

var _ upgrader.CFClient = new(FakeCFClient)
