# Morse code learning tools

Here are a selection of morse code learning tools.

## ncwtester

This measure and keep track of your morse code learning progress.

It sends morse characters for you to receive and times how quickly you receive each one.

It can send a group of characters and you can select which characters are sent.

```
$ ncwtester
Start test round with 41 letters? (y/n)> y
 1/41: u: reaction time  1002ms: OK
 2/41: s: reaction time   633ms: OK
 3/41: 5: reaction time  1315ms: OK
 4/41: m: reaction time   654ms: OK
 5/41: /: reaction time  1487ms: OK
 6/41: 9: reaction time   394ms: OK
 7/41: y: reaction time   596ms: OK
 8/41: h: reaction time   856ms: BAD s
 9/41: 6: reaction time   847ms: OK
10/41: 8: reaction time  1069ms: OK
11/41: =: reaction time   837ms: OK
12/41: w: reaction time   841ms: OK
13/41: i: reaction time   648ms: OK
14/41: o: reaction time  1175ms: OK
15/41: 3: reaction time   759ms: OK
16/41: g: reaction time   902ms: OK
17/41: l: reaction time   789ms: OK
18/41: .: reaction time  1080ms: OK
19/41: q: reaction time  1572ms: BAD ,
20/41: k: reaction time   907ms: OK
21/41: c: reaction time  1000ms: OK
22/41: e: reaction time   616ms: OK
23/41: n: reaction time   876ms: OK
24/41: j: reaction time   735ms: BAD 1
25/41: ?: reaction time   917ms: OK
26/41: t: reaction time   652ms: BAD n
27/41: z: reaction time  1527ms: OK
28/41: ,: reaction time   989ms: OK
29/41: b: reaction time  1237ms: OK
30/41: 1: reaction time   803ms: OK
31/41: v: reaction time  1032ms: OK
32/41: f: reaction time   895ms: OK
33/41: x: reaction time   922ms: OK
34/41: 4: reaction time   813ms: OK
35/41: d: reaction time   915ms: OK
36/41: a: reaction time   711ms: OK
37/41: 0: reaction time   644ms: OK
38/41: 7: reaction time  1667ms: OK
39/41: r: reaction time   692ms: OK
40/41: p: reaction time   874ms: OK
41/41: 2: reaction time   689ms: OK
Round stats
total: min  0.394 p50  0.896 p75  1.032 p90  1.315 errors   9.8%
```

The statistics show the 50%/75%/90% percentile timings.

After a round it complete it will show the cumulative statistics.

It is configured with command line flags.

```
$ ncwtester -h
Usage of ncwtester:
  -cutoff duration
    	If set, ignore stats older than this
  -frequency float
    	HZ of morse (default 600)
  -letters string
    	Letters to test (default "abcdefghijklmnopqrstuvwxyz0123456789.=/,?")
  -log string
    	CSV file to log attempts (default "ncwtesterstats.csv")
  -pipeline int
    	number of characters to send ahead
  -samplerate int
    	sample rate (default 44100)
  -wpm float
    	WPM to send at (default 25)
```

All stats are stored in the file specified by `-log` for analysis.

Setting `-group` can send multiple characters at once - it waits for
them all to be received before carrying on.

## keymorse

This is a program which listens to keyboard input and plays it as
morse code.

It is indented as a learning aid, or possibly an assistive aid.

It needs the `sudo` program to be installed and it will run a
subprocess using `sudo` to read the keys as this is a privileged
operation. It will listen to all keyboards it finds (technically input
devices with an `A` button!).

```
$ keymorse -wpm 25
Starting keyboard listener as root.
[sudo] password for user: 
ID  Device               Name                                Phys
--------------------------------------------------------------------------------------------------
1   /dev/input/event17   Microsoft Natural® Ergonomic Keyboard 4000 usb-0000:08:00.3-1.1.1/input0
2   /dev/input/event18   Microsoft Natural® Ergonomic Keyboard 4000 usb-0000:08:00.3-1.1.1/input1
3   /dev/input/event4    AT Translated Set 2 keyboard        isa0060/serio0/input0
Keyboard listener started.
Listening for keys pressed to send Morse.
```

This has the following flags

```
$ keymorse -h
Usage of keymorse:
  -frequency float
    	HZ of morse (default 600)
  -logger
    	Set this to start the logger (done automatically)
  -samplerate int
    	sample rate (default 44100)
  -v	Verbose debugging
  -wpm float
    	WPM to send at (default 25)
```

Use `-wpm` to set the words per minute of the Morse code generated.
The `-logger` is used internally when running the root process to read
the keypresses.

## Installation

Each binary is self contained.

### Build requirements

### Pre-built binaries
