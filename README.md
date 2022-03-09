# Thumbsup

Helper command to pull useful Terraform AWS Provider specific data from the Github API. Output goes to STDOUT.

To avoid rate limiting, requires credentials set by environment variables:

```
GITHUB_USER
GITHUB_TOKEN
```

Note that even with a token set, you can get rate limited in a large repo.

## Query Types

There are several query types available.

### Service Stats

`thumbsup service-stats`

This query aggregates reaction types by service type, defined by any GitHub label starting with `service/`. You can optionally sort the results by supplying a sort parameter eg `thumbsup service-stats eyes`

### Top Core Issues

`thumbsup top-core-issues`

This is basically the same as the `service-stats` query, except it is limited to a set of service labels (currently hardcoded). You can optionally sort the results by supplying a sort parameter eg `thumbsup top-core-issues eyes`

### Aggregated Issue Reactions

`thumbsup aggregated-issue-reactions`

This query goes through all issues in the repo and examines the timeline for each issue looking for [cross-referenced](https://docs.github.com/en/developers/webhooks-and-events/events/issue-event-types#cross-referenced) events. The reactions on the issue are aggregated with the reactions on the cross referenced event. This query can be useful for identifying issue clusters. These are related issues which by themselves may not have a large amount of upvotes, but in aggregate do.

### Issue PullRequest Reactions

`thumbsup issue-pullrequest-reactions`

This query goes through all open pull requests, and aggregates reactions for the pull request, and any issue linked to it by a closes keyword (ie fixes, resolves, closes etc...)
