package audioctl

import "github.com/moutend/go-wca/pkg/wca"

func RegisterDeviceEvents(callback wca.IMMNotificationClientCallback) error {
	var (
		mmde *wca.IMMDeviceEnumerator
		mmd  *wca.IMMDevice
		err  error
	)

	if err = wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return err
	}
	defer mmde.Release()

	if err = mmde.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &mmd); err != nil {
		return err
	}

	return mmde.RegisterEndpointNotificationCallback(wca.NewIMMNotificationClient(callback))
}
