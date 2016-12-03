package common

import "testing"

func TestCurrentTimeMillis(t *testing.T) {
	if CurrentTimeMillis() < 1480805318583 {
		t.Errorf("CurrentTimeMillis gave %v, but should be higher than 1480805318583 which was on 2016-12-03", CurrentTimeMillis())
	}
}
