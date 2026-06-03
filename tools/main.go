package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
	hid "github.com/sstallion/go-hid"
)

var errFound = errors.New("found")

const (
	vid               = 0x1209
	pid               = 0x2641
	vendorUsagePage   = 0xFF00
	configIfaceNumber = 1

	reportIDConfig  = 1
	reportIDCommand = 2

	configReportSize = 14
	configMagic      = 0x564B4333

	cmdSave     = 1
	cmdLoad     = 2
	cmdDefaults = 3
	cmdBootsel  = 4

	keyTypeConsumer = 0
	keyTypeKeyboard = 1
)

type keyDef struct {
	keyType uint8
	code    uint16
}

// allKeys maps symbolic names to (type, code) pairs.
// Consumer keys from TinyUSB HID_USAGE_CONSUMER_*, keyboard keys from HID_KEY_*.
// Conflicts (power, stop, mute) resolved: consumer wins unqualified, keyboard gets key_ prefix.
var allKeys = map[string]keyDef{
	// Consumer control keys
	"consumer_control":                  {keyTypeConsumer, 0x0001},
	"power":                             {keyTypeConsumer, 0x0030},
	"reset":                             {keyTypeConsumer, 0x0031},
	"sleep":                             {keyTypeConsumer, 0x0032},
	"brightness_increment":              {keyTypeConsumer, 0x006F},
	"brightness_decrement":              {keyTypeConsumer, 0x0070},
	"wireless_radio_controls":           {keyTypeConsumer, 0x000C},
	"wireless_radio_buttons":            {keyTypeConsumer, 0x00C6},
	"wireless_radio_led":                {keyTypeConsumer, 0x00C7},
	"wireless_radio_slider_switch":      {keyTypeConsumer, 0x00C8},
	"play_pause":                        {keyTypeConsumer, 0x00CD},
	"scan_next":                         {keyTypeConsumer, 0x00B5},
	"scan_previous":                     {keyTypeConsumer, 0x00B6},
	"stop":                              {keyTypeConsumer, 0x00B7},
	"volume":                            {keyTypeConsumer, 0x00E0},
	"mute":                              {keyTypeConsumer, 0x00E2},
	"bass":                              {keyTypeConsumer, 0x00E3},
	"treble":                            {keyTypeConsumer, 0x00E4},
	"bass_boost":                        {keyTypeConsumer, 0x00E5},
	"volume_increment":                  {keyTypeConsumer, 0x00E9},
	"volume_decrement":                  {keyTypeConsumer, 0x00EA},
	"bass_increment":                    {keyTypeConsumer, 0x0152},
	"bass_decrement":                    {keyTypeConsumer, 0x0153},
	"treble_increment":                  {keyTypeConsumer, 0x0154},
	"treble_decrement":                  {keyTypeConsumer, 0x0155},
	"al_consumer_control_configuration": {keyTypeConsumer, 0x0183},
	"al_email_reader":                   {keyTypeConsumer, 0x018A},
	"al_calculator":                     {keyTypeConsumer, 0x0192},
	"al_local_browser":                  {keyTypeConsumer, 0x0194},
	"ac_search":                         {keyTypeConsumer, 0x0221},
	"ac_home":                           {keyTypeConsumer, 0x0223},
	"ac_back":                           {keyTypeConsumer, 0x0224},
	"ac_forward":                        {keyTypeConsumer, 0x0225},
	"ac_stop":                           {keyTypeConsumer, 0x0226},
	"ac_refresh":                        {keyTypeConsumer, 0x0227},
	"ac_bookmarks":                      {keyTypeConsumer, 0x022A},
	"ac_pan":                            {keyTypeConsumer, 0x0238},

	// Keyboard keys
	"a":               {keyTypeKeyboard, 0x04},
	"b":               {keyTypeKeyboard, 0x05},
	"c":               {keyTypeKeyboard, 0x06},
	"d":               {keyTypeKeyboard, 0x07},
	"e":               {keyTypeKeyboard, 0x08},
	"f":               {keyTypeKeyboard, 0x09},
	"g":               {keyTypeKeyboard, 0x0A},
	"h":               {keyTypeKeyboard, 0x0B},
	"i":               {keyTypeKeyboard, 0x0C},
	"j":               {keyTypeKeyboard, 0x0D},
	"k":               {keyTypeKeyboard, 0x0E},
	"l":               {keyTypeKeyboard, 0x0F},
	"m":               {keyTypeKeyboard, 0x10},
	"n":               {keyTypeKeyboard, 0x11},
	"o":               {keyTypeKeyboard, 0x12},
	"p":               {keyTypeKeyboard, 0x13},
	"q":               {keyTypeKeyboard, 0x14},
	"r":               {keyTypeKeyboard, 0x15},
	"s":               {keyTypeKeyboard, 0x16},
	"t":               {keyTypeKeyboard, 0x17},
	"u":               {keyTypeKeyboard, 0x18},
	"v":               {keyTypeKeyboard, 0x19},
	"w":               {keyTypeKeyboard, 0x1A},
	"x":               {keyTypeKeyboard, 0x1B},
	"y":               {keyTypeKeyboard, 0x1C},
	"z":               {keyTypeKeyboard, 0x1D},
	"1":               {keyTypeKeyboard, 0x1E},
	"2":               {keyTypeKeyboard, 0x1F},
	"3":               {keyTypeKeyboard, 0x20},
	"4":               {keyTypeKeyboard, 0x21},
	"5":               {keyTypeKeyboard, 0x22},
	"6":               {keyTypeKeyboard, 0x23},
	"7":               {keyTypeKeyboard, 0x24},
	"8":               {keyTypeKeyboard, 0x25},
	"9":               {keyTypeKeyboard, 0x26},
	"0":               {keyTypeKeyboard, 0x27},
	"enter":           {keyTypeKeyboard, 0x28},
	"escape":          {keyTypeKeyboard, 0x29},
	"backspace":       {keyTypeKeyboard, 0x2A},
	"tab":             {keyTypeKeyboard, 0x2B},
	"space":           {keyTypeKeyboard, 0x2C},
	"minus":           {keyTypeKeyboard, 0x2D},
	"equal":           {keyTypeKeyboard, 0x2E},
	"bracket_left":    {keyTypeKeyboard, 0x2F},
	"bracket_right":   {keyTypeKeyboard, 0x30},
	"backslash":       {keyTypeKeyboard, 0x31},
	"semicolon":       {keyTypeKeyboard, 0x33},
	"apostrophe":      {keyTypeKeyboard, 0x34},
	"grave":           {keyTypeKeyboard, 0x35},
	"comma":           {keyTypeKeyboard, 0x36},
	"period":          {keyTypeKeyboard, 0x37},
	"slash":           {keyTypeKeyboard, 0x38},
	"caps_lock":       {keyTypeKeyboard, 0x39},
	"f1":              {keyTypeKeyboard, 0x3A},
	"f2":              {keyTypeKeyboard, 0x3B},
	"f3":              {keyTypeKeyboard, 0x3C},
	"f4":              {keyTypeKeyboard, 0x3D},
	"f5":              {keyTypeKeyboard, 0x3E},
	"f6":              {keyTypeKeyboard, 0x3F},
	"f7":              {keyTypeKeyboard, 0x40},
	"f8":              {keyTypeKeyboard, 0x41},
	"f9":              {keyTypeKeyboard, 0x42},
	"f10":             {keyTypeKeyboard, 0x43},
	"f11":             {keyTypeKeyboard, 0x44},
	"f12":             {keyTypeKeyboard, 0x45},
	"print_screen":    {keyTypeKeyboard, 0x46},
	"scroll_lock":     {keyTypeKeyboard, 0x47},
	"pause":           {keyTypeKeyboard, 0x48},
	"insert":          {keyTypeKeyboard, 0x49},
	"home":            {keyTypeKeyboard, 0x4A},
	"page_up":         {keyTypeKeyboard, 0x4B},
	"delete":          {keyTypeKeyboard, 0x4C},
	"end":             {keyTypeKeyboard, 0x4D},
	"page_down":       {keyTypeKeyboard, 0x4E},
	"arrow_right":     {keyTypeKeyboard, 0x4F},
	"arrow_left":      {keyTypeKeyboard, 0x50},
	"arrow_down":      {keyTypeKeyboard, 0x51},
	"arrow_up":        {keyTypeKeyboard, 0x52},
	"num_lock":        {keyTypeKeyboard, 0x53},
	"keypad_divide":   {keyTypeKeyboard, 0x54},
	"keypad_multiply": {keyTypeKeyboard, 0x55},
	"keypad_subtract": {keyTypeKeyboard, 0x56},
	"keypad_add":      {keyTypeKeyboard, 0x57},
	"keypad_enter":    {keyTypeKeyboard, 0x58},
	"keypad_1":        {keyTypeKeyboard, 0x59},
	"keypad_2":        {keyTypeKeyboard, 0x5A},
	"keypad_3":        {keyTypeKeyboard, 0x5B},
	"keypad_4":        {keyTypeKeyboard, 0x5C},
	"keypad_5":        {keyTypeKeyboard, 0x5D},
	"keypad_6":        {keyTypeKeyboard, 0x5E},
	"keypad_7":        {keyTypeKeyboard, 0x5F},
	"keypad_8":        {keyTypeKeyboard, 0x60},
	"keypad_9":        {keyTypeKeyboard, 0x61},
	"keypad_0":        {keyTypeKeyboard, 0x62},
	"keypad_decimal":  {keyTypeKeyboard, 0x63},
	"application":     {keyTypeKeyboard, 0x65},
	"key_power":       {keyTypeKeyboard, 0x66},
	"keypad_equal":    {keyTypeKeyboard, 0x67},
	"f13":             {keyTypeKeyboard, 0x68},
	"f14":             {keyTypeKeyboard, 0x69},
	"f15":             {keyTypeKeyboard, 0x6A},
	"f16":             {keyTypeKeyboard, 0x6B},
	"f17":             {keyTypeKeyboard, 0x6C},
	"f18":             {keyTypeKeyboard, 0x6D},
	"f19":             {keyTypeKeyboard, 0x6E},
	"f20":             {keyTypeKeyboard, 0x6F},
	"f21":             {keyTypeKeyboard, 0x70},
	"f22":             {keyTypeKeyboard, 0x71},
	"f23":             {keyTypeKeyboard, 0x72},
	"f24":             {keyTypeKeyboard, 0x73},
	"execute":         {keyTypeKeyboard, 0x74},
	"help":            {keyTypeKeyboard, 0x75},
	"menu":            {keyTypeKeyboard, 0x76},
	"select":          {keyTypeKeyboard, 0x77},
	"key_stop":        {keyTypeKeyboard, 0x78},
	"again":           {keyTypeKeyboard, 0x79},
	"undo":            {keyTypeKeyboard, 0x7A},
	"cut":             {keyTypeKeyboard, 0x7B},
	"copy":            {keyTypeKeyboard, 0x7C},
	"paste":           {keyTypeKeyboard, 0x7D},
	"find":            {keyTypeKeyboard, 0x7E},
	"key_mute":        {keyTypeKeyboard, 0x7F},
	"volume_up":       {keyTypeKeyboard, 0x80},
	"volume_down":     {keyTypeKeyboard, 0x81},
	"control_left":    {keyTypeKeyboard, 0xE0},
	"shift_left":      {keyTypeKeyboard, 0xE1},
	"alt_left":        {keyTypeKeyboard, 0xE2},
	"gui_left":        {keyTypeKeyboard, 0xE3},
	"control_right":   {keyTypeKeyboard, 0xE4},
	"shift_right":     {keyTypeKeyboard, 0xE5},
	"alt_right":       {keyTypeKeyboard, 0xE6},
	"gui_right":       {keyTypeKeyboard, 0xE7},
}

