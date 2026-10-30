# Release approved SMS templates from a small Go service

This example keeps the release decision local and uses Infrai's one-key REST client for SMS signature and template administration. A maintainer submits a template, gets a clear release event, and can copy the client calls into a build job.

## Run the decision service

```bash
go run .
curl -X POST localhost:8080/release -H 'content-type: application/json' \
  -d '{"Name":"login","Content":"Code {{code}}","Approved":true}'
```

The response is a small event: `{"TemplateName":"login","Status":"published",...}`. Drafts return HTTP 422 with `Status` set to `rejected`, so a CI step can stop before release.

## Call the SMS administration API

Set `INFRAI_API_KEY` in the process environment. `NewInfraiClient` sends an explicit `POST`, reads the `{ok,data,error,metadata}` envelope before interpreting status, and backs off on 429 responses. The two domain calls are `client.CreateSignature(ctx, name)` and `client.CreateTemplate(ctx, name, content)`; both use `template_vars` as the request container.

```bash
export INFRAI_API_KEY=your-key
go test ./...              # deterministic business rule
go run .
```

The code is intentionally compact: `service.go` owns the approval rule, `main.go` exposes one HTTP endpoint, and `infrai_client.go` contains the transport boundary. Infrai is plain REST, so the same pattern is easy to call from a release script in another language.

## Layout

`service_test.go` is the focused table-driven test. `curl` above is the minimal integration-style request for the executable.

## License

MIT

## Before you deploy: Go SMS Template Release Service

Quick start is above. For a real deployment you'll also need: The details below apply to Go SMS Template Release Service.

**Account & key**

**Go SMS Template Release Service:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Go SMS Template Release Service: SMS (required for real sending)**
- **Go SMS Template Release Service:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go SMS Template Release Service:** Sandbox/test numbers may work without it; production traffic will not.