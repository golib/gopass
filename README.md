# gopass

Manage credentials among `Chrome Passwords`, `Lastpass` and `Keychain`.

### Install

```bash
$ go get github.com/golib/gopass
```

### Migrations

#### Migrate credentials from Lastpass to Chrome

- exports credentials from lastpass

    click *Lastpass Icon* -> *More Options* -> *Advanced* -> *Export* -> *LastPass CSV File*

- converts csv file exported from lastpass to chrome import format

    ```bash
    $ gopass chrome --csv lastpass.csv
    ```

- imports csv file

    > NOTE: You have to enable `Password import` feature before you can import credentials from csv file. You can enable it by visiting [chrome://flags](chrome://flags)

    visit [chrome://settings](chrome://settings) -> *Passwords* -> *(More actions)*

#### Migrate credentials from csv file to macos Keychain

- exports credentials from your provider

    > Such as LastPass or Chrome

- write credentials to keychain

    ```bash
    $ gopass keychain --csv credentials.csv
    ```
