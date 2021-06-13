package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/robindiddams/emojidict"
)

var paddingRunes = []rune{
	0x269C,
	0x1F3CD,
	0x1F4D1,
	0x1F64B,
}

func parseMapping(buf []byte) ([]rune, error) {
	re := regexp.MustCompile(`\temojis\[\d+\] = 0x([0-9A-Z]+)\n`)
	matches := re.FindAllSubmatch(buf, -1)
	var emojis []rune
	for _, match := range matches {
		hexStr := string(match[1])
		n, err := strconv.ParseInt(hexStr, 16, 64)
		if err != nil {
			return nil, err
		}
		emojis = append(emojis, rune(n))
	}
	return emojis, nil
}

func getMapping() ([]byte, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/keith-turner/ecoji/master/mapping.go")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func checkEmoji(maybe rune) bool {
	for _, e := range emojidict.All {
		if len(e) == 1 && e[0] == maybe {
			return true
		}
	}
	return false
}

type sortableRunes []rune

func (s sortableRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortableRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortableRunes) Len() int {
	return len(s)
}

func sortRunes(runes []rune) []rune {
	// sort them
	s := sortableRunes(runes)
	sort.Sort(s)
	return s
}

func main() {
	fmt.Fprintln(os.Stderr, "fetching mapping from keith-turner/ecoji")
	buf, err := getMapping()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ecojiset, err := parseMapping(buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var newSet []rune
	usedRunes := make(map[rune]bool)

	var unsortedReplacers []rune
	for _, emoji := range emojidict.All {
		if len(emoji) == 1 && int(emoji[0]) > 0x1f004 {
			unsortedReplacers = append(unsortedReplacers, emoji[0])
		}
	}
	replacers := sortRunes(unsortedReplacers)

	for _, original := range ecojiset {
		usedRunes[original] = true
	}
	for _, padding := range paddingRunes {
		usedRunes[padding] = true
	}

	getComputedRune := func() rune {
		for _, emoji := range replacers {
			if !usedRunes[emoji] {
				return emoji
			}
		}
		panic("no emoji found!")
	}

	fmt.Fprintf(os.Stderr, "| Invalid Emoji (hex) | Replacement (hex) |\n")
	fmt.Fprintf(os.Stderr, "|---------------------|-------------------|\n")
	for _, original := range ecojiset {
		if checkEmoji(original) {
			newSet = append(newSet, original)
			usedRunes[original] = true
		} else {
			str := fmt.Sprintf("! %s (0x%x) is invalid", string(original), original)
			rando := getComputedRune()
			str += fmt.Sprintf(", using auto-selected rune: %s (0x%x)", string(rando), rando)
			newSet = append(newSet, rando)
			usedRunes[rando] = true
			fmt.Fprintf(os.Stderr, "| %c (%x) | %c (%x) |\n", original, original, rando, rando)
		}
	}
	builder := strings.Builder{}
	for _, r := range ecojiset {
		builder.WriteString(fmt.Sprintf("%x\n", r))
	}
	ioutil.WriteFile("emojisv1.txt", []byte(builder.String()), 0644)
	sorted := sortRunes(newSet)
	v2builder := strings.Builder{}
	for _, r := range sorted {
		v2builder.WriteString(fmt.Sprintf("%x\n", r))
	}
	ioutil.WriteFile("emojis.txt", []byte(v2builder.String()), 0644)
}
