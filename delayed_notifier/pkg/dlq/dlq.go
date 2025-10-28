package dlq

// DLQ is a generic type that stores objects and errors
//
//	DLQ := dlq.NewDLQ[int](len(nums) / 10)
//	go func() {
//		for _, n := range nums {
//			err = operation(n)
//			if err != nil {
//				DLQ.Put(n, fmt.Errorf("bad result: %w", err)
//			}
//		}
//		DLQ.Close()
//	}
//
//	// somewhere
//
//	for failedMessage := range DLQ.Items() {
//		zlog.Logger.Error(failedMessage.Error()).Msg(fmt.Sprintf("failed to send int: %d", failedMessage.Value()))
//	}
type DLQ[T any] struct {
	// items field is private so no one can put without knowing
	items chan *Item[T]
}

// Item is basically an item of DLQ
//
//	item.Value()
//	item.Error()
type Item[T any] struct {
	value T
	error error
}

// Value is getter for value field
func (i *Item[T]) Value() T {
	return i.value
}

// Error is getter for error field
func (i *Item[T]) Error() error {
	return i.error
}

// NewDLQ creates an empty DLQ
//
//	bufferSize - channel buffer, 0 => method Put blocks everything
func NewDLQ[T any](bufferSize int) *DLQ[T] {
	return &DLQ[T]{
		items: make(chan *Item[T], bufferSize),
	}
}

// Put inserts a new item into DLQ
//
// warning: blocking operation because of chan inside!
func (d *DLQ[T]) Put(value T, err error) {
	d.items <- &Item[T]{value: value, error: err}
}

// Items returns a read-only channel with inserted values
func (d *DLQ[T]) Items() <-chan *Item[T] {
	return d.items
}

// Close closes DLQ channel
func (d *DLQ[T]) Close() {
	close(d.items)
}
