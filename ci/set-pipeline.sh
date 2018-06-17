#!/bin/sh

fly -t wings set-pipeline -p boshspecs -c pipeline.yml
