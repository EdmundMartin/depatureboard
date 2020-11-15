package client_test

import (
	"testing"

	c "github.com/EdmundMartin/depatureboard/pkg/client"
)

func TestItCleansAString(t *testing.T) {
	dirtyString := `
                            
	Portsmouth Harbour&nbsp;via Horsham  

	&nbsp;&amp; Bognor Regis&nbsp;via Horsham  

`
	cleanString := c.CleanString(dirtyString)
	expected := "Portsmouth Harbour via Horsham & Bognor Regis via Horsham"
	if cleanString != expected {
		t.Fatalf(`Cleaning strings are now working.
Expected:
%v

Actual:
%v
`, expected, cleanString)
	}
}
