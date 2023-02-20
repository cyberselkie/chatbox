package chat

import (
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/justinian/dice"
)

const (
	Bel = "\007"
)

func CommandsOutput(msg string) string {
	msg = DiceCommands(msg)
	return msg
}

// Need to remove this now that it renders in markdown but I'm sentimental
/*func TextCommands(msg string) string {
	// change color of inline text
	msg = ColorText(msg, "/", "/")
	//adding pseudo markdown
	//italics
	msg = TextStyles(msg, "*", "*", "italics")
	//bold
	msg = TextStyles(msg, "+", "+", "bold")
	//whisper
	msg = TextStyles(msg, "{", "}", "whisper")
	//underline
	msg = TextStyles(msg, "_", "_", "underline")

	return msg
}*/

func DiceCommands(msg string) string {
	//shadowrun dice command
	if strings.Contains(msg, "sr[") {
		msg = shadowroll(msg)
	}
	//standard dice command
	if strings.Contains(msg, "roll[") {
		msg = standardroll(msg)
	}

	return msg
}

// Meat of the Commands
func standardroll(ro string) string {
	roll := GetStringInBetween(ro, "roll[", "]")
	name := strings.TrimRight(ro, ":")
	var me string
	if roll == "" {
		me = "ERROR: No roll! \n"
	} else {
		result, _, err := dice.Roll(roll)
		if err != nil {
			return "ERROR: Incorrect roll! \n"
		}
		r := result.(dice.StdResult)
		sort.Ints(r.Rolls)
		sNums := make([]string, len(r.Rolls))
		for i, x := range r.Rolls {
			sNums[i] = strconv.Itoa(x)
		}
		list := strings.Join(sNums, ", ")
		total := "`" + strconv.Itoa(r.Total) + "`"
		me = name + "\n total = " + total + "\n dice_rolled = " + list + "\n"
		if r.Dropped != nil {
			sort.Ints(r.Dropped)
			sNums = make([]string, len(r.Dropped))
			for i, x := range r.Dropped {
				sNums[i] = strconv.Itoa(x)
			}
			list = strings.Join(sNums, ", ")
			me += "\n dice_dropped = " + list + "\n"
		}
	}
	return me
}

func shadowroll(ro string) string {
	roll := GetStringInBetween(ro, "sr[", "]")
	name := strings.TrimRight(ro, ":")
	var me string
	if roll == "" {
		me = "ERROR: No roll! \n"
	} else {
		res := roll + "d6rv5"
		result, _, err := dice.Roll(res)
		if err != nil {
			return "ERROR: Incorrect roll! \n"
		}
		r := result.(dice.VsResult)
		sort.Ints(r.Rolls)
		sNums := make([]string, len(r.Rolls))
		for i, x := range r.Rolls {
			sNums[i] = strconv.Itoa(x)
		}
		hits := "`" + strconv.Itoa(r.Successes) + "`"
		list := strings.Join(sNums, ", ")
		me = name + " \n total_hits = " + hits + "\n total_results = " + list + "\n"
		count := strings.Count(list, "1")
		if count >= len(r.Rolls)/2 {
			me = me + "\n `**GLITCH**` \n "
		}
	}
	return me
}
func GetStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}
	return str[s : s+e]
}

// color text!
func ColorText(msg string, start string, end string) string {
	var edited string
	for {
		boring := GetStringInBetween(msg, start, end)
		if boring != "" {
			colorpick := GetStringInBetween(msg, "/", "=")
			txt := GetStringInBetween(msg, "=", "/")
			prettify := lipgloss.NewStyle().Foreground(lipgloss.Color(colorpick))
			pretty_words := prettify.Render(txt)
			edited = start + boring + end
			msg = strings.Replace(msg, edited, pretty_words, -1)
		} else {
			break
		}
	}
	return msg
}

// pseudo markdown
var bold = lipgloss.NewStyle().Bold(true)
var italic = lipgloss.NewStyle().Italic(true)
var whisper = lipgloss.NewStyle().Faint(true)
var underline = lipgloss.NewStyle().Underline(true)

func TextStyles(msg string, start string, end string, style string) string {
	badwords := GetStringInBetween(msg, start, end)
	var edited string
	pretty_words := bold.Render(badwords)
	switch {
	case style == "bold":
		pretty_words = bold.Render(badwords)
	case style == "italics":
		pretty_words = italic.Render(badwords)
	case style == "whisper":
		pretty_words = whisper.Render(badwords)
	case style == "underline":
		pretty_words = underline.Render(badwords)
	}
	edited = start + badwords + end
	msg = strings.Replace(msg, edited, pretty_words, -1)
	return msg
}
