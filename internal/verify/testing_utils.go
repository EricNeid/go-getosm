// SPDX-License-Identifier: MIT
// Copyright (c) 2021 Eric Neidhardt

// Package verify provides utility functions for tests.
package verify

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

// Assert fails the test if the condition is false and stops its execution.
func Assert(t *testing.T, condition bool, msg string, v ...interface{}) {
	t.Helper()

	if !condition {
		t.Fatalf(msg, v...)
	}
}

// Ok fails the test if an err is not nil and stops its execution.
func Ok(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

// Equals fails the test if exp is not equal to act and stops its execution.
func Equals(t *testing.T, exp, act any) {
	t.Helper()

	if !deepEqual(exp, act) {
		t.Fatalf("exp: %#v\ngot: %#v", exp, act)
	}
}

// AlmostEquals fails the test if exp is not almost equal to act and stops its execution.
func AlmostEquals(t *testing.T, exp, act float64) {
	t.Helper()

	if math.Abs(exp-act) > 0.000000001 {
		t.Fatalf("exp: %#v\ngot: %#v", exp, act)
	}
}

func deepEqual(v1, v2 interface{}) bool {
	if reflect.DeepEqual(v1, v2) {
		return true
	}
	var x1 interface{}
	bytesA, _ := json.Marshal(v1)
	_ = json.Unmarshal(bytesA, &x1)
	var x2 interface{}
	bytesB, _ := json.Marshal(v2)
	_ = json.Unmarshal(bytesB, &x2)
	return reflect.DeepEqual(x1, x2)
}

// NotNil fails the test if act is nil and stops its execution.
func NotNil(t *testing.T, act interface{}, msg string) {
	t.Helper()

	if act == nil {
		t.Fatalf("%s", msg)
	}
}
