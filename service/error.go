package service

import (
	"errors"
	"fmt"
)

type ServiceError struct {
	Service string
	Message string
	Err     error
}

func (e *ServiceError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s service: %s: %s", e.Service, e.Message, e.Err)
	}

	return fmt.Sprintf("%s service: %s", e.Service, e.Err)
}

func (e *ServiceError) Unwrap() error {
    return e.Err
}

type ServiceErrCreator struct {
	Service string
}

func NewServiceErrCreator(service string) ServiceErrCreator {
	return ServiceErrCreator{
		Service: service,
	}
}

func (s *ServiceErrCreator) Wrap(message string, err error) error {
	return &ServiceError{
		Service: s.Service,
		Message: message,
		Err:     err,
	}
}

func (s *ServiceErrCreator) New(text string) error {
	return &ServiceError{
		Service: s.Service,
		Err:     errors.New(text),
	}
}

func (s *ServiceErrCreator) Newf(format string, a ...any) error {
	return &ServiceError{
		Service: s.Service,
		Err:     fmt.Errorf(format, a...),
	}
}

type InvalidFilterError struct {
	Service string
	Message string
}

func (e *InvalidFilterError) Error() string {
	return fmt.Sprintf("%s: invalid filter: %s", e.Service, e.Message)
}

type InvalidSortError struct {
	Service string
	Message string
}

func (e *InvalidSortError) Error() string {
	return fmt.Sprintf("%s: invalid sort: %s", e.Service, e.Message)
}
