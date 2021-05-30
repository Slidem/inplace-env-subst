# In place environment variable substitution

This is a simple go module, to replace environment variables placeholders in files in place.

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

## Default values

You can use default values by following the following format:

```
${ENV:-defaultValue}
```

`:-` is considered the separator between the env key, and the default value.
The default value is considered only when the env variable is not found.

## Whitelist variables

You can whitelist environment variables with the WhitelistEnv value.
This will ensure that all other env variables will be ignored.

```go
inplaceenvsubst.ProcessFiles([]string{"/documents/test.txt"}, &inplaceenvsubst.Config{
    FailOnMissingVariables: false,
    RunInParallel:          false,
    ErrorListener:          nil,
    WhitelistEnvs:          inplaceenvsubst.NewStringSet("whitelisted")
})
```

## Blacklist Variables

You can blacklist environment variables with the BlacklistEnv value.
This will ensure that all other env variables will be ignored.

```go
inplaceenvsubst.ProcessFiles([]string{"/documents/test.txt"}, &inplaceenvsubst.Config{
    FailOnMissingVariables: false,
    RunInParallel:          false,
    ErrorListener:          nil,
    BlacklistEnvs:          inplaceenvsubst.NewStringSet("blacklisted")
})
```

## Notes

- Only accepts environment variables placeholders format `${ENV_KEY}`
- `$ENV_KEY` will not work
- Nested env variables not supported (Ex : `${ENV_KEY ${ENV_NESTED}}`, in this case, the env variable will be `ENV_KEY ${ENV_NESTED}`)
- You cannot use both blacklists and whitelist envs