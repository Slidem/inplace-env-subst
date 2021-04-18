# In place environment variable substitution

This is a simple go module, to replace environment variables placeholders in files in place.

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

After calling 

```go
inplaceenvsubst.ProcessFiles([]string{"/documents/test.txt"}, &inplaceenvsubst.Config{
    FailOnMissingVariables: false,
    RunInParallel:          false,
    ErrorListener:          nil,
})
```

`test.txt` will have the following content

```
Test this a
"b" don't forget about this one
```

## Notes

- Only accepts environment variables placeholders format `${ENV_KEY}`
- `$ENV_KEY` will not work
- Nested env variables not supported (Ex : `${ENV_KEY ${ENV_NESTED}}`, in this case, the env variable will be `ENV_KEY ${ENV_NESTED}`)