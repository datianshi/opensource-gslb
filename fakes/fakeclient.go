// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"
	"time"

	"github.com/datianshi/simple-cf-gtm"
	"github.com/miekg/dns"
)

type FakeDNSClient struct {
	ExchangeStub        func(m *dns.Msg, address string) (r *dns.Msg, rtt time.Duration, err error)
	exchangeMutex       sync.RWMutex
	exchangeArgsForCall []struct {
		m       *dns.Msg
		address string
	}
	exchangeReturns struct {
		result1 *dns.Msg
		result2 time.Duration
		result3 error
	}
	exchangeReturnsOnCall map[int]struct {
		result1 *dns.Msg
		result2 time.Duration
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeDNSClient) Exchange(m *dns.Msg, address string) (r *dns.Msg, rtt time.Duration, err error) {
	fake.exchangeMutex.Lock()
	ret, specificReturn := fake.exchangeReturnsOnCall[len(fake.exchangeArgsForCall)]
	fake.exchangeArgsForCall = append(fake.exchangeArgsForCall, struct {
		m       *dns.Msg
		address string
	}{m, address})
	fake.recordInvocation("Exchange", []interface{}{m, address})
	fake.exchangeMutex.Unlock()
	if fake.ExchangeStub != nil {
		return fake.ExchangeStub(m, address)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.exchangeReturns.result1, fake.exchangeReturns.result2, fake.exchangeReturns.result3
}

func (fake *FakeDNSClient) ExchangeCallCount() int {
	fake.exchangeMutex.RLock()
	defer fake.exchangeMutex.RUnlock()
	return len(fake.exchangeArgsForCall)
}

func (fake *FakeDNSClient) ExchangeArgsForCall(i int) (*dns.Msg, string) {
	fake.exchangeMutex.RLock()
	defer fake.exchangeMutex.RUnlock()
	return fake.exchangeArgsForCall[i].m, fake.exchangeArgsForCall[i].address
}

func (fake *FakeDNSClient) ExchangeReturns(result1 *dns.Msg, result2 time.Duration, result3 error) {
	fake.ExchangeStub = nil
	fake.exchangeReturns = struct {
		result1 *dns.Msg
		result2 time.Duration
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeDNSClient) ExchangeReturnsOnCall(i int, result1 *dns.Msg, result2 time.Duration, result3 error) {
	fake.ExchangeStub = nil
	if fake.exchangeReturnsOnCall == nil {
		fake.exchangeReturnsOnCall = make(map[int]struct {
			result1 *dns.Msg
			result2 time.Duration
			result3 error
		})
	}
	fake.exchangeReturnsOnCall[i] = struct {
		result1 *dns.Msg
		result2 time.Duration
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeDNSClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.exchangeMutex.RLock()
	defer fake.exchangeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeDNSClient) recordInvocation(key string, args []interface{}) {
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

var _ gtm.DNSClient = new(FakeDNSClient)