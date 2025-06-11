package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/45ai/backend/internal/model"
	"github.com/45ai/backend/internal/repository"
)

// Job represents a generation job in the queue
type Job struct {
	JobID      string
	UserID     int64
	TemplateID int
	ImageData  []byte
}

// QueueService interface for managing generation queue
type QueueService interface {
	AddJob(ctx context.Context, userID int64, templateID int, imageData []byte) (string, error)
	StartWorkers(ctx context.Context, workerCount int, generationService GenerationService)
	Shutdown(ctx context.Context) error
}

// channelQueueService implements QueueService using Go channels
type channelQueueService struct {
	jobChannel    chan *Job
	generationRepo repository.GenerationRepository
	workerDone    chan bool
	shutdown      chan bool
}

// NewChannelQueueService creates a new channel-based queue service
func NewChannelQueueService(generationRepo repository.GenerationRepository) QueueService {
	return &channelQueueService{
		jobChannel:     make(chan *Job, 100), // Buffer for 100 jobs
		generationRepo: generationRepo,
		workerDone:     make(chan bool),
		shutdown:       make(chan bool),
	}
}

func (s *channelQueueService) AddJob(ctx context.Context, userID int64, templateID int, imageData []byte) (string, error) {
	jobID, err := generateJobID()
	if err != nil {
		return "", fmt.Errorf("failed to generate job ID: %w", err)
	}

	// Create database record
	generation := &model.Generation{
		JobID:      jobID,
		UserID:     userID,
		TemplateID: templateID,
		Status:     "pending",
		Progress:   0,
	}

	if err := s.generationRepo.Create(ctx, generation); err != nil {
		return "", fmt.Errorf("failed to create generation record: %w", err)
	}

	// Create job and add to channel
	job := &Job{
		JobID:      jobID,
		UserID:     userID,
		TemplateID: templateID,
		ImageData:  imageData,
	}

	select {
	case s.jobChannel <- job:
		log.Printf("Job %s queued successfully", jobID)
		return jobID, nil
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return "", fmt.Errorf("queue is full")
	}
}

func (s *channelQueueService) StartWorkers(ctx context.Context, workerCount int, generationService GenerationService) {
	log.Printf("Starting %d queue workers", workerCount)

	for i := 0; i < workerCount; i++ {
		go s.worker(ctx, i+1, generationService)
	}
}

func (s *channelQueueService) worker(ctx context.Context, workerID int, generationService GenerationService) {
	log.Printf("Worker %d started", workerID)

	for {
		select {
		case job := <-s.jobChannel:
			s.processJob(ctx, workerID, job, generationService)

		case <-s.shutdown:
			log.Printf("Worker %d shutting down", workerID)
			s.workerDone <- true
			return

		case <-ctx.Done():
			log.Printf("Worker %d stopped due to context cancellation", workerID)
			return
		}
	}
}

func (s *channelQueueService) processJob(ctx context.Context, workerID int, job *Job, generationService GenerationService) {
	jobID := job.JobID
	log.Printf("Worker %d processing job %s", workerID, jobID)

	// Update status to processing
	if err := s.generationRepo.UpdateStatus(ctx, jobID, "processing", 10); err != nil {
		log.Printf("Worker %d failed to update job %s status to processing: %v", workerID, jobID, err)
		return
	}

	// Get template for generation
	templateRepo := generationService.(*generationServiceImpl).templateRepo
	template, err := templateRepo.GetByID(ctx, job.TemplateID)
	if err != nil {
		log.Printf("Worker %d failed to get template for job %s: %v", workerID, jobID, err)
		s.generationRepo.UpdateWithError(ctx, jobID, fmt.Sprintf("Template not found: %v", err), time.Now())
		return
	}

	// Update progress
	if err := s.generationRepo.UpdateStatus(ctx, jobID, "processing", 50); err != nil {
		log.Printf("Worker %d failed to update job %s progress: %v", workerID, jobID, err)
	}

	// Process the image generation
	comfyuiRepo := generationService.(*generationServiceImpl).comfyuiRepo
	imageURLs, err := comfyuiRepo.GenerateImage(ctx, job.TemplateID, &bytesReader{data: job.ImageData})
	if err != nil {
		log.Printf("Worker %d failed to generate image for job %s: %v", workerID, jobID, err)
		s.generationRepo.UpdateWithError(ctx, jobID, fmt.Sprintf("Generation failed: %v", err), time.Now())
		return
	}

	// Update with result
	if len(imageURLs) > 0 {
		if err := s.generationRepo.UpdateWithResult(ctx, jobID, "completed", imageURLs[0], time.Now()); err != nil {
			log.Printf("Worker %d failed to update job %s with result: %v", workerID, jobID, err)
			return
		}

		// Handle credit deduction and transaction recording
		userRepo := generationService.(*generationServiceImpl).userRepo
		transactionRepo := generationService.(*generationServiceImpl).transactionRepo

		if err := userRepo.UpdateCredits(ctx, job.UserID, -template.CreditCost); err != nil {
			log.Printf("Worker %d failed to deduct credits for job %s: %v", workerID, jobID, err)
		}

		transaction := &model.Transaction{
			UserID:            job.UserID,
			Type:              "generation",
			Amount:            -template.CreditCost,
			Description:       fmt.Sprintf("Used '%s' template", template.Name),
			RelatedTemplateID: &template.ID,
		}
		if err := transactionRepo.Create(ctx, transaction); err != nil {
			log.Printf("Worker %d failed to create transaction for job %s: %v", workerID, jobID, err)
		}

		log.Printf("Worker %d completed job %s successfully", workerID, jobID)
	} else {
		log.Printf("Worker %d: No images generated for job %s", workerID, jobID)
		s.generationRepo.UpdateWithError(ctx, jobID, "No images generated", time.Now())
	}
}

func (s *channelQueueService) Shutdown(ctx context.Context) error {
	log.Println("Shutting down queue service...")
	close(s.shutdown)
	
	// Wait for workers to finish (with timeout)
	timeout := time.After(30 * time.Second)
	for {
		select {
		case <-s.workerDone:
			// Worker finished
		case <-timeout:
			log.Println("Queue shutdown timed out")
			return fmt.Errorf("shutdown timed out")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// generateJobID creates a unique job identifier
func generateJobID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// bytesReader implements io.Reader for []byte
type bytesReader struct {
	data []byte
	pos  int
}

func (r *bytesReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, fmt.Errorf("EOF")
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
} 