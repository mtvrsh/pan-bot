package crappers //(s)crappers

import (
	"testing"
)

func TestQuerySjp(t *testing.T) {
	t.Parallel()

	tests := []struct{ str, want string }{
		{" przestać    \t", "przestać nie ma opisu"},
		{"prześladować", "prześladować nie ma opisu"},
		{"dupę", "wulgarnie:\n1. pośladki, tyłek; zad, rzyć\n2. atrakcyjna se" +
			"ksualnie kobieta; laska, dziwa, dupencja, towar\n3. człowiek nie" +
			"zaradny, oferma; niezdara, ślamazara, gamoń\n4. tył czegoś, tyln" +
			"a część czegoś; zad"},
		//{"burgundowy", ""}, //needs following links support
		{"sławomira", "imię żeńskie\nimię męskie"},
		{"moje", "zaimek dzierżawczy będący odpowiednikiem zaimka osobowego " +
			"\"ja\"\npotocznie: to, co stanowi własność mówiącego lub ma z ni" +
			"m związek\npotocznie: kobieta, która wraz z mówiącym stanowi par" +
			"ę (np. małżeńską)\npotocznie: osoba, która wraz z mówiącym stano" +
			"wi parę (np. małżeńską)"},
		{"jajo", "1. jajko ptaków domowych wykorzystywane do przygotowywania " +
			"potraw\n2. żeńska komórka rozrodcza, komórka jajowa\n3. potoczni" +
			"e: przedmiot o owalnym kształcie\n4. potocznie: jądro męskie; jajco"},
		{"niewystepuje2137", "Nie występuje w słowniku"},
		{"", emptyQuery},
		{" \n \t ", emptyQuery},
	}

	for _, tt := range tests {
		str, want := tt.str, tt.want
		t.Run(str, func(t *testing.T) {
			t.Parallel()
			got, err := QuerySjp(str)
			if err != nil {
				t.Error(err)
			}
			if got != want {
				t.Errorf("got:\n%v\nwant:\n%v", got, want)
			}
		})
	}
}

func BenchmarkQuerySjpEmpty(b *testing.B) {
	if _, err := QuerySjp(""); err != nil {
		b.Error(err)
	}
}

func BenchmarkQuerySjpNonExistent(b *testing.B) {
	if _, err := QuerySjp("123"); err != nil {
		b.Error(err)
	}
}

func BenchmarkQuerySjpReal(b *testing.B) {
	if _, err := QuerySjp("test"); err != nil {
		b.Error(err)
	}
}
