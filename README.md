# OneTimeSecret.

### Description:
Server for One-Time Secret messages

### Design pattern (MVC):
- ### Model
    - `User`
    - `Message`
    - `AccountData`
    - `MessageInfo`
- ### View (Templates)
    - `Account`
    - `Messages`
    - `Not found page`
- ### Controller
    - Control beetwen Model, View and DB(Postgres with ORM (GORM) )

### Security:
- `HTTP Basic Authentication`

### Needed to add:
- Unit-tests
- `TaskFile` build
- Continuous Integration
- Review code, code is so bad, and this code need to improve
