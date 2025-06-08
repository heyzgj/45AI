package service

import (
	"context"
	"sync"
)

type Job struct {
	UserID     int64
	TemplateID int
	ImageData  []byte
}

type QueueService interface {
	AddJob(ctx context.Context, job *Job) error
	GetJob(ctx context.Context) (*Job, error)
}

type inMemoryQueueService struct {
	queue []*Job
	mutex sync.Mutex
}

func NewInMemoryQueueService() QueueService {
	return &inMemoryQueueService{
		queue: make([]*Job, 0),
	}
}

func (s *inMemoryQueueService) AddJob(ctx context.Context, job *Job) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.queue = append(s.queue, job)
	return nil
}

func (s *inMemoryQueueService) GetJob(ctx context.Context) (*Job, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.queue) == 0 {
		return nil, nil
	}
	job := s.queue[0]
	s.queue = s.queue[1:]
	return job, nil
} 