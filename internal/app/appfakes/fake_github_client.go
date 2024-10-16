// Code generated by counterfeiter. DO NOT EDIT.
package appfakes

import (
	"context"
	"sync"

	"github.com/Mgla96/snappr/internal/adapters/clients"
	"github.com/google/go-github/v66/github"
)

type FakeGithubClient struct {
	AddCommentToPullRequestReviewStub        func(context.Context, string, string, int, string, string, string, int, int, clients.Side, clients.Side) (*github.PullRequestComment, error)
	addCommentToPullRequestReviewMutex       sync.RWMutex
	addCommentToPullRequestReviewArgsForCall []struct {
		arg1  context.Context
		arg2  string
		arg3  string
		arg4  int
		arg5  string
		arg6  string
		arg7  string
		arg8  int
		arg9  int
		arg10 clients.Side
		arg11 clients.Side
	}
	addCommentToPullRequestReviewReturns struct {
		result1 *github.PullRequestComment
		result2 error
	}
	addCommentToPullRequestReviewReturnsOnCall map[int]struct {
		result1 *github.PullRequestComment
		result2 error
	}
	AddCommitToBranchStub        func(context.Context, string, string, string, string, string, []byte) error
	addCommitToBranchMutex       sync.RWMutex
	addCommitToBranchArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
		arg5 string
		arg6 string
		arg7 []byte
	}
	addCommitToBranchReturns struct {
		result1 error
	}
	addCommitToBranchReturnsOnCall map[int]struct {
		result1 error
	}
	CreateBranchStub        func(context.Context, string, string, string, string) error
	createBranchMutex       sync.RWMutex
	createBranchArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
		arg5 string
	}
	createBranchReturns struct {
		result1 error
	}
	createBranchReturnsOnCall map[int]struct {
		result1 error
	}
	CreatePullRequestStub        func(context.Context, string, string, string, string, string, string) (*github.PullRequest, error)
	createPullRequestMutex       sync.RWMutex
	createPullRequestArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
		arg5 string
		arg6 string
		arg7 string
	}
	createPullRequestReturns struct {
		result1 *github.PullRequest
		result2 error
	}
	createPullRequestReturnsOnCall map[int]struct {
		result1 *github.PullRequest
		result2 error
	}
	GetCommitCodeStub        func(context.Context, string, string, string, clients.CodeFilter) (map[string]string, error)
	getCommitCodeMutex       sync.RWMutex
	getCommitCodeArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
		arg5 clients.CodeFilter
	}
	getCommitCodeReturns struct {
		result1 map[string]string
		result2 error
	}
	getCommitCodeReturnsOnCall map[int]struct {
		result1 map[string]string
		result2 error
	}
	GetLatestCommitFromBranchStub        func(context.Context, string, string, string) (string, error)
	getLatestCommitFromBranchMutex       sync.RWMutex
	getLatestCommitFromBranchArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
	}
	getLatestCommitFromBranchReturns struct {
		result1 string
		result2 error
	}
	getLatestCommitFromBranchReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	GetPRCodeStub        func(context.Context, string, string, int, *github.ListOptions, clients.CodeFilter) (map[string]string, error)
	getPRCodeMutex       sync.RWMutex
	getPRCodeArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int
		arg5 *github.ListOptions
		arg6 clients.CodeFilter
	}
	getPRCodeReturns struct {
		result1 map[string]string
		result2 error
	}
	getPRCodeReturnsOnCall map[int]struct {
		result1 map[string]string
		result2 error
	}
	GetPRPatchStub        func(context.Context, string, string, int) (string, error)
	getPRPatchMutex       sync.RWMutex
	getPRPatchArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int
	}
	getPRPatchReturns struct {
		result1 string
		result2 error
	}
	getPRPatchReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	ListPullRequestsStub        func(context.Context, string, string, *github.PullRequestListOptions) ([]*github.PullRequest, error)
	listPullRequestsMutex       sync.RWMutex
	listPullRequestsArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 *github.PullRequestListOptions
	}
	listPullRequestsReturns struct {
		result1 []*github.PullRequest
		result2 error
	}
	listPullRequestsReturnsOnCall map[int]struct {
		result1 []*github.PullRequest
		result2 error
	}
	MergePullRequestStub        func(context.Context, string, string, int, string) (*github.PullRequestMergeResult, error)
	mergePullRequestMutex       sync.RWMutex
	mergePullRequestArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int
		arg5 string
	}
	mergePullRequestReturns struct {
		result1 *github.PullRequestMergeResult
		result2 error
	}
	mergePullRequestReturnsOnCall map[int]struct {
		result1 *github.PullRequestMergeResult
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeGithubClient) AddCommentToPullRequestReview(arg1 context.Context, arg2 string, arg3 string, arg4 int, arg5 string, arg6 string, arg7 string, arg8 int, arg9 int, arg10 clients.Side, arg11 clients.Side) (*github.PullRequestComment, error) {
	fake.addCommentToPullRequestReviewMutex.Lock()
	ret, specificReturn := fake.addCommentToPullRequestReviewReturnsOnCall[len(fake.addCommentToPullRequestReviewArgsForCall)]
	fake.addCommentToPullRequestReviewArgsForCall = append(fake.addCommentToPullRequestReviewArgsForCall, struct {
		arg1  context.Context
		arg2  string
		arg3  string
		arg4  int
		arg5  string
		arg6  string
		arg7  string
		arg8  int
		arg9  int
		arg10 clients.Side
		arg11 clients.Side
	}{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11})
	stub := fake.AddCommentToPullRequestReviewStub
	fakeReturns := fake.addCommentToPullRequestReviewReturns
	fake.recordInvocation("AddCommentToPullRequestReview", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11})
	fake.addCommentToPullRequestReviewMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGithubClient) AddCommentToPullRequestReviewCallCount() int {
	fake.addCommentToPullRequestReviewMutex.RLock()
	defer fake.addCommentToPullRequestReviewMutex.RUnlock()
	return len(fake.addCommentToPullRequestReviewArgsForCall)
}

