// +build linux
// Input device event monitor.
package main

// Could try pure go evdev module github.com/holoplot/go-evdev

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	//evdev "github.com/holoplot/go-evdev"
	evdev "github.com/gvalkov/golang-evdev"
	"github.com/ncw/ncwtester/cwgenerator"
	"github.com/ncw/ncwtester/cwplayer"
)

var (
	sampleRate = flag.Int("samplerate", 44100, "sample rate")
	wpm        = flag.Float64("wpm", 25.0, "WPM to send at")
	frequency  = flag.Float64("frequency", 600.0, "HZ of morse")
	logger     = flag.Bool("logger", false, "Set this to start the logger (done automatically)")
	debug      = flag.Bool("v", false, "Verbose debugging")
)

const (
	device_glob     = "/dev/input/event*"
	channelNum      = 2
	bitDepthInBytes = 2
	maxSampleValue  = 32767
)

func debugf(format string, a ...interface{}) {
	if !*debug {
		return
	}
	if *logger {
		os.Stderr.WriteString("K: ")
	} else {
		os.Stderr.WriteString("M: ")
	}
	fmt.Fprintf(os.Stderr, format, a...)
	os.Stderr.WriteString("\n")
}

// Select a device from a list of accessible input devices.
func selectDevices() ([]*evdev.InputDevice, error) {
	devices, err := evdev.ListInputDevices(device_glob)
	if err != nil {
		return nil, fmt.Errorf("failed to list input devices: %w", err)
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("no accessible input devices found by %q", device_glob)
	}

	// Filter out everything except keyboards
	newDevices := devices[:0]
	for _, dev := range devices {
		// Work out if this is a keyboard - does it have an A button
		for _, caps := range dev.Capabilities {
			for _, cap := range caps {
				if cap.Code == evdev.KEY_A {
					newDevices = append(newDevices, dev)
					break
				}
			}
		}
	}
	devices = newDevices
	if len(devices) == 0 {
		return nil, fmt.Errorf("no accessible keyboard devices found by %q", device_glob)
	}

	lines := make([]string, len(devices))
	max := 0
	for i := range devices {
		dev := devices[i]
		str := fmt.Sprintf("%-3d %-20s %-35s %s", i+1, dev.Fn, dev.Name, dev.Phys)
		if len(str) > max {
			max = len(str)
		}
		lines[i] = str
	}
	fmt.Fprintf(os.Stderr, "%-3s %-20s %-35s %s\n", "ID", "Device", "Name", "Phys")
	fmt.Fprintf(os.Stderr, strings.Repeat("-", max)+"\n")
	fmt.Fprintf(os.Stderr, strings.Join(lines, "\n")+"\n")
	return devices, nil
}

type device struct {
	dev  *evdev.InputDevice
	down [evdev.KEY_MAX + 1]bool
}

func newDevice(dev *evdev.InputDevice) *device {
	return &device{
		dev: dev,
	}
}

// func (d *device) open() error {
// 	dev, err = evdev.Open(args[0])
// }

// Synchronise loggers writing to stdout
var outMu sync.Mutex

func (d *device) handleKeyDown(code int, codeName string) error {
	// log.Printf("%s (%d)", codeName, code)
	modifierHeld := (d.down[evdev.KEY_LEFTCTRL] || d.down[evdev.KEY_RIGHTCTRL] ||
		d.down[evdev.KEY_LEFTALT] || d.down[evdev.KEY_RIGHTALT] ||
		d.down[evdev.KEY_LEFTMETA] || d.down[evdev.KEY_RIGHTMETA])
	if modifierHeld {
		return nil
	}
	shiftHeld := d.down[evdev.KEY_LEFTSHIFT] || d.down[evdev.KEY_RIGHTSHIFT]
	keyMap := normal_map
	if shiftHeld {
		keyMap = shift_map
	}
	c, found := keyMap[code]
	if !found {
		return nil
	}
	outMu.Lock()
	defer outMu.Unlock()
	fmt.Printf("%c\n", c)
	debugf("Tx: %s: %c", d.dev.Name, c)
	return nil
}

