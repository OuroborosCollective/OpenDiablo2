package d2util

import (
"reflect"
"testing"
)

func TestAsterToEmpty(t *testing.T) {
tests := []struct {
name     string
input    string
expected string
}{
{"asterix", "*test", ""},
{"no asterix", "test", "test"},
{"empty", "", ""},
{"space", " ", " "},
{"asterix in middle", "te*st", "te*st"},
{"Starts with asterix", "*test", ""},
{"Asterix at end", "test*", "test*"},
{"Multiple asterixes", "**test", ""},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
if got := AsterToEmpty(tt.input); got != tt.expected {
t.Errorf("AsterToEmpty(%q) = %v, want %v", tt.input, got, tt.expected)
}
})
}
}

func TestEmptyToZero(t *testing.T) {
tests := []struct {
name     string
input    string
expected string
}{
{"empty string", "", "0"},
{"single space", " ", "0"},
{"already zero", "0", "0"},
{"positive number", "1", "1"},
{"alpha string", "abc", "abc"},
{"double space", "  ", "  "},
{"leading space", " 1", " 1"},
{"trailing space", "1 ", "1 "},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
if got := EmptyToZero(tt.input); got != tt.expected {
t.Errorf("EmptyToZero(%q) = %v, want %v", tt.input, got, tt.expected)
}
})
}
}

func TestStringToInt(t *testing.T) {
tests := []struct {
name     string
input    string
expected int
}{
{"Positive number", "123", 123},
{"Negative number", "-123", -123},
{"Zero", "0", 0},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
if got := StringToInt(tt.input); got != tt.expected {
t.Errorf("StringToInt(%q) = %v, want %v", tt.input, got, tt.expected)
}
})
}
}

func TestStringToUint(t *testing.T) {
tests := []struct {
name     string
input    string
expected uint
}{
{"Positive number", "123", 123},
{"Zero", "0", 0},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
if got := StringToUint(tt.input); got != tt.expected {
t.Errorf("StringToUint(%q) = %v, want %v", tt.input, got, tt.expected)
}
})
}
}

func TestStringToUint8(t *testing.T) {
tests := []struct {
name     string
input    string
expected uint8
}{
{"Positive number", "123", 123},
{"Zero", "0", 0},
{"Max uint8", "255", 255},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
if got := StringToUint8(tt.input); got != tt.expected {
t.Errorf("StringToUint8(%q) = %v, want %v", tt.input, got, tt.expected)
}
})
}
}

func TestStringToInt8(t *testing.T) {
tests := []struct {
name     string
input    string
expected int8
}{
{"Positive number", "122", 122},
{"Negative number", "-128", -128},
{"Zero", "0", 0},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
if got := StringToInt8(tt.input); got != tt.expected {
t.Errorf("StringToInt8(%q) = %v, want %v", tt.input, got, tt.expected)
}
})
}
}

func TestUtf16BytesToString(t *testing.T) {
tests := []struct {
name     string
input    []byte
expected string
wantErr  bool
}{
{"Simple ASCII", []byte{'H', 0, 'e', 0, 'l', 0, 'l', 0, 'o', 0}, "Hello", false},
{"Odd length", []byte{'H'}, "", true},
{"Empty", []byte{}, "", false},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got, err := Utf16BytesToString(tt.input)
if (err != nil) != tt.wantErr {
t.Errorf("Utf16BytesToString() error = %v, wantErr %v", err, tt.wantErr)
return
}
if got != tt.expected {
t.Errorf("Utf16BytesToString() = %q, want %q", got, tt.expected)
}
})
}
}

func TestSplitIntoLinesWithMaxWidth(t *testing.T) {
tests := []struct {
name     string
input    string
maxWidth int
expected []string
}{
{
"Simple split",
"The quick brown fox",
10,
[]string{" The quick", "brown fox"},
},
{
"Single long word",
"Supercalifragilisticexpialidocious",
5,
[]string{"Supe", "rcali", "fragi", "listi", "cexpi", "alido", "cious"},
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
if got := SplitIntoLinesWithMaxWidth(tt.input, tt.maxWidth); !reflect.DeepEqual(got, tt.expected) {
t.Errorf("SplitIntoLinesWithMaxWidth() = %v, want %v", got, tt.expected)
}
})
}
}
