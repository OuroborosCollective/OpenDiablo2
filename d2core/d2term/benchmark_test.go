package d2term

import (
	"testing"
)

func BenchmarkParseCommand(b *testing.B) {
	command := "hello world \"quoted string\" \\escaped\\ space"
	for i := 0; i < b.N; i++ {
		parseCommand(command)
	}
}

func BenchmarkParseCommandLong(b *testing.B) {
	command := "cmd "
	for i := 0; i < 100; i++ {
		command += "arg" + string(rune(i)) + " "
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseCommand(command)
	}
}
