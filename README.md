# resequencer

![CI](https://github.com/teivah/resequencer/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/teivah/resequencer)](https://goreportcard.com/report/github.com/teivah/resequencer)

Resequencing in Go

## Introduction

`resequencer` is a Go library that implements the [resequencer pattern](https://www.enterpriseintegrationpatterns.com/Resequencer.html).

One use case, for example, is when using Sarama with a consumer group, and we distribute each message to a set of workers. If we want to make sure the offsets are committed in sequence, we can use a resequencer per partition.

## How to Use

First, we need to create a new `Handler`. Then we have to use the two methods:
* `Push` to add new sequence IDs
* `Messages` that returns a `<-chan []int` that contains the ordered sequence IDs

```go
ctx, cancel := context.WithCancel(context.Background())
handler := resequencer.NewHandler(ctx, -1) // Create a resequencer and initialize the first sequence ID to -1

max := 10
for i := 0; i < max; i++ {
	i := i
	go func() {
		handler.Push(i) // Push a new sequence ID
	}()
}

for sequenceIDs := range handler.Messages() { // Read all the sequence IDs (sequenceIDs is an []int).
	for _, sequenceID := range sequenceIDs {
		fmt.Println(sequenceID)
		if sequenceID == max-1 {
			cancel()
			return
		}
	}
}
```
