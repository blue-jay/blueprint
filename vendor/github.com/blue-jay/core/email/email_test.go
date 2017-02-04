package email_test

import (
	"testing"

	"github.com/blue-jay/core/email"
)

// TestEmailFail ensures email fails.
func TestEmailFail(t *testing.T) {
	config := email.Info{
		Username: "",
		Password: "",
		Hostname: "127.0.0.1",
		Port:     25,
		From:     "from@example.com",
	}

	err := config.Send("to@example.com", "Subject", "Body")
	if err == nil {
		t.Errorf("Expected an error: %v", err)
	}
}

// TestEmailSuccess ensures email succeeds.
// Add your own working code here.
/*func TestEmailSuccess(t *testing.T) {
	config := email.Info{
		Username: "",
		Password: "",
		Hostname: "127.0.0.1",
		Port:     25,
		From:     "from@example.com",
	}

	err := config.Send("to@example.com", "Subject", "Body")
	if err != nil {
		t.Errorf("Error not expected: %v", err)
	}
}*/
