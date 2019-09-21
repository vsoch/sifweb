// Copyright 2019 Vanessa Sochat. All rights reserved.
// Use of this source code is governed by the Polyform Strict license
// that can be found in the LICENSE file and available at
// https://polyformproject.org/licenses/noncommercial/1.0.0

package main
 
import (
	"fmt"
	"regexp"
	"time"
)

// readableSize returns the size in human readable format.
// https://github.com/sylabs/sif/blob/master/pkg/sif/fmt.go
func readableSize(size uint64) string {
	var divs int
	var conversion string

	for ; size != 0; size >>= 10 {
		if size < 1024 {
			break
		}
		divs++
	}

	switch divs {
	case 0:
		conversion = fmt.Sprintf("%d", size)
	case 1:
		conversion = fmt.Sprintf("%dKB", size)
	case 2:
		conversion = fmt.Sprintf("%dMB", size)
	case 3:
		conversion = fmt.Sprintf("%dGB", size)
	case 4:
		conversion = fmt.Sprintf("%dTB", size)
	}
	return conversion
}

// Replace newlines with another character (e.g., <br>) 
func replaceNewLine(input string, replacement string) string {
	re := regexp.MustCompile(`\r?\n`)
	return re.ReplaceAllString(input, replacement)
}

func addFileName(fileName string, s string) string {
	s = fmt.Sprintln("File:    ", fileName) + s
	return s
}


// FmtHeader formats the output of a SIF file global header.
func (fimg *FileImage) FmtHeader() string {
	s := fmt.Sprintln("Launch:  ", trimZeroBytes(fimg.Header.Launch[:]))
	s += fmt.Sprintln("Magic:   ", trimZeroBytes(fimg.Header.Magic[:]))
	s += fmt.Sprintln("Version: ", trimZeroBytes(fimg.Header.Version[:]))
	s += fmt.Sprintln("Arch:    ", GetGoArch(trimZeroBytes(fimg.Header.Arch[:])))
	s += fmt.Sprintln("ID:      ", fimg.Header.ID)
	s += fmt.Sprintln("Ctime:   ", time.Unix(fimg.Header.Ctime, 0))
	s += fmt.Sprintln("Mtime:   ", time.Unix(fimg.Header.Mtime, 0))
	s += fmt.Sprintln("Dfree:   ", fimg.Header.Dfree)
	s += fmt.Sprintln("Dtotal:  ", fimg.Header.Dtotal)
	s += fmt.Sprintln("Descoff: ", fimg.Header.Descroff)
	s += fmt.Sprintln("Descrlen:", readableSize(uint64(fimg.Header.Descrlen)))
	s += fmt.Sprintln("Dataoff: ", fimg.Header.Dataoff)
	s += fmt.Sprintln("Datalen: ", readableSize(uint64(fimg.Header.Datalen)))

	return s
}
