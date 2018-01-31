package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"sort"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
)

var (
	repo           string
	startStr       string
	endStr         string
	owner          string
	integrationId  int
	installationId int
	keyPath        string
)

func main() {
	flag.StringVar(&startStr, "start", "", "start date in format YYYY-MM-DD")
	flag.StringVar(&endStr, "end", "", "end date in format YYYY-MM-DD")
	flag.StringVar(&owner, "owner", "skycoin", "repo owner. `skycoin` by default")
	flag.StringVar(&repo, "repo", "skycoin", "repo name. `skycoin` by default")
	flag.IntVar(&integrationId, "integrationID", 0, "github app integrationID")
	flag.IntVar(&installationId, "installationID", 0, "github app installationID")
	flag.StringVar(&keyPath, "key", "", "path to the api key *.pem")
	flag.Parse()
	var start, end time.Time
	var err error
	ctx := context.Background()

	if start, err = time.Parse("2006-01-02", startStr); err != nil {
		fmt.Printf("Failed parse `start` param. err: %s\n", err)
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
			log.Printf("Failed create http Transport. err: %v", err)
			return
		}
		httpClient = &http.Client{
			Transport: itr,
		}
	}

	client := github.NewClient(httpClient)
	getClosedIssues(ctx, client, start, end)
}

func getClosedIssues(ctx context.Context, client *github.Client, start, end time.Time) {
	var page = 0
	var perPage = 100
	var seen = make(map[int]bool, 0)
	var issues = issuesSlice{}
	var pullRequests = issuesSlice{}
	for {
		options := github.IssueListByRepoOptions{
			Sort:  "updated",
			State: "closed",
			Since: start,
			Direction: "asc",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
		}
		batchIssues, _, err := client.Issues.ListByRepo(ctx, owner, repo, &options)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				log.Printf("hit rate limit.\n message: %s", err.Error())
			} else {
				log.Printf("failed get closed issues list. err: %v", err)
			}
		}
		if len(batchIssues) == 0 {
			break
		}
		for _, issue := range batchIssues {
			if !seen[*issue.Number] && (issue.ClosedAt.Equal(start) || issue.ClosedAt.After(start)) && (issue.ClosedAt.Equal(end) || issue.ClosedAt.Before(end)) {
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
	fmt.Println("CLOSED PULL REQUESTS")
	fmt.Println("--------------------")
	pullRequests.report()
	fmt.Println("\n\nCLOSED ISSUES")
	fmt.Println("-------------")
	issues.report()
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
