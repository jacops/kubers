package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/go-hclog"
)

type testChval struct {
	c   chan *Secret
	val *Secret
}

type testChErrval struct {
	c   chan error
	val error
}

func TestMergeSecretsChannelsCanTransmitSecrets(t *testing.T) {
	tests := []struct {
		name           string
		chValPairSlice []testChval
	}{
		{name: "no-secrets"},
		{name: "one-secret", chValPairSlice: []testChval{{val: &Secret{Name: "secret1"}}}},
		{name: "one-secret", chValPairSlice: []testChval{{val: &Secret{Name: "secret1"}}, {val: &Secret{Name: "secret2"}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipeline := &Pipeline{
				provider: &DummyProvider{},
				writer:   NewDummyPathWriter(),
				logger: hclog.New(&hclog.LoggerOptions{
					Name: "handler",
				}),
			}

			var clist []<-chan *Secret
			var secrets []*Secret

			for _, pair := range tt.chValPairSlice {
				pair.c = make(chan *Secret)
				clist = append(clist, pair.c)
				secrets = append(secrets, pair.val)

				go func(pair testChval) {
					defer close(pair.c)
					pair.c <- pair.val
				}(pair)
			}

			merged := pipeline.mergeSecretsChannels(clist)

			msgCounter := 0
			for sc := range merged {
				msgCounter++
				found := false
				for _, b := range secrets {
					if b.Name == sc.Name {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("Could not found transmitted value: %s", sc.Name)
				}
			}

			if msgCounter != len(secrets) {
				t.Errorf("Not all secrets have been transmitted")
			}
		})
	}
}

func TestMergeErrorChannelsCanTransmitErrors(t *testing.T) {
	tests := []struct {
		name string
		err  []testChErrval
	}{
		{name: "no-errors"},
		{name: "one-error", err: []testChErrval{{val: errors.New("error1")}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipeline := &Pipeline{
				provider: &DummyProvider{},
				writer:   NewDummyPathWriter(),
				logger: hclog.New(&hclog.LoggerOptions{
					Name: "handler",
				}),
			}

			var clist []<-chan error
			var errs []error

			for _, pair := range tt.err {
				pair.c = make(chan error)
				clist = append(clist, pair.c)
				errs = append(errs, pair.val)

				go func(pair testChErrval) {
					defer close(pair.c)
					pair.c <- pair.val
				}(pair)
			}

			merged := pipeline.mergeErrorChannels(clist)

			msgCounter := 0
			for sc := range merged {
				msgCounter++
				found := false
				for _, b := range errs {
					if b.Error() == sc.Error() {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("Could not found transmitted error: %s", sc.Error())
				}
			}

			if msgCounter != len(errs) {
				t.Errorf("Not all errors have been transmitted")
			}
		})
	}
}

func TestFanOut(t *testing.T) {
	tests := []struct {
		secrets []*Secret
	}{
		{secrets: []*Secret{{Name: "secret1"}}},
	}

	msgCounter := 0
	for _, tt := range tests {
		t.Run("test", func(t *testing.T) {
			pipeline := &Pipeline{
				provider: &DummyProvider{},
				writer:   NewDummyPathWriter(),
				logger: hclog.New(&hclog.LoggerOptions{
					Name: "handler",
				}),
			}
			fanoutc := pipeline.fanOut(tt.secrets)
			for s := range fanoutc {
				msgCounter++
				found := false
				for _, b := range tt.secrets {
					if b.Name == s.Name {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("Could not found fanned out secret: %s", s.Name)
				}
			}
			if msgCounter != len(tt.secrets) {
				t.Errorf("Not all secrets have been transmitted")
			}
		})
	}
}

func TestSecretsCanGoThroughPipeline(t *testing.T) {
	tests := []struct {
		secrets       []*Secret
		workerNumbers int
	}{
		{[]*Secret{{Name: "secret1"}}, 1},
		{[]*Secret{{Name: "secret1"}, {Name: "secret2"}}, 1},
		{[]*Secret{{Name: "secret1"}}, 5},
		{[]*Secret{{Name: "secret1"}, {Name: "secret2"}}, 5},
	}

	for _, tt := range tests {
		t.Run("test", func(t *testing.T) {
			pipeline := &Pipeline{
				provider: &DummyProvider{},
				writer:   NewDummyPathWriter(),
				logger: hclog.New(&hclog.LoggerOptions{
					Name: "handler",
				}),
				workersNumber: tt.workerNumbers,
			}
			err := pipeline.Do(context.Background(), tt.secrets)

			if err != nil {
				t.Errorf("did not expect any other errors: %s", err)
			}
		})
	}
}

type DummyProvider struct{}
type DummyPathWriter struct{}

func (dd *DummyProvider) GetSecret(ctx context.Context, secretURL string) (string, error) {
	return "dummy-secret", nil
}

func (w *DummyPathWriter) WriteSecret(secret *Secret) error {
	return nil
}

func NewDummyPathWriter() *DummyPathWriter {
	return &DummyPathWriter{}
}
