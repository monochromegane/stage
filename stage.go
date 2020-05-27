package stage

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sync/errgroup"
)

type Stage struct {
	baseDir string
	startAt time.Time
}

func New(outDir string) *Stage {
	return &Stage{
		baseDir: outDir,
	}
}

func (s *Stage) Run(iter int, newActorFn NewActorFn, scenario Scenario) error {
	s.startAt = time.Now()
	err := s.ensureOutDir()
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < iter; i++ {
		i := i
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				actor := newActorFn()
				err := s.runWithLogFile(actor, scenario, iter, i)
				if err != nil {
					return err
				}
				return nil
			}
		})
	}
	return eg.Wait()
}

func (s *Stage) runWithLogFile(actor Actor, scenario Scenario, iter, i int) error {
	w, err := s.createLogFile(iter, i)
	if err != nil {
		return err
	}
	defer w.Close()

	err = s.run(actor, scenario, w)
	if err != nil {
		return err
	}
}

func (s *Stage) run(actor Actor, scenario Scenario, w io.Writer) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	for scenario.Scan() {
		action, err := actor.Act(scenario.Line())
		if err != nil {
			return err
		}
		_, err = bw.WriteString(action)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Stage) outDirName() string {
	return filepath.Join(s.baseDir, s.startAt.Format("20060102150405"))
}

func (s *Stage) ensureOutDir() error {
	return ensureDir(s.outDirName())
}

func (s *Stage) createLogFile(max, i int) (io.Writer, error) {
	format := fmt.Sprintf("iter-%%0%dd.log", numberOfDigit(max))
	name := fmt.Sprintf(format, i)
	return os.Create(filepath.Join(s.outDirName(), name))
}