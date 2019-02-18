package gotabcmd

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// tableau actions
const (
	actLogin           = "login"
	actLogout          = "logout"
	actRefreshExtracts = "refreshextracts"
)

// TabRunner is an interface for tableau command runner
type TabRunner interface {
	run(action string, args ...string) (out string, err error)
}

// Tableau keeps current connection paramters.  Depending on the TabExecer
// could be thread usafe.  If TabExecer is Tabcmd (which it is) then certainly
// thread unsafe!
type Tableau struct {
	loggedIn bool
	user     string
	server   string

	e TabRunner
}

// Singleton
var (
	gTb   *Tableau
	gOnce sync.Once
)

// Errors
var (
	ErrAlreadyLoggedIn = errors.New("already logged in")
	ErrNotLoggedIn     = errors.New("not logged in")
	ErrEmptyDs         = errors.New("empty dataset, nothing to do")
)

// ServerList is a list of servers that will be used to attempt the login
// taken from the following URL:
// https://onlinehelp.tableau.com/current/pro/desktop/en-us/publish_tableau_online_ip_authorization.htm
var ServerList = []string{
	"dub01.online.tableau.com",
	"eu-west-1a.online.tableau.com",
	"10ax.online.tableau.com",
	"10ay.online.tableau.com",
	"10az.online.tableau.com",
	"us-east-1.online.tableau.com",
	"us-west-2b.online.tableau.com",
}

// NewTableau returns a local Tableau Instance
func NewTableau(commandTimeout time.Duration) *Tableau {
	// make sure that there could be only one
	gOnce.Do(func() { gTb = &Tableau{e: newTabcmd(commandTimeout)} })
	return gTb
}

func (t Tableau) String() string {
	return fmt.Sprintf("Tableau: %s@%s", t.user, t.server)
}

// login attempts to login onto the server
func (t *Tableau) login(server, username, password string) (string, error) {
	if t.loggedIn {
		return "", ErrAlreadyLoggedIn
	}

	server = addHTTPS(server)
	t.server = server
	t.user = username

	log.Printf("attempting to log in as %s", t)

	out, err := t.e.run(actLogin, "-s", server, "-u", username, "-p", password)

	t.loggedIn = (err == nil)

	return out, err
}

// addHTTPS adds https:// prefix if it's not present
func addHTTPS(uri string) (addr string) {
	const (
		http       = "http://"
		securehttp = "https://"
	)
	switch {
	case strings.HasPrefix(uri, http):
		addr = strings.Replace(uri, http, securehttp, 1)
	case !strings.HasPrefix(uri, securehttp):
		addr = securehttp + uri
	default:
		addr = uri
	}
	return
}

// Login attempts to login onto the specified server
func (t *Tableau) Login(server, username, password string) (string, error) {
	return t.login(server, username, password)
}

// LoginOnline logs into online instance.  Attemts to login into each
// online server in ServerList in order, until finally logs in (or fails)
// completely.
func (t *Tableau) LoginOnline(username, password string) (string, error) {
	var (
		out string
		err error
		buf bytes.Buffer
	)
	// iterate trough the server list
	for _, srv := range ServerList {
		out, err = t.login(srv, username, password)
		buf.WriteString(out)
		if err == nil {
			break
		}
	}
	return buf.String(), err
}

// Logout does what is should.
func (t *Tableau) Logout() (string, error) {
	log.Print("logout")
	if !t.loggedIn {
		return "", nil
	}

	out, err := t.e.run(actLogout)
	if err != nil {
		return out, err
	}
	t.loggedIn = false

	return out, nil
}

// RefreshExtracts launches extract refresh of one or more datasources.
func (t *Tableau) RefreshExtracts(datasets ...string) (string, error) {
	if !t.loggedIn {
		return "", ErrNotLoggedIn
	}
	if len(datasets) == 0 {
		return "", ErrEmptyDs
	}

	var buf bytes.Buffer
	for _, ds := range datasets {
		log.Printf("starting refresh of dataset %s", ds)
		out, err := t.e.run(actRefreshExtracts, "--datasource", ds)
		buf.WriteString(out)
		if err != nil {
			return buf.String(), err
		}
	}
	return buf.String(), nil
}
