#!/usr/bin/env bash
set -e

./jenkins-get-commenters quota
./jenkins-get-commenters test-data/submissions-2023-08.csv -a
./jenkins-get-commenters quota
./jenkins-get-commenters test-data/submissions-2023-07.csv -a
./jenkins-get-commenters quota
./jenkins-get-commenters test-data/submissions-2023-06.csv -a
./jenkins-get-commenters quota
./jenkins-get-commenters test-data/submissions-2023-05.csv -a
./jenkins-get-commenters quota
./jenkins-get-commenters test-data/submissions-2023-04.csv -a
./jenkins-get-commenters quota
