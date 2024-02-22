package gdit_test

import (
	"fmt"
	"testing"

	"github.com/saweima12/gdit"
)

// Define testConfig
type testConfig struct {
	DomainUrl string
}

// Define testRepository
type testRepo struct{}

func NewTestRepo(ctx *gdit.Context) (*testRepo, error) {
	fmt.Println("OnTestRepo create.")
	return &testRepo{}, nil
}

// Define TestServ
type TestService interface {
	Run()
}
type testService struct {
	repo *testRepo
}

func (t *testService) Run() {

}

func NewTestServ(ctx *gdit.Context) (TestService, error) {
	serv := &testService{}
	fmt.Println("OnTestService Create")
	serv.repo = gdit.MustInject[*testRepo](ctx)
	return serv, nil
}

func TestGdit(t *testing.T) {
	app := gdit.New()

	// Test provideValue.
	gdit.ProvideValue[*testConfig](&testConfig{
		DomainUrl: "http://example.com",
	}).Attach(app)

	// Test provide
	gdit.Provide[*testRepo](NewTestRepo).Attach(app)

	// Test provide with name
	gdit.Provide[TestService](NewTestServ).
		WithName("TestServ").
		Attach(app)

	// Invoke
	gdit.Invoke(app, NewTestServ)
	gdit.InvokeFunc(app, Runner)

	app.Run()
	fmt.Println(app)
}

func Runner(ctx *gdit.Context) error {
	cfg := gdit.MustInject[*testConfig](ctx)
	fmt.Println(cfg)
	return nil
}