var modifierNames = map[string]uint8{
	"ctrl":   0x01,
	"shift":  0x02,
	"alt":    0x04,
	"gui":    0x08,
	"rctrl":  0x10,
	"rshift": 0x20,
	"ralt":   0x40,
	"rgui":   0x80,
}

var modifierBitNames = [8]string{
	"ctrl", "shift", "alt", "gui",
	"rctrl", "rshift", "ralt", "rgui",
}

var keyNames map[keyDef]string

func init() {
	// Prepare a reverse map of keycodes to names
	keyNames = make(map[keyDef]string, len(allKeys))
	for name, def := range allKeys {
		keyNames[def] = name
	}
}

// Given a key type and code, return the string representation.
func keyName(kt uint8, code uint16) string {
	if name, ok := keyNames[keyDef{kt, code}]; ok {
		return name
	}
	prefix := "consumer"
	if kt == keyTypeKeyboard {
		prefix = "keyboard"
	}
	return fmt.Sprintf("%s:0x%04X", prefix, code)
}

func parseKey(s string) (keyDef, error) {
	name := strings.ToLower(s)
	if def, ok := allKeys[name]; ok {
		return def, nil
	}
	v, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return keyDef{}, fmt.Errorf("unknown key name %q (use 'list-keys' to see valid names)", s)
	}
	return keyDef{keyTypeConsumer, uint16(v)}, nil
}

