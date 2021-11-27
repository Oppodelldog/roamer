package sound

import (
	"path/filepath"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

type Session struct {
	asc  *wca.IAudioSessionControl
	asc2 *wca.IAudioSessionControl2
	sv   *wca.ISimpleAudioVolume
}

func (s *Session) Release() {
	s.sv.Release()
	s.asc2.Release()
	s.asc.Release()
}

func (s *Session) GetDisplayName() (string, error) {
	var (
		v   string
		err = s.asc2.GetDisplayName(&v)
	)

	return v, err
}

func (s *Session) GetDisplayNameEnhanced() (string, error) {
	if s.IsSystemSoundsSession() {
		pid, err := s.GetProcessId()
		if err != nil {
			return "", err
		}

		v := ProcessName(int(pid))
		v = strings.TrimSuffix(v, filepath.Ext(v))

		if len(v) > 0 {
			v = strings.Title(v)
		}

		return v, nil
	}

	v, err := s.GetDisplayName()
	if err != nil {
		return "", err
	}

	if v == "@%SystemRoot%\\System32\\AudioSrv.Dll,-202" {
		return "Systemsounds", nil
	}

	v = strings.TrimSuffix(v, filepath.Ext(v))

	return v, err
}

func (s *Session) GetGroupingParam() (ole.GUID, error) {
	var (
		v   ole.GUID
		err = s.asc2.GetGroupingParam(&v)
	)

	return v, err
}

func (s *Session) GetIconPath() (string, error) {
	var (
		v   string
		err = s.asc2.GetIconPath(&v)
	)

	return v, err
}

func (s *Session) GetMasterVolume() (float32, error) {
	var (
		v   float32
		err = s.sv.GetMasterVolume(&v)
	)

	return v, err
}

func (s *Session) GetMute() (bool, error) {
	var (
		v   bool
		err = s.sv.GetMute(&v)
	)

	return v, err
}

func (s *Session) GetSessionIdentifier() (string, error) {
	var (
		v   string
		err = s.asc2.GetSessionIdentifier(&v)
	)

	return v, err
}

func (s *Session) GetSessionInstanceIdentifier() (string, error) {
	var (
		v   string
		err = s.asc2.GetSessionInstanceIdentifier(&v)
	)

	return v, err
}

func (s *Session) GetProcessId() (uint32, error) {
	var (
		v   uint32
		err = s.asc2.GetProcessId(&v)
	)

	return v, err
}

func (s *Session) GetState() (uint32, error) {
	var (
		v   uint32
		err = s.asc2.GetState(&v)
	)

	return v, err
}

func (s *Session) IsSystemSoundsSession() bool {
	return s.asc2.IsSystemSoundsSession() != nil
}

func (s *Session) SetDisplayName(v string) error {
	return s.asc2.SetDisplayName(&v, nil)
}

func (s *Session) SetDuckingPreference(v bool) error {
	return s.asc2.SetDuckingPreference(v)
}

func (s *Session) SetGroupingParam(v ole.GUID) error {
	return s.asc2.SetGroupingParam(&v, nil)
}

func (s *Session) SetIconPath(v string) error {
	return s.asc2.SetIconPath(&v, nil)
}

func (s *Session) SetMasterVolume(v float32) error {
	return s.sv.SetMasterVolume(v, nil)
}

func (s *Session) SetMute(v bool) error {
	return s.sv.SetMute(v, nil)
}
