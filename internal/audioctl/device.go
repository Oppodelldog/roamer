package sound

import (
	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

type Device struct {
	Sessions []*Session
	mmd      *wca.IMMDevice
	aev      *wca.IAudioEndpointVolume
	ps       *wca.IPropertyStore
}

func (d *Device) Release() {
	for i := range d.Sessions {
		d.Sessions[i].Release()
	}

	d.ps.Release()
	d.aev.Release()
	d.mmd.Release()
	ole.CoUninitialize()
}

func (d *Device) GetId() (string, error) {
	var (
		v   string
		err = d.mmd.GetId(&v)
	)

	return v, err
}

func (d *Device) GetDeviceFriendlyName() (string, error) {
	var (
		pv  wca.PROPVARIANT
		err = d.ps.GetValue(&wca.PKEY_Device_FriendlyName, &pv)
	)

	return pv.String(), err
}

func (d *Device) GetMasterVolume() (float32, error) {
	var (
		v   float32
		err = d.aev.GetMasterVolumeLevel(&v)
	)

	return v, err
}

func (d *Device) GetMasterVolumeLevelScalar() (float32, error) {
	var (
		v   float32
		err = d.aev.GetMasterVolumeLevelScalar(&v)
	)

	return v, err
}

func (d *Device) GetChannelCount() (uint32, error) {
	var (
		v   uint32
		err = d.aev.GetChannelCount(&v)
	)

	return v, err
}

func (d *Device) GetMute() (bool, error) {
	var (
		v   bool
		err = d.aev.GetMute(&v)
	)

	return v, err
}

func (d *Device) GetVolumeRange() (float32, float32, float32, error) {
	var (
		minDB       float32
		maxDB       float32
		incrementDB float32
		err         = d.aev.GetVolumeRange(&minDB, &maxDB, &incrementDB)
	)

	return minDB, maxDB, incrementDB, err
}

func (d *Device) SetMasterVolume(v float32) error {
	return d.aev.SetMasterVolumeLevel(v, nil)
}

func (d *Device) SetMasterVolumeLevelScalar(v float32) error {
	return d.aev.SetMasterVolumeLevelScalar(v, nil)
}

func (d *Device) SetMute(v bool) error {
	return d.aev.SetMute(v, nil)
}
