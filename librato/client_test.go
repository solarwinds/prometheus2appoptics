package librato

var paramErrorFixture = `{
    "measurements": {
        "summary": {
            "total": 6,
            "accepted": 0,
            "failed": 6
        }
    },
    "errors": [
        {
            "param": "time",
            "value": 1507056682391,
            "reason": "Is too far in the future (>30 minutes ahead). Check for local clock drift or enable NTP."
        }
    ]
}`

var requestErrorFixture = `{
  "errors": {
    "request": [
      "Please use secured connection through https!",
      "Please provide credentials for authentication."
    ]
  }
}`

// TODO: this probably represents a bug, as I got this from the same request that produced "paramErrorFixture" above
// but simply altered the URL scheme to `http` instead of `https`
var genericErrorSliceFixture = `{
    "errors": [
        "must specify metric name or compose parameter"
    ]
}`

var textErrorFixture = `Credentials are required to access this resource.`
