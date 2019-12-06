package fakes

import "sync"

type Decompressor struct {
	UnTarCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Destination string
		}
		Returns struct {
			Error error
		}
		Stub func(string) error
	}
}

func (f *Decompressor) UnTar(param1 string) error {
	f.UnTarCall.Lock()
	defer f.UnTarCall.Unlock()
	f.UnTarCall.CallCount++
	f.UnTarCall.Receives.Destination = param1
	if f.UnTarCall.Stub != nil {
		return f.UnTarCall.Stub(param1)
	}
	return f.UnTarCall.Returns.Error
}
