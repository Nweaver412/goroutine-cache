package main

import (
	"fmt"
	"kvstore/store"
	"time"
)

func main() {
	kv := store.NewTTLStore()
	kv.Set("key", "val")
	value, exists := kv.Get("key")
	fmt.Println("Get key:", value, "exists:", exists)

	kv.SetWithTTL("temp", "expire soon", 2*time.Second)
	value, exists = kv.Get("temp")
	fmt.Println("Get temp:", value, "exists:", exists)

	time.Sleep(3 * time.Second)
	value, exists = kv.Get("temp")
	fmt.Println("Get temp after TTL:", value, "exists:", exists)
}
