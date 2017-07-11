# errors

Build upon [pkg/errors](https://github.com/pkg/errors), with some customization for my customized [yyscamper/logrus](https?//github.com/yyscamper/logrus)

## Usage
```go
import "github.com/yyscamper/errors"

var err error

//errors.New style
err = errors.New("some error")
if err != nil {
    ...
}

//fmt.Errorf style
err = errors.Errorf("some error, path=%s", "/home/me/foo.txt")
if err != nil {
    ...
}

//Wrap the existing error
path := "/home/me/foo.txt"
file, err = os.Open(path)
if err != nil {
    return errors.Wrap(err) //The callstack will be attached
}
```

## Compatiblity
My error wrapper extends the golang error interface, so it can be casted into the golang error

## Stacktrace
Stacktrace is automatically collected when the error is created

## Fields
Designed for logrus, logrus can extract additional fileds from the error objects
```go
errors.New("some error").WithField("path": "/home/me/foo.txt")
errors.New("some error").WithFields(errors.Fileds{
    "path": "/home/me/foo.txt",
    "mode": "readwrite",
    "time": time.Now(),
})
errors.New("some error").With("path": "/home/me/foo.txt", "time", time.Now())
```

## Generator
Using the generator can automatically set a name for the error, it is useful for logrus module logging

```go
errgen := errors.NewGenerator("testmodule")
err := errgen.New("some error")
err = errgen.Errorf("some error with path=%s", "/home/me/foo.txt")
```

## Logrus Integration
I have a customized [yyscamper/logrus](http://github.com/yyscamper/logrus), it extends the support for this error wrapper, so that the error stacktrace and all fields information can be automatically logged by logrus
```go
import (
    log "github.com/yyscamper/logrus"
    "github.com/yyscamper/errors"
)

err := errors.New("some error").WithField("path", "/home/me/foo.txt")
log.WithError(err).Error("some happens")
```
