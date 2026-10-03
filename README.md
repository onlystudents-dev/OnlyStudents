# OnlyStudents
This is a student diary app, mimicking e-kreta.hu, but modern, fast and just better.

# WARNING
This is pre-alpha software, untested, and does have a ton of breaking API and other changes constantly. There is no versioning system yet either.
For testing/development, `DEMO_MODE=true` can be used right now which registers demo accounts to OPAQUE. See `Demo Mode Credentials`

# Installation
- Copy `.env.example` to `.env`
- Run `go run ./cmd/opaque-setup` to get OPAQUE secrets and put them inside `.env`
- Populate the other variables in `.env`
- Run `docker compose up -d` or `podman compose up -d` depending on your configuration. Rootless Docker/Podman is supported with no additional changes needed.

# Features (some planned, some implemented)
- Latest protocols, encryption and standards
- Post-Quantum E2EE Messaging (with [FS](https://en.wikipedia.org/wiki/Forward_secrecy) and [PCS](https://crypto.stackexchange.com/a/84455)) with reporting/moderation possible through client side encryption key sharing (actually trustable!)
- OPAQUE + Argon2ID Login/Registration which means the server never sees the password but can prove you are the same person
- Admin panel for governments, who can add new schools, new accounts, etc
- Student/Guardian/Teacher Timetable, Grades, Homework, Absences, Exams, etc
- Multiple schools for a student, and multiple children for a guardian instead of many accounts
- Insanely fast Fiber Go server which serves both the backend and a static Preact (instead of react) frontend with lazy loading.
- Custom subjects, grades, permissions
- Moderation panel for moderators
- Custom Themes

# Demo Mode Credentials ( https://onlystudents.hu )
- Student/Teacher/Guardian ID 1, Password: test
- Email Change/Password reset does not work to combat spam!