func parseKeyCombo(s string) (keyDef, uint8, error) {
	parts := strings.Split(strings.ToLower(s), "+")

	if len(parts) == 1 {
		k, err := parseKey(parts[0])
		return k, 0, err
	}

	var mod uint8
	for _, p := range parts[:len(parts)-1] {
		if p == "" {
			return keyDef{}, 0, fmt.Errorf("empty modifier name")
		}
		bit, ok := modifierNames[p]
		if !ok {
			return keyDef{}, 0, fmt.Errorf("unknown modifier %q (valid: ctrl, shift, alt, gui, rctrl, rshift, ralt, rgui)", p)
		}
		if mod&bit != 0 {
			return keyDef{}, 0, fmt.Errorf("duplicate modifier %q", p)
		}
		mod |= bit
	}

	keyPart := parts[len(parts)-1]
	if keyPart == "" {
		return keyDef{}, 0, fmt.Errorf("missing key name after modifier(s)")
	}

	k, err := parseKey(keyPart)
	if err != nil {
		return keyDef{}, 0, err
	}
	if k.keyType == keyTypeConsumer {
		return keyDef{}, 0, fmt.Errorf("modifiers cannot be used with consumer keys; use a keyboard key instead")
	}

	return k, mod, nil
}

func keyComboName(kt uint8, code uint16, mod uint8) string {
	name := keyName(kt, code)
	if mod == 0 {
		return name
	}
	var parts []string
	for i, n := range modifierBitNames {
		if mod&(1<<i) != 0 {
			parts = append(parts, n)
		}
	}
	parts = append(parts, name)
	return strings.Join(parts, "+")
}

