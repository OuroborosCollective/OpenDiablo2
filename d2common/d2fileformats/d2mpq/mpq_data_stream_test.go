package d2mpq

import (
	"io"
	"testing"
)

func TestMpqDataStream_Read(t *testing.T) {
	data := []byte("hello world")
	stream := &Stream{
		Data: data,
		Size: 4096,
	}
	stream.Block = &Block{
		UncompressedFileSize: uint32(len(data)),
		Flags:                FileSingleUnit,
	}
	m := &MpqDataStream{
		stream: stream,
	}

	buf := make([]byte, 5)
	n, err := m.Read(buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if n != 5 {
		t.Errorf("Read returned n=%d, want 5", n)
	}
	if string(buf) != "hello" {
		t.Errorf("Read returned %q, want \"hello\"", string(buf))
	}

	n, err = m.Read(buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if n != 5 {
		t.Errorf("Read returned n=%d, want 5", n)
	}
	if string(buf) != " worl" {
		t.Errorf("Read returned %q, want \" worl\"", string(buf))
	}

	buf = make([]byte, 5)
	n, err = m.Read(buf)
	if n != 1 {
		t.Errorf("Read returned n=%d, want 1", n)
	}
	if string(buf[:n]) != "d" {
		t.Errorf("Read returned %q, want \"d\"", string(buf[:n]))
	}
}

func TestMpqDataStream_Seek(t *testing.T) {
	data := []byte("hello world")
	stream := &Stream{
		Data: data,
		Size: 4096,
	}
	stream.Block = &Block{
		UncompressedFileSize: uint32(len(data)),
		Flags:                FileSingleUnit,
	}
	m := &MpqDataStream{
		stream: stream,
	}

	tests := []struct {
		name    string
		offset  int64
		whence  int
		wantPos int64
		wantErr bool
	}{
		{"SeekStart", 6, io.SeekStart, 6, false},
		{"SeekCurrent", 2, io.SeekCurrent, 8, false},
		{"SeekEnd", -1, io.SeekEnd, 10, false},
		{"SeekStartAgain", 0, io.SeekStart, 0, false},
		{"InvalidWhence", 0, 99, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPos, err := m.Seek(tt.offset, tt.whence)
			if (err != nil) != tt.wantErr {
				t.Errorf("Seek() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if gotPos != tt.wantPos {
				t.Errorf("Seek() gotPos = %d, want %d", gotPos, tt.wantPos)
			}
			if int64(m.stream.Position) != tt.wantPos {
				t.Errorf("stream.Position = %d, want %d", m.stream.Position, tt.wantPos)
			}
		})
	}
}

func TestMpqDataStream_Close(t *testing.T) {
	stream := &Stream{}
	m := &MpqDataStream{
		stream: stream,
	}

	if err := m.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}

	if m.stream != nil {
		t.Error("Expected stream to be nil after Close")
	}
}