func (fake *FakeGithubClient) AddCommentToPullRequestReviewCalls(stub func(context.Context, string, string, int, string, string, string, int, int, clients.Side, clients.Side) (*github.PullRequestComment, error)) {
	fake.addCommentToPullRequestReviewMutex.Lock()
	defer fake.addCommentToPullRequestReviewMutex.Unlock()
	fake.AddCommentToPullRequestReviewStub = stub
}

func (fake *FakeGithubClient) AddCommentToPullRequestReviewArgsForCall(i int) (context.Context, string, string, int, string, string, string, int, int, clients.Side, clients.Side) {
	fake.addCommentToPullRequestReviewMutex.RLock()
	defer fake.addCommentToPullRequestReviewMutex.RUnlock()
	argsForCall := fake.addCommentToPullRequestReviewArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6, argsForCall.arg7, argsForCall.arg8, argsForCall.arg9, argsForCall.arg10, argsForCall.arg11
}

func (fake *FakeGithubClient) AddCommentToPullRequestReviewReturns(result1 *github.PullRequestComment, result2 error) {
	fake.addCommentToPullRequestReviewMutex.Lock()
	defer fake.addCommentToPullRequestReviewMutex.Unlock()
	fake.AddCommentToPullRequestReviewStub = nil
	fake.addCommentToPullRequestReviewReturns = struct {
		result1 *github.PullRequestComment
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) AddCommentToPullRequestReviewReturnsOnCall(i int, result1 *github.PullRequestComment, result2 error) {
	fake.addCommentToPullRequestReviewMutex.Lock()
	defer fake.addCommentToPullRequestReviewMutex.Unlock()
	fake.AddCommentToPullRequestReviewStub = nil
	if fake.addCommentToPullRequestReviewReturnsOnCall == nil {
		fake.addCommentToPullRequestReviewReturnsOnCall = make(map[int]struct {
			result1 *github.PullRequestComment
			result2 error
		})
	}
	fake.addCommentToPullRequestReviewReturnsOnCall[i] = struct {
		result1 *github.PullRequestComment
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) AddCommitToBranch(arg1 context.Context, arg2 string, arg3 string, arg4 string, arg5 string, arg6 string, arg7 []byte) error {
	var arg7Copy []byte
	if arg7 != nil {
		arg7Copy = make([]byte, len(arg7))
		copy(arg7Copy, arg7)
	}
	fake.addCommitToBranchMutex.Lock()
	ret, specificReturn := fake.addCommitToBranchReturnsOnCall[len(fake.addCommitToBranchArgsForCall)]
	fake.addCommitToBranchArgsForCall = append(fake.addCommitToBranchArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
		arg5 string
		arg6 string
		arg7 []byte
	}{arg1, arg2, arg3, arg4, arg5, arg6, arg7Copy})
	stub := fake.AddCommitToBranchStub
	fakeReturns := fake.addCommitToBranchReturns
	fake.recordInvocation("AddCommitToBranch", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6, arg7Copy})
	fake.addCommitToBranchMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGithubClient) AddCommitToBranchCallCount() int {
	fake.addCommitToBranchMutex.RLock()
	defer fake.addCommitToBranchMutex.RUnlock()
	return len(fake.addCommitToBranchArgsForCall)
}

