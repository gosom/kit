package rollbar

import (
	"context"

	"github.com/gosom/kit/core"
	"github.com/rollbar/rollbar-go"
)

type rollbarErrorReporter struct {
}

func NewRollbarErrorReporter(token, environment, codeVersion, serverHost, serverRoot string) *rollbarErrorReporter {
	rollbar.SetToken(token)
	if len(environment) > 0 {
		rollbar.SetEnvironment(environment) // defaults to "development"
	}
	if len(codeVersion) > 0 {
		rollbar.SetCodeVersion(codeVersion) // optional Git hash/branch/tag (required for GitHub integration)
	}
	if len(serverHost) > 0 {
		rollbar.SetServerHost(serverHost) // optional override; defaults to hostname
	}
	if len(serverRoot) > 0 {
		rollbar.SetServerRoot(serverRoot) // path of project (required for GitHub integration and non-project stacktrace collapsing)
	}
	ans := rollbarErrorReporter{}
	return &ans
}

func (r rollbarErrorReporter) ReportError(ctx context.Context, args ...any) {
	var reportArgs []any
	custom := map[string]any{}
	requestID := core.RequestIDFromContext(ctx)
	if len(requestID) > 0 {
		custom["request_id"] = requestID
	}
	requestIP := core.IPFromContext(ctx)
	if requestIP != "" {
		custom["request_ip"] = requestIP
	}
	user := core.UserFromContext(ctx)
	if user != nil {
		ctx := rollbar.NewPersonContext(
			ctx,
			&rollbar.Person{
				Id: user.GetID(), Extra: user.GetExtra(),
			},
		)
		reportArgs = append(reportArgs, ctx)
	}
	for _, arg := range args {
		switch v := arg.(type) {
		case error:
			reportArgs = append(reportArgs, v)
			if st, ok := v.(core.StackTracer); ok {
				custom["stacktrace"] = st.StackTrace()
			}
		default:
			reportArgs = append(reportArgs, v)
		}
	}
	if len(custom) > 0 {
		reportArgs = append(reportArgs, custom)
	}
	rollbar.Error(reportArgs...)
}

func (r rollbarErrorReporter) ReportPanic(ctx context.Context, args ...any) {
	rollbar.Critical(args...)
	rollbar.Wait()
}

func (r rollbarErrorReporter) Close() {
	rollbar.Close()
}
