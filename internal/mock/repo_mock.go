// Code generated by MockGen. DO NOT EDIT.
// Source: ocp-video-api/internal/repo (interfaces: Repo)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	models "ocp-video-api/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// AddVideo mocks base method.
func (m *MockRepo) AddVideo(arg0 context.Context, arg1 models.Video) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddVideo", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddVideo indicates an expected call of AddVideo.
func (mr *MockRepoMockRecorder) AddVideo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddVideo", reflect.TypeOf((*MockRepo)(nil).AddVideo), arg0, arg1)
}

// AddVideos mocks base method.
func (m *MockRepo) AddVideos(arg0 context.Context, arg1 []models.Video) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddVideos", arg0, arg1)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddVideos indicates an expected call of AddVideos.
func (mr *MockRepoMockRecorder) AddVideos(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddVideos", reflect.TypeOf((*MockRepo)(nil).AddVideos), arg0, arg1)
}

// GetVideo mocks base method.
func (m *MockRepo) GetVideo(arg0 context.Context, arg1 uint64) (*models.Video, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVideo", arg0, arg1)
	ret0, _ := ret[0].(*models.Video)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVideo indicates an expected call of GetVideo.
func (mr *MockRepoMockRecorder) GetVideo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVideo", reflect.TypeOf((*MockRepo)(nil).GetVideo), arg0, arg1)
}

// GetVideos mocks base method.
func (m *MockRepo) GetVideos(arg0 context.Context, arg1, arg2 uint64) ([]models.Video, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVideos", arg0, arg1, arg2)
	ret0, _ := ret[0].([]models.Video)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVideos indicates an expected call of GetVideos.
func (mr *MockRepoMockRecorder) GetVideos(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVideos", reflect.TypeOf((*MockRepo)(nil).GetVideos), arg0, arg1, arg2)
}

// RemoveVideo mocks base method.
func (m *MockRepo) RemoveVideo(arg0 context.Context, arg1 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveVideo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveVideo indicates an expected call of RemoveVideo.
func (mr *MockRepoMockRecorder) RemoveVideo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveVideo", reflect.TypeOf((*MockRepo)(nil).RemoveVideo), arg0, arg1)
}

// UpdateVideo mocks base method.
func (m *MockRepo) UpdateVideo(arg0 context.Context, arg1 models.Video) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVideo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateVideo indicates an expected call of UpdateVideo.
func (mr *MockRepoMockRecorder) UpdateVideo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVideo", reflect.TypeOf((*MockRepo)(nil).UpdateVideo), arg0, arg1)
}
