package main

import "testing"

func TestEncodePCM16LE(t *testing.T) {
	buffer, err := encodePCM16LE([]int{0, 1, -1, 32767, -32768})
	if err != nil {
		t.Fatalf("encodePCM16LE returned error: %v", err)
	}

	expected := []byte{0, 0, 1, 0, 255, 255, 255, 127, 0, 128}
	if len(buffer) != len(expected) {
		t.Fatalf("unexpected buffer size: got %d want %d", len(buffer), len(expected))
	}
	for index := range expected {
		if buffer[index] != expected[index] {
			t.Fatalf("byte %d: got %d want %d", index, buffer[index], expected[index])
		}
	}
}

func TestEncodePCM16LERejectsOutOfRangeSamples(t *testing.T) {
	if _, err := encodePCM16LE([]int{32768}); err == nil {
		t.Fatal("expected out-of-range sample error")
	}
}