package queue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {

	t.Parallel()

	var tests = []struct {
		len int
	}{
		{0},
		{1000000},
		{401},
		{12},
		{1212},
		{23},
		{1},
	}
	for _, test := range tests {
		q := NewQueue[string]()

		require.Equal(t, q.Count, 0, "err: q.Count = %d, but must be 0", q.Count)

		for i := 0; i < test.len; i++ {
			q.Add(fmt.Sprintf("str number %d", i))
		}

		assert.Equal(t, q.Count, test.len, fmt.Sprintf("err: q.Count = %d, but must be %d", q.Count, test.len))
	}

}

func TestPop(t *testing.T) {

	t.Parallel()

	q := NewQueue[string]()

	str := make([]string, 0)
	var newStr string
	for i := 0; i < 10; i++ {
		newStr = fmt.Sprintf("str number %d", i)
		q.Add(newStr)
		str = append(str, newStr)
	}

	for i := 0; i < 11; i++ {
		val, err := q.Pop()

		if err != nil {
			assert.Equal(t, i, 10, fmt.Sprintf("err: %s\nnum: %d", err, i))
		}

		if i != 10 && val != str[i] {
			t.Errorf("err: %s != %s", val, str[i])
		}
	}

	assert.Equal(t, q.Count, 0, fmt.Sprintf("err: q.Count = %d, but must be 0", q.Count))
}

func BenchmarkAdd(b *testing.B) {
	q := NewQueue[int]()
	for i := 0; i < b.N; i++ {
		q.Add(i)
	}

	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}
