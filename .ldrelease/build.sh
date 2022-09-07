#!/bin/bash

set -ue

helm lint .
helm package .
