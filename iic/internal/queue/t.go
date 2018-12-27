package queue

// T represets a queue (FIFO) of ints.
type T struct {
	D          []int
	start, end int
}

// New creates a new queue with a capacity hint
// which can be used to reduce re-allocations and
// copying.
func New(capHint int) *T {
	return &T{D: make([]int, capHint)}
}

// Push pushes v onto the end of the queue.
func (q *T) Push(v int) {
	ne := q.end + 1
	if ne == len(q.D) {
		ne = 0
	}
	if ne == q.start {
		q.grow()
		q.Push(v)
		return
	}
	q.D[q.end] = v
	q.end = ne
}

// Pop pops the first pushed element which has not yet been
// popped.  Pop panics if the queue is empty.
func (q *T) Pop() int {
	if q.start == q.end {
		panic("oob")
	}
	r := q.D[q.start]
	q.start++
	if q.start == len(q.D) {
		q.start = 0
	}
	return r
}

// Len returns the length of the queue.
func (q *T) Len() int {
	if q.end > q.start {
		return q.end - q.start
	}
	if q.end < q.start {
		return (len(q.D) - q.start) + q.end
	}
	return 0
}

func (q *T) grow() {
	N := len(q.D)
	M := 2 * N
	if M == 0 {
		M = 13
	}
	tmp := make([]int, M)
	if q.start < q.end {
		copy(tmp, q.D[q.start:q.end])
		q.D = tmp
		q.end = q.end - q.start
		q.start = 0
		return
	}
	copy(tmp, q.D[q.start:])
	copy(tmp[N-q.start:], q.D[:q.start])
	q.D = tmp
	q.end = (N - q.start) + q.end
	q.start = 0
}
