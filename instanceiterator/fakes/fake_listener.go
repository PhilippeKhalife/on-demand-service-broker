// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"
	"time"

	"github.com/pivotal-cf/on-demand-service-broker/broker/services"
	"github.com/pivotal-cf/on-demand-service-broker/config"
	"github.com/pivotal-cf/on-demand-service-broker/instanceiterator"
	"github.com/pivotal-cf/on-demand-service-broker/service"
)

type FakeListener struct {
	FailedToRefreshInstanceInfoStub        func(instance string)
	failedToRefreshInstanceInfoMutex       sync.RWMutex
	failedToRefreshInstanceInfoArgsForCall []struct {
		instance string
	}
	StartingStub        func(maxInFlight int)
	startingMutex       sync.RWMutex
	startingArgsForCall []struct {
		maxInFlight int
	}
	RetryAttemptStub        func(num, limit int)
	retryAttemptMutex       sync.RWMutex
	retryAttemptArgsForCall []struct {
		num   int
		limit int
	}
	RetryCanariesAttemptStub        func(num, limit, remainingCanaries int)
	retryCanariesAttemptMutex       sync.RWMutex
	retryCanariesAttemptArgsForCall []struct {
		num               int
		limit             int
		remainingCanaries int
	}
	InstancesToProcessStub        func(instances []service.Instance)
	instancesToProcessMutex       sync.RWMutex
	instancesToProcessArgsForCall []struct {
		instances []service.Instance
	}
	InstanceOperationStartingStub        func(instance string, index int, totalInstances int, isCanary bool)
	instanceOperationStartingMutex       sync.RWMutex
	instanceOperationStartingArgsForCall []struct {
		instance       string
		index          int
		totalInstances int
		isCanary       bool
	}
	InstanceOperationStartResultStub        func(instance string, status services.BOSHOperationType)
	instanceOperationStartResultMutex       sync.RWMutex
	instanceOperationStartResultArgsForCall []struct {
		instance string
		status   services.BOSHOperationType
	}
	InstanceOperationFinishedStub        func(instance string, result string)
	instanceOperationFinishedMutex       sync.RWMutex
	instanceOperationFinishedArgsForCall []struct {
		instance string
		result   string
	}
	WaitingForStub        func(instance string, boshTaskId int)
	waitingForMutex       sync.RWMutex
	waitingForArgsForCall []struct {
		instance   string
		boshTaskId int
	}
	ProgressStub        func(pollingInterval time.Duration, orphanCount, processedCount, toRetryCount, deletedCount int)
	progressMutex       sync.RWMutex
	progressArgsForCall []struct {
		pollingInterval time.Duration
		orphanCount     int
		processedCount  int
		toRetryCount    int
		deletedCount    int
	}
	FinishedStub        func(orphanCount, finishedCount, deletedCount int, busyInstances, failedInstances []string)
	finishedMutex       sync.RWMutex
	finishedArgsForCall []struct {
		orphanCount     int
		finishedCount   int
		deletedCount    int
		busyInstances   []string
		failedInstances []string
	}
	CanariesStartingStub        func(canaries int, filter config.CanarySelectionParams)
	canariesStartingMutex       sync.RWMutex
	canariesStartingArgsForCall []struct {
		canaries int
		filter   config.CanarySelectionParams
	}
	CanariesFinishedStub        func()
	canariesFinishedMutex       sync.RWMutex
	canariesFinishedArgsForCall []struct{}
	invocations                 map[string][][]interface{}
	invocationsMutex            sync.RWMutex
}

func (fake *FakeListener) FailedToRefreshInstanceInfo(instance string) {
	fake.failedToRefreshInstanceInfoMutex.Lock()
	fake.failedToRefreshInstanceInfoArgsForCall = append(fake.failedToRefreshInstanceInfoArgsForCall, struct {
		instance string
	}{instance})
	fake.recordInvocation("FailedToRefreshInstanceInfo", []interface{}{instance})
	fake.failedToRefreshInstanceInfoMutex.Unlock()
	if fake.FailedToRefreshInstanceInfoStub != nil {
		fake.FailedToRefreshInstanceInfoStub(instance)
	}
}

