package main

import (
	"testing"
	"time"
)

func TestResultset_execCommand(t *testing.T) {

	type fields struct {
		OSCommand            string
		CommandArgs          []string
		Stdout               string
		Stderr               string
		CmdStarttime         time.Time
		CMDStoptime          time.Time
		DurationSecounds     int
		SuccessfullExecution bool
		ErrorStr             string
	}
	tests := []struct {
		name    string
		fields  fields
		compare fields
	}{
		{"return 0", fields{OSCommand: "/usr/bin/true"}, fields{SuccessfullExecution: true}},
		{"return 1", fields{OSCommand: "/usr/bin/false"}, fields{SuccessfullExecution: false}},
		{"stdout", fields{OSCommand: "/bin/echo", CommandArgs: []string{"oh", "yeah"}},
			fields{SuccessfullExecution: true, Stdout: "oh yeah\n"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resultset{
				OSCommand:            tt.fields.OSCommand,
				CommandArgs:          tt.fields.CommandArgs,
				Stdout:               tt.fields.Stdout,
				Stderr:               tt.fields.Stderr,
				SuccessfullExecution: tt.fields.SuccessfullExecution,
				ErrorStr:             tt.fields.ErrorStr,
			}
			r.execCommand()
			// switch debug on

			if r.SuccessfullExecution != tt.compare.SuccessfullExecution {
				t.Errorf("exec command failed: got: %v want: %v", tt.fields.SuccessfullExecution, tt.compare.SuccessfullExecution)
			}

			if r.Stdout != tt.compare.Stdout {
				t.Errorf("execCommand() = got: ->%v<-, want ->%v<-", r.Stdout, tt.compare.Stdout)
			}

		},
		)
	}

}
