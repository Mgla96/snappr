package clients

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

// gitService is an interface for interacting with github git service
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . gitService
type gitService interface {
	CreateBlob(ctx context.Context, owner string, repo string, blob *github.Blob) (*github.Blob, *github.Response, error)
	CreateCommit(ctx context.Context, owner string, repo string, commit *github.Commit) (*github.Commit, *github.Response, error)
	CreateRef(ctx context.Context, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error)
	CreateTag(ctx context.Context, owner string, repo string, tag *github.Tag) (*github.Tag, *github.Response, error)
	CreateTree(ctx context.Context, owner string, repo string, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error)
	DeleteRef(ctx context.Context, owner string, repo string, ref string) (*github.Response, error)
	GetBlob(ctx context.Context, owner string, repo string, sha string) (*github.Blob, *github.Response, error)
	GetBlobRaw(ctx context.Context, owner string, repo string, sha string) ([]byte, *github.Response, error)
	GetCommit(ctx context.Context, owner string, repo string, sha string) (*github.Commit, *github.Response, error)
	GetRef(ctx context.Context, owner string, repo string, ref string) (*github.Reference, *github.Response, error)
	GetTag(ctx context.Context, owner string, repo string, sha string) (*github.Tag, *github.Response, error)
	GetTree(ctx context.Context, owner string, repo string, sha string, recursive bool) (*github.Tree, *github.Response, error)
	ListMatchingRefs(ctx context.Context, owner string, repo string, opts *github.ReferenceListOptions) ([]*github.Reference, *github.Response, error)
	UpdateRef(ctx context.Context, owner string, repo string, ref *github.Reference, force bool) (*github.Reference, *github.Response, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

// pullRequestService is an interface for interacting with github pull request service
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . pullRequestService
type pullRequestService interface {
	Create(ctx context.Context, owner string, repo string, pull *github.NewPullRequest) (*github.PullRequest, *github.Response, error)
	CreateComment(ctx context.Context, owner string, repo string, number int, comment *github.PullRequestComment) (*github.PullRequestComment, *github.Response, error)
	CreateCommentInReplyTo(ctx context.Context, owner string, repo string, number int, body string, commentID int64) (*github.PullRequestComment, *github.Response, error)
	CreateReview(ctx context.Context, owner string, repo string, number int, review *github.PullRequestReviewRequest) (*github.PullRequestReview, *github.Response, error)
	DeleteComment(ctx context.Context, owner string, repo string, commentID int64) (*github.Response, error)
	DeletePendingReview(ctx context.Context, owner string, repo string, number int, reviewID int64) (*github.PullRequestReview, *github.Response, error)
	DismissReview(ctx context.Context, owner string, repo string, number int, reviewID int64, review *github.PullRequestReviewDismissalRequest) (*github.PullRequestReview, *github.Response, error)
	Edit(ctx context.Context, owner string, repo string, number int, pull *github.PullRequest) (*github.PullRequest, *github.Response, error)
	EditComment(ctx context.Context, owner string, repo string, commentID int64, comment *github.PullRequestComment) (*github.PullRequestComment, *github.Response, error)
	Get(ctx context.Context, owner string, repo string, number int) (*github.PullRequest, *github.Response, error)
	GetComment(ctx context.Context, owner string, repo string, commentID int64) (*github.PullRequestComment, *github.Response, error)
	GetRaw(ctx context.Context, owner string, repo string, number int, opts github.RawOptions) (string, *github.Response, error)
	GetReview(ctx context.Context, owner string, repo string, number int, reviewID int64) (*github.PullRequestReview, *github.Response, error)
	IsMerged(ctx context.Context, owner string, repo string, number int) (bool, *github.Response, error)
	List(ctx context.Context, owner string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
	ListComments(ctx context.Context, owner string, repo string, number int, opts *github.PullRequestListCommentsOptions) ([]*github.PullRequestComment, *github.Response, error)
	ListCommits(ctx context.Context, owner string, repo string, number int, opts *github.ListOptions) ([]*github.RepositoryCommit, *github.Response, error)
	ListFiles(ctx context.Context, owner string, repo string, number int, opts *github.ListOptions) ([]*github.CommitFile, *github.Response, error)
	ListPullRequestsWithCommit(ctx context.Context, owner string, repo string, sha string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
	ListReviewComments(ctx context.Context, owner string, repo string, number int, reviewID int64, opts *github.ListOptions) ([]*github.PullRequestComment, *github.Response, error)
	ListReviewers(ctx context.Context, owner string, repo string, number int, opts *github.ListOptions) (*github.Reviewers, *github.Response, error)
	ListReviews(ctx context.Context, owner string, repo string, number int, opts *github.ListOptions) ([]*github.PullRequestReview, *github.Response, error)
	Merge(ctx context.Context, owner string, repo string, number int, commitMessage string, options *github.PullRequestOptions) (*github.PullRequestMergeResult, *github.Response, error)
	RemoveReviewers(ctx context.Context, owner string, repo string, number int, reviewers github.ReviewersRequest) (*github.Response, error)
	RequestReviewers(ctx context.Context, owner string, repo string, number int, reviewers github.ReviewersRequest) (*github.PullRequest, *github.Response, error)
	SubmitReview(ctx context.Context, owner string, repo string, number int, reviewID int64, review *github.PullRequestReviewRequest) (*github.PullRequestReview, *github.Response, error)
	UpdateBranch(ctx context.Context, owner string, repo string, number int, opts *github.PullRequestBranchUpdateOptions) (*github.PullRequestBranchUpdateResponse, *github.Response, error)
	UpdateReview(ctx context.Context, owner string, repo string, number int, reviewID int64, body string) (*github.PullRequestReview, *github.Response, error)
}

type GithubClient struct {
	ghGitClient         gitService
	ghPullRequestClient pullRequestService
	log                 zerolog.Logger
}

// NewGithubClient creates a new instance of the GithubClient.
func NewGithubClient(token string, logger zerolog.Logger) *GithubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &GithubClient{
		ghGitClient:         client.Git,
		ghPullRequestClient: client.PullRequests,
		log:                 logger,
	}
}

// CreateBranch creates a new branch from the specified base branch.
func (gc *GithubClient) CreateBranch(ctx context.Context, owner, repo, newBranch, baseBranch string) error {
	baseRef, _, err := gc.ghGitClient.GetRef(ctx, owner, repo, "refs/heads/"+baseBranch)
	if err != nil {
		return fmt.Errorf("failed to get base branch reference: %w", err)
	}

	newRef := &github.Reference{
		Ref:    github.String("refs/heads/" + newBranch),
		Object: &github.GitObject{SHA: baseRef.Object.SHA},
	}

	_, _, err = gc.ghGitClient.CreateRef(ctx, owner, repo, newRef)
	if err != nil {
		return fmt.Errorf("failed to create new branch reference: %w", err)
	}

	return nil
}

// GetLatestCommitFromBranch gets the latest commit SHA from a branch.
func (gc *GithubClient) GetLatestCommitFromBranch(ctx context.Context, owner, repo, branch string) (string, error) {
	ref, _, err := gc.ghGitClient.GetRef(ctx, owner, repo, "refs/heads/"+branch)
	if err != nil {
		return "", fmt.Errorf("failed to get branch reference: %w", err)
	}

	commitSHA := ref.GetObject().GetSHA()
	return commitSHA, nil
}

// AddCommitToBranch adds a commit to the specified branch.
func (gc *GithubClient) AddCommitToBranch(ctx context.Context, owner, repo, branch, filePath, commitMessage string, fileContent []byte) error {
	ref, _, err := gc.ghGitClient.GetRef(ctx, owner, repo, "refs/heads/"+branch)
	if err != nil {
		return err
	}

	blob := &github.Blob{
		Content:  github.String(string(fileContent)),
		Encoding: github.String("utf-8"),
	}
	blobRes, _, err := gc.ghGitClient.CreateBlob(ctx, owner, repo, blob)
	if err != nil {
		return err
	}

	treeEntry := &github.TreeEntry{
		Path: github.String(filePath),
		Mode: github.String("100644"),
		Type: github.String("blob"),
		SHA:  blobRes.SHA,
	}
	baseTree, _, err := gc.ghGitClient.GetTree(ctx, owner, repo, *ref.Object.SHA, false)
	if err != nil {
		return err
	}
	tree, _, err := gc.ghGitClient.CreateTree(ctx, owner, repo, *baseTree.SHA, []*github.TreeEntry{treeEntry})
	if err != nil {
		return err
	}

	parentCommit, _, err := gc.ghGitClient.GetCommit(ctx, owner, repo, *ref.Object.SHA)
	if err != nil {
		return err
	}
	commit := &github.Commit{
		Message: github.String(commitMessage),
		Tree:    tree,
		Parents: []*github.Commit{parentCommit},
	}
	newCommit, _, err := gc.ghGitClient.CreateCommit(ctx, owner, repo, commit)
	if err != nil {
		return err
	}

	ref.Object.SHA = newCommit.SHA
	_, _, err = gc.ghGitClient.UpdateRef(ctx, owner, repo, ref, false)
	if err != nil {
		return err
	}

	return nil
}

// CodeFilter is used to filter files based on a regex pattern.
type CodeFilter struct {
	FileRegexPattern string
}

func (gc *GithubClient) processEntry(entry *github.TreeEntry, codeFilter CodeFilter, context context.Context, owner, repo string, files map[string]string) error {
	if entry.GetType() != "blob" {
		return nil
	}

	re, err := regexp.Compile(codeFilter.FileRegexPattern)
	if err != nil {
		return fmt.Errorf("error compiling regex: %w", err)
	}
	if re.MatchString(entry.GetPath()) {
		blob, resp, err := gc.ghGitClient.GetBlob(context, owner, repo, entry.GetSHA())
		if err != nil {
			return fmt.Errorf("error getting blob: %w", err)
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("unexpected status code getting blob: %d", resp.StatusCode)
		}
		content, err := base64.StdEncoding.DecodeString(blob.GetContent())
		if err != nil {
			return fmt.Errorf("error base64 decoding blob string: %w", err)
		}
		// if first line has substring DO NOT EDIT then skip
		if IsDoNotEditFile(content) {
			return nil
		}

		files[entry.GetPath()] = string(content)
	}
	return nil
}

// GetCommitCode gets the code from a commit.
func (gc *GithubClient) GetCommitCode(context context.Context, owner, repo, commitSHA string, codeFilter CodeFilter) (map[string]string, error) {
	commit, _, err := gc.ghGitClient.GetCommit(context, owner, repo, commitSHA)
	if err != nil {
		return nil, err
	}

	treeSHA := commit.Tree.GetSHA()
	tree, _, err := gc.ghGitClient.GetTree(context, owner, repo, treeSHA, true)
	if err != nil {
		return nil, fmt.Errorf("error getting git tree: %w", err)
	}

	files := make(map[string]string)
	// temporary so we are not processing same files
	Shuffle(tree.Entries)
	for _, entry := range tree.Entries {
		if len(files) >= 5 {
			break // Stop processing once we have 5 entries. Temporary fix for context length limit reached
		}
		err := gc.processEntry(entry, codeFilter, context, owner, repo, files)
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

// AddCommentToPullRequestReview adds a comment to a pull request review.
//
// Parameters:
//   - ctx: The context for the API request.
//   - owner: The owner of the repository.
//   - repo: The repository name.
//   - prNumber: The pull request number.
//   - commentBody: The content of the comment to be added.
//   - commitID: The SHA of the commit to comment on.
//   - path: The file path in the repository where the comment should be added.
//   - position: The position in the diff where the comment should be added.
//
// Returns:
//   - The created PullRequestComment object.
//   - An error if any occurred during the API request.
func (gc *GithubClient) AddCommentToPullRequestReview(ctx context.Context, owner, repo string, prNumber int, commentBody, commitID, path string, startLine, line int) (*github.PullRequestComment, error) {
	comment := &github.PullRequestComment{
		// Text content of the comment
		Body: github.String(commentBody),
		// SHA of the commit to comment on
		CommitID: github.String(commitID),
		// Filepath which the comment applies
		Path: github.String(path),
		// Position in the diff where the comment should be applied
		// Position: github.Int(position),
		// First line of range you want to comment on
		StartLine: github.Int(startLine),
		// Last line of range you want to comment on
		Line: github.Int(line),
	}

	prComment, _, err := gc.ghPullRequestClient.CreateComment(ctx, owner, repo, prNumber, comment)
	if err != nil {
		return nil, fmt.Errorf("error creating comment: %w", err)
	}

	return prComment, nil
}

// CreatePullRequest creates a new pull request.
func (gc *GithubClient) CreatePullRequest(ctx context.Context, owner, repo, title, head, base, body string) (*github.PullRequest, error) {
	pr := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(head),
		Base:  github.String(base),
		Body:  github.String(body),
	}

	pullRequest, _, err := gc.ghPullRequestClient.Create(ctx, owner, repo, pr)
	if err != nil {
		return nil, fmt.Errorf("error creating pull request: %w", err)
	}

	return pullRequest, nil
}

// MergePullRequest merges a pull request.
func (gc *GithubClient) MergePullRequest(ctx context.Context, owner, repo string, prNumber int, commitMessage string) (*github.PullRequestMergeResult, error) {
	mergeResult, _, err := gc.ghPullRequestClient.Merge(ctx, owner, repo, prNumber, commitMessage, &github.PullRequestOptions{})
	if err != nil {
		return nil, fmt.Errorf("error merging pull request: %w", err)
	}

	return mergeResult, nil
}

// ListPullRequests lists all pull requests in a repository.
func (gc *GithubClient) ListPullRequests(ctx context.Context, owner, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, error) {
	pullRequests, _, err := gc.ghPullRequestClient.List(ctx, owner, repo, opts)
	if err != nil {
		return nil, fmt.Errorf("error listing pull requests: %w", err)
	}

	return pullRequests, nil
}

// GetPRCode gets the code from a pull request.
func (gc *GithubClient) GetPRCode(ctx context.Context, owner, repo string, prNumber int, opts *github.ListOptions) (map[string]string, error) {
	commitFiles, _, err := gc.ghPullRequestClient.ListFiles(ctx, owner, repo, prNumber, opts)
	if err != nil {
		return nil, fmt.Errorf("error listing PR files: %w", err)
	}

	files := make(map[string]string)
	for _, commitFile := range commitFiles {
		// TODO(mgottlieb): Use custom regex pattern instead of hardcoding to .go
		if strings.HasSuffix(commitFile.GetFilename(), ".go") {
			blob, resp, err := gc.ghGitClient.GetBlob(ctx, owner, repo, commitFile.GetSHA())
			if err != nil {
				return nil, fmt.Errorf("error getting blob: %w", err)
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return nil, fmt.Errorf("unexpected status code getting blob: %d", resp.StatusCode)
			}
			content, err := base64.StdEncoding.DecodeString(blob.GetContent())
			if err != nil {
				return nil, fmt.Errorf("error base64 decoding blob string: %w", err)
			}

			files[commitFile.GetFilename()] = string(content)
		}
	}

	return files, nil
}

func (gc *GithubClient) GetPRDiff(ctx context.Context, owner, repo string, prNumber int) (string, error) {
	diffOpts := &github.RawOptions{Type: github.Diff}
	diff, _, err := gc.ghPullRequestClient.GetRaw(ctx, owner, repo, prNumber, *diffOpts)
	if err != nil {
		return "", fmt.Errorf("error retrieving PR diff: %w", err)
	}

	return diff, nil
}

func (gc *GithubClient) GetPRPatch(ctx context.Context, owner, repo string, prNumber int) (string, error) {
	diffOpts := &github.RawOptions{Type: github.Patch}
	diff, _, err := gc.ghPullRequestClient.GetRaw(ctx, owner, repo, prNumber, *diffOpts)
	if err != nil {
		return "", fmt.Errorf("error retrieving PR diff: %w", err)
	}

	return diff, nil
}

// Shuffle shuffles a slice of any type
func Shuffle(slice interface{}) {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		panic("Shuffle: not a slice")
	}

	rand.Seed(time.Now().UnixNano())
	n := rv.Len()
	swap := reflect.Swapper(slice)

	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		swap(i, j)
	}
}

// IsDoNotEditFile checks if the file has a DO NOT EDIT comment.
func IsDoNotEditFile(data []byte) bool {
	lines := bytes.SplitN(data, []byte("\n"), 2)
	if len(lines) > 0 && bytes.Contains(lines[0], []byte("DO NOT EDIT")) {
		return true
	}
	return false
}
