# cwtool for Morse code learning by M5NCW

cwtool is a command line program to help you learn and practice Morse code.

It has 4 main modes which are described below:

- [cwtool keymorse](#cwtool_keymorse) - play keystrokes from all applications as Morse code (Linux only)
- [cwtool ncwtester](#cwtool_ncwtester) - measure your reaction times when decoding Morse code
- [cwtool play](#cwtool_play) - play morse code from the command line or a file
- [cwtool rss](#cwtool_rss) - turn an RSS feed into morse code

## Installation

`cwtool` is a self contained binary. It needs no installation. Just
download it and run it.

Download the relevant binary from

- https://github.com/ncw/cwtool/releases

Or alternatively if you have Go installed use

- `go get github.com/ncw/cwtool`

and this will build the binary in `$GOPATH/bin`.

Or if you want to hack on the source code, use `git` to clone the
repository then `go build` to build the code from within that
directory.

## License

This is free software under the terms of MIT the license (check the [LICENSE](LICENSE) file).

## Contact and support

The project website is at:

https://github.com/ncw/cwtool

There you can file bug reports, ask for help or contribute patches.

## Authors

- Nick Craig-Wood M5NCW nick@craig-wood.com

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
* [cwtool keymorse](#cwtool_keymorse)	 - Snoop on all keypresses and turn into Morse code
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


Snoop on all keypresses and turn into Morse code

### Synopsis



This is a program which listens to keyboard input and plays it as
morse code.

It is intended as an aid for learning Morse code, or possibly an
assistive aid.

It needs the `sudo` program to be installed and it will run a
subprocess using `sudo` to read the keys as this is a privileged
operation. It will listen to all keyboards it finds (technically input
devices with an `A` button!).

This command installs a listener to listen to all key presses and
turns them into morse code.

Since it snoops key presses from all applications, it requires root
privileges. It will use sudo to start a keylistener subprocess so
expect a sudo prompt.

Use `-wpm` to set the words per minute of the Morse code generated.
The `-logger` is used internally when running the root process to read
the keypresses.

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

The statistics show the 50%/75%/90% percentile timings.

After a round it complete it will show the cumulative statistics.

It is configured with command line flags.

All stats are stored in the file specified by `--log` for analysis.

Setting `--group` can send multiple characters at once - it waits for
them all to be received before carrying on.


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
`--file` flag or from stdin with the `--stdin` flag.



```
cwtool play [flags]
```

### Options

```
  -c, --channels int       channels to generate (default 1)
      --farnsworth float   Increase character spacing to match this WPM
      --file string        File to play Morse from (optional)
      --frequency float    HZ of morse (default 600)
  -h, --help               help for play
      --out string         WAV file for output instead of speaker
  -s, --samplerate int     sample rate in samples/s (default 8000)
      --stdin              If set play Morse from stdin
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



This fetches an RSS feed parses it and plays the items as Morse code.

Here are some examples of RSS feeds which can be passed to the `--url`
flag. These are easy to find - just search for the name of the
publication + RSS.

    http://feeds.bbci.co.uk/news/uk/rss.xml
    http://www.nytimes.com/services/xml/rss/nyt/HomePage.xml
    http://rss.cnn.com/rss/cnn_topstories.rss
    https://news.ycombinator.com/rss
    http://feeds.wired.com/wired/index

Most RSS, Atom and JSON feed types are supported.

The following info is played from the feed. `BT` is the morse prosign.

- Title `BT`
- Description of feed `BT` (if `--description` is used)

Then for each item

- NR # `BT` (# is a number starting from 1 and incrementing)
- Title of item `BT`
- Description of item `BT` (if `--description` is used)

This plays the title of the RSS feed and then the titles of each item
in the feed.

Use `--description` to add the descriptions of each link in as well as
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


