package bot

import "testing"

func TestPudzianConverter(t *testing.T) {
	t.Parallel()
	errMsg := "Błędne dane: podaj liczbę i jednostkę (KM,PM,W)\n" +
		"Na przykład: 1337 KM"

	tests := []struct{ str, want string }{
		{"  120   KM  ", "120 KM to 4.46 PM"},
		{"  120  ", errMsg},
		{"", errMsg},
		{"120 121 122 123", errMsg},
		{"120 121 PM", errMsg},
		{"120 121", errMsg},
		{"120 HP", "120 HP to 4.46 PM"},
		{"120 KM", "120 KM to 4.46 PM"},
		{"120 PM", "120 PM to 3227.16 KM"},
		{"120 W", "120 W to 0.01 PM"},
		{"120", errMsg},
		{"120HP", "120 KM to 4.46 PM"},
		{"120KM", "120 KM to 4.46 PM"},
		{"120PM", "120 PM to 3227.16 KM"},
		{"120W", "120 W to 0.01 PM"},
		{"HP", errMsg},
		{"KM", errMsg},
		{"PM", errMsg},
		{"W", errMsg},
		{"\n121KM \t ", "121 KM to 4.50 PM"},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			got := pudzianConverter(tt.str)
			if got != tt.want {
				t.Errorf("got:\n%v\nwant:\n%v", got, tt.want)
			}
		})
	}
}
