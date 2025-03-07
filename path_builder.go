package main

type PathBuilder interface {
	ConfigLocations() []string
	LogLocations() string
}
