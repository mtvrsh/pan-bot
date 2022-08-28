package main

import (
	"testing"
)

func TestSjpQuery(t *testing.T) {
	tests := make(map[string]string)

	tests[" przestać    \t"] = "przestać nie ma opisu"
	tests["prześladować"] = "prześladować nie ma opisu"
	tests["dupę"] = "wulgarnie:\n1. pośladki, tyłek; zad, rzyć\n2. atrakcyjna" +
		" seksualnie kobieta; laska, dziwa, dupencja, towar\n3. człowiek niez" +
		"aradny, oferma; niezdara, ślamazara, gamoń\n4. tył czegoś, tylna czę" +
		"ść czegoś; zad"
	tests["burgundowy"] = ""
	tests["sławomira"] = ""
	tests["moje"] = ""
	tests["jajo"] = "1. jajko ptaków domowych wykorzystywane do przygotowywan" +
		"ia potraw\n2. żeńska komórka rozrodcza, komórka jajowa\n3. potocznie" +
		": przedmiot o owalnym kształcie\n4. potocznie: jądro męskie; jajco"
	tests["niewystepuje2137"] = "Nie występuje w słowniku"
	tests[""] = "Puste zapytanie, pusta odpowiedź :)"
	tests["  "] = "Puste zapytanie, pusta odpowiedź :)"

	for tt, want := range tests {
		t.Run(tt, func(t *testing.T) {
			resp, err := sjpQuery(tt)
			if err != nil {
				t.Error(err)
			}
			if resp != want {
				t.Errorf("got:\n%v\nwant:\n%v", resp, want)
			}
		})
	}
}

func TestAsMdCode(t *testing.T) {
	s := asMdCode("test")
	want := "```text\ntest\n```"
	if s != want {
		t.Errorf("got:\n%v\nwant:\n%v", s, want)
	}
}

func BenchmarkSjpQueryEmpty(b *testing.B) {
	_, err := sjpQuery("")
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkSjpQueryNonExistent(b *testing.B) {
	_, err := sjpQuery("123")
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkSjpQueryReal(b *testing.B) {
	_, err := sjpQuery("test")
	if err != nil {
		b.Error(err)
	}
}
