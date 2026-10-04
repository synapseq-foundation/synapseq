// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package sbg

import "testing"

func TestParseVoiceCompatibility(t *testing.T) {
	tests := []struct {
		token     string
		want      voice
		wantError string
	}{
		{token: "-", want: voice{kind: voiceOff}},
		{token: "pink/.5", want: voice{kind: voicePink, amplitude: .5}},
		{token: "mix/25", want: voice{kind: voiceMix, amplitude: 25}},
		{token: "200/10", want: voice{kind: voiceTone, carrier: 200, amplitude: 10}},
		{token: "200-4/10", want: voice{kind: voiceBinaural, carrier: 200, beat: 4, reverse: true, amplitude: 10}},
		{token: "spin:300+4.2/10", want: voice{kind: voiceSpin, carrier: 300, beat: 4.2, amplitude: 10}},

		{token: "200+4/10", want: voice{kind: voiceBinaural, carrier: 200, beat: 4, amplitude: 10}},
		{token: "200-4.2/10", want: voice{kind: voiceBinaural, carrier: 200, beat: 4.2, reverse: true, amplitude: 10}},
		{token: "200+4.2/10", want: voice{kind: voiceBinaural, carrier: 200, beat: 4.2, amplitude: 10}},
		{token: "200+0/10", want: voice{kind: voiceBinaural, carrier: 200, amplitude: 10}},
		{token: "200-0/10", want: voice{kind: voiceBinaural, carrier: 200, reverse: true, amplitude: 10}},
		{token: "spin:300-4.2/10", want: voice{kind: voiceSpin, carrier: 300, beat: 4.2, amplitude: 10}},
		{token: "200", wantError: "expected carrier[+|-beat]/amplitude"},
		{token: "200/1/2", wantError: "expected carrier[+|-beat]/amplitude"},
		{token: "200+4-2/10", wantError: "multiple beat signs"},
		{token: "spin:200/10", wantError: "missing beat frequency"},
		{token: "200+/10", wantError: "missing beat frequency"},
		{token: "200/1..0", wantError: "invalid amplitude: more than one decimal point"},
		{token: "200+a/10", wantError: "invalid beat: unexpected character 'a'"},
		{token: "bell", wantError: "bell voices are not supported by SynapSeq"},
		{token: "bell/20", wantError: "bell voices are not supported by SynapSeq"},
		{token: "bell+10/20", wantError: "bell voices are not supported by SynapSeq"},
		{token: "bell-10/20", wantError: "bell voices are not supported by SynapSeq"},
	}
	for _, test := range tests {
		t.Run(test.token, func(t *testing.T) {
			got, err := parseVoice(test.token)
			if test.wantError != "" {
				if err == nil || err.Error() != test.wantError {
					t.Fatalf("error = %v, want %s", err, test.wantError)
				}
				return
			}
			if err != nil || got != test.want {
				t.Fatalf("voice = %#v, error = %v, want %#v", got, err, test.want)
			}
		})
	}
}
