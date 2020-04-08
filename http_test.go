package recaptcha_test

import (
	"github.com/lithdew/recaptcha"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var Request = recaptcha.Request{Secret: "shhhh", Response: "hhhhhh"}

func TestDo(t *testing.T) {
	res, err := recaptcha.Do(Request)
	require.NoError(t, err)
	require.False(t, res.Success)
	require.Len(t, res.ErrorCodes, 2)

	var client recaptcha.Client

	res, err = client.Do(Request)
	require.NoError(t, err)
	require.False(t, res.Success)
	require.Len(t, res.ErrorCodes, 2)
}

func TestDoTimeout(t *testing.T) {
	res, err := recaptcha.DoTimeout(Request, 3*time.Second)
	require.NoError(t, err)
	require.False(t, res.Success)
	require.Len(t, res.ErrorCodes, 2)

	var client recaptcha.Client

	res, err = client.DoTimeout(Request, 3*time.Second)
	require.NoError(t, err)
	require.False(t, res.Success)
	require.Len(t, res.ErrorCodes, 2)
}

func BenchmarkDo(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := recaptcha.Do(Request)
		require.NoError(b, err)
	}
}

func BenchmarkDoTimeout(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := recaptcha.DoTimeout(Request, 3*time.Second)
		require.NoError(b, err)
	}
}

func BenchmarkParallelDo(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := recaptcha.Do(Request)
			require.NoError(b, err)
		}
	})
}

func BenchmarkParallelDoTimeout(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := recaptcha.DoTimeout(Request, 3*time.Second)
			require.NoError(b, err)
		}
	})
}
