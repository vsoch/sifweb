// Copyright 2019 Vanessa Sochat. All rights reserved.
// Use of this source code is governed by the Polyform Strict license
// that can be found in the LICENSE file and available at
// https://polyformproject.org/licenses/noncommercial/1.0.0

package main
 
import (
	"encoding/binary"
	"bytes"
	"fmt"
	"syscall/js"
	"time"
)


// loadBytes loads an imageString from the browser and populates FileImage with data.
func (fimg *FileImage) loadBytes(value js.Value, size int) error {

	// We can use CopyBytesToGo, need golang 1.13+
	sif := make([]byte, size)
	fmt.Println(value)
	howmany := js.CopyBytesToGo(sif, value)
	fmt.Println(howmany)

	// Read in the string to bytes, n should equal size
	reader := bytes.NewReader(sif)
        n, _ := reader.Read(sif)
	fmt.Println(sif)

	// Save the data and size to the FileImage
	fimg.Filesize = int64(n)
	fimg.Filedata = sif
	fimg.Reader = bytes.NewReader(fimg.Filedata)

	return nil
}

// Read the global header from the container file.
// https://github.com/sylabs/sif/blob/master/pkg/sif/load.go#L20
func (fimg *FileImage) readHeader() error {
	if err := binary.Read(fimg.Reader, binary.LittleEndian, &fimg.Header); err != nil {
		return fmt.Errorf("reading global header from container file: %s", err)
	}
	return nil
}

// A valid sif has SIFMAGIC at the top
func (fimg *FileImage) isValidSif() error {

	// check various header fields
	if trimZeroBytes(fimg.Header.Magic[:]) != HdrMagic {
		return fmt.Errorf("invalid SIF file: Magic |%s| want |%s|", fimg.Header.Magic, HdrMagic)
	}
	if trimZeroBytes(fimg.Header.Version[:]) > HdrVersion {
		return fmt.Errorf("invalid SIF file: Version %s want <= %s", fimg.Header.Version, HdrVersion)
	}

	return nil
}

// Seek to a particular spot in the Reader, exit on error
func (fimg *FileImage) seek(offset int64) error {

	_, err := fimg.Reader.Seek(offset, 0); 
	if err != nil {
		return fmt.Errorf("seek() to offset:%s %s", offset, err)
	}
	return nil

}


// Read descriptors from the SIF
// https://github.com/sylabs/sif/blob/master/pkg/sif/load.go#L29
func (fimg *FileImage) readDescriptors() error {


	fmt.Println("fimg.Header.Descroff", fimg.Header.Descroff)
	fimg.seek(fimg.Header.Descroff)

	fmt.Println("fimg.Header.Dtotal", fimg.Header.Dtotal)

	descr, _, err := fimg.GetPartPrimSys()
	if err == nil {
		fimg.PrimPartID = descr.ID
	}

	fmt.Println("fimg.PrimPartID", fimg.PrimPartID)

	// Initialize descriptor array (slice) and read them all from file
	// This seems to be too much for the browser	
	fimg.DescrArr = make([]Descriptor, DescrNumEntries)// fimg.Header.Dtotal)
	if err := binary.Read(fimg.Reader, binary.LittleEndian, &fimg.DescrArr); err != nil {
		fimg.DescrArr = nil
		return fmt.Errorf("reading descriptor array from container file: %s", err)
	}

	return nil
}

// returnResult back to the browser, in the innerHTML of the result element
func returnResult(output string) {
	js.Global().Get("document").
		Call("getElementById", "result").
		Set("innerHTML", output)
}

// loadContainer is linked with the JavaScript function of the same name.
// It takes as input the binary data from the SIF image, and attempts
// to read the header. This has to be modified to compile with wasm.
func loadContainer(this js.Value, val []js.Value) interface{} {
	fmt.Println("The container binary is:", val[0])
        fmt.Println("Size:", val[2].Int())
	fmt.Println("ArrayBuffer:", val[1])
	
	fimg := FileImage{}

	// read the string of given size to bytes from the SIF file
	if err := fimg.loadBytes(val[1], val[2].Int()); err != nil {
		returnResult("Error loading bytes.")
		return nil
	}

	// read global header from SIF file
	if err := fimg.readHeader(); err != nil {
		returnResult("Error reading header.")
		return nil
	}

	// validate global header
	if err := fimg.isValidSif(); err != nil {
		returnResult("This is not a valid sif")
		return nil
	}

	// read descriptor data
	if err := fimg.readDescriptors(); err != nil {
		fmt.Println("Skipping reading descriptors: ", err)
	}

	// header with newlines
	header := fimg.FmtHeader()

	// Add file info
	header = addFileName(val[0].String(), header)

	// Print header, and descriptors to console
	fmt.Print(header)

	// Replace with breaks
	header = replaceNewLine(header, "<br>")

	fmt.Println("Container id:", fimg.Header.ID)
	fmt.Println("Created on:  ", time.Unix(fimg.Header.Ctime, 0))
	fmt.Println("Modified on: ", time.Unix(fimg.Header.Mtime, 0))
	fmt.Println("----------------------------------------------------")

	// Send result back to browser
	returnResult(header)
	
	return nil
}

func trimZeroBytes(str []byte) string {
	return string(bytes.TrimRight(str, "\x00"))
}
