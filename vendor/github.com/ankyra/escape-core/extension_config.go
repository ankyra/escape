package core

/*

A package can extend another package and inherit its build scripts and input
and output variables. This makes it possible to reuse build and deployment
patterns.

## Which fields get extended?

## How do scripts get executed?

## When to use dependencies, extensions, providers?

## Shortcomings

## Escape Plan

Extensions are configured in the [`extends`](/docs/escape-plan/#extends)
field of the Escape Plan.
*/
type ExtensionConfig struct {
	ReleaseId string `json:"release_id"`
}

func NewExtensionConfig(releaseId string) *ExtensionConfig {
	return &ExtensionConfig{releaseId}
}
