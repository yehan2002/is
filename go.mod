module github.com/yehan2002/is/v2

go 1.16

retract v2.2.0 // this version deadlocks on parallel tests

require github.com/google/go-cmp v0.5.9
