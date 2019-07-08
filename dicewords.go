// Copyright 2018 Timothy Ham
package dicewords

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var Debug bool
var NumWords int
var VersionString string

var EFFLargeWordList []string
var EFFShortWordList []string
var EFFShortWordUniqPrefix []string

type Dictionary int

const (
	Large Dictionary = iota
	Short
	Short2
)

func init() {
	EFFLargeWordList = strings.Split(EFFLargeWordListRaw, "\n")
	EFFLargeWordList = EFFLargeWordList[1 : len(EFFLargeWordList)-1]

	EFFShortWordList = strings.Split(EFFShortWordListRaw, "\n")
	EFFShortWordList = EFFShortWordList[1 : len(EFFShortWordList)-1]

	EFFShortWordUniqPrefix = strings.Split(EFFShortWordUniqPrefixRaw, "\n")
	EFFShortWordUniqPrefix = EFFShortWordUniqPrefix[1 : len(EFFShortWordUniqPrefix)-1]
}

// GetLargeWord needs 5 digit rolls
func GetLargeWord(rolls int) (string, error) {
	return getWord(EFFLargeWordList, rolls, 11111)
}

func GetShortWord(rolls int) (string, error) {
	return getWord(EFFShortWordList, rolls, 1111)
}

func GetShortUniqueWord(rolls int) (string, error) {
	return getWord(EFFShortWordUniqPrefix, rolls, 1111)
}

func getWord(list []string, rolls, smallest int) (string, error) {
	if rolls < smallest {
		return "", errors.New(fmt.Sprintf("Roll smaller than %d. Got %d\n", smallest, rolls))
	}
	// convert rolls into row index
	idx := 0
	factor := 1
	_ = idx
	rollsCopy := rolls
	for {
		digit := rolls % 10
		if digit > 6 {
			return "", errors.New(fmt.Sprintf("Bad roll input %d\n", rollsCopy))
		}
		idx += factor * (digit - 1)
		factor = factor * 6
		rolls = rolls / 10

		if rolls == 0 {
			break
		}
	}

	if idx < 0 || idx > len(list)-1 {
		return "", errors.New(fmt.Sprintf("roll outside range %d\n", rollsCopy))
	}
	row := list[idx]
	fields := strings.Split(row, "\t")
	return fields[1], nil
}

func PrintStats(stats Stats) string {
	return fmt.Sprintf("%.1f bits; %d long, %d non space chars", stats.NumBits, stats.Length, stats.NumChars)
}

type Stats struct {
	NumBits  float64
	Length   int
	NumChars int
}

func getStats(phrase string, dict Dictionary) Stats {
	phrase = strings.TrimSpace(phrase)
	words := strings.Split(phrase, " ")

	numChars := 0
	for _, word := range words {
		numChars += len(word)
	}

	stats := Stats{}
	stats.Length = len(phrase)
	stats.NumChars = numChars
	stats.NumBits = EstimateBits(len(words), dict)

	return stats
}

func GetPhrase(numWords int, dict Dictionary) (string, Stats) {
	// dice has six sides
	six := big.NewInt(6)

	// large words = 5 rolls, short words = 4 rolls
	count := 5
	if dict != Large {
		count = 4
	}

	nums := make([]int, count)
	res := ""

	for j := 0; j < numWords; j++ {
		for i := 0; i < count; i++ {
			n, err := rand.Int(rand.Reader, six)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return "", Stats{}
			}
			nums[i] = int(n.Int64())
		}

		factor := 1
		rolls := 0
		for i := count - 1; i > -1; i-- {
			digit := nums[i] + 1
			rolls = rolls + factor*digit
			factor = factor * 10
		}
		var word string
		var err error
		switch dict {
		case Large:
			word, err = GetLargeWord(rolls)
		case Short:
			word, err = GetShortWord(rolls)
		case Short2:
			word, err = GetShortUniqueWord(rolls)
		default:
			word, err = GetLargeWord(rolls)
		}

		if err != nil {
			panic(err.Error())
		}
		res = res + word + " "
	}

	res = strings.TrimSpace(res)
	return res, getStats(res, dict)
}

func EstimateBits(numWords int, dict Dictionary) float64 {
	baseBits := 12.9 // long words have 12.9 bits per word
	if dict == Short || dict == Short2 {
		baseBits = 10.3 // short words have 10.3 bits per word
	}

	numBits := baseBits*float64(numWords) + 0.4
	return numBits
}
