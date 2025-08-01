// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	models "Events-Service/internal/models"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// GetEvents is an autogenerated mock type for the GetEvents type
type GetEvents struct {
	mock.Mock
}

// GetEventsByDay provides a mock function with given fields: userID, date
func (_m *GetEvents) GetEventsByDay(userID int64, date string) ([]models.Event, error) {
	ret := _m.Called(userID, date)

	if len(ret) == 0 {
		panic("no return value specified for GetEventsByDay")
	}

	var r0 []models.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, string) ([]models.Event, error)); ok {
		return rf(userID, date)
	}
	if rf, ok := ret.Get(0).(func(int64, string) []models.Event); ok {
		r0 = rf(userID, date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, string) error); ok {
		r1 = rf(userID, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEventsByMonth provides a mock function with given fields: userID, year, month
func (_m *GetEvents) GetEventsByMonth(userID int64, year int, month time.Month) ([]models.Event, error) {
	ret := _m.Called(userID, year, month)

	if len(ret) == 0 {
		panic("no return value specified for GetEventsByMonth")
	}

	var r0 []models.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, int, time.Month) ([]models.Event, error)); ok {
		return rf(userID, year, month)
	}
	if rf, ok := ret.Get(0).(func(int64, int, time.Month) []models.Event); ok {
		r0 = rf(userID, year, month)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, int, time.Month) error); ok {
		r1 = rf(userID, year, month)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEventsByWeek provides a mock function with given fields: userID, date
func (_m *GetEvents) GetEventsByWeek(userID int64, date time.Time) ([]models.Event, error) {
	ret := _m.Called(userID, date)

	if len(ret) == 0 {
		panic("no return value specified for GetEventsByWeek")
	}

	var r0 []models.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, time.Time) ([]models.Event, error)); ok {
		return rf(userID, date)
	}
	if rf, ok := ret.Get(0).(func(int64, time.Time) []models.Event); ok {
		r0 = rf(userID, date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, time.Time) error); ok {
		r1 = rf(userID, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewGetEvents creates a new instance of GetEvents. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGetEvents(t interface {
	mock.TestingT
	Cleanup(func())
}) *GetEvents {
	mock := &GetEvents{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
