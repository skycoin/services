package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"sort"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
)

var (
	startStr       string
	endStr         string
	repo           string
	owner          string
	org            string
	integrationId  int
	installationId int
	keyPath        string
)

func main() {
	flag.StringVar(&startStr, "start", "", "start date in format YYYY-MM-DD")
	flag.StringVar(&endStr, "end", "", "end date in format YYYY-MM-DD")
	flag.StringVar(&owner, "owner", "skycoin", "repo owner. `skycoin` by default")
	flag.StringVar(&repo, "repo", "skycoin", "repo name. `skycoin` by default")
	flag.StringVar(&org, "org", "", "org name. empty by default")
	flag.IntVar(&integrationId, "integrationID", 0, "github app integrationID")
	flag.IntVar(&installationId, "installationID", 0, "github app installationID")
	flag.StringVar(&keyPath, "key", "", "path to the api key *.pem")
	flag.Parse()
	var start, end time.Time
	var err error
	ctx := context.Background()
	if startStr != "" {
		if start, err = time.Parse("2006-01-02", startStr); err != nil {
			fmt.Printf("Failed parse `start` param. err: %s\n", err)
			flag.Usage()
			return
		}
	} else {
		flag.Usage()
		return
	}

	if endStr != "" {
		if end, err = time.Parse("2006-01-02", endStr); err != nil {
			fmt.Printf("Failed parse `end` param. err: %s\n", err)
			flag.Usage()
			return
		}
		end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second + 999*time.Millisecond)
	} else {
		end = time.Now()
	}

	var httpClient *http.Client
	if keyPath != "" {
		itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, integrationId, installationId, keyPath)
		if err != nil {
			fmt.Printf("Failed create http Transport. err: %v", err)
			return
		}
		httpClient = &http.Client{
			Transport: itr,
		}
	}

	client := github.NewClient(httpClient)
	c := newClientWrapper(ctx, client, start, end, owner, repo, org)
	if c.org != "" {
		c.getClosedIssuesByOrg()
	} else {
		pullRequests, issues := c.getClosedIssues()
		c.report(pullRequests, issues)
	}
}

type clientWrapper struct {
	ctx    context.Context
	client *github.Client
	start  time.Time
	end    time.Time
	owner  string
	repo   string
	org    string
}

func newClientWrapper(ctx context.Context, client *github.Client, start, end time.Time, owner, repo, org string) *clientWrapper {
	return &clientWrapper{
		ctx:    ctx,
		client: client,
		start:  start,
		end:    end,
		owner:  owner,
		repo:   repo,
		org:    org,
	}
}

func (c *clientWrapper) getClosedIssuesByOrg() {
	repos, _, err := c.client.Repositories.ListByOrg(c.ctx, c.org, &github.RepositoryListByOrgOptions{})
	if err != nil {
		fmt.Printf("failed get organisation. err: %s", err)
	}
	for _, repo := range repos {
		c.repo = *repo.Name
		c.owner = *repo.Owner.Login
		pullRequests, issues := c.getClosedIssues()
		c.report(pullRequests, issues)
	}
}

func (c *clientWrapper) getClosedIssues() (issuesSlice, issuesSlice) {
	var page = 0
	var perPage = 100
	var seen = make(map[int]bool, 0)
	var issues = issuesSlice{}
	var pullRequests = issuesSlice{}
	var batchIssues []*github.Issue
	var err error
	for {
		batchIssues, err = c.getBatch(page, perPage)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				fmt.Printf("hit rate limit.\n message: %s", err.Error())
			} else {
				fmt.Printf("failed get closed issues list. err: %v", err)
			}
			break
		}
		if len(batchIssues) == 0 {
			break
		}
		for _, issue := range batchIssues {
			if !seen[*issue.Number] &&
				(issue.ClosedAt.Equal(c.start) || issue.ClosedAt.After(c.start)) &&
				(issue.ClosedAt.Equal(c.end) || issue.ClosedAt.Before(c.end)) {
				seen[*issue.Number] = true
				if issue.IsPullRequest() {
					pullRequests.addIssue(issue)
				} else {
					issues.addIssue(issue)
				}
			}
		}
		page++
	}
	return pullRequests, issues
}

func (c *clientWrapper) report(pullRequests, issues issuesSlice) {
	if len(issues.issues) == 0 && len(pullRequests.issues) == 0 {
		return
	}
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")
	fmt.Printf("%s/%s\n", c.owner, c.repo)
	if len(pullRequests.issues) > 0 {
		fmt.Println("\nCLOSED PULL REQUESTS")
		fmt.Println("--------------------")
		pullRequests.report()
	}

	if len(issues.issues) > 0 {
		fmt.Println("\nCLOSED ISSUES")
		fmt.Println("-------------")
		issues.report()
	}
	fmt.Println("------------------------------------------------------------------------------------------------------------------------\n\n")
}

func (c *clientWrapper) getBatch(page, perPage int) ([]*github.Issue, error) {
	var batchIssues []*github.Issue
	var err error
	options := github.IssueListByRepoOptions{
		Sort:      "updated",
		State:     "closed",
		Since:     c.start,
		Direction: "asc",
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}
	batchIssues, _, err = c.client.Issues.ListByRepo(c.ctx, c.owner, c.repo, &options)
	return batchIssues, err
}

type issuesSlice struct {
	longestTitle int
	longestUrl   int
	longestLogin int
	issues       []github.Issue
}

func (p *issuesSlice) checkForReport(issue github.Issue) {
	if p.longestUrl < len(*issue.HTMLURL) {
		p.longestUrl = len(*issue.HTMLURL)
	}
	if p.longestTitle < len(*issue.Title) {
		p.longestTitle = len(*issue.Title)
	}
	if p.longestLogin < len(*issue.User.Login) {
		p.longestLogin = len(*issue.User.Login)
	}
}

func (p *issuesSlice) addIssue(issue *github.Issue) {
	p.issues = append(p.issues, *issue)
	p.checkForReport(*issue)
}

func (p *issuesSlice) report() {
	p.sort()
	template := fmt.Sprintf("|%%4d | %%%ds | %%%ds | %%%ds |\n", p.longestTitle, p.longestUrl, p.longestLogin)
	for _, issue := range p.issues {
		fmt.Printf(template, *issue.Number, *issue.Title, *issue.HTMLURL, *issue.User.Login)
	}
}

func (p *issuesSlice) sort() {
	sort.Sort(p)
}

func (p *issuesSlice) Len() int {
	return len(p.issues)
}

func (p *issuesSlice) Less(i, j int) bool {
	return p.issues[i].ClosedAt.Before(*p.issues[j].ClosedAt)
}

func (p *issuesSlice) Swap(i, j int) {
	p.issues[i], p.issues[j] = p.issues[j], p.issues[i]
}
