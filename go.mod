module github.com/yehan2002/is/v2

go 1.16

retract v2.2.0 // this version deadlocks on parallel tests

require github.com/go-test/deep v1.1.0