// Open the device and return an hid.Device to the caller.
func openDevice() (*hid.Device, error) {
	var path string
	err := hid.Enumerate(vid, pid, func(info *hid.DeviceInfo) error {
		if info.UsagePage == vendorUsagePage {
			path = info.Path
			return errFound
		}
		return nil
	})
	if err != nil && !errors.Is(err, errFound) {
		return nil, fmt.Errorf("enumerating HID devices: %w", err)
	}
	if path == "" {
		err := hid.Enumerate(vid, pid, func(info *hid.DeviceInfo) error {
			if info.InterfaceNbr == configIfaceNumber {
				path = info.Path
				return errFound
			}
			return nil
		})
		if err != nil && !errors.Is(err, errFound) {
			return nil, fmt.Errorf("enumerating HID devices: %w", err)
		}
	}
	if path == "" {
		return nil, fmt.Errorf("volume knob not found")
	}

	dev, err := hid.OpenPath(path)
	if err != nil {
		return nil, err
	}
	return dev, nil
}

type config struct {
	typeCW  uint8
	typeCCW uint8
	keyCW   uint16
	keyCCW  uint16
	divider uint16
	modCW   uint8
	modCCW  uint8
}

func configFromBuf(buf []byte) config {
	if len(buf) != configReportSize+1 {
		log.Fatalf("device configuration is invalid")
	}
	magic := binary.LittleEndian.Uint32(buf[1:5])
	if magic != configMagic {
		log.Fatalf("firmware version mismatch (expected magic 0x%08X, got 0x%08X); please update firmware or vkcfg", configMagic, magic)
	}
	return config{
		typeCW:  buf[5],
		typeCCW: buf[6],
		keyCW:   binary.LittleEndian.Uint16(buf[7:9]),
		keyCCW:  binary.LittleEndian.Uint16(buf[9:11]),
		divider: binary.LittleEndian.Uint16(buf[11:13]),
		modCW:   buf[13],
		modCCW:  buf[14],
	}
}

func configToBuf(reportId uint8, cfg config) []byte {
	buf := make([]byte, configReportSize+1)
	buf[0] = reportIDConfig
	binary.LittleEndian.PutUint32(buf[1:5], configMagic)
	buf[5] = cfg.typeCW
	buf[6] = cfg.typeCCW
	binary.LittleEndian.PutUint16(buf[7:9], cfg.keyCW)
	binary.LittleEndian.PutUint16(buf[9:11], cfg.keyCCW)
	binary.LittleEndian.PutUint16(buf[11:13], cfg.divider)
	buf[13] = cfg.modCW
	buf[14] = cfg.modCCW
	return buf
}

// Read configuration from the device.
func getConfig(dev *hid.Device) (config, error) {
	buf := make([]byte, configReportSize+1)
	buf[0] = reportIDConfig
	_, err := dev.GetFeatureReport(buf)
	if err != nil {
		return config{}, err
	}
	cfg := configFromBuf(buf)
	return cfg, nil
}

