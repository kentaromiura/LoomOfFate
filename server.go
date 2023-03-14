package main

import (
	"html/template"
	"io"
	"math/bits"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	// === force template renderer
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.Renderer = renderer
	// ===
	e.GET("/next/:surge", func(c echo.Context) error {

		surge, err := strconv.ParseInt(c.Param("surge"), 10, 0)
		if err != nil {
			return c.JSON(400, map[string]interface{}{
				"error": "surge is wrong.",
			})
		}
		random := nextRandom()
		if random > 49 {
			random += int(surge)
		} else {
			random -= int(surge)
		}
		if random < 0 {
			random = 0
		}
		if random > 99 {
			random = 99
		}

		next := loom{value: random}
		var unexpected string = ""
		if next.value <= 4 || next.value >= 95 {
			nextUn := nextUnexpected()
			unexpected = UnexpectedResult(nextUn).String()
		}
		return c.JSON(200, map[string]interface{}{
			"knowledge":  next.knowledge().String(),
			"conflict":   next.conflict().String(),
			"endings":    next.endings().String(),
			"unexpected": unexpected,
		})
	})

	e.GET("/", func(c echo.Context) error {
		next := loom{value: nextRandom()}
		var unexpected string = ""
		if next.value <= 4 || next.value >= 95 {
			nextUn := nextUnexpected()
			unexpected = UnexpectedResult(nextUn).String()
		}
		cssText, _ := os.ReadFile("out.css")
		return c.Render(http.StatusOK, "template.html", map[string]interface{}{
			"knowledge":  next.knowledge().String(),
			"conflict":   next.conflict().String(),
			"endings":    next.endings().String(),
			"unexpected": unexpected,
			"css":        template.CSS(cssText),
		})
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func nextRandom() int {
	return rand.Intn(100)
}

func Mod32(n uint32, d uint32) uint32 {
	_, r := bits.Div32(0, n, d)
	return r
}

func nextUnexpected() uint32 {
	var seventeen uint32 = 17
	return Mod32(rand.Uint32(), seventeen)
}

type loom struct {
	value int
}
type LoomResult uint

const (
	Undefined LoomResult = iota
	Yes
	YesUnexpectedly
	YesBut
	YesAnd
	No
	NoAnd
	NoBut
	NoUnexpectedly
)

func (s LoomResult) String() string {
	switch s {
	case Yes:
		return "Yes"
	case YesUnexpectedly:
		return "Yes, and unexpectedly"
	case YesBut:
		return "Yes, but"
	case YesAnd:
		return "Yes, and"
	case No:
		return "No"
	case NoAnd:
		return "No, and"
	case NoBut:
		return "No, but"
	case NoUnexpectedly:
		return "No, and unexpectedly"
	}
	return "===UNDEFINED==="
}

type UnexpectedResult uint32

const (
	Foreshadowing UnexpectedResult = iota
	TyingOff
	ToConflict
	CostumeChange
	KeyGrip
	ToKnowledge
	Framing
	SetChange
	Upstaged
	PatternChange
	Limelit
	EnteringTheRed
	ToEndings
	Montage
	EnterStageLeft
	CrossStitch
	SixDegrees
)

func (s UnexpectedResult) String() string {
	switch s {
	case Foreshadowing:
		return "Wrap up current scene, set a new thread for the next one."
	case TyingOff:
		return "Fully resolve or substantially advance the main thread in this scene. Possibly create follow-up threads."
	case ToConflict:
		return "Next scene is a conflict one, move to it."
	case CostumeChange:
		return "one NPC either positively or negatively flips their position regarging something."
	case KeyGrip:
		return "move to a new scene with a new location or new goal."
	case ToKnowledge:
		return "move to a new scene centered on lore or investigation."
	case Framing:
		return "focus the main thread either on an existing or new NPC or on a specific object."
	case SetChange:
		return "move to another location in this scene"
	case Upstaged:
		return "an NPC do something to substatially move torwards their goal"
	case PatternChange:
		return "drastically change the main thread."
	case Limelit:
		return "next scene greatly benefit the PC."
	case EnteringTheRed:
		return "Move torwards a dangerous situation, leave the PC the choice of what to do."
	case ToEndings:
		return "give closure to current scene and start next."
	case Montage:
		return "skip forward using short descriptions to summarize what happens next."
	case EnterStageLeft:
		return "Someone appears in the scene."
	case CrossStitch:
		return "the rest of the scene revolves around another thread."
	case SixDegrees:
		return "a not necessarily positive connection is estabilished between PC and an NPC."
	}
	return ""
}

func (v loom) knowledge() LoomResult {
	switch check := v.value; {
	case check < 5:
		return NoUnexpectedly
	case check < 15:
		return NoBut
	case check < 20:
		return NoAnd
	case check < 50:
		return No
	case check < 80:
		return Yes
	case check < 85:
		return YesAnd
	case check < 95:
		return YesBut
	}
	return YesUnexpectedly
}

func (v loom) conflict() LoomResult {
	switch check := v.value; {
	case check < 2:
		return NoUnexpectedly
	case check < 6:
		return NoBut
	case check < 16:
		return NoAnd
	case check < 50:
		return No
	case check < 84:
		return Yes
	case check < 94:
		return YesAnd
	case check < 98:
		return YesBut
	}
	return YesUnexpectedly
}

func (v loom) endings() LoomResult {
	switch check := v.value; {
	case check == 0:
		return NoUnexpectedly
	case check == 1:
		return NoBut
	case check < 20:
		return NoAnd
	case check < 50:
		return No
	case check < 80:
		return Yes
	case check < 98:
		return YesAnd
	case check == 98:
		return YesBut
	}
	return YesUnexpectedly
}
