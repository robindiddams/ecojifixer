package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/robindiddams/emojidict"
)

var paddingRunes = []rune{
	0x2615,
	0x269C,
	0x1F3CD,
	0x1F4D1,
	0x1F64B,
}

func getName(r rune) (string, bool) {
	resp, err := http.Get(fmt.Sprintf("https://emojipedia.org/emoji/%c/", r))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	str := string(buf)
	match := regexp.MustCompile(`<title>(.*)</title>`).FindStringSubmatch(str)
	return match[1], false
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

// these are a part of future emoji spec (14)
var newEmojis = []rune{
	emojidict.MeltingFace[0],
	emojidict.FaceWithOpenEyesAndHandOverMouth[0],
	emojidict.FaceWithPeekingEye[0],
	emojidict.SalutingFace[0],
	emojidict.DottedLineFace[0],
	emojidict.FaceWithDiagonalMouth[0],
	emojidict.FaceHoldingBackTears[0],
}

// these ones keith didnt really like
var redundantRunes = []rune{
	emojidict.WhiteCircle[0],
	emojidict.BlackCircle[0],
	emojidict.CrossMark[0],              // x
	emojidict.CrossMarkButton[0],        // negative_squared_cross_mark
	emojidict.RedQuestionMark[0],        // question
	emojidict.WhiteQuestionMark[0],      // grey_question
	emojidict.WhiteExclamationMark[0],   // grey_exclamation
	emojidict.RedExclamationMark[0],     // exclamation
	emojidict.Plus[0],                   // heavy_plus_sign
	emojidict.Minus[0],                  // heavy_minus_sign
	emojidict.Divide[0],                 // heavy_division_sign
	emojidict.OrangeCircle[0],           // orange_circle
	emojidict.YellowCircle[0],           // yellow_circle
	emojidict.GreenCircle[0],            // green_circle
	emojidict.PurpleCircle[0],           // purple_circle
	emojidict.BrownCircle[0],            // brown_circle
	emojidict.RedSquare[0],              // red_square
	emojidict.BlueSquare[0],             // blue_square
	emojidict.OrangeSquare[0],           // orange_square
	emojidict.YellowSquare[0],           // yellow_square
	emojidict.GreenSquare[0],            // green_square
	emojidict.PurpleSquare[0],           // purple_square
	emojidict.BrownSquare[0],            // brown_square
	emojidict.BlackLargeSquare[0],       // black_large_square
	emojidict.WhiteLargeSquare[0],       // white_large_square
	emojidict.WhiteMediumSmallSquare[0], // white_medium_small_square
	emojidict.BlackMediumSmallSquare[0], // black_medium_small_square
	emojidict.CheckMarkButton[0],        // white_check_mark

	emojidict.Watch[0],            // watch
	emojidict.HourglassDone[0],    // hourglass
	emojidict.AlarmClock[0],       // alarm_clock
	emojidict.HourglassNotDone[0], // hourglass_flowing_sand
	emojidict.WhiteHeart[0],       // white_heart
	emojidict.BrownHeart[0],       // brown_heart
	emojidict.OrangeHeart[0],      // orange_heart
}

func main() {
	fmt.Println("fetching mapping from keith-turner/ecoji")
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
	checkRune := func(r rune) bool {
		for _, emoji := range emojidict.All {
			if len(emoji) == 1 && emoji[0] == r {
				return true
			}
		}
		return false
	}

	// singlePointRunes := make(map[rune]bool)
	var singlePointRunesStack []rune
	for _, emoji := range emojidict.All {
		if len(emoji) == 1 {
			// singlePointRunes[emoji[0]] = true
			singlePointRunesStack = append(singlePointRunesStack, emoji[0])
		}
	}

	removeRune := func(r rune) {
		for i, rr := range singlePointRunesStack {
			if rr == r {
				singlePointRunesStack = append(singlePointRunesStack[:i], singlePointRunesStack[i+1:]...)
				return
			}
		}
	}

	for _, original := range ecojiset {
		removeRune(original)
	}

	for _, originalPadding := range paddingRunes {
		removeRune(originalPadding)
	}
	for _, redundant := range redundantRunes {
		removeRune(redundant)
	}
	for _, new := range newEmojis {
		removeRune(new)
	}

	fmt.Println("remaining:", len(singlePointRunesStack))

	var index int
	getReplacement := func() rune {
		if len(singlePointRunesStack) == 0 {
			return 'x'
		}
		next := singlePointRunesStack[index]
		index++
		return next
	}

	fmt.Fprintf(os.Stderr, "## Padding \n\n")

	fmt.Fprintf(os.Stderr, "| index | Emoji (hex) | Replacement (hex) |\n")
	fmt.Fprintf(os.Stderr, "|-------|-------------|-------------------|\n")

	for i, original := range paddingRunes {
		if !checkRune(original) {
			replacement := getReplacement()
			name, draft := getName(replacement)

			fmt.Printf("replacement padding emoji (%c), using %x ( %c )  %s  (draft: %t)\n", original, replacement, replacement, name, draft)
			fmt.Fprintf(os.Stderr, "| %d | %c (%x) | %c (%x) |\n", i, original, original, replacement, replacement)
		} else {
			fmt.Fprintf(os.Stderr, "| %d | %c (%x) | - |\n", i, original, original)
		}
	}

	fmt.Fprintf(os.Stderr, "\n## Emojis \n\n")

	fmt.Fprintf(os.Stderr, "| index | Emoji (hex) | Replacement (hex) |\n")
	fmt.Fprintf(os.Stderr, "|-------|-------------|-------------------|\n")

	var finalSet []rune
	for i, original := range ecojiset {
		if !checkRune(original) {
			replacement := getReplacement()
			finalSet = append(finalSet, replacement)
			name, draft := getName(replacement)
			fmt.Printf("replacemed emoji %d (%c), with %x ( %c )  %s  (draft: %t)\n", i, original, replacement, replacement, name, draft)
			fmt.Fprintf(os.Stderr, "| %d | %c (%x) | %c (%x) |\n", i, original, original, replacement, replacement)
		} else {
			finalSet = append(finalSet, original)
			fmt.Fprintf(os.Stderr, "| %d | %c (%x) | - |\n", i, original, original)
		}
	}

	fmt.Fprintf(os.Stderr, "\n## Unused/remaining \n\n")

	fmt.Fprintf(os.Stderr, "| index | Emoji (hex) | Replacement (hex) |\n")
	fmt.Fprintf(os.Stderr, "|-------|-------------|-------------------|\n")

	for i := index; i < len(singlePointRunesStack); i++ {
		fmt.Fprintf(os.Stderr, "| - | %c (%x) | - |\n", singlePointRunesStack[i], singlePointRunesStack[i])
	}
	fmt.Println("unused:", len(singlePointRunesStack)-index+1)

	fmt.Println("writing final set")
	var str string
	for _, r := range finalSet {
		str = fmt.Sprintf("%s%x\n", str, r)

	}

	if err := os.WriteFile("emojis.txt", []byte(str), 0644); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// defer f.Close()

}
