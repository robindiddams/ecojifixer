package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/robindiddams/emojidict"
)

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

var paddingRunes = []rune{
	0x269C,
	0x1F3CD,
	0x1F4D1,
	0x1F64B,
}

var replacements = map[rune]rune{
	// block letters
	0x1f170: '🩸',
	0x1f171: '❌',
	0x1f17e: '⭕',
	0x1f17f: '🧊',
	0x1f202: '🫐',
	0x1f237: '🪙',

	// regional indicators
	0x1f1e6: '♈',
	0x1f1e7: '♉',
	0x1f1e8: '♊',
	0x1f1e9: '♋',
	0x1f1ea: '♌',
	0x1f1eb: '♍',
	0x1f1ec: '♎',
	0x1f1ed: '♏',
	0x1f1ee: '♐',
	0x1f1ef: '♑',
	0x1f1f0: '♒',
	0x1f1f1: '♓',
	0x1f1f2: '⛎',
	0x1f1f3: '🪐',
	0x1f1f4: '⭐',

	// Skin tones
	0x1f3fb: '🟧',
	0x1f3fc: '🟪',
	0x1f3fd: '🟦',
	0x1f3fe: '🟩',
	0x1f3ff: '🟫',
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

	for _, original := range ecojiset {
		usedRunes[original] = true
	}
	for _, replacements := range replacements {
		usedRunes[replacements] = true
	}
	for _, padding := range paddingRunes {
		usedRunes[padding] = true
	}

	getComputedRune := func() rune {
		for _, emoji := range emojidict.All {
			if len(emoji) == 1 && !usedRunes[emoji[0]] {
				return emoji[0]
			}
		}
		panic("no emoji found!")
	}

	for _, original := range ecojiset {
		if checkEmoji(original) {
			newSet = append(newSet, original)
			usedRunes[original] = true
		} else {
			str := fmt.Sprintf("! %s (0x%x) is invalid", string(original), original)

			if replacement, ok := replacements[original]; ok {
				str += fmt.Sprintf(", using replacement: %s (0x%x)", string(replacement), replacement)
				newSet = append(newSet, replacement)
				usedRunes[replacement] = true
			} else {
				rando := getComputedRune()
				str += fmt.Sprintf(", using auto-selected rune: %s (0x%x)", string(rando), rando)
				newSet = append(newSet, rando)
				usedRunes[rando] = true
			}
			fmt.Fprintln(os.Stderr, str)
		}
	}
	builder := strings.Builder{}
	for _, r := range ecojiset {
		builder.WriteString(fmt.Sprintf("%x\n", r))
	}
	ioutil.WriteFile("emojis.txt", []byte(builder.String()), 0644)
	v2builder := strings.Builder{}
	for _, r := range newSet {
		v2builder.WriteString(fmt.Sprintf("%x\n", r))
	}
	ioutil.WriteFile("emojisv2.txt", []byte(v2builder.String()), 0644)
}
