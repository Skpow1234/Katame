package cli

import (
	"testing"

	"github.com/Skpow1234/Katame/internal/cloud"
)

func TestIsRemoteURI_Azure(t *testing.T) {
	if !isRemoteURI("az://container/blob") {
		t.Error("expected az:// to be remote")
	}
	if isRemoteURI("./local/path") {
		t.Error("local path should not be remote")
	}
	if cloud.DetectScheme("az://c/b") != cloud.SchemeAzure {
		t.Error("expected SchemeAzure")
	}
	if cloud.DetectScheme("s3://bucket/key") == cloud.SchemeAzure {
		t.Error("s3 should not be Azure")
	}
}
