# `auth-core`

> Identity, credentials, and permissions — one service, done right.

---

`auth-core` is a Visea's own self-hosted Single Sign-On service built to be the single source of truth for authentication and authorization across an ecosystem of applications. It handles who you are, how you prove it, and what you're allowed to do — then gets out of the way.

---

## What it does

### Identity Management
Users, organizations, and memberships are first-class concepts. A user exists independently of any organization and can belong to many. Organizations can enforce SSO and MFA at the policy level, not just the application level.

### Authentication
Two login paths, both secure:

- **Email & Password** — Credentials verified against Argon2id hashes. Rate-limited at the IP level via Redis. MFA-gated when enabled.
- **Passwordless OTP** — A time-limited, single-use code delivered over email or SMS, stored only as a hash. Resistant to user enumeration by design.

All authentication events — successes and failures — are written to an immutable audit log.

### Multi-Factor Authentication
TOTP-based MFA (compatible with any authenticator app). The TOTP secret is AES-256-GCM encrypted before storage — the encryption key lives in the environment, never the database. Enrollment requires verification before activation. Ten single-use backup codes are generated at enrollment and shown exactly once.

### Session Management
Every login creates a `session` tied to a device, IP, and user agent. Sessions can be revoked individually or all at once. Refresh tokens use a family-based rotation scheme — if a token that has already been used is ever presented again, the entire family is immediately revoked and a security alert is sent to the user.

### Token Infrastructure
Access tokens are short-lived JWTs (15 minutes), signed with ES256 using a private key loaded from the environment at startup. The private key never touches the database. Public keys are served via a standard `/.well-known/jwks.json` endpoint so downstream services can verify tokens locally — no round-trip to `auth-core` on every request.

### Permissions & RBAC
A full role-based access control system scoped to organizations. Permissions are namespaced by service using a `service:resource:action` key format (e.g. `payments:invoices:write`). Roles are compositions of permissions, assigned to users within an org. Resource servers check claims embedded in the JWT — no permission lookup at runtime.

### Security Operations
Every sensitive action in the system produces an entry in `audit_logs` — an append-only record of what happened, who did it, from where, and to what. Login attempts are tracked separately to power rate limiting and account lockout logic without polluting the audit trail.

---

*Designed to be the last auth service we build.*
