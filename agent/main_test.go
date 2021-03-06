package agent

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/recipe"
)

var (
	iniFile      = "../tests/configs/config_regular.ini"
	path         = &iniFile
	cfgWrap      = config.NewConfigs()
	defaultAgent = New("test", cfgWrap.New("test", *path), false)
)

func TestGetterers(t *testing.T) {
	defaultAgent.GetName()
	defaultAgent.GetContext()
	defaultAgent.GetLogs()
	defaultAgent.GetConfig()
	defaultAgent.GetMutex()

	defaultAgent.GetStdout()
	defaultAgent.GetStdin()

	defaultAgent.SetStdin(os.Stdin)
	defaultAgent.SetStdout(os.Stdout)

	h := defaultAgent.DefaultAPIHandles()
	if len(h) > 2 {
		t.Fail()
	}

}

func TestCommands(t *testing.T) {
	req, err := http.NewRequest("GET", "/client/workers/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	defaultAgent.WorkersList(rr, req)
	r, err := recipe.Build("../tests/etc/leprechaun/recipes/schedule.yml")
	if err != nil {
		t.Fail()
	}
	// create worker
	_, err = defaultAgent.CreateWorker(&r)
	if err != nil {
		t.Fail()
	}

	defaultAgent.WorkersList(rr, req)
	// not existent worker

	req, err = http.NewRequest("GET", "/client/workers/kill?name=schedule", nil)
	if err != nil {
		t.Fatal(err)
	}

	defaultAgent.KillWorker(rr, req)

	req, err = http.NewRequest("GET", "/client/workers/kill?name=jobber", nil)
	if err != nil {
		t.Fatal(err)
	}
	defaultAgent.KillWorker(rr, req)
}
