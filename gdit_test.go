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

func NewTestConfig(ctx gdit.Context) (*testConfig, error) {
	ctx.OnStart(func(ctx gdit.Context) error {
		fmt.Println("On TestConfig Start")
		return nil
	})

	ctx.OnStop(func(ctx gdit.Context) error {
		fmt.Println("On TestConfig stop")
		return nil
	})

	return &testConfig{
		DomainUrl: "http://example.com",
	}, nil
}

// Define testRepository
type testRepo struct {
	cfg *testConfig
}

func NewTestRepo(ctx gdit.Context) (*testRepo, error) {
	fmt.Println("OnTestRepo create.")

	ctx.OnStart(func(ctx gdit.Context) error {
		fmt.Println("On TestRepo Start.")
		return nil
	})

	ctx.OnStop(func(ctx gdit.Context) error {
		fmt.Println("On TestRepo Stop.")
		return nil
	})

	gdit.MustInject[*testConfig](ctx)
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

func NewTestServ(ctx gdit.Context) (TestService, error) {
	serv := &testService{}
	fmt.Println("OnTestService Create")

	ctx.OnStart(func(ctx gdit.Context) error {
		fmt.Println("On TestServ Start.")
		return nil
	})

	ctx.OnStop(func(ctx gdit.Context) error {
		fmt.Println("On TestServ Stop.")
		return nil
	})
	serv.repo = gdit.MustInject[*testRepo](ctx)

	return serv, nil
}

func TestGdit(t *testing.T) {
	app := gdit.New()

	// Test provideValue.
	gdit.ProvideValue[*testConfig](&testConfig{
		DomainUrl: "http://example.com",
	}).Attach(app)

	gdit.Provide[*testConfig](NewTestConfig).Attach(app)

	// Test provide
	gdit.Provide[*testRepo](NewTestRepo).Attach(app)

	// Test provide with name
	helloScope := app.GetScope("Hello")
	gdit.Provide[TestService](NewTestServ).
		Attach(app)

	// Invoke

	app.Setup()
	fmt.Println(app.Scope)
	fmt.Println(helloScope)

}

func Runner(ctx gdit.Context) error {

	return nil
}
