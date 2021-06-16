package helper

import (
	"testing"
)

func TestGeoIP(t *testing.T) {
	ip, err := GeoIP()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(ip)
}
