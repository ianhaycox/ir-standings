// Package irsdk originally from https://github.com/quimcalpe/iracing-sdk and https://github.com/Sj-Si/iracing-sdk and
package irsdk

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/hidez8891/shm"
	"github.com/ianhaycox/ir-standings/irsdk/events"
	"github.com/ianhaycox/ir-standings/irsdk/iryaml"
)

const (
	exportFileMode = 0600
)

type SDK interface {
	RefreshSession()
	WaitForData(timeout time.Duration) bool
	GetVars() (map[string]Variable, error)
	GetVar(name string) (Variable, error)
	GetVarValue(name string) (interface{}, error)
	GetVarValues(name string) (interface{}, error)
	GetSession() iryaml.IRSession
	SessionChanged() bool
	GetLastVersion() int
	IsConnected() bool
	ExportIbtTo(fileName string)
	ExportSessionTo(fileName string)
	GetYaml() string
	BroadcastMsg(msg Msg)
	Close()
}

// IRSDK is the main SDK object clients must use
type IRSDK struct {
	SDK
	r                 reader
	h                 *header
	session           iryaml.IRSession
	s                 []string
	tVars             *TelemetryVars
	lastValidData     int64
	lastSessionUpdate int
}

func (sdk *IRSDK) RefreshSession() {
	if sdk.SessionChanged() {
		sRaw := readSessionData(sdk.r, sdk.h)

		err := yaml.Unmarshal([]byte(sRaw), &sdk.session)
		if err != nil {
			log.Println(err)
		}

		sdk.s = strings.Split(sRaw, "\n")
	}
}

func (sdk *IRSDK) WaitForData(timeout time.Duration) bool {
	if !sdk.IsConnected() {
		initIRSDK(sdk)
	}

	if events.WaitForSingleObject(timeout) {
		sdk.RefreshSession()
		return readVariableValues(sdk)
	}

	return false
}

func (sdk *IRSDK) GetVars() (map[string]Variable, error) {
	results := make(map[string]Variable, 0)

	if !sessionStatusOK(sdk.h.status) {
		return results, fmt.Errorf("session is not active")
	}

	sdk.tVars.mux.Lock()
	defer sdk.tVars.mux.Unlock()

	for _, variable := range sdk.tVars.vars {
		results[variable.Name] = variable
	}

	return results, nil
}

func (sdk *IRSDK) GetVar(name string) (Variable, error) {
	if !sessionStatusOK(sdk.h.status) {
		return Variable{}, fmt.Errorf("session is not active")
	}

	sdk.tVars.mux.Lock()
	defer sdk.tVars.mux.Unlock()

	if v, ok := sdk.tVars.vars[name]; ok {
		return v, nil
	}

	return Variable{}, fmt.Errorf("telemetry variable %q not found", name)
}

func (sdk *IRSDK) GetVarValue(name string) (interface{}, error) {
	var (
		r   Variable
		err error
	)

	if r, err = sdk.GetVar(name); err == nil {
		return r.Value, nil
	}

	return r, err
}

func (sdk *IRSDK) GetVarValues(name string) (interface{}, error) {
	var (
		r   Variable
		err error
	)

	if r, err = sdk.GetVar(name); err == nil {
		return r.Values, nil
	}

	return r, err
}

func (sdk *IRSDK) GetSession() iryaml.IRSession {
	return sdk.session
}

func (sdk *IRSDK) SessionChanged() bool {
	log.Println("session status:", sdk.h.status)

	if !sessionStatusOK(sdk.h.status) {
		log.Println("session status not ok")

		return false
	}

	if sdk.lastSessionUpdate != sdk.h.sessionInfoUpdate {
		log.Println("Session changed", sdk.lastSessionUpdate, sdk.h.sessionInfoUpdate)
		sdk.lastSessionUpdate = sdk.h.sessionInfoUpdate

		return true
	}

	log.Println("Session ", sdk.lastSessionUpdate, "head session", sdk.h.sessionInfoUpdate)

	return false
}

func (sdk *IRSDK) GetLastVersion() int {
	if !sessionStatusOK(sdk.h.status) {
		return -1
	}

	sdk.tVars.mux.Lock()
	defer sdk.tVars.mux.Unlock()

	last := sdk.tVars.lastVersion

	return last
}

func (sdk *IRSDK) GetSessionData(path string) (string, error) {
	if !sessionStatusOK(sdk.h.status) {
		return "", fmt.Errorf("session not connected")
	}

	return getSessionDataPath(sdk.s, path)
}

func (sdk *IRSDK) IsConnected() bool {
	if sdk.h != nil {
		if sessionStatusOK(sdk.h.status) && (sdk.lastValidData+connTimeout > time.Now().Unix()) {
			return true
		}
	}

	return false
}

// ExportIbtTo exports current memory data to a file
func (sdk *IRSDK) ExportIbtTo(fileName string) {
	rbuf := make([]byte, fileMapSize)

	_, err := sdk.r.ReadAt(rbuf, 0)
	if err != nil {
		log.Fatal(err)
	}

	_ = os.WriteFile(fileName, rbuf, exportFileMode)
}

// ExportSessionTo exports current session yaml data to a file
func (sdk *IRSDK) ExportSessionTo(fileName string) {
	y := strings.Join(sdk.s, "\n")

	_ = os.WriteFile(fileName, []byte(y), exportFileMode)
}

func (sdk *IRSDK) GetYaml() string {
	return strings.Join(sdk.s, "\n")
}

func (sdk *IRSDK) BroadcastMsg(msg Msg) {
	if msg.P2 == nil {
		msg.P2 = 0
	}

	events.BroadcastMsg(broadcastMsgName, msg.Cmd, msg.P1, msg.P2, msg.P3)
}

// Close clean up sdk resources
func (sdk *IRSDK) Close() {
	_ = sdk.r.Close()
}

// Init creates a SDK instance to operate with
func Init(r reader) SDK {
	if r == nil {
		var err error

		r, err = shm.Open(fileMapName, fileMapSize)
		if err != nil {
			log.Fatal(err)
		}
	}

	sdk := &IRSDK{r: r, lastValidData: 0}

	events.OpenEvent(dataValidEventName)
	initIRSDK(sdk)

	return sdk
}

func initIRSDK(sdk *IRSDK) {
	h := readHeader(sdk.r)
	sdk.h = &h
	sdk.s = nil
	sdk.lastSessionUpdate = -1

	if sdk.tVars != nil {
		sdk.tVars.vars = nil
	}

	if sessionStatusOK(h.status) {
		sRaw := readSessionData(sdk.r, &h)

		err := yaml.Unmarshal([]byte(sRaw), &sdk.session)
		if err != nil {
			log.Println(err)
		}

		sdk.s = strings.Split(sRaw, "\n")
		sdk.tVars = readVariableHeaders(sdk.r, &h)

		readVariableValues(sdk)
	}
}

func sessionStatusOK(status int) bool {
	return (status & stConnected) > 0
}
