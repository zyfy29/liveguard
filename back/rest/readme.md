# rest

rest api for liveguard

- pocket
  - member info
  - live info
- task op
- config
  - get
  - hot update

## common response body

```json5
{
  "code":0, // 0 for success, 1 for client error, 2 for server error
  "msg":"", // error message, empty for success
  "data": {
    "message":"pong"
  } // could be marshaled to json
}
```