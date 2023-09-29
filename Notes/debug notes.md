*** Debug mode enabled ***
PR/comment
- jenkinsci/docker/1711, dduportal, 2023-09-19 13:44:12 +0000 UTC, 1330155191, https://github.com/jenkinsci/docker/pull/1711#discussion_r1330155191
- jenkinsci/docker/1711, dduportal, 2023-09-23 17:26:43 +0000 UTC, 1335046154, https://github.com/jenkinsci/docker/pull/1711#discussion_r1335046154
- jenkinsci/docker/1711, dduportal, 2023-09-23 17:26:57 +0000 UTC, 1335046170, https://github.com/jenkinsci/docker/pull/1711#discussion_r1335046170
- jenkinsci/docker/1711, dduportal, 2023-09-23 17:27:04 +0000 UTC, 1335046196, https://github.com/jenkinsci/docker/pull/1711#discussion_r1335046196
- jenkinsci/docker/1711, dduportal, 2023-09-23 17:27:11 +0000 UTC, 1335046200, https://github.com/jenkinsci/docker/pull/1711#discussion_r1335046200
- jenkinsci/docker/1711, gounthar, 2023-09-23 17:54:26 +0000 UTC, 1335048928, https://github.com/jenkinsci/docker/pull/1711#discussion_r1335048928

ISSUE/comment
https://pkg.go.dev/github.com/google/go-github/v55@v55.0.0/github#IssuesService.ListComments
GitHub API docs: https://docs.github.com/en/rest/issues/comments#list-issue-comments GitHub API docs: https://docs.github.com/en/rest/issues/comments#list-issue-comments-for-a-repository

- jenkinsci/docker/1711, gounthar, 2023-09-19 17:09:21 +0000 UTC, https://github.com/jenkinsci/docker/pull/1711#issuecomment-1726106746, %!s(MISSING)
- jenkinsci/docker/1711, gounthar, 2023-09-24 16:49:30 +0000 UTC, https://github.com/jenkinsci/docker/pull/1711#issuecomment-1732618145, %!s(MISSING)
- jenkinsci/docker/1711, gounthar, 2023-09-24 17:03:13 +0000 UTC, https://github.com/jenkinsci/docker/pull/1711#issuecomment-1732620853, %!s(MISSING)
- jenkinsci/docker/1711, dduportal, 2023-09-24 20:03:57 +0000 UTC, https://github.com/jenkinsci/docker/pull/1711#issuecomment-1732657642, %!s(MISSING)

-----

*** Debug mode enabled ***
See "debug.log" for the trace

Processing "data/submissions-2023-06.csv"
   5% |██████████                                                                                                                                                                                                | (76/1385, 2 it/s) [47s:13m13s]panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x13a81ff]

goroutine 1 [running]:
github.com/jmMeessen/jenkins-get-commenters/cmd.load_reviewComments({0xc000284840, 0x9}, {0xc00028484a, 0x1c}, {0x14af07b, 0x2}, {0xc000190dc0, 0x9, 0x0?})
        /home/runner/work/jenkins-get-commenters/jenkins-get-commenters/cmd/get.go:264 +0x1bf
github.com/jmMeessen/jenkins-get-commenters/cmd.getCommenters({0xc000284840, 0x29}, 0x0?, 0x68?, {0x149c8c7, 0x1b})
        /home/runner/work/jenkins-get-commenters/jenkins-get-commenters/cmd/get.go:142 +0x4d9
github.com/jmMeessen/jenkins-get-commenters/cmd.performAction({0x7ff7bfeff527, 0x1c})
        /home/runner/work/jenkins-get-commenters/jenkins-get-commenters/cmd/root.go:280 +0x3bd
github.com/jmMeessen/jenkins-get-commenters/cmd.glob..func5(0xc00012e200?, {0xc000091740, 0x1, 0x1492750?})
        /home/runner/work/jenkins-get-commenters/jenkins-get-commenters/cmd/root.go:76 +0x21f
github.com/spf13/cobra.(*Command).execute(0x1884920, {0xc0000be050, 0x3, 0x3})
        /home/runner/go/pkg/mod/github.com/spf13/cobra@v1.7.0/command.go:944 +0x863
github.com/spf13/cobra.(*Command).ExecuteC(0x1884920)
        /home/runner/go/pkg/mod/github.com/spf13/cobra@v1.7.0/command.go:1068 +0x3a5
github.com/spf13/cobra.(*Command).Execute(...)
        /home/runner/go/pkg/mod/github.com/spf13/cobra@v1.7.0/command.go:992
github.com/jmMeessen/jenkins-get-commenters/cmd.Execute()
        /home/runner/work/jenkins-get-commenters/jenkins-get-commenters/cmd/root.go:88 +0x1a