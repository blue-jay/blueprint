package asset

import (
	"testing"
)

// TestReset ensures config works properly.
func TestReset(t *testing.T) {
	folder := "asset"

	info := Info{
		Folder: folder,
	}

	SetConfig(info)

	c := Config()

	if c.Folder != folder {
		t.Errorf("Config returned unexpected body: got %v want %v",
			c.Folder, folder)
	}

	ResetConfig()

	c = Config()

	if c.Folder != "" {
		t.Errorf("Config returned unexpected body: got %v want %v",
			c.Folder, "''")
	}
}

// BenchmarkRace detects race conditions.
func BenchmarkRace(b *testing.B) {
	for n := 0; n < b.N; n++ {
		go func() {
			info := Info{
				Folder: "asset",
			}

			SetConfig(info)
			Config()
			ResetConfig()
		}()
	}
}
