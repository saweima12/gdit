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

func NewTestConfig(ctx gdit.InvokeCtx) (*testConfig, error) {
	return &testConfig{
		DomainUrl: "http://example.com",
	}, nil
}

// Define testRepository
type testRepo struct {
	cfg *testConfig
}

func NewTestRepo(ctx gdit.InvokeCtx) (*testRepo, error) {
	fmt.Println("OnTestRepo create.")

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

func NewTestServ(ctx gdit.InvokeCtx) (TestService, error) {
	serv := &testService{}
	fmt.Println("OnTestService Create")

	serv.repo = gdit.MustInject[*testRepo](ctx)
	return serv, nil
}

func TestProvide(t *testing.T) {
	app := getTestApp()

	gdit.ProvideValue[*testConfig](&testConfig{
		DomainUrl: "http://example.com",
	}).Attach(app)

	gdit.InvokeFunc(app, func(ctx gdit.InvokeCtx) error {
		// Test inject provide value
		cfg := gdit.MustInject[*testConfig](ctx)
		if cfg.DomainUrl != "http://example.com" {
			t.Fail()
		}

		// Test inject named lazy provider with interface.
		serv := gdit.MustInjectNamed[TestService](ctx, "TestService")
		serv2 := gdit.MustInjectNamed[TestService](ctx, "TestService")
		t.Run("The serv and serv2 should be equals", func(ct *testing.T) {
			if serv != serv2 {
				ct.Fail()
			}
		})

		// test factory provider.
		repo1 := gdit.MustInject[*testRepo](ctx)
		repo2 := gdit.MustInject[*testRepo](ctx)
		t.Run("The repositories repo1 and repo2 should be distinct.", func(ct *testing.T) {
			if &repo1 == &repo2 {
				ct.Fail()
			}
		})

		app.Startup()
		app.Teardown()
		return nil
	})
}

func TestOverwriteScope(t *testing.T) {
}

func getTestApp() gdit.App {
	app := gdit.New()

	gdit.ProvideFactory[*testRepo](NewTestRepo).
		Attach(app)

	gdit.Provide[TestService](NewTestServ).
		WithName("TestService").
		Attach(app)

	return app
}
