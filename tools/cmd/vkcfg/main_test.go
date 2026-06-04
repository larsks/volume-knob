package main

import (
	"testing"
)

func TestParseKeyCombo(t *testing.T) {
	tests := []struct {
		input   string
		wantKey keyDef
		wantMod uint8
	}{
		{"page_up", keyDef{keyTypeKeyboard, 0x4B}, 0x00},
		{"volume_increment", keyDef{keyTypeConsumer, 0x00E9}, 0x00},
		{"shift+page_up", keyDef{keyTypeKeyboard, 0x4B}, 0x02},
		{"ctrl+page_up", keyDef{keyTypeKeyboard, 0x4B}, 0x01},
		{"ctrl+shift+page_up", keyDef{keyTypeKeyboard, 0x4B}, 0x03},
		{"ctrl+shift+alt+gui+a", keyDef{keyTypeKeyboard, 0x04}, 0x0F},
		{"rctrl+a", keyDef{keyTypeKeyboard, 0x04}, 0x10},
		{"rshift+ralt+a", keyDef{keyTypeKeyboard, 0x04}, 0x60},
		// Standalone modifier key (no +) should be treated as a regular key
		{"shift_left", keyDef{keyTypeKeyboard, 0xE1}, 0x00},
		{"control_left", keyDef{keyTypeKeyboard, 0xE0}, 0x00},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			k, mod, err := parseKeyCombo(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if k != tt.wantKey {
				t.Errorf("key = %+v, want %+v", k, tt.wantKey)
			}
			if mod != tt.wantMod {
				t.Errorf("mod = 0x%02X, want 0x%02X", mod, tt.wantMod)
			}
		})
	}
}

func TestParseKeyComboErrors(t *testing.T) {
	tests := []struct {
		input   string
		wantMsg string
	}{
		{"ctrl+volume_increment", "modifiers cannot be used with consumer keys"},
		{"shift+mute", "modifiers cannot be used with consumer keys"},
		{"foo+a", "unknown modifier"},
		{"ctrl+ctrl+a", "duplicate modifier"},
		{"ctrl+", "missing key name after modifier"},
		{"+a", "empty modifier name"},
		{"ctrl+nonexistent", "unknown key name"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, _, err := parseKeyCombo(tt.input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !contains(err.Error(), tt.wantMsg) {
				t.Errorf("error = %q, want substring %q", err.Error(), tt.wantMsg)
			}
		})
	}
}

func TestKeyComboName(t *testing.T) {
	tests := []struct {
		name string
		kt   uint8
		code uint16
		mod  uint8
		want string
	}{
		{"no modifier", keyTypeKeyboard, 0x4B, 0x00, "page_up"},
		{"consumer no modifier", keyTypeConsumer, 0x00E9, 0x00, "volume_increment"},
		{"shift", keyTypeKeyboard, 0x4B, 0x02, "shift+page_up"},
		{"ctrl+shift", keyTypeKeyboard, 0x4B, 0x03, "ctrl+shift+page_up"},
		{"all left modifiers", keyTypeKeyboard, 0x04, 0x0F, "ctrl+shift+alt+gui+a"},
		{"right ctrl", keyTypeKeyboard, 0x04, 0x10, "rctrl+a"},
		{"unknown key with modifier", keyTypeKeyboard, 0xFF, 0x02, "shift+keyboard:0x00FF"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := keyComboName(tt.kt, tt.code, tt.mod)
			if got != tt.want {
				t.Errorf("keyComboName(%d, 0x%04X, 0x%02X) = %q, want %q", tt.kt, tt.code, tt.mod, got, tt.want)
			}
		})
	}
}

func TestConfigRoundTrip(t *testing.T) {
	cfg := config{
		typeCW:  keyTypeKeyboard,
		typeCCW: keyTypeKeyboard,
		keyCW:   0x4B,
		keyCCW:  0x4E,
		divider: 256,
		modCW:   0x03,
		modCCW:  0x01,
	}
	buf := configToBuf(reportIDConfig, cfg)
	got := configFromBuf(buf)
	if got != cfg {
		t.Errorf("round-trip failed:\n  got  %+v\n  want %+v", got, cfg)
	}
}

func TestConfigRoundTripNoModifiers(t *testing.T) {
	cfg := config{
		typeCW:  keyTypeConsumer,
		typeCCW: keyTypeConsumer,
		keyCW:   0x00E9,
		keyCCW:  0x00EA,
		divider: 256,
		modCW:   0,
		modCCW:  0,
	}
	buf := configToBuf(reportIDConfig, cfg)
	got := configFromBuf(buf)
	if got != cfg {
		t.Errorf("round-trip failed:\n  got  %+v\n  want %+v", got, cfg)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
