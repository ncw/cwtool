# Morse code learning tools

Here are a selection of morse code learning tools.

## ncwtester

This measures and keep track of your morse code learning progress.

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

# Command docs

Here is the full documentation for each command of the cwtool.

<!-- auto generated from here on -->

## cwtool {#cwtool}


Show help for cwtool commands.

### Synopsis


Cwtool provides a suite of morse code tools.


### Options

```
  -h, --help      help for cwtool
  -v, --verbose   Verbose debugging
```

### SEE ALSO

* [cwtool completion](#cwtool_completion)	 - Generate the autocompletion script for the specified shell
* [cwtool keymorse](#cwtool_keymorse)	 - Snoop on all keypresses and turn into morse
* [cwtool ncwtester](#cwtool_ncwtester)	 - See how your morse receiving is going
* [cwtool play](#cwtool_play)	 - Play morse code from the command line or file
* [cwtool rss](#cwtool_rss)	 - Fetch RSS and turn into morse code


## cwtool completion {#cwtool_completion}


Generate the autocompletion script for the specified shell

### Synopsis

Generate the autocompletion script for cwtool for the specified shell.
See each sub-command's help for details on how to use the generated script.


### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
  -v, --verbose   Verbose debugging
```

### SEE ALSO

* [cwtool](#cwtool)	 - Show help for cwtool commands.
* [cwtool completion bash](#cwtool_completion_bash)	 - Generate the autocompletion script for bash
* [cwtool completion fish](#cwtool_completion_fish)	 - Generate the autocompletion script for fish
* [cwtool completion powershell](#cwtool_completion_powershell)	 - Generate the autocompletion script for powershell
* [cwtool completion zsh](#cwtool_completion_zsh)	 - Generate the autocompletion script for zsh


## cwtool keymorse {#cwtool_keymorse}


Snoop on all keypresses and turn into morse

### Synopsis

This command installs a listener to listen to all key presses and
turns them into morse code.

Since it snoops key presses from all applications, it requires root
privileges. It will use sudo to start a keylistener subprocess so
expect a sudo prompt.

For example to play all keypresses at 30 WPM

    cwtool keymorse --wpm 30



```
cwtool keymorse [flags]
```

### Options

```
  -c, --channels int       channels to generate (default 1)
      --farnsworth float   Increase character spacing to match this WPM
      --frequency float    HZ of morse (default 600)
  -h, --help               help for keymorse
      --out string         WAV file for output instead of speaker
  -s, --samplerate int     sample rate in samples/s (default 8000)
      --wpm float          WPM to send at (default 25)
```

### Options inherited from parent commands

```
  -v, --verbose   Verbose debugging
```

### SEE ALSO

* [cwtool](#cwtool)	 - Show help for cwtool commands.


## cwtool ncwtester {#cwtool_ncwtester}


See how your morse receiving is going

### Synopsis

This measures and keep track of your morse code learning progress.

It sends morse characters for you to receive and times how quickly you
receive each one.

It can send a group of characters and you can select which characters
are sent.


```
cwtool ncwtester [flags]
```

### Options

```
  -c, --channels int       channels to generate (default 1)
      --cutoff duration    If set, ignore stats older than this
      --farnsworth float   Increase character spacing to match this WPM
      --frequency float    HZ of morse (default 600)
      --group int          Send letters in groups this big (default 1)
  -h, --help               help for ncwtester
      --letters string     Letters to test (default "abcdefghijklmnopqrstuvwxyz0123456789.=/,?")
      --log string         CSV file to log attempts (default "ncwtesterstats.csv")
      --out string         WAV file for output instead of speaker
  -s, --samplerate int     sample rate in samples/s (default 8000)
      --wpm float          WPM to send at (default 25)
```

### Options inherited from parent commands

```
  -v, --verbose   Verbose debugging
```

### SEE ALSO

* [cwtool](#cwtool)	 - Show help for cwtool commands.


## cwtool play {#cwtool_play}


Play morse code from the command line or file

### Synopsis



This plays morse code from the command line or from a file with the
--file flag or from stdin with the --stdin flag.


```
cwtool play [flags]
```

### Options

```
  -c, --channels int       channels to generate (default 1)
      --farnsworth float   Increase character spacing to match this WPM
      --file string        File to play morse from (optional)
      --frequency float    HZ of morse (default 600)
  -h, --help               help for play
      --out string         WAV file for output instead of speaker
  -s, --samplerate int     sample rate in samples/s (default 8000)
      --stdin              If set play morse from stdin
      --wpm float          WPM to send at (default 25)
```

### Options inherited from parent commands

```
  -v, --verbose   Verbose debugging
```

### SEE ALSO

* [cwtool](#cwtool)	 - Show help for cwtool commands.


## cwtool rss {#cwtool_rss}


Fetch RSS and turn into morse code

### Synopsis

Fetch RSS and turn it into morse code

This fetches an RSS feed parses it and sends the items as morse code.

Here are some examples of RSS feeds. These are easy to find - just
search for the name of the publication + RSS.

    http://feeds.bbci.co.uk/news/uk/rss.xml
    http://www.nytimes.com/services/xml/rss/nyt/HomePage.xml
    http://rss.cnn.com/rss/cnn_topstories.rss
    https://news.ycombinator.com/rss
    http://feeds.wired.com/wired/index

Most RSS, Atom and JSON feed types are supported.

The following info is played from the feed. <BT> is the morse prosign.

- Title <BT>
- Description of feed <BT> (if --description is used)

Then for each item

- NR # <BT> (# is number starting from 1 and incrementing)
- Title of item <BT>
- Description of item <BT> (if --description is used)

This plays the title of the RSS feed and then the titles of each item
in the feed.

Use --description to add the descriptions of each link in as well as
their titles.

For example to play the BBC UK News to a file at 20 WPM but with 8 WPM
Farnsworth spacing:

    cwtool rss -v --url http://feeds.bbci.co.uk/news/uk/rss.xml --wpm 20 --farnsworth 8 --out bbc.wav



```
cwtool rss [flags]
```

### Options

```
  -c, --channels int       channels to generate (default 1)
      --description        If set add the description too
      --farnsworth float   Increase character spacing to match this WPM
      --frequency float    HZ of morse (default 600)
  -h, --help               help for rss
      --out string         WAV file for output instead of speaker
  -s, --samplerate int     sample rate in samples/s (default 8000)
      --url string         URL to fetch RSS from
      --wpm float          WPM to send at (default 25)
```

### Options inherited from parent commands

```
  -v, --verbose   Verbose debugging
```

### SEE ALSO

* [cwtool](#cwtool)	 - Show help for cwtool commands.


## cwtool completion bash {#cwtool_completion_bash}


Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(cwtool completion bash)

To load completions for every new session, execute once:

#### Linux:

	cwtool completion bash > /etc/bash_completion.d/cwtool

#### macOS:

	cwtool completion bash > $(brew --prefix)/etc/bash_completion.d/cwtool

You will need to start a new shell for this setup to take effect.


```
cwtool completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -v, --verbose   Verbose debugging
```

### SEE ALSO

* [cwtool completion](#cwtool_completion)	 - Generate the autocompletion script for the specified shell


## cwtool completion fish {#cwtool_completion_fish}


Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	cwtool completion fish | source

To load completions for every new session, execute once:

	cwtool completion fish > ~/.config/fish/completions/cwtool.fish

You will need to start a new shell for this setup to take effect.


```
cwtool completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -v, --verbose   Verbose debugging
```

### SEE ALSO

* [cwtool completion](#cwtool_completion)	 - Generate the autocompletion script for the specified shell


## cwtool completion powershell {#cwtool_completion_powershell}


Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	cwtool completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
cwtool completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -v, --verbose   Verbose debugging
```

### SEE ALSO

* [cwtool completion](#cwtool_completion)	 - Generate the autocompletion script for the specified shell


## cwtool completion zsh {#cwtool_completion_zsh}


Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(cwtool completion zsh)

To load completions for every new session, execute once:

#### Linux:

	cwtool completion zsh > "${fpath[1]}/_cwtool"

#### macOS:

	cwtool completion zsh > $(brew --prefix)/share/zsh/site-functions/_cwtool

You will need to start a new shell for this setup to take effect.


```
cwtool completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -v, --verbose   Verbose debugging
```

### SEE ALSO

* [cwtool completion](#cwtool_completion)	 - Generate the autocompletion script for the specified shell