func (d *device) handleEvent(ev *evdev.InputEvent) error {
	switch ev.Type {
	case evdev.EV_KEY:
		var codeName string
		var kev evdev.KeyEvent
		kev.New(ev)
		code := int(kev.Scancode)
		val, haskey := evdev.KEY[code]
		if haskey {
			codeName = val
		} else {
			val, haskey := evdev.BTN[code]
			if haskey {
				codeName = val
			} else {
				codeName = "?"
			}
		}
		switch kev.State {
		case evdev.KeyUp:
			d.down[code] = false
		case evdev.KeyDown:
			d.down[code] = true
			d.handleKeyDown(code, codeName)
		case evdev.KeyHold:
			d.down[code] = true
		default:
		}
	}
	return nil
}

func (d *device) read() error {
	for {
		events, err := d.dev.Read()
		if err != nil {
			return fmt.Errorf("event read failed: %w", err)
		}
		for i := range events {
			err = d.handleEvent(&events[i])
			if err != nil {
				return fmt.Errorf("event handle failed: %w", err)
			}
		}
	}
}

// Starts the keyboard logger logging to stdout
func runLogger(args []string) error {
	var wg sync.WaitGroup
	devs, err := selectDevices()
	if err != nil {
		return fmt.Errorf("select devices failed: %w", err)
	}
	for _, dev := range devs {
		d := newDevice(dev)
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := d.read()
			if err != nil {
				log.Printf("keyboard listener failed: %v", err)
			}
		}()
	}
	fmt.Println("Keyboard listener started.")
	wg.Wait()
	return nil
}

// Read keys from in and send morse
func runMorser(in io.Reader) error {
	bufIn := bufio.NewReader(in)
	opt := cwgenerator.Options{
		WPM:             *wpm,
		Frequency:       *frequency,
		SampleRate:      *sampleRate,
		ChannelNum:      channelNum,
		BitDepthInBytes: bitDepthInBytes,
		MaxSampleValue:  maxSampleValue,
		Continuous:      true,
	}
	cw, err := cwplayer.New(opt)
	if err != nil {
		return fmt.Errorf("failed to make cw generator: %w", err)
	}
	cw.String(" vvv")
	cw.Play()

	line, _, err := bufIn.ReadLine()
	if err != nil {
		return fmt.Errorf("failed to start keyboard listener: %w", err)
	}
	_, _ = os.Stderr.Write(line)

	fmt.Fprintf(os.Stderr, "\nListening for keys pressed to send Morse.\n")
	for {
		c, _, err := bufIn.ReadRune()
		if c < 0x20 {
			continue
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		debugf("Rx: %c", c)
		cw.Rune(c)
	}
	return nil
}

func start() error {
	fmt.Println("Starting keyboard listener as root.")
	binary, err := exec.LookPath(os.Args[0])
	if err != nil {
		return fmt.Errorf("failed to find path for %q: %w", os.Args[0], err)
	}
	binary, err = filepath.Abs(binary)
	if err != nil {
		return fmt.Errorf("failed find absolute path for %q: %w", os.Args[0], err)
	}
	args := []string{binary, "-logger"}
	args = append(args, os.Args[1:]...)
	cmd := exec.Command("sudo", args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	what := cmd.Path + " " + strings.Join(cmd.Args, " ")
	debugf("Starting: %s", what)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to make stdout pipe: %s: %w", what, err)
	}
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start keyboard listener: %s: %w", what, err)
	}
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Fatalf("keyboard listener failed: %s: %v", what, err)
		}

	}()
	return runMorser(stdout)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()
	var err error
	if *logger {
		err = runLogger(flag.Args())
	} else {
		err = start()
	}
	if err != nil {
		log.Fatal(err)
	}
}
