package server_test

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/stretchr/testify/require"
)

type SafeMap struct {
	rw sync.Mutex
	m  map[string]bool
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		rw: sync.Mutex{},
		m:  map[string]bool{}}
}

func (m *SafeMap) Add(key string) {
	m.rw.Lock()
	defer m.rw.Unlock()
	if _, ok := m.m[key]; !ok {
		m.m[key] = true
	}
}

func TestFreeTCPAddr(t *testing.T) {
	var m sync.Map
	nports := 1000
	t.Parallel()

	var wg sync.WaitGroup
	loadedCh := make(chan bool, 0)
	done := make(chan bool, 1)
	var conflicts int

	go func(loaded <-chan bool, conflicts *int, done chan bool) {
		for p := range loaded {
			if p {
				*conflicts++
			}
		}
		done <- true
		close(done)
	}(loadedCh, &conflicts, done)

	for i := 0; i < nports; i++ {
		wg.Add(1)
		go func(m *sync.Map, loaded chan<- bool) {
			_, port, err := server.FreeTCPAddr()
			actual, load := m.LoadOrStore(port, true)
			loaded <- load
			fmt.Fprintf(os.Stderr, "port: %s actual: %v loaded: %v\n", port, actual, load)
			require.NoError(t, err)
			wg.Done()
		}(&m, loadedCh)
	}
	wg.Wait()
	close(loadedCh)
	<-done

	require.Equal(t, 0, conflicts)
	//	require.Equal(t, nports, len(m))
}
