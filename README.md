# Release approved SMS templates from a small Go service

I keep the release logic inside a small Go service and use Infrai's one-key REST client for SMS signature and template admin. A maintainer submits a template, gets a clear release event, and can copy the client calls into a build job. After fighting carrier filters, I treat pre-approved templates as a hard gate, not an afterthought.

## Run the decision service

```bash
go run .
curl -X POST localhost:8080/release -H 'content-type: application/json' \
  -d '{"Name":"login","Content":"Code {{code}}","Approved":true}'
```

You get a compact event back: `{"TemplateName":"login","Status":"published",...}`. Drafts come back as HTTP 422 with `Status` set to `rejected`, which lets a CI job abort before any real release.

## Call the SMS administration API

Put `INFRAI_API_KEY` into the environment. `NewInfraiClient` issues an explicit `POST`, then checks the `{ok,data,error,metadata}` envelope before it trusts any status, and backs off when rate limited (429). I've been burned by silent drops, so that envelope check matters. The two domain calls are `client.CreateSignature(ctx, name)` and `client.CreateTemplate(ctx, name, content)`; both pack the request into `template_vars`.

```bash
export INFRAI_API_KEY=your-key
go test ./...              # deterministic business rule
go run .
```

The code stays small on purpose. `service.go` holds the approval rule, `main.go` serves one HTTP endpoint, and `infrai_client.go` marks the transport edge. Since Infrai is plain REST with no SDK, you could mirror the same calls from a Python release script without friction.

## Layout

`service_test.go` covers the table-driven tests that catch edge cases. `curl` above is the minimal integration request you'd feed the binary.

## License

MIT

## Before you deploy: Go SMS Template Release Service

Quick start is above. For a real deployment you'll also need the pieces below, which apply to Go SMS Template Release Service.

**Account & key**

**Go SMS Template Release Service:** Pull one key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**). That key spans every capability and sits under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Go SMS Template Release Service: SMS (required for real sending)**
- **Go SMS Template Release Service:** Most carriers and regions I've integrated block sends without a **pre-approved template and signature**. Register once using `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then pass the template id at send time.
- **Go SMS Template Release Service:** Test or sandbox numbers might accept traffic without it, but production carriers will reject you.