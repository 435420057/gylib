package handler

import (
	"testing"
	"external/test"
)

func TestFirmwareHandler(t *testing.T) {
	req := make(map[string]interface{}, 1)
	req["firmware_version"] = "111"
	resp := CheckFirmwareVersion(req)
	t.Log(test.ToJsonString(resp))
}