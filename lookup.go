// Copyright 2019 Vanessa Sochat. All rights reserved.
// Use of this source code is governed by the Polyform Strict license
// that can be found in the LICENSE file and available at
// https://polyformproject.org/licenses/noncommercial/1.0.0

package main
 
import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// GetGoArch returns the go runtime arch code from the SIF arch code.
// https://github.com/sylabs/sif/blob/master/pkg/sif/lookup.go#L48
func GetGoArch(sifarch string) (goarch string) {
	var ok bool

	archMap := map[string]string{
		HdrArch386:      "386",
		HdrArchAMD64:    "amd64",
		HdrArchARM:      "arm",
		HdrArchARM64:    "arm64",
		HdrArchPPC64:    "ppc64",
		HdrArchPPC64le:  "ppc64le",
		HdrArchMIPS:     "mips",
		HdrArchMIPSle:   "mipsle",
		HdrArchMIPS64:   "mips64",
		HdrArchMIPS64le: "mips64le",
		HdrArchS390x:    "s390x",
	}

	if goarch, ok = archMap[sifarch]; !ok {
		goarch = "unknown"
	}
	return goarch
}

// GetPartPrimSys returns the primary system partition if present. There should
// be only one primary system partition in a SIF file.
// https://github.com/sylabs/sif/blob/master/pkg/sif/lookup.go#L400
func (fimg *FileImage) GetPartPrimSys() (*Descriptor, int, error) {
	var descr *Descriptor
	index := -1

	for i, v := range fimg.DescrArr {
		if !v.Used {
			continue
		} else {
			if v.Datatype == DataPartition {
				ptype, err := v.GetPartType()
				if err != nil {
					return nil, -1, err
				}
				if ptype == PartPrimSys {
					if index != -1 {
						return nil, -1, ErrMultValues
					}
					index = i
					descr = &fimg.DescrArr[i]
				}
			}
		}
	}

	if index == -1 {
		return nil, -1, ErrNotFound
	}

	return descr, index, nil
}

// GetPartType extracts the Parttype field from the Extra field of a Partition Descriptor.
// https://github.com/sylabs/sif/blob/master/pkg/sif/lookup.go#L300
func (d *Descriptor) GetPartType() (Parttype, error) {
	if d.Datatype != DataPartition {
		return -1, fmt.Errorf("expected DataPartition, got %v", d.Datatype)
	}

	var pinfo Partition
	b := bytes.NewReader(d.Extra[:])
	if err := binary.Read(b, binary.LittleEndian, &pinfo); err != nil {
		return -1, fmt.Errorf("while extracting Partition extra info: %s", err)
	}

	return pinfo.Parttype, nil
}
