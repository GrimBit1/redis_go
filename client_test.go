package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestClient(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := range 10 {
		go func(it int) {
			client, err := NewClient("localhost:6379")
			if err != nil {
				t.Error(err)
			}

			err = client.Set(context.Background(), fmt.Sprintf("key_%d", it), fmt.Sprintf("value_%d", it))
			if err != nil {
				t.Error(err)
			}

			val, err := client.Get(context.Background(), fmt.Sprintf("key_%d", it))
			if err != nil {
				t.Error(err)
			}

			fmt.Println("val", val)
			if val != fmt.Sprintf("value_%d", it) {
				t.Error("val not match")
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

}
