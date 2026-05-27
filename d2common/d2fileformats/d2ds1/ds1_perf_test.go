package d2ds1

import (
	"testing"
)

func BenchmarkLoadFileList(b *testing.B) {
	ds1 := exampleData()
	// Create a DS1 with many files
	files := make([]string, 1000)
	for i := range files {
		files[i] = "some/path/to/a/file/with/a/reasonably/long/name/and/extension.dt1"
	}
	ds1.Files = files
	ds1.version = 18 // version that has file list

	data := ds1.Marshal()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Unmarshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
