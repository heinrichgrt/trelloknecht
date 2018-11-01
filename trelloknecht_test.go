package main

import "testing"

func Test_shortenStringIfToLong(t *testing.T) {
	configuration["headLineCharsSkip"] = "20"
	//configuration["headLineMaxChars"] = "25"
	type args struct {
		instring string
	}
	//
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Test no shortening", args{"Ein kurzer String"}, "Ein kurzer String"},
		{"Max reached", args{"vier Vier viEr vieRX"}, "vier Vier viEr vieRX"},
		{"Short one word", args{"vier Vier viEr vieR VIERX"}, "vier Vier viEr vieR..."},

		// TODO: Add test case
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shortenStringIfToLong(tt.args.instring); got != tt.want {
				t.Errorf("shortenStringIfToLong() = ->%v<-, want ->%v<-", got, tt.want)
			}
		})
	}
}
