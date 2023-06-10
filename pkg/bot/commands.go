package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/m3tav3rse/pan-bot/pkg/pudzianmechaniczny"
)

const helpMessage = `Dostępne komendy:
!echo <tekst>   Wypisuje tekst
!pm             Konwerter Pudzianów Mechanicznyh (PM) do watów i koni mechanicznych (W/KM)
!sjp <wyraz>    Znaczenie wyrazu z sjp.pl
!wpc            Aktualizuje kurs WPC/PLN
!h|!help        Na pewno nie wyświetli tej listy komend
!v|!version     Wyświetla wersję bota
`

func pudzianConverter(s string) string { // refactor
	var unit string
	failMsg := "Błędne dane: podaj liczbę i jednostkę (KM,PM,W)\n" +
		"Na przykład: 1337 KM"
	words := strings.Fields(s)

	// TODO rm HP?
	if len(words) == 1 {
		switch {
		case strings.HasSuffix(words[0], "HP"):
			words[0] = strings.TrimSuffix(words[0], "HP")
			fallthrough
		case strings.HasSuffix(words[0], "KM"):
			unit = "KM"
		case strings.HasSuffix(words[0], "PM"):
			unit = "PM"
		case strings.HasSuffix(words[0], "W"):
			unit = "W"
		}
		words[0] = strings.TrimSuffix(words[0], unit)
		words = append(words, unit)
	}

	if len(words) == 2 {
		var result float64
		number, err := strconv.ParseFloat(words[0], 64)
		if err != nil {
			return failMsg
		}
		switch words[1] {
		case "HP":
			fallthrough
		case "KM":
			unit = "PM"
			result = pudzianmechaniczny.HPToPM(number)
		case "PM":
			unit = "KM"
			result = pudzianmechaniczny.PMToHP(number)
		case "W":
			unit = "PM"
			result = pudzianmechaniczny.WToPM(number)
		default:
			return failMsg
		}
		return fmt.Sprintf("%v %v to %.2f %v", words[0], words[1], result, unit)
	}

	return failMsg
}
