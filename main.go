package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/robindiddams/emojidict"
)

var paddingRunes = []rune{
	0x2615,
	0x269C,
	0x1F3CD,
	0x1F4D1,
	0x1F64B,
}

func getCachedName(r rune) (string, bool) {
	buf, err := os.ReadFile(fmt.Sprintf("cache/%x", r))
	if err != nil {
		if os.IsNotExist(err) {
			return "", false
		}
		panic(err)
	}
	return string(buf), true
}

func saveNameToCache(r rune, name string) {
	if err := os.WriteFile(fmt.Sprintf("cache/%x", r), []byte(name), 0644); err != nil {
		fmt.Fprintln(os.Stderr, err, os.IsExist(err))
	}
}

func getName(r rune) string {
	if cached, found := getCachedName(r); found {
		return cached
	}
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
	name := strings.TrimSpace(strings.Replace(strings.Replace(match[1], string(r), "", 1), "Emoji", "", 1))
	saveNameToCache(r, name)
	return name
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
	buf, err := ioutil.ReadFile("mapping.txt")
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// these are a part of future emoji spec (14)
var newEmojis = [][]rune{
	emojidict.MeltingFace,
	emojidict.FaceWithOpenEyesAndHandOverMouth,
	emojidict.FaceWithPeekingEye,
	emojidict.SalutingFace,
	emojidict.DottedLineFace,
	emojidict.FaceWithDiagonalMouth,
	emojidict.FaceHoldingBackTears,
	emojidict.RightwardsHand,
	emojidict.LeftwardsHand,
	emojidict.PalmDownHand,
	emojidict.PalmUpHand,
	emojidict.HandWithIndexFingerAndThumbCrossed,
	emojidict.IndexPointingAtTheViewer,
	emojidict.HeartHands,
	emojidict.BitingLip,
	emojidict.PregnantMan,
	emojidict.Coral,
	emojidict.Lotus,
	emojidict.EmptyNest,
	emojidict.NestWithEggs,
	emojidict.Beans,
	emojidict.PouringLiquid,
	emojidict.Jar,
	emojidict.PlaygroundSlide,
	emojidict.Wheel,
	emojidict.RingBuoy,
	emojidict.Hamsa,
	emojidict.MirrorBall,
	emojidict.LowBattery,
	emojidict.Crutch,
	emojidict.XRay,
	emojidict.HeavyEqualsSign,
	emojidict.Bubbles,
}

var peopleRunes = [][]rune{
	emojidict.DeafPerson,
	emojidict.Ninja,
	emojidict.PersonWithCrown,
	emojidict.PregnantPerson,
	emojidict.Mage,
	emojidict.Fairy,
	emojidict.Vampire,
	emojidict.Merperson,
	emojidict.Elf,
	emojidict.Genie,
	emojidict.Zombie,
	emojidict.Troll,
	emojidict.PersonStanding,
	emojidict.PersonKneeling,
	emojidict.PersonInSteamyRoom,
	emojidict.PersonInLotusPosition,
	emojidict.PersonClimbing,
	emojidict.PeopleHugging,

	// things that skin tone modifiers attatch to
	emojidict.RaisedHand,
	emojidict.PinchedFingers,
	emojidict.PinchingHand,
	emojidict.RaisedFist,
}

var redundantRunes = [][]rune{
	// these ones keith didnt really like
	emojidict.WhiteCircle,
	emojidict.BlackCircle,
	emojidict.CrossMark,              // x
	emojidict.CrossMarkButton,        // negative_squared_cross_mark
	emojidict.RedQuestionMark,        // question
	emojidict.WhiteQuestionMark,      // grey_question
	emojidict.WhiteExclamationMark,   // grey_exclamation
	emojidict.RedExclamationMark,     // exclamation
	emojidict.Plus,                   // heavy_plus_sign
	emojidict.Minus,                  // heavy_minus_sign
	emojidict.Divide,                 // heavy_division_sign
	emojidict.OrangeCircle,           // orange_circle
	emojidict.YellowCircle,           // yellow_circle
	emojidict.GreenCircle,            // green_circle
	emojidict.PurpleCircle,           // purple_circle
	emojidict.BrownCircle,            // brown_circle
	emojidict.RedSquare,              // red_square
	emojidict.BlueSquare,             // blue_square
	emojidict.OrangeSquare,           // orange_square
	emojidict.YellowSquare,           // yellow_square
	emojidict.GreenSquare,            // green_square
	emojidict.PurpleSquare,           // purple_square
	emojidict.BrownSquare,            // brown_square
	emojidict.BlackLargeSquare,       // black_large_square
	emojidict.WhiteLargeSquare,       // white_large_square
	emojidict.WhiteMediumSmallSquare, // white_medium_small_square
	emojidict.BlackMediumSmallSquare, // black_medium_small_square
	emojidict.CheckMarkButton,        // white_check_mark

	emojidict.Watch,            // watch
	emojidict.HourglassDone,    // hourglass
	emojidict.AlarmClock,       // alarm_clock
	emojidict.HourglassNotDone, // hourglass_flowing_sand
	emojidict.WhiteHeart,       // white_heart
	emojidict.BrownHeart,       // brown_heart
	emojidict.OrangeHeart,      // orange_heart

	// these are ones I dont really like
	emojidict.Elevator,
}

var selectionOverrides = map[int][]rune{
	859: emojidict.SmilingFaceWithTear,
	860: emojidict.DisguisedFace,
	664: emojidict.YawningFace,
}

var paddingSelectionOverrides = map[int][]rune{
	1: emojidict.PottedPlant,
	2: emojidict.RollerSkate,
}

func main() {
	os.Mkdir("cache", 0777)
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
		removeRune(redundant[0])
	}
	for _, new := range newEmojis {
		removeRune(new[0])
	}
	for _, personEmoji := range peopleRunes {
		removeRune(personEmoji[0])
	}
	for _, override := range selectionOverrides {
		removeRune(override[0])
	}
	for _, override := range paddingSelectionOverrides {
		removeRune(override[0])
	}

	fmt.Fprintln(os.Stderr, "remaining:", len(singlePointRunesStack))

	var index int
	getReplacement := func(isPadding bool, setIndex int) rune {
		if isPadding {
			if override, ok := paddingSelectionOverrides[setIndex]; ok {
				return override[0]
			}
		} else {
			if override, ok := selectionOverrides[setIndex]; ok {
				return override[0]
			}
		}
		if len(singlePointRunesStack) == 0 {
			return 'x'
		}
		next := singlePointRunesStack[index]
		index++
		return next
	}

	fmt.Printf("## Padding \n\n")

	fmt.Printf("| index | V1 Emoji (hex) | Replacement (hex) (name) |\n")
	fmt.Printf("|-------|-------------|-------------------|\n")

	for i, original := range paddingRunes {
		if !checkRune(original) {
			replacement := getReplacement(true, i)
			name := getName(replacement)

			fmt.Fprintf(os.Stderr, "replacement padding emoji (%c), using %x ( %c )  %s\n", original, replacement, replacement, name)
			fmt.Printf("| %d | %c (%x) | %c (%x) (%s) |\n", i, original, original, replacement, replacement, name)
		} else {
			fmt.Printf("| %d | %c (%x) | - |\n", i, original, original)
		}
	}

	fmt.Printf("\n## Emojis \n\n")

	fmt.Printf("| index | V1 Emoji (hex) | Replacement (hex) (name) |\n")
	fmt.Printf("|-------|-------------|-------------------|\n")

	var finalSet []rune
	for i, original := range ecojiset {
		if !checkRune(original) {
			replacement := getReplacement(false, i)
			finalSet = append(finalSet, replacement)
			name := getName(replacement)
			fmt.Fprintf(os.Stderr, "replacemed emoji %d (%c), with %x ( %c )  %s\n", i, original, replacement, replacement, name)
			fmt.Printf("| %d | %c (%x) | %c (%x) (%s) |\n", i, original, original, replacement, replacement, name)
		} else {
			finalSet = append(finalSet, original)
			fmt.Printf("| %d | %c (%x) | - |\n", i, original, original)
		}
	}

	fmt.Printf("\n## Unused/remaining \n\n")

	fmt.Printf("| index | V1 Emoji (hex) | Replacement (hex) (name) |\n")
	fmt.Printf("|-------|-------------|-------------------|\n")

	for i := index; i < len(singlePointRunesStack); i++ {
		name := getName(singlePointRunesStack[i])
		fmt.Printf("| - | %c (%x) (%s) | - |\n", singlePointRunesStack[i], singlePointRunesStack[i], name)
	}
	fmt.Fprintln(os.Stderr, "unused:", len(singlePointRunesStack)-index+1)

	fmt.Fprintln(os.Stderr, "writing final set")
	var str string
	for _, r := range finalSet {
		str = fmt.Sprintf("%s%x\n", str, r)

	}

	if err := os.WriteFile("emojis.txt", []byte(str), 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// defer f.Close()

}
