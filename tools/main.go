package main

import (
	"encoding/binary"
	"errors"
	"fmt"
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
	configIfaceNumber = 3

	reportIDConfig  = 1
	reportIDCommand = 2

	configReportSize = 8

	cmdSave     = 1
	cmdLoad     = 2
	cmdDefaults = 3

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

var keyNames map[keyDef]string

func init() {
	keyNames = make(map[keyDef]string, len(allKeys))
	for name, def := range allKeys {
		keyNames[def] = name
	}
}

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

func openDevice() (*hid.Device, error) {
	var path string
	hid.Enumerate(vid, pid, func(info *hid.DeviceInfo) error {
		if info.UsagePage == vendorUsagePage {
			path = info.Path
			return errFound
		}
		return nil
	})
	if path == "" {
		hid.Enumerate(vid, pid, func(info *hid.DeviceInfo) error {
			if info.InterfaceNbr == configIfaceNumber {
				path = info.Path
				return errFound
			}
			return nil
		})
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
}

func getConfig(dev *hid.Device) (config, error) {
	buf := make([]byte, configReportSize+1)
	buf[0] = reportIDConfig
	_, err := dev.GetFeatureReport(buf)
	if err != nil {
		return config{}, err
	}
	var cfg config
	cfg.typeCW = buf[1]
	cfg.typeCCW = buf[2]
	cfg.keyCW = binary.LittleEndian.Uint16(buf[3:5])
	cfg.keyCCW = binary.LittleEndian.Uint16(buf[5:7])
	cfg.divider = binary.LittleEndian.Uint16(buf[7:9])
	return cfg, nil
}

func setConfig(dev *hid.Device, cfg config) error {
	buf := make([]byte, configReportSize+1)
	buf[0] = reportIDConfig
	buf[1] = cfg.typeCW
	buf[2] = cfg.typeCCW
	binary.LittleEndian.PutUint16(buf[3:5], cfg.keyCW)
	binary.LittleEndian.PutUint16(buf[5:7], cfg.keyCCW)
	binary.LittleEndian.PutUint16(buf[7:9], cfg.divider)
	_, err := dev.SendFeatureReport(buf)
	return err
}

func sendCommand(dev *hid.Device, cmd byte) error {
	buf := []byte{reportIDCommand, cmd}
	_, err := dev.SendFeatureReport(buf)
	return err
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: vkcfg <command> [options]

Commands:
  get        show current configuration
  set        update configuration values
  save       persist current config to flash
  load       reload config from flash
  defaults   reset to compiled-in defaults
  list-keys  list known key names`)
	os.Exit(1)
}

func cmdGet(dev *hid.Device) {
	cfg, err := getConfig(dev)
	if err != nil {
		fatal(err)
	}
	fmt.Printf("key_cw   = %s\n", keyName(cfg.typeCW, cfg.keyCW))
	fmt.Printf("key_ccw  = %s\n", keyName(cfg.typeCCW, cfg.keyCCW))
	fmt.Printf("divider  = %d\n", cfg.divider)
}

func cmdSet(dev *hid.Device, args []string) {
	fs := flag.NewFlagSet("set", flag.ExitOnError)
	cwStr := fs.String("cw", "", "clockwise key (e.g. volume_increment, a, f1, 0xE9)")
	ccwStr := fs.String("ccw", "", "counter-clockwise key (e.g. volume_decrement, b, f2, 0xEA)")
	dividerStr := fs.String("divider", "", "encoder divider")
	fs.Parse(args)

	cfg, err := getConfig(dev)
	if err != nil {
		fatal(err)
	}

	if *cwStr != "" {
		k, err := parseKey(*cwStr)
		if err != nil {
			fatal(err)
		}
		cfg.typeCW = k.keyType
		cfg.keyCW = k.code
	}
	if *ccwStr != "" {
		k, err := parseKey(*ccwStr)
		if err != nil {
			fatal(err)
		}
		cfg.typeCCW = k.keyType
		cfg.keyCCW = k.code
	}
	if *dividerStr != "" {
		v, err := strconv.ParseUint(*dividerStr, 0, 16)
		if err != nil {
			fatal(fmt.Errorf("invalid divider: %v", err))
		}
		cfg.divider = uint16(v)
	}

	if err := setConfig(dev, cfg); err != nil {
		fatal(err)
	}
	fmt.Println("OK")
}

func cmdSimple(dev *hid.Device, cmd byte) {
	if err := sendCommand(dev, cmd); err != nil {
		fatal(err)
	}
	fmt.Println("OK")
}

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
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	cmd := os.Args[1]

	if cmd == "list-keys" {
		cmdListKeys()
		return
	}

	if err := hid.Init(); err != nil {
		fatal(err)
	}
	defer hid.Exit()

	dev, err := openDevice()
	if err != nil {
		fatal(err)
	}
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
	default:
		usage()
	}
}
