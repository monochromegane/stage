package stage

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"golang.org/x/sync/errgroup"
)

type Stage struct {
	concurrency int
	seed        int64
	baseDir     string
	startAt     time.Time
}

func New(outDir string, concurrency int, seed int64) *Stage {
	if concurrency < 1 {
		concurrency = runtime.NumCPU()
	}
	return &Stage{
		concurrency: concurrency,
		seed:        seed,
		baseDir:     outDir,
	}
}

func (s *Stage) Run(iter int, newActorFn NewActorFn, newScenarioFn NewScenarioFn, callbackFn CallbackFn) error {
	s.startAt = time.Now()
	err := s.ensureOutDir()
	if err != nil {
		return err
	}
	sem := make(chan struct{}, s.concurrency)
	defer close(sem)
	progressCh := make(chan int, s.concurrency)
	doneCh := make(chan struct{})
	defer close(doneCh)

	go func() {
		for i := range progressCh {
			callbackFn(i)
		}
		doneCh <- struct{}{}
	}()

	rnd := rand.New(rand.NewSource(s.seed))
	eg, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < iter; i++ {
		sem <- struct{}{}
		aSeed := rnd.Int63()
		sSeed := rnd.Int63()
		i := i
		eg.Go(func() error {
			defer func() {
				progressCh <- i
				<-sem
			}()

			select {
			case <-ctx.Done():
				return nil
			default:
				actor := newActorFn(aSeed)
				scenario := newScenarioFn(sSeed)
				err := s.runWithLogFile(actor, scenario, iter, i, aSeed, sSeed)
				if err != nil {
					return err
				}
				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	close(progressCh)
	<-doneCh
	return nil
}

func (s *Stage) runWithLogFile(actor Actor, scenario Scenario, iter, i int, aSeed, sSeed int64) error {
	w, err := s.createLogFile(iter, i, aSeed, sSeed)
	if err != nil {
		return err
	}
	defer w.Close()

	err = s.run(actor, scenario, w)
	if err != nil {
		return err
	}
	return nil
}

func (s *Stage) run(actor Actor, scenario Scenario, w io.Writer) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	for scenario.Scan() {
		action, err := actor.Act(scenario.Line())
		if err != nil {
			return err
		}
		_, err = bw.WriteString(action.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Stage) outDirName() string {
	return filepath.Join(s.baseDir, fmt.Sprintf("%s-%d", s.startAt.Format("20060102150405"), s.seed))
}

func (s *Stage) ensureOutDir() error {
	return ensureDir(s.outDirName())
}

func (s *Stage) createLogFile(max, i int, aSeed, sSeed int64) (*os.File, error) {
	format := fmt.Sprintf("iter_%%0%dd-a_%%d-s_%%d.log", numberOfDigit(max))
	name := fmt.Sprintf(format, i, aSeed, sSeed)
	return os.Create(filepath.Join(s.outDirName(), name))
}
