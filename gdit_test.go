package gdit_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

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

	ctx.OnStart(func(startCtx gdit.StartCtx) error {
		fmt.Println("OnRepoStart")
		return nil
	})

	ctx.OnStop(func(stopCtx gdit.StopCtx) error {
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

func NewTestServ(ctx gdit.InvokeCtx) (TestService, error) {
	serv := &testService{}
	fmt.Println("OnTestService Create")

	serv.repo = gdit.MustInject[*testRepo](ctx)
	return serv, nil
}

func TestProvide(t *testing.T) {
	app := getTestApp()

	gdit.InvokeFunc(app, func(ctx gdit.InvokeCtx) error {

		// Test inject provide value
		t.Run("cfg should have value", func(t *testing.T) {
			cfg := gdit.MustInject[*testConfig](ctx)
			if cfg.DomainUrl != "http://example.com" {
				t.Fail()
			}
		})

		t.Run("TestService not registed to a type, that will trigger panic", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			serv := gdit.MustInject[TestService](ctx)
			if serv != nil {
				t.Fail()
			}
		})

		t.Run("TestServicess not exists, that will trigger panic", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			serv := gdit.MustInjectNamed[TestService](ctx, "TestServicess")
			if serv != nil {
				t.Fail()
			}
		})

		t.Run("TestServ use type should be nil", func(ct *testing.T) {
			serv, err := gdit.Inject[TestService](ctx)
			if err == nil || serv != nil {
				t.Fail()
			}
		})

		t.Run("TestServ use name should have instance", func(ct *testing.T) {
			serv, err := gdit.InjectNamed[TestService](ctx, "TestService")
			if err != nil || serv == nil {
				t.Fail()
			}
		})

		t.Run("The serv and serv2 should be equals", func(ct *testing.T) {
			// Test inject named lazy provider with interface.
			serv := gdit.MustInjectNamed[TestService](ctx, "TestService")
			serv2 := gdit.MustInjectNamed[TestService](ctx, "TestService")
			if serv != serv2 {
				ct.Fail()
			}
		})

		app.Startup()

		// test factory provider.
		t.Run("The repositories repo1 and repo2 should be distinct.", func(ct *testing.T) {
			repo1 := gdit.MustInject[*testRepo](ctx)
			repo2 := gdit.MustInject[*testRepo](ctx)
			if &repo1 == &repo2 {
				ct.Fail()
			}
		})

		t.Run("Invoke should be success", func(t *testing.T) {
			repo, err := gdit.Invoke[*testRepo](app, func(ic gdit.InvokeCtx) (*testRepo, error) {
				return &testRepo{}, nil
			})
			if repo == nil || err != nil {
				t.Fail()
			}
		})

		t.Run("InvokeProvide should be success", func(t *testing.T) {
			repo, err := gdit.InvokeProvide[*testRepo](app, func(ic gdit.InvokeCtx) (*testRepo, error) {
				return &testRepo{}, nil
			})

			if repo == nil || err != nil {
				t.Fail()
			}
		})

		t.Run("InvokeProvide return a error should be failed", func(t *testing.T) {
			repo, err := gdit.InvokeProvide[*testRepo](app, func(ic gdit.InvokeCtx) (*testRepo, error) {
				return nil, errors.New("failed")
			})

			if repo != nil || err == nil {
				t.Fail()
			}
		})

		app.Teardown()
		return nil
	})
}

func TestInvokeWithHook(t *testing.T) {
	app := getTestApp()
	fmt.Println("Start to test InvokeWithHook")

	testCh := make(chan struct{})

	gdit.InvokeFunc(app, func(ic gdit.InvokeCtx) error {
		ic.OnStart(func(startCtx gdit.StartCtx) error {
			testCh <- struct{}{}
			return nil
		})
		return nil
	})
	go app.Startup()
	ok := false
	select {
	case <-testCh:
		ok = true
	case <-time.After(time.Second * 2):
		ok = false
	}

	if !ok {
		t.Fail()
	}

	fmt.Println("Test InvokeWithHook end")

}

func getTestApp() gdit.App {
	app := gdit.New().SetLogLevel(gdit.LOG_DEBUG)

	gdit.ProvideValue[*testConfig](&testConfig{
		DomainUrl: "http://example.com",
	}).Attach(app)

	gdit.ProvideFactory[*testRepo](NewTestRepo).
		Attach(app)

	gdit.Provide[TestService](NewTestServ).
		WithName("TestService").
		Attach(app)

	return app
}
