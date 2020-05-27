package stage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func NewTestScenarioFn(seed int64) Scenario {
	return &TestScenario{}
}

type TestScenario struct {
	counter int
}

func (s *TestScenario) Scan() bool {
	if s.counter < 10 {
		s.counter++
		return true
	}
	return false
}

func (s *TestScenario) Line() Line {
	return Line{"counter": s.counter}
}

func NewTestActorFn(seed int64) Actor {
	return TestActor{}
}

type TestActor struct{}

func (a TestActor) Act(line Line) (Action, error) {
	counter := line["counter"].(int)
	return TestAction{action: counter}, nil
}

type TestAction struct {
	action int
}

func (a TestAction) String() string { return fmt.Sprintf("%d", a.action) }

func TestStageRunOnce(t *testing.T) {
	stage := New("", 2, 20)
	actor := TestActor{}
	scenario := &TestScenario{}
	out := new(bytes.Buffer)
	err := stage.run(actor, scenario, out)
	if err != nil {
		t.Errorf("stage.run shoud not happen the error %v", err)
	}

	expected := "12345678910"
	if actual := out.String(); actual != expected {
		t.Errorf("stage.run shoud output %s, but %s", expected, actual)
	}
}

func TestStageRunWitLogFileOnce(t *testing.T) {
	dir, _ := ioutil.TempDir("", "test")
	defer os.RemoveAll(dir)

	stage := New(dir, 2, 20)
	startAt, _ := time.Parse("2006/01/02 15:04:05", "2020/05/27 15:14:13")
	stage.startAt = startAt
	stage.ensureOutDir()

	actor := TestActor{}
	scenario := &TestScenario{}
	err := stage.runWithLogFile(actor, scenario, 1000, 10, 30, 40)
	if err != nil {
		t.Errorf("stage.runWithLogFile shoud not happen the error %v", err)
	}

	log := filepath.Join(dir, "20200527151413-20", "iter_0010-a_30-s_40.log")
	expected := "12345678910"

	b, err := ioutil.ReadFile(log)
	if err != nil {
		t.Errorf("stage.runWithLogFile should create log %s, but not exist", log)
	}
	actual := string(b)
	if actual != expected {
		t.Errorf("stage.run shoud output %s, but %s", expected, actual)
	}
}

func TestStageRun(t *testing.T) {
	dir, _ := ioutil.TempDir("", "test")
	defer os.RemoveAll(dir)

	stage := New(dir, 2, 20)
	err := stage.Run(10, NewTestActorFn, NewTestScenarioFn, NoOpeCallbackFn)
	if err != nil {
		t.Errorf("stage.Run shoud not happen the error %v", err)
	}

	expected := "12345678910"
	timestamps, _ := ioutil.ReadDir(dir)
	logs, _ := ioutil.ReadDir(filepath.Join(dir, timestamps[0].Name()))
	for _, log := range logs {
		b, err := ioutil.ReadFile(filepath.Join(dir, timestamps[0].Name(), log.Name()))
		if err != nil {
			t.Errorf("stage.Run should create log %s, but not exist", log.Name())
		}
		actual := string(b)
		if actual != expected {
			t.Errorf("stage.Run shoud output %s, but %s", expected, actual)
		}
	}
}
