package chat

import (
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/justinian/dice"
)

func standardroll(ro string) string {
	roll := GetStringInBetween(ro, "roll[", "]")
	name := strings.TrimRight(ro, ":")
	var me string
	if roll == "" {
		me = "ERROR: No roll!"
	} else {
		result, _, err := dice.Roll(roll)
		if err != nil {
			return "ERROR: Incorrect roll!"
		}
		r := result.(dice.StdResult)
		sort.Ints(r.Rolls)
		sNums := make([]string, len(r.Rolls))
		for i, x := range r.Rolls {
			sNums[i] = strconv.Itoa(x)
		}
		list := strings.Join(sNums, ", ")
		total := roll_style.Render(strconv.Itoa(r.Total))
		me = name + "\n total = " + total + "\n dice_rolled = " + list
		if r.Dropped != nil {
			sort.Ints(r.Dropped)
			sNums = make([]string, len(r.Dropped))
			for i, x := range r.Dropped {
				sNums[i] = strconv.Itoa(x)
			}
			list = strings.Join(sNums, ", ")
			me += "\n dice_dropped = " + list
		}
	}
	return me
}

func shadowroll(ro string) string {
	roll := GetStringInBetween(ro, "sr[", "]")
	name := strings.TrimRight(ro, ":")
	var me string
	if roll == "" {
		me = "ERROR: No roll!"
	} else {
		res := roll + "d6rv5"
		result, _, err := dice.Roll(res)
		if err != nil {
			return "ERROR: Incorrect roll!"
		}
		r := result.(dice.VsResult)
		sort.Ints(r.Rolls)
		sNums := make([]string, len(r.Rolls))
		for i, x := range r.Rolls {
			sNums[i] = strconv.Itoa(x)
		}
		hits := roll_style.Render(strconv.Itoa(r.Successes))
		list := strings.Join(sNums, ", ")
		me = name + " \n total_hits = " + hits + "\n total_results = " + list
		count := strings.Count(list, "1")
		if count >= len(r.Rolls)/2 {
			me = me + " **GLITCH** \n "
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
	boring := GetStringInBetween(msg, start, end)
	colorpick := GetStringInBetween(msg, "/", "=")
	txt := GetStringInBetween(msg, "=", "/")
	prettify := lipgloss.NewStyle().Foreground(lipgloss.Color(colorpick))
	pretty_words := prettify.Render(txt)
	edited = start + boring + end
	msg = strings.Replace(msg, edited, pretty_words, -1)
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
