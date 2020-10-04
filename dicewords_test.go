// Copyright 2018 Timothy Ham
package dicewords

import (
	"strings"
	"testing"
)

func TestGetWords(t *testing.T) {
	word, err := GetLargeWord(11111)
	if word != "abacus" {
		t.Errorf("expected abacus, got %v\n", word)
	}
	word, err = GetLargeWord(66666)
	if word != "zoom" {
		t.Errorf("expected zoom, got %v\n", word)
	}

	_, err = GetLargeWord(12734)
	if err == nil {
		t.Errorf("expected error, got none")
	}
	_, err = GetLargeWord(3)
	if err == nil {
		t.Errorf("expected error, got none")
	}
	_, err = GetLargeWord(266666)
	if err == nil {
		t.Errorf("expected error, got none")
	}

	word, err = GetShortWord(6666)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if word != "zoom" {
		t.Errorf("unexpected word %s", word)
	}

	word, err = GetShortUniqueWord(6666)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if word != "zucchini" {
		t.Errorf("unexpected word %s", word)
	}
}

func TestGetPhrases(t *testing.T) {
	NumWords = 5

	phrase, stats := GetPhrase(6, Large)
	words := strings.Split(phrase, " ")
	if len(words) != 6 {
		t.Errorf("unexpected %v", len(words))
	}

	stats = getStats("hello there", Large)
	if stats.Length != 11 {
		t.Errorf("expected 11, got %d", stats.Length)
	}
	if stats.NumChars != 10 {
		t.Errorf("expected 10, got %d", stats.Length)
	}
	if stats.NumBits != 26.2 {
		t.Errorf("unexpected %v", stats.NumBits)
	}

	stats = getStats("hello there", Short)
	if stats.NumBits != 21 {
		t.Errorf("unexpected %v", stats.NumBits)
	}
}

func TestMakeWords(t *testing.T) {
	conf := MakeConfig()
	w, _ := MakeWords(conf)
	if len(w) != 5 {
		t.Errorf("unexpected %v", len(w))
	}
	fields := strings.Split(w[0], " ")
	if len(fields) != 5 {
		t.Errorf("unexpected %v", len(fields))
	}

	conf.NumWords = 0
	conf.NumBits = 120
	w, _ = MakeWords(conf)
	fields = strings.Split(w[0], " ")
	if len(fields) != 10 {
		t.Errorf("unexpected %v", len(fields))
	}
}
