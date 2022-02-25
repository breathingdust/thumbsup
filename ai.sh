#!/bin/bash
#

services=("comprehend" "codeguruprofiler" "codegurureviewer" "lexmodelbuildingservice" "forecastservice" "kendra" "rekognition" "personalize" "polly" "translate" "transcribeservice" "frauddetector")

for t in ${services[@]}; do
  go run main.go issues-by-service $t >> ai.csv
done