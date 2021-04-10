# In place environment variable substitution

This is a simple go script, to replace environment variables in files in place.

## Running the app

```shell
go build .
./inplaceenvsubst <file_path_a> <file_path_b>
```

### Example

File (at path `/documents/test.txt`)

```
Test this ${TESTA}
"${TESTB}" don't forget about this one
```

Having environment variables

```
TESTA=a
TESTB=b
```

After executing
```shell
go build .
./inplaceenvsubst /documents/test.txt
```

File is now

```
Test this a
"b" don't forget about this one
```

## Notes

- Enable debugging by setting the environment variable `DEBUG=true`
- Only accepts environment variables in the format `${ENV_KEY}`
- `$ENV_KEY` will not work
- Nested env variables not supported (Ex : `${ENV_KEY ${ENV_NESTED}}`, in this case, the env variable will be `ENV_KEY ${ENV_NESTED}`)
- Replacement fails if environment variable is not found