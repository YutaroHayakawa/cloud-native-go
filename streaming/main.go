package main

import (
  "fmt"
  "time"
  "context"
)

type Value int

func SlowOperation(ctx context.Context) (Value, error) {
  select {
  case <-ctx.Done():
    return 0, ctx.Err()
  case <-time.After(time.Second * 10):
    return 0, nil
  }
}

func Stream(ctx context.Context, out chan Value) error {
  dctx, cancel := context.WithTimeout(ctx, time.Second * 5)

  defer cancel()

  res, err := SlowOperation(dctx)
  if err != nil {
    return err
  }

  for {
    select {
    case out <- res:
    case <-ctx.Done():
      return ctx.Err()
    }
  }
}

func main() {
  out := make(chan Value)
  ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)

  defer cancel()

  go Stream(ctx, out)

  select {
  case val := <-out:
    fmt.Println(val)
  case <-ctx.Done():
    fmt.Println(ctx.Err())
  }
}