func (fake *FakeListener) FailedToRefreshInstanceInfoCallCount() int {
	fake.failedToRefreshInstanceInfoMutex.RLock()
	defer fake.failedToRefreshInstanceInfoMutex.RUnlock()
	return len(fake.failedToRefreshInstanceInfoArgsForCall)
}

func (fake *FakeListener) FailedToRefreshInstanceInfoArgsForCall(i int) string {
	fake.failedToRefreshInstanceInfoMutex.RLock()
	defer fake.failedToRefreshInstanceInfoMutex.RUnlock()
	return fake.failedToRefreshInstanceInfoArgsForCall[i].instance
}

func (fake *FakeListener) Starting(maxInFlight int) {
	fake.startingMutex.Lock()
	fake.startingArgsForCall = append(fake.startingArgsForCall, struct {
		maxInFlight int
	}{maxInFlight})
	fake.recordInvocation("Starting", []interface{}{maxInFlight})
	fake.startingMutex.Unlock()
	if fake.StartingStub != nil {
		fake.StartingStub(maxInFlight)
	}
}

func (fake *FakeListener) StartingCallCount() int {
	fake.startingMutex.RLock()
	defer fake.startingMutex.RUnlock()
	return len(fake.startingArgsForCall)
}

func (fake *FakeListener) StartingArgsForCall(i int) int {
	fake.startingMutex.RLock()
	defer fake.startingMutex.RUnlock()
	return fake.startingArgsForCall[i].maxInFlight
}

func (fake *FakeListener) RetryAttempt(num int, limit int) {
	fake.retryAttemptMutex.Lock()
	fake.retryAttemptArgsForCall = append(fake.retryAttemptArgsForCall, struct {
		num   int
		limit int
	}{num, limit})
	fake.recordInvocation("RetryAttempt", []interface{}{num, limit})
	fake.retryAttemptMutex.Unlock()
	if fake.RetryAttemptStub != nil {
		fake.RetryAttemptStub(num, limit)
	}
}

func (fake *FakeListener) RetryAttemptCallCount() int {
	fake.retryAttemptMutex.RLock()
	defer fake.retryAttemptMutex.RUnlock()
	return len(fake.retryAttemptArgsForCall)
}

func (fake *FakeListener) RetryAttemptArgsForCall(i int) (int, int) {
	fake.retryAttemptMutex.RLock()
	defer fake.retryAttemptMutex.RUnlock()
	return fake.retryAttemptArgsForCall[i].num, fake.retryAttemptArgsForCall[i].limit
}

func (fake *FakeListener) RetryCanariesAttempt(num int, limit int, remainingCanaries int) {
	fake.retryCanariesAttemptMutex.Lock()
	fake.retryCanariesAttemptArgsForCall = append(fake.retryCanariesAttemptArgsForCall, struct {
		num               int
		limit             int
		remainingCanaries int
	}{num, limit, remainingCanaries})
	fake.recordInvocation("RetryCanariesAttempt", []interface{}{num, limit, remainingCanaries})
	fake.retryCanariesAttemptMutex.Unlock()
	if fake.RetryCanariesAttemptStub != nil {
		fake.RetryCanariesAttemptStub(num, limit, remainingCanaries)
	}
}

func (fake *FakeListener) RetryCanariesAttemptCallCount() int {
	fake.retryCanariesAttemptMutex.RLock()
	defer fake.retryCanariesAttemptMutex.RUnlock()
	return len(fake.retryCanariesAttemptArgsForCall)
}

func (fake *FakeListener) RetryCanariesAttemptArgsForCall(i int) (int, int, int) {
	fake.retryCanariesAttemptMutex.RLock()
	defer fake.retryCanariesAttemptMutex.RUnlock()
	return fake.retryCanariesAttemptArgsForCall[i].num, fake.retryCanariesAttemptArgsForCall[i].limit, fake.retryCanariesAttemptArgsForCall[i].remainingCanaries
}

