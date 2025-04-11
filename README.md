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
<<<<<<< HEAD
- Add Message Password
- Add ExpirationDate
=======
- Unit-tests
- `TaskFile` build
- Continuous Integration
- Review code, code is so bad, and this code need to improve
>>>>>>> parent of 1c18ca3 (Merge branch 'main' of https://github.com/kvm101/OneTimeSecret)
