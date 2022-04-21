package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func mvvmCheckWatch(t *testing.T, x *mvvm, model, view string) {
	t.Helper()

	assert.Contains(t, x.m2v, model)
	assert.True(t, x.m2v[model][view])
	assert.Contains(t, x.v2m, view)
	assert.True(t, x.v2m[view][model])

	key := x.rel(model, view)
	assert.Contains(t, x.funcs, key)
}

// -----------------------------------------------------------------------------

func mvvmCheckUnwatch(t *testing.T, x *mvvm, model, view string) {
	t.Helper()

	assert.False(t, x.m2v[model][view])
	assert.False(t, x.v2m[view][model])

	key := x.rel(model, view)
	assert.NotContains(t, x.funcs, key)
}

// -----------------------------------------------------------------------------

func TestWatch(t *testing.T) {
	x := newMVVM()

	model := "session"
	view := "view"
	sum := 0
	x.Watch(model, view, func(a, b int) {
		sum = a + b
	})

	mvvmCheckWatch(t, x, model, view)
	x.trigger(model, view, 1, 2)
	assert.Equal(t, sum, 1+2)
}

// -----------------------------------------------------------------------------

func TestUnwatch(t *testing.T) {

	x := newMVVM()

	model1 := "session1"
	view1 := "view1"

	x.Watch(model1, view1, func(a, b int) {
		_ = a + b
	})
	mvvmCheckWatch(t, x, model1, view1)

	model2 := "session2"
	view2 := "view2"

	x.Watch(model2, view2, func(a, b int) {
		_ = a + b
	})

	mvvmCheckWatch(t, x, model2, view2)

	x.Unwatch(model2, view2)

	mvvmCheckUnwatch(t, x, model2, view2)
}