// Send configuration to the device.
func setConfig(dev *hid.Device, cfg config) error {
	buf := configToBuf(reportIDConfig, cfg)
	_, err := dev.SendFeatureReport(buf)
	return err
}

func sendCommand(dev *hid.Device, cmd byte) error {
	buf := []byte{reportIDCommand, cmd}
	_, err := dev.SendFeatureReport(buf)
	return err
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: vkcfg <command> [options]

Commands:
  get        show current configuration
  set        update configuration values
  save       persist current config to flash
  load       reload config from flash
  defaults   reset to compiled-in defaults
  bootsel    reboot device into BOOTSEL mode
  list-keys  list known key names`)
	os.Exit(1)
}

// Get and display the current configuration.
func cmdGet(dev *hid.Device) {
	cfg, err := getConfig(dev)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	fmt.Printf("key_cw   = %s\n", keyComboName(cfg.typeCW, cfg.keyCW, cfg.modCW))
	fmt.Printf("key_ccw  = %s\n", keyComboName(cfg.typeCCW, cfg.keyCCW, cfg.modCCW))
	fmt.Printf("divider  = %d\n", cfg.divider)
}

// Update the configuration. First get the configuration from the device, then
// update it with any new values before sending it back to the device.
func cmdSet(dev *hid.Device, args []string) {
	fs := flag.NewFlagSet("set", flag.ExitOnError)
	cwStr := fs.String("cw", "", "clockwise key (e.g. volume_increment, shift+page_up, ctrl+alt+a)")
	ccwStr := fs.String("ccw", "", "counter-clockwise key (e.g. volume_decrement, shift+page_down)")
	dividerStr := fs.String("divider", "", "encoder divider")
	if err := fs.Parse(args); err != nil {
		log.Fatalf("ERROR: failed to parse arguments")
	}

	cfg, err := getConfig(dev)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	if *cwStr != "" {
		k, mod, err := parseKeyCombo(*cwStr)
		if err != nil {
			log.Fatalf("ERROR: %v", err)
		}
		cfg.typeCW = k.keyType
		cfg.keyCW = k.code
		cfg.modCW = mod
	}
	if *ccwStr != "" {
		k, mod, err := parseKeyCombo(*ccwStr)
		if err != nil {
			log.Fatalf("ERROR: %v", err)
		}
		cfg.typeCCW = k.keyType
		cfg.keyCCW = k.code
		cfg.modCCW = mod
	}
	if *dividerStr != "" {
		v, err := strconv.ParseUint(*dividerStr, 0, 16)
		if err != nil {
			log.Fatalf("ERROR: invalid divider: %v", err)
		}
		cfg.divider = uint16(v)
	}

	if err := setConfig(dev, cfg); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	fmt.Println("OK")
}

// Send a simple command (one that does not require client-side logic) to the device.
func cmdSimple(dev *hid.Device, cmd byte) {
	if err := sendCommand(dev, cmd); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	fmt.Println("OK")
}

// Show a list of available key names.
func cmdListKeys() {
	names := make([]string, 0, len(allKeys))
	for name := range allKeys {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		def := allKeys[name]
		kind := "consumer"
		if def.keyType == keyTypeKeyboard {
			kind = "keyboard"
		}
		fmt.Printf("  %-40s %-10s 0x%04X\n", name, kind, def.code)
	}

	modNames := make([]string, 0, len(modifierNames))
	for name := range modifierNames {
		modNames = append(modNames, name)
	}
	sort.Strings(modNames)
	fmt.Println("\nModifiers (keyboard keys only):")
	fmt.Println("  Use modifier+key syntax, e.g. shift+page_up, ctrl+alt+a")
	for _, name := range modNames {
		fmt.Printf("  %s\n", name)
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("vkcfg: ")
	if len(os.Args) < 2 {
		usage()
	}

	cmd := os.Args[1]

	if cmd == "list-keys" {
		cmdListKeys()
		return
	}

	if err := hid.Init(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	// nolint:errcheck
	defer hid.Exit()

	dev, err := openDevice()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	// nolint:errcheck
	defer dev.Close()

	switch cmd {
	case "get":
		cmdGet(dev)
	case "set":
		cmdSet(dev, os.Args[2:])
	case "save":
		cmdSimple(dev, cmdSave)
	case "load":
		cmdSimple(dev, cmdLoad)
	case "defaults":
		cmdSimple(dev, cmdDefaults)
	case "bootsel":
		cmdSimple(dev, cmdBootsel)
	case "version":
		cmdVersion()
	default:
		usage()
	}
}
