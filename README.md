# OnlyStudents
This is a student diary app, mimicking ekreta.hu, but modern, fast and just better.

# Installation
- Copy `.env.example` to `.env`
- Populate the variables in `.env`
- Run `docker compose up -d` or `podman compose up -d` depending on your configuration. Rootless Docker/Podman is supported with no additional changes needed.

# Features (some planned, some implemented)
- Latest protocols, encryption and standards
- Post-Quantum E2EE Messaging (with [FS](https://en.wikipedia.org/wiki/Forward_secrecy) and [PCS](https://crypto.stackexchange.com/a/84455)) with reporting/moderation possible through client side encryption key sharing (even trustable by students!)
- OPAQUE + Argon2ID Login/Registration which means the server never sees the password
- Admin panel for governments, who can add new schools, new accounts, etc
- Moderation panel for moderators
- Student/Guardian/Teacher Timetable, Grades, Homework, Absences, Exams, etc
- Custom subjects
- Granular permission system
- Guardian children selector instead of multiple children accounts
- Custom Themes
- Insanely fast Fiber Go server, Node.JS frontend is served as static files
