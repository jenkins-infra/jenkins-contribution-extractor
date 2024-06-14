# notes on the CLI

## Verbs

- **none**
    - used to read a previously generated list of PRs with the first three fields being org, project, PR_nbr
        - first (and required) parameter is the input filename
    - define output file `-o or --output`
    - append or overwrite output
    - define env_var containing the Github token
    - generated CSV with or without header.
    - verbose processing

- **get**
    - retrieves the comment count for a given org - project - PR combination
        - required params:
            - `--org` : organization
            - `--project` : project
            - `--pr` : Pull Request Reference
            - `--full_ref` : combination of above as "org/project/pr"

## New interface

### **`jenkins-contribution-extractor`**
- **root** (displays help)
    - **version**
        - -d : detailed
    - **quota**
    - **get**
        - **commenters**
            -- debug
            - **for_pr** \[pr_spec\]
                --debug
        - **pr**
            org
            mois
    - **test**