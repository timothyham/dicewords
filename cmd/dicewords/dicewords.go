// Copyright 2018 Timothy Ham
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/timothyham/dicewords"
)

var numPhrases = flag.Int("p", 5, "Number of phrases to generate")
var numWords = flag.Int("w", 0, "Number of words per passphrase")
var numBits = flag.Int("b", 64, "Number of bits to generates")
var short = flag.Bool("short", false, "Short words")
var shortUniq = flag.Bool("short2", false, "Short words with unique beginning")
var verbose = flag.Bool("v", false, "Print additional info")
var version = flag.Bool("version", false, "Print version")
var help = flag.Bool("h", false, "Print help")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		printHelp()
	}
	flag.Parse()
	if *help {
		printHelp()
		return
	}

	if *version {
		printVersion()
		return
	}

	dict := dicewords.Large
	if *short {
		dict = dicewords.Short
	} else if *shortUniq {
		dict = dicewords.Short2
	}

	if *numWords == 0 {
		// use numBits to determine numWords
		for i := 1; i < 20; i++ {
			estBits := dicewords.EstimateBits(i, dict)
			if estBits >= float64(*numBits) {
				*numWords = i
				break
			}
		}
	}

	for i := 0; i < *numPhrases; i++ {
		phrase, stats := dicewords.GetPhrase(*numWords, dict)
		fmt.Printf("%s\n", phrase)
		if *verbose {
			statStr := dicewords.PrintStats(stats)
			fmt.Printf("    %s\n", statStr)
		}
	}
}

func printHelp() {
	helpText := `
dicewords - print EFF dicewords

options:
-version 
    Show version.
-help 
    Show this help.
-b
    Target number of bits. Default is 64 bits.
-p 
    Number of passphrases to generate. Default is 5.
-w
    Number of words per passphrase. Overrides -b.
-short
    Use eff short words list.
-short2
    Use eff short unique 3 letter beginning words list.
-v
    Show additional information.
`
	fmt.Printf(helpText)
}

func printVersion() {
	fmt.Println(dicewords.VersionString)
}
