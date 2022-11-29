# Test dependencies import for zap implementation of the logger

The purpose of this test is to evaluate the dependency list when using only parts of the implementation.

For example, in this test, only zap should be included in the dependency list, even if logrus is listed as dependency of [go-belt](https://github.com/facebookincubator/go-belt).

## Listing the dependencies

To list the dependencies for this module:
```
$ go list -f '{{ join .Deps "\n" }}'
# Checking if zap is included (it should)
$ go list -f '{{ join .Deps "\n" }}' | grep zap
github.com/facebookincubator/go-belt/tool/logger/implementation/zap
go.uber.org/zap
go.uber.org/zap/buffer
go.uber.org/zap/internal
go.uber.org/zap/internal/bufferpool
go.uber.org/zap/internal/color
go.uber.org/zap/internal/exit
go.uber.org/zap/zapcore
# Checking if logrus is included (it should not)
$ go list -f '{{ join .Deps "\n" }}' | grep logrus
```

Also, the [go.sum](test/dependencies/zap/go.sum) does not contain any entry for logrus:
```
$ cat go.sum | grep logrus
```

## Pending

* Evaluate the possibility of adding a CI test to check the number of dependencies imported by each logger implementation
