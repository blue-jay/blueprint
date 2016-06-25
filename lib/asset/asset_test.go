package asset

import (
	"testing"
)

// BenchmarkRace detects race conditions
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