func (fake *FakeGithubClient) AddCommitToBranchCalls(stub func(context.Context, string, string, string, string, string, []byte) error) {
	fake.addCommitToBranchMutex.Lock()
	defer fake.addCommitToBranchMutex.Unlock()
	fake.AddCommitToBranchStub = stub
}

func (fake *FakeGithubClient) AddCommitToBranchArgsForCall(i int) (context.Context, string, string, string, string, string, []byte) {
	fake.addCommitToBranchMutex.RLock()
	defer fake.addCommitToBranchMutex.RUnlock()
	argsForCall := fake.addCommitToBranchArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6, argsForCall.arg7
}

func (fake *FakeGithubClient) AddCommitToBranchReturns(result1 error) {
	fake.addCommitToBranchMutex.Lock()
	defer fake.addCommitToBranchMutex.Unlock()
	fake.AddCommitToBranchStub = nil
	fake.addCommitToBranchReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGithubClient) AddCommitToBranchReturnsOnCall(i int, result1 error) {
	fake.addCommitToBranchMutex.Lock()
	defer fake.addCommitToBranchMutex.Unlock()
	fake.AddCommitToBranchStub = nil
	if fake.addCommitToBranchReturnsOnCall == nil {
		fake.addCommitToBranchReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addCommitToBranchReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGithubClient) CreateBranch(arg1 context.Context, arg2 string, arg3 string, arg4 string, arg5 string) error {
	fake.createBranchMutex.Lock()
	ret, specificReturn := fake.createBranchReturnsOnCall[len(fake.createBranchArgsForCall)]
	fake.createBranchArgsForCall = append(fake.createBranchArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
		arg5 string
	}{arg1, arg2, arg3, arg4, arg5})
	stub := fake.CreateBranchStub
	fakeReturns := fake.createBranchReturns
	fake.recordInvocation("CreateBranch", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.createBranchMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeGithubClient) CreateBranchCallCount() int {
	fake.createBranchMutex.RLock()
	defer fake.createBranchMutex.RUnlock()
	return len(fake.createBranchArgsForCall)
}

func (fake *FakeGithubClient) CreateBranchCalls(stub func(context.Context, string, string, string, string) error) {
	fake.createBranchMutex.Lock()
	defer fake.createBranchMutex.Unlock()
	fake.CreateBranchStub = stub
}

func (fake *FakeGithubClient) CreateBranchArgsForCall(i int) (context.Context, string, string, string, string) {
	fake.createBranchMutex.RLock()
	defer fake.createBranchMutex.RUnlock()
	argsForCall := fake.createBranchArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *FakeGithubClient) CreateBranchReturns(result1 error) {
	fake.createBranchMutex.Lock()
	defer fake.createBranchMutex.Unlock()
	fake.CreateBranchStub = nil
	fake.createBranchReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGithubClient) CreateBranchReturnsOnCall(i int, result1 error) {
	fake.createBranchMutex.Lock()
	defer fake.createBranchMutex.Unlock()
	fake.CreateBranchStub = nil
	if fake.createBranchReturnsOnCall == nil {
		fake.createBranchReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createBranchReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGithubClient) CreatePullRequest(arg1 context.Context, arg2 string, arg3 string, arg4 string, arg5 string, arg6 string, arg7 string) (*github.PullRequest, error) {
	fake.createPullRequestMutex.Lock()
	ret, specificReturn := fake.createPullRequestReturnsOnCall[len(fake.createPullRequestArgsForCall)]
	fake.createPullRequestArgsForCall = append(fake.createPullRequestArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
		arg5 string
		arg6 string
		arg7 string
	}{arg1, arg2, arg3, arg4, arg5, arg6, arg7})
	stub := fake.CreatePullRequestStub
	fakeReturns := fake.createPullRequestReturns
	fake.recordInvocation("CreatePullRequest", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6, arg7})
	fake.createPullRequestMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGithubClient) CreatePullRequestCallCount() int {
	fake.createPullRequestMutex.RLock()
	defer fake.createPullRequestMutex.RUnlock()
	return len(fake.createPullRequestArgsForCall)
}