func (fake *FakeListener) InstancesToProcess(instances []service.Instance) {
	var instancesCopy []service.Instance
	if instances != nil {
		instancesCopy = make([]service.Instance, len(instances))
		copy(instancesCopy, instances)
	}
	fake.instancesToProcessMutex.Lock()
	fake.instancesToProcessArgsForCall = append(fake.instancesToProcessArgsForCall, struct {
		instances []service.Instance
	}{instancesCopy})
	fake.recordInvocation("InstancesToProcess", []interface{}{instancesCopy})
	fake.instancesToProcessMutex.Unlock()
	if fake.InstancesToProcessStub != nil {
		fake.InstancesToProcessStub(instances)
	}
}

func (fake *FakeListener) InstancesToProcessCallCount() int {
	fake.instancesToProcessMutex.RLock()
	defer fake.instancesToProcessMutex.RUnlock()
	return len(fake.instancesToProcessArgsForCall)
}

func (fake *FakeListener) InstancesToProcessArgsForCall(i int) []service.Instance {
	fake.instancesToProcessMutex.RLock()
	defer fake.instancesToProcessMutex.RUnlock()
	return fake.instancesToProcessArgsForCall[i].instances
}

func (fake *FakeListener) InstanceOperationStarting(instance string, index int, totalInstances int, isCanary bool) {
	fake.instanceOperationStartingMutex.Lock()
	fake.instanceOperationStartingArgsForCall = append(fake.instanceOperationStartingArgsForCall, struct {
		instance       string
		index          int
		totalInstances int
		isCanary       bool
	}{instance, index, totalInstances, isCanary})
	fake.recordInvocation("InstanceOperationStarting", []interface{}{instance, index, totalInstances, isCanary})
	fake.instanceOperationStartingMutex.Unlock()
	if fake.InstanceOperationStartingStub != nil {
		fake.InstanceOperationStartingStub(instance, index, totalInstances, isCanary)
	}
}

func (fake *FakeListener) InstanceOperationStartingCallCount() int {
	fake.instanceOperationStartingMutex.RLock()
	defer fake.instanceOperationStartingMutex.RUnlock()
	return len(fake.instanceOperationStartingArgsForCall)
}

func (fake *FakeListener) InstanceOperationStartingArgsForCall(i int) (string, int, int, bool) {
	fake.instanceOperationStartingMutex.RLock()
	defer fake.instanceOperationStartingMutex.RUnlock()
	return fake.instanceOperationStartingArgsForCall[i].instance, fake.instanceOperationStartingArgsForCall[i].index, fake.instanceOperationStartingArgsForCall[i].totalInstances, fake.instanceOperationStartingArgsForCall[i].isCanary
}

func (fake *FakeListener) InstanceOperationStartResult(instance string, status services.BOSHOperationType) {
	fake.instanceOperationStartResultMutex.Lock()
	fake.instanceOperationStartResultArgsForCall = append(fake.instanceOperationStartResultArgsForCall, struct {
		instance string
		status   services.BOSHOperationType
	}{instance, status})
	fake.recordInvocation("InstanceOperationStartResult", []interface{}{instance, status})
	fake.instanceOperationStartResultMutex.Unlock()
	if fake.InstanceOperationStartResultStub != nil {
		fake.InstanceOperationStartResultStub(instance, status)
	}
}

func (fake *FakeListener) InstanceOperationStartResultCallCount() int {
	fake.instanceOperationStartResultMutex.RLock()
	defer fake.instanceOperationStartResultMutex.RUnlock()
	return len(fake.instanceOperationStartResultArgsForCall)
}

func (fake *FakeListener) InstanceOperationStartResultArgsForCall(i int) (string, services.BOSHOperationType) {
	fake.instanceOperationStartResultMutex.RLock()
	defer fake.instanceOperationStartResultMutex.RUnlock()
	return fake.instanceOperationStartResultArgsForCall[i].instance, fake.instanceOperationStartResultArgsForCall[i].status
}

func (fake *FakeListener) InstanceOperationFinished(instance string, result string) {
	fake.instanceOperationFinishedMutex.Lock()
	fake.instanceOperationFinishedArgsForCall = append(fake.instanceOperationFinishedArgsForCall, struct {
		instance string
		result   string
	}{instance, result})
	fake.recordInvocation("InstanceOperationFinished", []interface{}{instance, result})
	fake.instanceOperationFinishedMutex.Unlock()
	if fake.InstanceOperationFinishedStub != nil {
		fake.InstanceOperationFinishedStub(instance, result)
	}
}

