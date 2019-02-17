package gotabcmd

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

// fakeRunner is a runner used for testing
type fakeRunner struct {
	timeout   time.Duration
	requested []string // requested will contain the arguments passed

	returnOut   string // if not defined, will return asterisk: '*'. Useful to count number of executions
	returnError error
}

func (r *fakeRunner) run(action string, args ...string) (out string, err error) {
	r.requested = append(r.requested, action)
	r.requested = append(r.requested, args...)

	retval := "*"
	if r.returnOut != "" {
		retval = r.returnOut
	}
	return retval, r.returnError
}

func TestTableau_login(t *testing.T) {
	type fields struct {
		loggedIn bool
		user     string
		server   string
		e        TabRunner
	}
	type args struct {
		server   string
		username string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantCmd []string
		want    string
		wantErr bool
	}{
		{"ok",
			fields{
				loggedIn: false,
				e:        &fakeRunner{},
			},
			args{"123", "user", "pass"},
			[]string{"login", "-s", "https://123", "-u", "user", "-p", "pass"},
			"*", false,
		},
		{"not ok",
			fields{
				loggedIn: false,
				e: &fakeRunner{
					returnOut:   "something bad happened",
					returnError: errors.New("invalid login")},
			},
			args{"432", "user", "pass"},
			[]string{"login", "-s", "https://432", "-u", "user", "-p", "pass"},
			"something bad happened", true,
		},
		{"already logged in",
			fields{
				loggedIn: true,
				e:        &fakeRunner{},
			},
			args{"123", "user", "pass"},
			nil,
			"", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &Tableau{
				loggedIn: tt.fields.loggedIn,
				user:     tt.fields.user,
				server:   tt.fields.server,
				e:        tt.fields.e,
			}
			got, err := tb.login(tt.args.server, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tableau.login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Tableau.login() = %v, want %v", got, tt.want)
			}
			fake := tt.fields.e.(*fakeRunner)
			if !reflect.DeepEqual(fake.requested, tt.wantCmd) {
				t.Errorf("Tabrunner.requested() = %v, want %v", fake.requested, tt.wantCmd)
			}
		})
	}
}

func TestTableau_Logout(t *testing.T) {
	type fields struct {
		loggedIn bool
		user     string
		server   string
		e        TabRunner
	}
	tests := []struct {
		name    string
		fields  fields
		wantCmd []string
		want    string
		wantErr bool
	}{
		{"not logged int",
			fields{e: &fakeRunner{}},
			nil,
			"", false},
		{"ok",
			fields{loggedIn: true, e: &fakeRunner{}},
			[]string{"logout"},
			"*", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &Tableau{
				loggedIn: tt.fields.loggedIn,
				user:     tt.fields.user,
				server:   tt.fields.server,
				e:        tt.fields.e,
			}
			got, err := tb.Logout()
			if (err != nil) != tt.wantErr {
				t.Errorf("Tableau.Logout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Tableau.Logout() = %v, want %v", got, tt.want)
			}
			fake := tt.fields.e.(*fakeRunner)
			if !reflect.DeepEqual(fake.requested, tt.wantCmd) {
				t.Errorf("Tabrunner.requested() = %v, want %v", fake.requested, tt.wantCmd)
			}
		})
	}
}

func Test_addHTTPS(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no https", args{"myserver"}, "https://myserver"},
		{"https", args{"https://myserver"}, "https://myserver"},
		{"http", args{"http://myserver"}, "https://myserver"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addHTTPS(tt.args.uri); got != tt.want {
				t.Errorf("addHTTPS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableau_RefreshExtracts(t *testing.T) {
	type fields struct {
		loggedIn bool
		user     string
		server   string
		e        TabRunner
	}
	type args struct {
		datasets []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"not logged in",
			fields{
				loggedIn: false,
				e:        &fakeRunner{},
			},
			args{[]string{"dataset"}},
			"", true,
		},
		{"0 datasets",
			fields{
				loggedIn: true,
				e:        &fakeRunner{},
			},
			args{[]string{}},
			"", true,
		},
		{"1 dataset",
			fields{
				loggedIn: true,
				e:        &fakeRunner{},
			},
			args{[]string{"dataset1"}},
			"*", false,
		},
		{"2 datasets",
			fields{
				loggedIn: true,
				e:        &fakeRunner{},
			},
			args{[]string{"dataset1", "dataset2"}},
			"**", false,
		},
		{"2 datasets error",
			fields{
				loggedIn: true,
				e:        &fakeRunner{returnError: errors.New("error refreshing")},
			},
			args{[]string{"dataset1", "dataset2"}},
			"*", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &Tableau{
				loggedIn: tt.fields.loggedIn,
				user:     tt.fields.user,
				server:   tt.fields.server,
				e:        tt.fields.e,
			}
			got, err := tb.RefreshExtracts(tt.args.datasets...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tableau.RefreshExtracts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Tableau.RefreshExtracts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableau_LoginOnline(t *testing.T) {
	type fields struct {
		loggedIn bool
		user     string
		server   string
		e        TabRunner
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"login ok", fields{loggedIn: false, e: &fakeRunner{}},
			args{"user", "pass"},
			"*", false},
		{"login not ok", fields{loggedIn: false, e: &fakeRunner{returnError: errors.New("login error")}},
			args{"user", "pass"},
			// number of stars should match number of defined online servers
			"*******", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &Tableau{
				loggedIn: tt.fields.loggedIn,
				user:     tt.fields.user,
				server:   tt.fields.server,
				e:        tt.fields.e,
			}
			got, err := tb.LoginOnline(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tableau.LoginOnline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Tableau.LoginOnline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableau_Login(t *testing.T) {
	type fields struct {
		loggedIn bool
		user     string
		server   string
		e        TabRunner
	}
	type args struct {
		server   string
		username string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"ok",
			fields{
				loggedIn: false,
				e:        &fakeRunner{},
			},
			args{"123", "user", "pass"},
			"*", false,
		},
		{"not ok",
			fields{
				loggedIn: false,
				e: &fakeRunner{
					returnOut:   "something bad happened",
					returnError: errors.New("invalid login")},
			},
			args{"432", "user", "pass"},
			"something bad happened", true,
		},
		{"already logged in",
			fields{
				loggedIn: true,
				e:        &fakeRunner{},
			},
			args{"123", "user", "pass"},
			"", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &Tableau{
				loggedIn: tt.fields.loggedIn,
				user:     tt.fields.user,
				server:   tt.fields.server,
				e:        tt.fields.e,
			}
			got, err := tb.Login(tt.args.server, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tableau.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Tableau.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTableau(t *testing.T) {
	tb := NewTableau(400)
	tb.server = "old server"
	tb.user = "peter"

	// make sure that calling NewTableau returns the same instance
	newTb := NewTableau(500)
	newTb.server = "new server"
	newTb.user = "rupert"
	if !reflect.DeepEqual(newTb, tb) {
		t.Error("tb and newTb are not the same")
	}
}
