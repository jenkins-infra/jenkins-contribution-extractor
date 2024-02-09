## Releasing

- tag the main branch with `git tag v0.2.14 -m "V 0.2.14"`
- push the branch to GitHub `git push origin --tags`

Verify on github the release action that the delivery to HomeBrew worked.

### to retry releasing
- delete release on GitHub
- delete tag locally with `git tag -d v0.2.14`
- delete the tag on the remote with `git push --delete origin v0.2.14`