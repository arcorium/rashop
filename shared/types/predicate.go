package types

type UnaryPredicate[T any] func(T)

type UnaryPredicateWithReturn[T, R any] func(T) R

type UnaryPredicateError[T any] UnaryPredicateWithReturn[T, error]

type BinaryPredicate[T, U any] func(T, U)

type BinaryPredicateWithReturn[T, U, R any] func(T, U) R

type BinaryPredicateError[T, U any] BinaryPredicateWithReturn[T, U, error]
