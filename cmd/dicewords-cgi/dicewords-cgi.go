// Copyright 2020 Timothy Ham

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

	conf := dicewords.MakeConfig()
	if *short {
		conf.Dict = dicewords.Short
	} else if *shortUniq {
		conf.Dict = dicewords.Short2
	}

	conf.NumWords = *numWords
	conf.NumBits = *numBits
	conf.NumPhrases = *numPhrases

	phrases, stats := dicewords.MakeWords(conf)

	outWords := ""

	for i, words := range phrases {
		outWords += fmt.Sprintf("%s</br>", words)
		if *verbose {
			outWords += fmt.Sprintf("    %s</br>", dicewords.PrintStats(stats[i]))
		}
	}

	outApple := ""
	applePhrases, _ := dicewords.MakeApple(dicewords.Config{NumPhrases: 5, AppleStyle: true})
	for _, words := range applePhrases {
		outApple += fmt.Sprintf("%s</br>", words)
	}

	fmt.Printf(template, outWords, outApple)
}

var template string = `Content-type: text/html

<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
<title>Dicewords</title>
</head>
<body>
Each line of random words is about 65 bits</br></br>
%s
</br></br>
Apple style passwords with about 80 bits</br></br>
%s
</br>
</body>
</html>

`

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
-cgi
	Use html linebreak </br> instead of \n
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
