package crappers //(s)crapper

import (
	"testing"
)

func TestQuery(t *testing.T) {
	t.Parallel()

	tests := make(map[string]string)

	tests[" przestać    \t"] = "przestać nie ma opisu"
	tests["prześladować"] = "prześladować nie ma opisu"
	tests["dupę"] = "wulgarnie:\n1. pośladki, tyłek; zad, rzyć\n2. atrakcyjna" +
		" seksualnie kobieta; laska, dziwa, dupencja, towar\n3. człowiek niez" +
		"aradny, oferma; niezdara, ślamazara, gamoń\n4. tył czegoś, tylna czę" +
		"ść czegoś; zad"
	tests["burgundowy"] = "" //needs following links support
	tests["sławomira"] = "imię żeńskie\nimię męskie"
	tests["moje"] = "zaimek dzierżawczy będący odpowiednikiem zaimka osoboweg" +
		"o \"ja\"\npotocznie: to, co stanowi własność mówiącego lub ma z nim " +
		"związek\npotocznie: kobieta, która wraz z mówiącym stanowi parę (np." +
		" małżeńską)\npotocznie: osoba, która wraz z mówiącym stanowi parę (n" +
		"p. małżeńską)"
	tests["jajo"] = "1. jajko ptaków domowych wykorzystywane do przygotowywan" +
		"ia potraw\n2. żeńska komórka rozrodcza, komórka jajowa\n3. potocznie" +
		": przedmiot o owalnym kształcie\n4. potocznie: jądro męskie; jajco"
	tests["niewystepuje2137"] = "Nie występuje w słowniku"
	tests[""] = emptyQuery
	tests[" \n \t "] = emptyQuery

	for tt, want := range tests {
		tt, want := tt, want

		t.Run(tt, func(t *testing.T) {
			t.Parallel()

			resp, err := SjpQuery(tt)
			if err != nil {
				t.Error(err)
			}
			if resp != want {
				t.Errorf("got:\n%v\nwant:\n%v", resp, want)
			}
		})
	}
}

func BenchmarkQueryEmpty(b *testing.B) {
	if _, err := SjpQuery(""); err != nil {
		b.Error(err)
	}
}

func BenchmarkQueryNonExistent(b *testing.B) {
	if _, err := SjpQuery("123"); err != nil {
		b.Error(err)
	}
}

func BenchmarkQueryReal(b *testing.B) {
	if _, err := SjpQuery("test"); err != nil {
		b.Error(err)
	}
}
