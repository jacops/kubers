package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/go-hclog"
)

//Pipeline ...
type Pipeline struct {
	logger        hclog.Logger
	writer        SecretsWriter
	provider      Provider
	workersNumber int
}

//NewPipeline ...
func NewPipeline(writer SecretsWriter, provider Provider, logger hclog.Logger, workersNumber int) *Pipeline {
	return &Pipeline{
		writer:        writer,
		provider:      provider,
		logger:        logger,
		workersNumber: workersNumber,
	}
}

//Do ...
func (p *Pipeline) Do(ctx context.Context, secrets []*Secret) error {
	var errc <-chan error
	var errcList []<-chan error
	var pipe1, pipe2 <-chan *Secret
	var pipe2List []<-chan *Secret

	p.logger.Debug(fmt.Sprintf("Processing pipeline started with %d workers.", p.workersNumber))

	pipe1 = fanOut(secrets)

	for i := 1; i <= p.workersNumber; i++ {
		pipe2, errc := p.populate(ctx, pipe1)
		pipe2List = append(pipe2List, pipe2)
		errcList = append(errcList, errc)
	}

	pipe2 = mergeSecretsChannels(pipe2List)

	for i := 1; i <= p.workersNumber; i++ {
		errc = p.store(ctx, pipe2)
		errcList = append(errcList, errc)
	}

	return waitForPipeline(errcList)
}

func waitForPipeline(cs []<-chan error) error {
	errorc := mergeErrorChannels(cs)
	for err := range errorc {
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pipeline) populate(ctx context.Context, in <-chan *Secret) (<-chan *Secret, <-chan error) {
	out := make(chan *Secret)
	errc := make(chan error, 1)

	go func() {
		defer func() {
			close(out)
			close(errc)
		}()
		for secret := range in {
			value, err := p.provider.GetSecret(ctx, secret.URL)
			if err != nil {
				errc <- err
			}
			secret.Value = value

			select {
			case out <- secret:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, errc
}

func (p *Pipeline) store(ctx context.Context, in <-chan *Secret) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		for secret := range in {
			if err := p.writer.WriteSecret(secret); err != nil {
				errc <- err
				return
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	return errc
}

func fanOut(secrets []*Secret) <-chan *Secret {
	out := make(chan *Secret)
	go func() {
		for _, secret := range secrets {
			out <- secret
		}
		close(out)
	}()
	return out
}

func mergeErrorChannels(cs []<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error, len(cs))

	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func mergeSecretsChannels(cs []<-chan *Secret) <-chan *Secret {
	var wg sync.WaitGroup

	out := make(chan *Secret)
	send := func(c <-chan *Secret) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go send(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
