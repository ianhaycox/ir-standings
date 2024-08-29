//nolint:mnd // ok
package irsdk

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type VarBuffer struct {
	TickCount int // used to detect changes in data
	bufOffset int // offset from header
}

type VarType int

const (
	VarTypeChar     VarType = 0
	VarTypeBool     VarType = 1
	VarTypeInt      VarType = 2
	VarTypeBitField VarType = 3
	VarTypeFloat    VarType = 4
	VarTypeDouble   VarType = 5
	VarTypeETCount  VarType = 6
)

type Variable struct {
	VarType     VarType // irsdk_VarType
	offset      int     // offset fron start of buffer row
	Count       int     // number of entrys (array) so length in bytes would be irsdk_VarTypeBytes[type] * count
	countAsTime bool
	Name        string
	Desc        string
	Unit        string
	Value       interface{}
	Values      interface{}
	RawBytes    []byte
}

func (v Variable) String() string {
	var ret string

	switch v.VarType {
	case VarTypeChar:
		ret = fmt.Sprintf("%c", v.Value)
	case VarTypeBool:
		ret = fmt.Sprintf("%v", v.Value)
	case VarTypeInt:
		ret = fmt.Sprintf("%d", v.Value)
	case VarTypeBitField:
		ret = fmt.Sprintf("%s", v.Value)
	case VarTypeFloat:
		ret = fmt.Sprintf("%f", v.Value)
	case VarTypeDouble:
		ret = fmt.Sprintf("%f", v.Value)
	case VarTypeETCount:
		ret = fmt.Sprintf("Unknown (%d)", v.VarType)
	default:
		ret = fmt.Sprintf("Unknown (%d)", v.VarType)
	}

	return ret
}

// TelemetryVars holds all variables we can read from telemetry live
type TelemetryVars struct {
	lastVersion int
	vars        map[string]Variable
	mux         sync.Mutex
}

func findLatestBuffer(r reader, h *header) VarBuffer {
	var vb VarBuffer

	foundTickCount := 0

	for i := 0; i < h.numBuf; i++ {
		rbuf := make([]byte, 16)

		_, err := r.ReadAt(rbuf, int64(48+i*16))
		if err != nil {
			log.Fatal(err)
		}

		currentVb := VarBuffer{
			byte4ToInt(rbuf[0:4]),
			byte4ToInt(rbuf[4:8]),
		}

		// fmt.Printf("BUFF?: %+v\n", currentVb)

		if foundTickCount < currentVb.TickCount {
			foundTickCount = currentVb.TickCount
			vb = currentVb
		}
	}

	// fmt.Printf("BUFF: %+v\n", vb)

	return vb
}

func readVariableHeaders(r reader, h *header) *TelemetryVars {
	vars := TelemetryVars{vars: make(map[string]Variable, h.numVars)}

	for i := 0; i < h.numVars; i++ {
		rbuf := make([]byte, 144)

		_, err := r.ReadAt(rbuf, int64(h.headerOffset+i*144))
		if err != nil {
			log.Fatal(err)
		}

		v := Variable{
			VarType(byte4ToInt(rbuf[0:4])),
			byte4ToInt(rbuf[4:8]),
			byte4ToInt(rbuf[8:12]),
			int(rbuf[12]) > 0,
			bytesToString(rbuf[16:48]),
			bytesToString(rbuf[48:112]),
			bytesToString(rbuf[112:144]),
			nil,
			nil,
			nil,
		}

		vars.vars[v.Name] = v
	}

	return &vars
}

//nolint:funlen,gocognit,gocyclo // Engage brain
func readVariableValues(sdk *IRSDK) bool {
	newData := false

	if sessionStatusOK(sdk.h.status) {
		// find latest buffer for variables
		vb := findLatestBuffer(sdk.r, sdk.h)
		sdk.tVars.mux.Lock()

		if sdk.tVars.lastVersion < vb.TickCount {
			newData = true
			sdk.tVars.lastVersion = vb.TickCount
			sdk.lastValidData = time.Now().Unix()

			for varName, v := range sdk.tVars.vars {
				var rbuf []byte

				switch v.VarType {
				case VarTypeChar:
					values := make([]string, v.Count)

					for i := 0; i < v.Count; i++ {
						rbuf = make([]byte, 1)

						_, err := sdk.r.ReadAt(rbuf, int64(vb.bufOffset+v.offset+(1*i)))
						if err != nil {
							log.Fatal(err)
						}

						values[i] = string(rbuf[0])
					}

					v.Value = values[0]
					v.Values = values
				case VarTypeBool:
					values := make([]bool, v.Count)

					for i := 0; i < v.Count; i++ {
						rbuf = make([]byte, 1)

						_, err := sdk.r.ReadAt(rbuf, int64(vb.bufOffset+v.offset+(1*i)))
						if err != nil {
							log.Fatal(err)
						}

						values[i] = int(rbuf[0]) > 0
					}

					v.Value = values[0]
					v.Values = values
				case VarTypeInt:
					values := make([]int, v.Count)

					for i := 0; i < v.Count; i++ {
						rbuf = make([]byte, 4)

						_, err := sdk.r.ReadAt(rbuf, int64(vb.bufOffset+v.offset+(4*i)))
						if err != nil {
							log.Fatal(err)
						}

						values[i] = byte4ToInt(rbuf)
					}

					v.Value = values[0]
					v.Values = values
				case VarTypeBitField:
					values := make([]int, v.Count)

					for i := 0; i < v.Count; i++ {
						rbuf = make([]byte, 4)

						_, err := sdk.r.ReadAt(rbuf, int64(vb.bufOffset+v.offset+(4*i)))
						if err != nil {
							log.Fatal(err)
						}

						values[i] = byte4ToInt(rbuf)
					}

					v.Value = values[0]
					v.Values = values
				case VarTypeFloat:
					values := make([]float32, v.Count)

					for i := 0; i < v.Count; i++ {
						rbuf = make([]byte, 4)

						_, err := sdk.r.ReadAt(rbuf, int64(vb.bufOffset+v.offset+(4*i)))
						if err != nil {
							log.Fatal(err)
						}

						values[i] = byte4ToFloat(rbuf)
					}

					v.Value = values[0]
					v.Values = values
				case VarTypeDouble:
					values := make([]float64, v.Count)

					for i := 0; i < v.Count; i++ {
						rbuf = make([]byte, 8)

						_, err := sdk.r.ReadAt(rbuf, int64(vb.bufOffset+v.offset+(8*i)))
						if err != nil {
							log.Fatal(err)
						}

						values[i] = byte8ToFloat(rbuf)
					}

					v.Value = values[0]
					v.Values = values
				case VarTypeETCount:
					log.Printf("unknown var type: %d", v.VarType)
				default:
					log.Printf("unknown var type: %d", v.VarType)
				}

				v.RawBytes = rbuf
				sdk.tVars.vars[varName] = v
			}
		}
		sdk.tVars.mux.Unlock()
	}

	return newData
}