func (fake *FakeListener) InstanceOperationFinishedCallCount() int {
	fake.instanceOperationFinishedMutex.RLock()
	defer fake.instanceOperationFinishedMutex.RUnlock()
	return len(fake.instanceOperationFinishedArgsForCall)
}

func (fake *FakeListener) InstanceOperationFinishedArgsForCall(i int) (string, string) {
	fake.instanceOperationFinishedMutex.RLock()
	defer fake.instanceOperationFinishedMutex.RUnlock()
	return fake.instanceOperationFinishedArgsForCall[i].instance, fake.instanceOperationFinishedArgsForCall[i].result
}

func (fake *FakeListener) WaitingFor(instance string, boshTaskId int) {
	fake.waitingForMutex.Lock()
	fake.waitingForArgsForCall = append(fake.waitingForArgsForCall, struct {
		instance   string
		boshTaskId int
	}{instance, boshTaskId})
	fake.recordInvocation("WaitingFor", []interface{}{instance, boshTaskId})
	fake.waitingForMutex.Unlock()
	if fake.WaitingForStub != nil {
		fake.WaitingForStub(instance, boshTaskId)
	}
}

func (fake *FakeListener) WaitingForCallCount() int {
	fake.waitingForMutex.RLock()
	defer fake.waitingForMutex.RUnlock()
	return len(fake.waitingForArgsForCall)
}

func (fake *FakeListener) WaitingForArgsForCall(i int) (string, int) {
	fake.waitingForMutex.RLock()
	defer fake.waitingForMutex.RUnlock()
	return fake.waitingForArgsForCall[i].instance, fake.waitingForArgsForCall[i].boshTaskId
}

func (fake *FakeListener) Progress(pollingInterval time.Duration, orphanCount int, processedCount int, toRetryCount int, deletedCount int) {
	fake.progressMutex.Lock()
	fake.progressArgsForCall = append(fake.progressArgsForCall, struct {
		pollingInterval time.Duration
		orphanCount     int
		processedCount  int
		toRetryCount    int
		deletedCount    int
	}{pollingInterval, orphanCount, processedCount, toRetryCount, deletedCount})
	fake.recordInvocation("Progress", []interface{}{pollingInterval, orphanCount, processedCount, toRetryCount, deletedCount})
	fake.progressMutex.Unlock()
	if fake.ProgressStub != nil {
		fake.ProgressStub(pollingInterval, orphanCount, processedCount, toRetryCount, deletedCount)
	}
}

func (fake *FakeListener) ProgressCallCount() int {
	fake.progressMutex.RLock()
	defer fake.progressMutex.RUnlock()
	return len(fake.progressArgsForCall)
}

func (fake *FakeListener) ProgressArgsForCall(i int) (time.Duration, int, int, int, int) {
	fake.progressMutex.RLock()
	defer fake.progressMutex.RUnlock()
	return fake.progressArgsForCall[i].pollingInterval, fake.progressArgsForCall[i].orphanCount, fake.progressArgsForCall[i].processedCount, fake.progressArgsForCall[i].toRetryCount, fake.progressArgsForCall[i].deletedCount
}

func (fake *FakeListener) Finished(orphanCount int, finishedCount int, deletedCount int, busyInstances []string, failedInstances []string) {
	var busyInstancesCopy []string
	if busyInstances != nil {
		busyInstancesCopy = make([]string, len(busyInstances))
		copy(busyInstancesCopy, busyInstances)
	}
	var failedInstancesCopy []string
	if failedInstances != nil {
		failedInstancesCopy = make([]string, len(failedInstances))
		copy(failedInstancesCopy, failedInstances)
	}
	fake.finishedMutex.Lock()
	fake.finishedArgsForCall = append(fake.finishedArgsForCall, struct {
		orphanCount     int
		finishedCount   int
		deletedCount    int
		busyInstances   []string
		failedInstances []string
	}{orphanCount, finishedCount, deletedCount, busyInstancesCopy, failedInstancesCopy})
	fake.recordInvocation("Finished", []interface{}{orphanCount, finishedCount, deletedCount, busyInstancesCopy, failedInstancesCopy})
	fake.finishedMutex.Unlock()
	if fake.FinishedStub != nil {
		fake.FinishedStub(orphanCount, finishedCount, deletedCount, busyInstances, failedInstances)
	}
}

