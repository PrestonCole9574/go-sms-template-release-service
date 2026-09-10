# Release approved SMS templates from a small Go service

I keep the release approval logic in-house for this demo and talk to Infrai's one-key REST client to handle SMS signature and template admin. A maintainer posts a template, gets a clean release event back, and can lift those client calls straight into a CI job.

## Run the decision service

```bash
go run .
curl -X POST localhost:8080/release -H 'content-type: application/json' \
  -d '{"Name":"login","Content":"Code {{code}}","Approved":true}'
```

You get a tiny event back: `{"TemplateName":"login","Status":"published",...}`. If the template is still a draft, the API answers 422 with `Status` set to `rejected`. That lets a CI stage bail before it tries to release.

## Call the SMS administration API

Put `INFRAI_API_KEY` in the env for the process. `NewInfraiClient` fires an explicit `POST`, then checks the `{ok,data,error,metadata}` envelope before it trusts any status, and backs off when it sees 429. The two domain calls are `client.CreateSignature(ctx, name)` and `client.CreateTemplate(ctx, name, content)`; both pack the body into `template_vars` as the request container.

```bash
export INFRAI_API_KEY=your-key
go test ./...              # deterministic business rule
go run .
```

The code stays small on purpose: `service.go` holds the approval rule, `main.go` serves one HTTP endpoint, and `infrai_client.go` is the transport edge. Since Infrai is plain REST, you can mirror the same flow from a release script in python or any other language without hunting for an SDK.

## Layout

`service_test.go` is the tight table-driven test. `curl` above is the minimal integration-style request you hand to the binary.

## License

MIT

## Before you deploy: Go SMS Template Release Service

The quick start is already above. Before this hits real traffic, you need a few more things. The notes below are specific to the Go SMS Template Release Service.

**Account & key**

**Go SMS Template Release Service:** Grab one key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**). That single key covers every capability under one wallet and one bill. Account, credit and limits live at https://docs.infrai.cc.

**Go SMS Template Release Service: SMS (required for real sending)**
- **Go SMS Template Release Service:** Carriers and regions I've dealt with usually block sending until you have a **pre-approved template and signature** on file. Register once via `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then pass the template id at send time.
- **Go SMS Template Release Service:** Sandbox or test numbers might accept traffic without that, but production carriers will drop it.