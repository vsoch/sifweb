// Copyright 2019 Vanessa Sochat. All rights reserved.
// Use of this source code is governed by the Polyform Strict license
// that can be found in the LICENSE file and available at
// https://polyformproject.org/licenses/noncommercial/1.0.0

package main
 
import "errors"

// ErrNotFound is the code for when no search key is not found.
var ErrNotFound = errors.New("no match found")

// ErrMultValues is the code for when search key is not unique.
var ErrMultValues = errors.New("lookup would return more than one match")
