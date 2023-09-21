# notes on the CLI

## Verbs

- **none**
    - used to read a previously generated list of PRs with the first three fields being org, project, PR_nbr
        - first (and required) parameter is the input filename
    - define output file `-o or --output`
    - define env_var containing the Github token
    - generated CSV with or without header.

- **get**
    - retrieves the comment count for a given org - project - PR combination
        - required params:
            - `--org` : organisation
            - `--project` : project
            - `--pr` : Pull Request Reference
            - `--full_ref` : combination of above as "org/project/pr"
