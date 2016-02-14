#!/bin/bash
set -o errexit
set -o pipefail
set -o nounset

readonly CWD=$(pwd)
readonly SAMPLES_DIR="$CWD/samples"
readonly DB_USER=${DB_USER:-postgres}
readonly DB_NAME="integration_test"
readonly DB_SCHEMA="public"
readonly GITHUB_SAMPLE="$SAMPLES_DIR/2015-01-01-15.json"

function download_github_samples() {
    if [ ! -f "$GITHUB_SAMPLE" ]; then
        mkdir -p $SAMPLES_DIR
        cd $SAMPLES_DIR
        wget -nc http://data.githubarchive.org/2015-01-01-15.json.gz && gunzip -f 2015-01-01-15.json.gz
        cd $CWD
    fi
}

function recreate_db() {
  psql -U ${DB_USER} -c "drop database if exists ${DB_NAME};"
  psql -U ${DB_USER} -c "create database ${DB_NAME};"
}

function import_json() {
    local table=$1
    local filename=$2
    pgfutter --table "$table" --schema $DB_SCHEMA --db $DB_NAME --user $DB_USER json "$filename"
    if [ $? -ne 0 ]; then
        echo "pgfutter could not import $filename"
        exit 300
    else
        echo "Imported $filename into $table"
    fi
}

function import_github_samples() {
    download_github_samples
    import_json "github_events" "$GITHUB_SAMPLE"
}

function test_csv_export() {
    local query="SELECT * FROM public.github_events"
    local filename="github_events.csv"
    pgreport --db $DB_NAME --user $DB_USER csv "$query"
    echo "Exported CVV to $filename"
}

function main() {
    recreate_db
    import_github_samples
    test_csv_export
}

main
