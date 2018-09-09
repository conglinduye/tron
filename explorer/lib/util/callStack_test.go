package util

import "testing"

func TestGetCallStackInfo(t *testing.T) {
	var fileName string
	var codeLine int
	var functionName string

	var callStack = 1
	fileName, codeLine, functionName = GetCallStackInfo(callStack)
	t.Log(fileName, "    ", codeLine, "    ", functionName)
}

func TestGetCurrentCallStackInfo(t *testing.T) {
	var fileName string
	var codeLine int
	var functionName string

	fileName, codeLine, functionName = GetCurrentCallStackInfo()
	t.Log(fileName, "    ", codeLine, "    ", functionName)
}
