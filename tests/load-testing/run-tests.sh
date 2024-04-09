#!/bin/bash

k6 run user-banner-test.js &
k6 run get-filtered-banner.js