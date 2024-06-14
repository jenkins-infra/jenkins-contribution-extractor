# Random notes

## Creating a Cobra skeleton

```sh
go mod init github.com/jenkins-infra/jenkins-get-commenters
cobra-cli init --author "Jean-Marc Meessen jean-marc@meessen-web.org" --license MIT
cobra-cli add get --author "Jean-Marc Meessen jean-marc@meessen-web.org" --license MIT
```

## hide flags in Cobra (debug)
- https://stackoverflow.com/questions/46591225/how-to-mark-some-global-persistent-flags-as-hidden-for-some-cobra-commands

## Rate limits handling
- https://docs.github.com/en/rest/guides/best-practices-for-using-the-rest-api?apiVersion=2022-11-28

```
Nbr of PR without comments: 446
Nbr of PR with comments:    638
Total comments:            2163
➜  jenkins-get-commenters git:(quota) ✗ ./jenkins-get-commenters quota                            
Limit: 5000 
Remaining 4068 
➜  jenkins-get-commenters git:(quota) ✗ 
```

https://umarcor.github.io/cobra/#generating-markdown-docs-for-your-own-cobracommand