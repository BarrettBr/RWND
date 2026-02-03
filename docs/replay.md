# Replay

Replay mode lets you step through recorded requests and optionally resend
the current request to compare old and new responses.

## Step Flow

```text
Press Enter for next, r to replay, q to quit >
```

When you press `Enter`, the request is printed in a readable format.

## Replay Flow

When you press `r`, the current request is re-sent to its recorded URL.
The old response is printed then the new response.

## Output Shape

Example:

```text
Request #1
GET http://localhost:3000/test
Headers:
  Accept: */*

Old Response
Status: 200
Body:
  hello from upstream: /test
---
New Response
Status: 200
Body:
  hello from upstream: /test
```
