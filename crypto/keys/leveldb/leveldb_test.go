package leveldb_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/stretchr/testify/require"
)

func benchmarkTestDir(b *testing.B) (string, func()) {
	dir, err := ioutil.TempDir("", b.Name()+"_")
	require.NoError(b, err)
	return dir, func() { os.RemoveAll(dir) }
}

func BenchmarkCreateAccount(b *testing.B) {
	dir, cleanup := benchmarkTestDir(b)
	defer cleanup()
	kb := keys.New(b.Name(), dir)
	for i := 0; i < b.N; i++ {
		_, _, err := kb.CreateMnemonic(fmt.Sprintf("%s_%d", b.Name(), i), keys.English, "012345678", keys.Secp256k1)
		if err != nil {
			panic(err)
		}
	}
}
