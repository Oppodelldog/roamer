package sound

import (
	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

func DefaultDevice() (*Device, error) {
	var (
		device = &Device{}
		err    error
		mmde   *wca.IMMDeviceEnumerator
		sm     *wca.IAudioSessionManager2
	)

	if err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return nil, err
	}

	if err = wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return nil, err
	}
	defer mmde.Release()

	if err = mmde.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &device.mmd); err != nil {
		return nil, err
	}

	if err = device.mmd.OpenPropertyStore(wca.STGM_READ, &device.ps); err != nil {
		return nil, err
	}

	if err = device.mmd.Activate(wca.IID_IAudioEndpointVolume, wca.CLSCTX_ALL, nil, &device.aev); err != nil {
		return nil, err
	}

	err = device.mmd.Activate(wca.IID_IAudioSessionManager2, wca.CLSCTX_ALL, nil, &sm)
	if err != nil {
		return nil, err
	}
	defer sm.Release()

	device.Sessions, err = sessions(sm)
	if err != nil {
		return nil, err
	}

	return device, nil
}
func sessions(sm *wca.IAudioSessionManager2) ([]*Session, error) {
	var (
		err       error
		sessions  []*Session
		se        *wca.IAudioSessionEnumerator
		sessCount int
	)

	err = sm.GetSessionEnumerator(&se)
	if err != nil {
		return nil, err
	}

	err = se.GetCount(&sessCount)
	if err != nil {
		return nil, err
	}

	for sessId := 0; sessId < sessCount; sessId++ {
		s, errSess := session(se, sessId)
		if errSess != nil {
			return nil, errSess
		}

		sm.Release()

		sessions = append(sessions, s)
	}

	return sessions, err
}

func session(se *wca.IAudioSessionEnumerator, sessId int) (*Session, error) {
	var err error
	var s = Session{}

	err = se.GetSession(sessId, &s.asc)
	if err != nil {
		return nil, err
	}

	err = s.asc.PutQueryInterface(wca.IID_IAudioSessionControl2, &s.asc2)
	if err != nil {
		return nil, err
	}

	err = s.asc.PutQueryInterface(wca.IID_ISimpleAudioVolume, &s.sv)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
