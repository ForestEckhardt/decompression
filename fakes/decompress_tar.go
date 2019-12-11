package fakes

import (
	"io"
	"sync"
)

type DecompressTar struct {
	GetReaderCall struct {
		sync.Mutex
		CallCount int
		Returns   struct {
			Reader io.Reader
		}
		Stub func() io.Reader
	}
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

func (f *DecompressTar) GetReader() io.Reader {
	f.GetReaderCall.Lock()
	defer f.GetReaderCall.Unlock()
	f.GetReaderCall.CallCount++
	if f.GetReaderCall.Stub != nil {
		return f.GetReaderCall.Stub()
	}
	return f.GetReaderCall.Returns.Reader
}
func (f *DecompressTar) UnTar(param1 string) error {
	f.UnTarCall.Lock()
	defer f.UnTarCall.Unlock()
	f.UnTarCall.CallCount++
	f.UnTarCall.Receives.Destination = param1
	if f.UnTarCall.Stub != nil {
		return f.UnTarCall.Stub(param1)
	}
	return f.UnTarCall.Returns.Error
}