func (fake *FakeListener) FinishedCallCount() int {
	fake.finishedMutex.RLock()
	defer fake.finishedMutex.RUnlock()
	return len(fake.finishedArgsForCall)
}

func (fake *FakeListener) FinishedArgsForCall(i int) (int, int, int, []string, []string) {
	fake.finishedMutex.RLock()
	defer fake.finishedMutex.RUnlock()
	return fake.finishedArgsForCall[i].orphanCount, fake.finishedArgsForCall[i].finishedCount, fake.finishedArgsForCall[i].deletedCount, fake.finishedArgsForCall[i].busyInstances, fake.finishedArgsForCall[i].failedInstances
}

func (fake *FakeListener) CanariesStarting(canaries int, filter config.CanarySelectionParams) {
	fake.canariesStartingMutex.Lock()
	fake.canariesStartingArgsForCall = append(fake.canariesStartingArgsForCall, struct {
		canaries int
		filter   config.CanarySelectionParams
	}{canaries, filter})
	fake.recordInvocation("CanariesStarting", []interface{}{canaries, filter})
	fake.canariesStartingMutex.Unlock()
	if fake.CanariesStartingStub != nil {
		fake.CanariesStartingStub(canaries, filter)
	}
}

func (fake *FakeListener) CanariesStartingCallCount() int {
	fake.canariesStartingMutex.RLock()
	defer fake.canariesStartingMutex.RUnlock()
	return len(fake.canariesStartingArgsForCall)
}

func (fake *FakeListener) CanariesStartingArgsForCall(i int) (int, config.CanarySelectionParams) {
	fake.canariesStartingMutex.RLock()
	defer fake.canariesStartingMutex.RUnlock()
	return fake.canariesStartingArgsForCall[i].canaries, fake.canariesStartingArgsForCall[i].filter
}

func (fake *FakeListener) CanariesFinished() {
	fake.canariesFinishedMutex.Lock()
	fake.canariesFinishedArgsForCall = append(fake.canariesFinishedArgsForCall, struct{}{})
	fake.recordInvocation("CanariesFinished", []interface{}{})
	fake.canariesFinishedMutex.Unlock()
	if fake.CanariesFinishedStub != nil {
		fake.CanariesFinishedStub()
	}
}

func (fake *FakeListener) CanariesFinishedCallCount() int {
	fake.canariesFinishedMutex.RLock()
	defer fake.canariesFinishedMutex.RUnlock()
	return len(fake.canariesFinishedArgsForCall)
}

func (fake *FakeListener) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.failedToRefreshInstanceInfoMutex.RLock()
	defer fake.failedToRefreshInstanceInfoMutex.RUnlock()
	fake.startingMutex.RLock()
	defer fake.startingMutex.RUnlock()
	fake.retryAttemptMutex.RLock()
	defer fake.retryAttemptMutex.RUnlock()
	fake.retryCanariesAttemptMutex.RLock()
	defer fake.retryCanariesAttemptMutex.RUnlock()
	fake.instancesToProcessMutex.RLock()
	defer fake.instancesToProcessMutex.RUnlock()
	fake.instanceOperationStartingMutex.RLock()
	defer fake.instanceOperationStartingMutex.RUnlock()
	fake.instanceOperationStartResultMutex.RLock()
	defer fake.instanceOperationStartResultMutex.RUnlock()
	fake.instanceOperationFinishedMutex.RLock()
	defer fake.instanceOperationFinishedMutex.RUnlock()
	fake.waitingForMutex.RLock()
	defer fake.waitingForMutex.RUnlock()
	fake.progressMutex.RLock()
	defer fake.progressMutex.RUnlock()
	fake.finishedMutex.RLock()
	defer fake.finishedMutex.RUnlock()
	fake.canariesStartingMutex.RLock()
	defer fake.canariesStartingMutex.RUnlock()
	fake.canariesFinishedMutex.RLock()
	defer fake.canariesFinishedMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeListener) recordInvocation(key string, args []interface{}) {
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

var _ instanceiterator.Listener = new(FakeListener)