func (fake *FakeGithubClient) CreatePullRequestCalls(stub func(context.Context, string, string, string, string, string, string) (*github.PullRequest, error)) {
	fake.createPullRequestMutex.Lock()
	defer fake.createPullRequestMutex.Unlock()
	fake.CreatePullRequestStub = stub
}

func (fake *FakeGithubClient) CreatePullRequestArgsForCall(i int) (context.Context, string, string, string, string, string, string) {
	fake.createPullRequestMutex.RLock()
	defer fake.createPullRequestMutex.RUnlock()
	argsForCall := fake.createPullRequestArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6, argsForCall.arg7
}

func (fake *FakeGithubClient) CreatePullRequestReturns(result1 *github.PullRequest, result2 error) {
	fake.createPullRequestMutex.Lock()
	defer fake.createPullRequestMutex.Unlock()
	fake.CreatePullRequestStub = nil
	fake.createPullRequestReturns = struct {
		result1 *github.PullRequest
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) CreatePullRequestReturnsOnCall(i int, result1 *github.PullRequest, result2 error) {
	fake.createPullRequestMutex.Lock()
	defer fake.createPullRequestMutex.Unlock()
	fake.CreatePullRequestStub = nil
	if fake.createPullRequestReturnsOnCall == nil {
		fake.createPullRequestReturnsOnCall = make(map[int]struct {
			result1 *github.PullRequest
			result2 error
		})
	}
	fake.createPullRequestReturnsOnCall[i] = struct {
		result1 *github.PullRequest
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) GetCommitCode(arg1 context.Context, arg2 string, arg3 string, arg4 string, arg5 clients.CodeFilter) (map[string]string, error) {
	fake.getCommitCodeMutex.Lock()
	ret, specificReturn := fake.getCommitCodeReturnsOnCall[len(fake.getCommitCodeArgsForCall)]
	fake.getCommitCodeArgsForCall = append(fake.getCommitCodeArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
		arg5 clients.CodeFilter
	}{arg1, arg2, arg3, arg4, arg5})
	stub := fake.GetCommitCodeStub
	fakeReturns := fake.getCommitCodeReturns
	fake.recordInvocation("GetCommitCode", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.getCommitCodeMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGithubClient) GetCommitCodeCallCount() int {
	fake.getCommitCodeMutex.RLock()
	defer fake.getCommitCodeMutex.RUnlock()
	return len(fake.getCommitCodeArgsForCall)
}

func (fake *FakeGithubClient) GetCommitCodeCalls(stub func(context.Context, string, string, string, clients.CodeFilter) (map[string]string, error)) {
	fake.getCommitCodeMutex.Lock()
	defer fake.getCommitCodeMutex.Unlock()
	fake.GetCommitCodeStub = stub
}

func (fake *FakeGithubClient) GetCommitCodeArgsForCall(i int) (context.Context, string, string, string, clients.CodeFilter) {
	fake.getCommitCodeMutex.RLock()
	defer fake.getCommitCodeMutex.RUnlock()
	argsForCall := fake.getCommitCodeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *FakeGithubClient) GetCommitCodeReturns(result1 map[string]string, result2 error) {
	fake.getCommitCodeMutex.Lock()
	defer fake.getCommitCodeMutex.Unlock()
	fake.GetCommitCodeStub = nil
	fake.getCommitCodeReturns = struct {
		result1 map[string]string
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) GetCommitCodeReturnsOnCall(i int, result1 map[string]string, result2 error) {
	fake.getCommitCodeMutex.Lock()
	defer fake.getCommitCodeMutex.Unlock()
	fake.GetCommitCodeStub = nil
	if fake.getCommitCodeReturnsOnCall == nil {
		fake.getCommitCodeReturnsOnCall = make(map[int]struct {
			result1 map[string]string
			result2 error
		})
	}
	fake.getCommitCodeReturnsOnCall[i] = struct {
		result1 map[string]string
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) GetLatestCommitFromBranch(arg1 context.Context, arg2 string, arg3 string, arg4 string) (string, error) {
	fake.getLatestCommitFromBranchMutex.Lock()
	ret, specificReturn := fake.getLatestCommitFromBranchReturnsOnCall[len(fake.getLatestCommitFromBranchArgsForCall)]
	fake.getLatestCommitFromBranchArgsForCall = append(fake.getLatestCommitFromBranchArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	stub := fake.GetLatestCommitFromBranchStub
	fakeReturns := fake.getLatestCommitFromBranchReturns
	fake.recordInvocation("GetLatestCommitFromBranch", []interface{}{arg1, arg2, arg3, arg4})
	fake.getLatestCommitFromBranchMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGithubClient) GetLatestCommitFromBranchCallCount() int {
	fake.getLatestCommitFromBranchMutex.RLock()
	defer fake.getLatestCommitFromBranchMutex.RUnlock()
	return len(fake.getLatestCommitFromBranchArgsForCall)
}

func (fake *FakeGithubClient) GetLatestCommitFromBranchCalls(stub func(context.Context, string, string, string) (string, error)) {
	fake.getLatestCommitFromBranchMutex.Lock()
	defer fake.getLatestCommitFromBranchMutex.Unlock()
	fake.GetLatestCommitFromBranchStub = stub
}

func (fake *FakeGithubClient) GetLatestCommitFromBranchArgsForCall(i int) (context.Context, string, string, string) {
	fake.getLatestCommitFromBranchMutex.RLock()
	defer fake.getLatestCommitFromBranchMutex.RUnlock()
	argsForCall := fake.getLatestCommitFromBranchArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeGithubClient) GetLatestCommitFromBranchReturns(result1 string, result2 error) {
	fake.getLatestCommitFromBranchMutex.Lock()
	defer fake.getLatestCommitFromBranchMutex.Unlock()
	fake.GetLatestCommitFromBranchStub = nil
	fake.getLatestCommitFromBranchReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) GetLatestCommitFromBranchReturnsOnCall(i int, result1 string, result2 error) {
	fake.getLatestCommitFromBranchMutex.Lock()
	defer fake.getLatestCommitFromBranchMutex.Unlock()
	fake.GetLatestCommitFromBranchStub = nil
	if fake.getLatestCommitFromBranchReturnsOnCall == nil {
		fake.getLatestCommitFromBranchReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getLatestCommitFromBranchReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) GetPRCode(arg1 context.Context, arg2 string, arg3 string, arg4 int, arg5 *github.ListOptions, arg6 clients.CodeFilter) (map[string]string, error) {
	fake.getPRCodeMutex.Lock()
	ret, specificReturn := fake.getPRCodeReturnsOnCall[len(fake.getPRCodeArgsForCall)]
	fake.getPRCodeArgsForCall = append(fake.getPRCodeArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int
		arg5 *github.ListOptions
		arg6 clients.CodeFilter
	}{arg1, arg2, arg3, arg4, arg5, arg6})
	stub := fake.GetPRCodeStub
	fakeReturns := fake.getPRCodeReturns
	fake.recordInvocation("GetPRCode", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.getPRCodeMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGithubClient) GetPRCodeCallCount() int {
	fake.getPRCodeMutex.RLock()
	defer fake.getPRCodeMutex.RUnlock()
	return len(fake.getPRCodeArgsForCall)
}

func (fake *FakeGithubClient) GetPRCodeCalls(stub func(context.Context, string, string, int, *github.ListOptions, clients.CodeFilter) (map[string]string, error)) {
	fake.getPRCodeMutex.Lock()
	defer fake.getPRCodeMutex.Unlock()
	fake.GetPRCodeStub = stub
}

func (fake *FakeGithubClient) GetPRCodeArgsForCall(i int) (context.Context, string, string, int, *github.ListOptions, clients.CodeFilter) {
	fake.getPRCodeMutex.RLock()
	defer fake.getPRCodeMutex.RUnlock()
	argsForCall := fake.getPRCodeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6
}

func (fake *FakeGithubClient) GetPRCodeReturns(result1 map[string]string, result2 error) {
	fake.getPRCodeMutex.Lock()
	defer fake.getPRCodeMutex.Unlock()
	fake.GetPRCodeStub = nil
	fake.getPRCodeReturns = struct {
		result1 map[string]string
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) GetPRCodeReturnsOnCall(i int, result1 map[string]string, result2 error) {
	fake.getPRCodeMutex.Lock()
	defer fake.getPRCodeMutex.Unlock()
	fake.GetPRCodeStub = nil
	if fake.getPRCodeReturnsOnCall == nil {
		fake.getPRCodeReturnsOnCall = make(map[int]struct {
			result1 map[string]string
			result2 error
		})
	}
	fake.getPRCodeReturnsOnCall[i] = struct {
		result1 map[string]string
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) GetPRPatch(arg1 context.Context, arg2 string, arg3 string, arg4 int) (string, error) {
	fake.getPRPatchMutex.Lock()
	ret, specificReturn := fake.getPRPatchReturnsOnCall[len(fake.getPRPatchArgsForCall)]
	fake.getPRPatchArgsForCall = append(fake.getPRPatchArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int
	}{arg1, arg2, arg3, arg4})
	stub := fake.GetPRPatchStub
	fakeReturns := fake.getPRPatchReturns
	fake.recordInvocation("GetPRPatch", []interface{}{arg1, arg2, arg3, arg4})
	fake.getPRPatchMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGithubClient) GetPRPatchCallCount() int {
	fake.getPRPatchMutex.RLock()
	defer fake.getPRPatchMutex.RUnlock()
	return len(fake.getPRPatchArgsForCall)
}

func (fake *FakeGithubClient) GetPRPatchCalls(stub func(context.Context, string, string, int) (string, error)) {
	fake.getPRPatchMutex.Lock()
	defer fake.getPRPatchMutex.Unlock()
	fake.GetPRPatchStub = stub
}

func (fake *FakeGithubClient) GetPRPatchArgsForCall(i int) (context.Context, string, string, int) {
	fake.getPRPatchMutex.RLock()
	defer fake.getPRPatchMutex.RUnlock()
	argsForCall := fake.getPRPatchArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeGithubClient) GetPRPatchReturns(result1 string, result2 error) {
	fake.getPRPatchMutex.Lock()
	defer fake.getPRPatchMutex.Unlock()
	fake.GetPRPatchStub = nil
	fake.getPRPatchReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) GetPRPatchReturnsOnCall(i int, result1 string, result2 error) {
	fake.getPRPatchMutex.Lock()
	defer fake.getPRPatchMutex.Unlock()
	fake.GetPRPatchStub = nil
	if fake.getPRPatchReturnsOnCall == nil {
		fake.getPRPatchReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getPRPatchReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) ListPullRequests(arg1 context.Context, arg2 string, arg3 string, arg4 *github.PullRequestListOptions) ([]*github.PullRequest, error) {
	fake.listPullRequestsMutex.Lock()
	ret, specificReturn := fake.listPullRequestsReturnsOnCall[len(fake.listPullRequestsArgsForCall)]
	fake.listPullRequestsArgsForCall = append(fake.listPullRequestsArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 *github.PullRequestListOptions
	}{arg1, arg2, arg3, arg4})
	stub := fake.ListPullRequestsStub
	fakeReturns := fake.listPullRequestsReturns
	fake.recordInvocation("ListPullRequests", []interface{}{arg1, arg2, arg3, arg4})
	fake.listPullRequestsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGithubClient) ListPullRequestsCallCount() int {
	fake.listPullRequestsMutex.RLock()
	defer fake.listPullRequestsMutex.RUnlock()
	return len(fake.listPullRequestsArgsForCall)
}

func (fake *FakeGithubClient) ListPullRequestsCalls(stub func(context.Context, string, string, *github.PullRequestListOptions) ([]*github.PullRequest, error)) {
	fake.listPullRequestsMutex.Lock()
	defer fake.listPullRequestsMutex.Unlock()
	fake.ListPullRequestsStub = stub
}

func (fake *FakeGithubClient) ListPullRequestsArgsForCall(i int) (context.Context, string, string, *github.PullRequestListOptions) {
	fake.listPullRequestsMutex.RLock()
	defer fake.listPullRequestsMutex.RUnlock()
	argsForCall := fake.listPullRequestsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeGithubClient) ListPullRequestsReturns(result1 []*github.PullRequest, result2 error) {
	fake.listPullRequestsMutex.Lock()
	defer fake.listPullRequestsMutex.Unlock()
	fake.ListPullRequestsStub = nil
	fake.listPullRequestsReturns = struct {
		result1 []*github.PullRequest
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) ListPullRequestsReturnsOnCall(i int, result1 []*github.PullRequest, result2 error) {
	fake.listPullRequestsMutex.Lock()
	defer fake.listPullRequestsMutex.Unlock()
	fake.ListPullRequestsStub = nil
	if fake.listPullRequestsReturnsOnCall == nil {
		fake.listPullRequestsReturnsOnCall = make(map[int]struct {
			result1 []*github.PullRequest
			result2 error
		})
	}
	fake.listPullRequestsReturnsOnCall[i] = struct {
		result1 []*github.PullRequest
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) MergePullRequest(arg1 context.Context, arg2 string, arg3 string, arg4 int, arg5 string) (*github.PullRequestMergeResult, error) {
	fake.mergePullRequestMutex.Lock()
	ret, specificReturn := fake.mergePullRequestReturnsOnCall[len(fake.mergePullRequestArgsForCall)]
	fake.mergePullRequestArgsForCall = append(fake.mergePullRequestArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 int
		arg5 string
	}{arg1, arg2, arg3, arg4, arg5})
	stub := fake.MergePullRequestStub
	fakeReturns := fake.mergePullRequestReturns
	fake.recordInvocation("MergePullRequest", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.mergePullRequestMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGithubClient) MergePullRequestCallCount() int {
	fake.mergePullRequestMutex.RLock()
	defer fake.mergePullRequestMutex.RUnlock()
	return len(fake.mergePullRequestArgsForCall)
}

func (fake *FakeGithubClient) MergePullRequestCalls(stub func(context.Context, string, string, int, string) (*github.PullRequestMergeResult, error)) {
	fake.mergePullRequestMutex.Lock()
	defer fake.mergePullRequestMutex.Unlock()
	fake.MergePullRequestStub = stub
}

func (fake *FakeGithubClient) MergePullRequestArgsForCall(i int) (context.Context, string, string, int, string) {
	fake.mergePullRequestMutex.RLock()
	defer fake.mergePullRequestMutex.RUnlock()
	argsForCall := fake.mergePullRequestArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *FakeGithubClient) MergePullRequestReturns(result1 *github.PullRequestMergeResult, result2 error) {
	fake.mergePullRequestMutex.Lock()
	defer fake.mergePullRequestMutex.Unlock()
	fake.MergePullRequestStub = nil
	fake.mergePullRequestReturns = struct {
		result1 *github.PullRequestMergeResult
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) MergePullRequestReturnsOnCall(i int, result1 *github.PullRequestMergeResult, result2 error) {
	fake.mergePullRequestMutex.Lock()
	defer fake.mergePullRequestMutex.Unlock()
	fake.MergePullRequestStub = nil
	if fake.mergePullRequestReturnsOnCall == nil {
		fake.mergePullRequestReturnsOnCall = make(map[int]struct {
			result1 *github.PullRequestMergeResult
			result2 error
		})
	}
	fake.mergePullRequestReturnsOnCall[i] = struct {
		result1 *github.PullRequestMergeResult
		result2 error
	}{result1, result2}
}

func (fake *FakeGithubClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addCommentToPullRequestReviewMutex.RLock()
	defer fake.addCommentToPullRequestReviewMutex.RUnlock()
	fake.addCommitToBranchMutex.RLock()
	defer fake.addCommitToBranchMutex.RUnlock()
	fake.createBranchMutex.RLock()
	defer fake.createBranchMutex.RUnlock()
	fake.createPullRequestMutex.RLock()
	defer fake.createPullRequestMutex.RUnlock()
	fake.getCommitCodeMutex.RLock()
	defer fake.getCommitCodeMutex.RUnlock()
	fake.getLatestCommitFromBranchMutex.RLock()
	defer fake.getLatestCommitFromBranchMutex.RUnlock()
	fake.getPRCodeMutex.RLock()
	defer fake.getPRCodeMutex.RUnlock()
	fake.getPRPatchMutex.RLock()
	defer fake.getPRPatchMutex.RUnlock()
	fake.listPullRequestsMutex.RLock()
	defer fake.listPullRequestsMutex.RUnlock()
	fake.mergePullRequestMutex.RLock()
	defer fake.mergePullRequestMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeGithubClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
