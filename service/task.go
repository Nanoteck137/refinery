package service

import (
	"context"
	"errors"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/maruel/natural"
	"github.com/nanoteck137/refinery/tools/broker"
	"github.com/nanoteck137/refinery/tools/utils"
	"github.com/robfig/cron/v3"
)

type UpdateFunc func()

type Task interface {
	Name() string
	Schedule() string
	Run(ctx context.Context) error
}

var _ (broker.Event) = (*TaskSyncStateEvent)(nil)

type TaskSyncStateEventTask struct {
	Name      string `json:"name"`
	IsRunning bool   `json:"isRunning"`
}

type TaskSyncStateEvent struct {
	Tasks []TaskSyncStateEventTask `json:"tasks"`
}

func (e TaskSyncStateEvent) GetEventType() string {
	return "task-sync-state"
}

type taskEntry struct {
	task    Task
	entryId *cron.EntryID

	isRunning bool

	lastRunSuccess bool
	lastRunTime    *time.Duration
	lastRunError   error
}

type TaskService struct {
	logger *slog.Logger

	cron  *cron.Cron
	tasks map[string]*taskEntry

	wg sync.WaitGroup
	mu sync.RWMutex

	updateFunc UpdateFunc
}

func NewTaskService(logger *slog.Logger) *TaskService {
	return &TaskService{
		logger: logger,
		cron:   cron.New(),
		tasks:  map[string]*taskEntry{},
	}
}

func (s *TaskService) GetSyncStateEvent() TaskSyncStateEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := TaskSyncStateEvent{
		Tasks: []TaskSyncStateEventTask{},
	}

	for _, task := range s.tasks {
		res.Tasks = append(res.Tasks, TaskSyncStateEventTask{
			Name:      task.task.Name(),
			IsRunning: task.isRunning,
		})
	}

	sort.SliceStable(res.Tasks, func(i, j int) bool {
		return natural.Less(res.Tasks[i].Name, res.Tasks[j].Name)
	})

	return res
}

func (s *TaskService) SetUpdateFunc(f UpdateFunc) {
	s.updateFunc = f
}

func (s *TaskService) update() {
	if s.updateFunc != nil {
		s.updateFunc()
	}
}

func (s *TaskService) AddTask(task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.tasks[task.Name()]
	if exists {
		s.logger.Error("task with name already exists",
			slog.String("name", task.Name()),
		)

		return errors.New("task with name already exists: " + task.Name())
	}

	schedule := task.Schedule()

	var entryId *cron.EntryID

	if schedule != "" {
		id, err := s.cron.AddFunc(task.Schedule(), func() {
			s.RunTask(context.Background(), task.Name())
		})
		if err != nil {
			s.logger.Error("failed to add task to cron instance", "err", err)
			return err
		}

		entryId = &id
	}

	s.tasks[task.Name()] = &taskEntry{
		task:    task,
		entryId: entryId,
	}

	s.logger.Info("added task",
		slog.String("name", task.Name()),
	)

	return nil
}

func (s *TaskService) DisplayTasks() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, task := range s.tasks {
		s.logger.Info("registered task",
			"name", task.task.Name(),
			"schedule", task.task.Schedule(),
		)
	}
}

func (s *TaskService) Start() {
	s.cron.Start()
}

func (s *TaskService) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()

	s.Wait()
}

func (s *TaskService) Wait() {
	s.wg.Wait()
}

func (s *TaskService) RunTask(ctx context.Context, name string) {
	s.mu.Lock()
	entry, exists := s.tasks[name]
	if !exists {
		s.mu.Unlock()
		s.logger.Error("no task with name", "name", name)
		return
	}

	if entry.isRunning {
		s.mu.Unlock()
		s.logger.Error("task already running", "name", name)
		return
	}

	entry.isRunning = true
	s.mu.Unlock()

	s.update()

	s.wg.Add(1)
	defer func() {
		s.wg.Done()

		s.mu.Lock()
		entry.isRunning = false
		s.mu.Unlock()

		s.update()
	}()

	s.logger.Info("running task", "task", entry.task.Name())

	timer := utils.SimpleTimer{}
	timer.Start()

	err := entry.task.Run(ctx)

	s.mu.Lock()
	dur := timer.Stop()
	entry.lastRunTime = &dur

	if err != nil {
		s.logger.Error("task returned error", "err", err, "task", name, "duration", dur)
		entry.lastRunError = err
		entry.lastRunSuccess = false
	} else {
		s.logger.Info("task was successfully", "task", name, "duration", dur)
		entry.lastRunError = nil
		entry.lastRunSuccess = true
	}
	s.mu.Unlock()
}
