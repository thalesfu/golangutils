package logging

import "context"

const LogContextKey = "log"
const LogStoreKey = "log-store"

func InitializeContextLogStore(ctx context.Context, name string) (context.Context, *LogStore) {
	if ctx == nil {
		return nil, nil
	}

	store, ok := ctx.Value(LogStoreKey).(*LogStore)

	if !ok {
		store = NewTopLogStore(name)
		ctx = context.WithValue(ctx, LogStoreKey, store)
		return ctx, store
	}

	child := store.CreateChild(name)

	ctx = context.WithValue(ctx, LogStoreKey, child)

	return ctx, child
}

func GetContextLogStore(ctx context.Context) (*LogStore, bool) {
	if ctx == nil {
		return nil, false
	}

	store, ok := ctx.Value(LogStoreKey).(*LogStore)

	if !ok {
		return nil, false
	}

	return store, true
}

type LogStore struct {
	name   string
	data   map[string]string
	parent *LogStore
}

func NewLogStore(name string, parent *LogStore) *LogStore {
	return &LogStore{
		name:   name,
		data:   make(map[string]string),
		parent: parent,
	}
}

func NewTopLogStore(name string) *LogStore {
	return &LogStore{
		name: name,
		data: make(map[string]string),
	}
}

func (l *LogStore) Set(key, value string) {
	l.data[key] = value
}

func (l *LogStore) Get(key string) string {
	return l.data[key]
}

func (l *LogStore) Delete(key string) {
	delete(l.data, key)
}

func (l *LogStore) GetAll() map[string]string {
	result := make(map[string]string)

	if l.parent != nil {
		for k, v := range l.parent.GetAll() {
			result[k] = v
		}
	}

	for k, v := range l.data {
		result[k] = v
	}

	return result
}

func (l *LogStore) GetParent() *LogStore {
	return l.parent
}

func (l *LogStore) GetRoot() *LogStore {
	if l.parent == nil {
		return l
	}
	return l.parent.GetRoot()
}

func (l *LogStore) IsTop() bool {
	return l.parent == nil
}

func (l *LogStore) GetPath() string {
	if l.parent == nil {
		return l.name
	}
	return l.parent.GetPath() + "." + l.name
}

func (l *LogStore) CreateChild(name string) *LogStore {
	return NewLogStore(name, l)
}
