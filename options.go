package eagle

import "time"

type Options struct {
	capacity int32
	expiryDuration  time.Duration
}

type option func(*Options) 

func WithCapacity(cap int32) option {
	return func (o *Options)  {
		o.capacity = cap
	}
}

func WithExpiry(t time.Duration) option {
	return func (o *Options)  {
		o.expiryDuration = t
	}
}
