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
readonly MONTGOMERY_SALARIES_SAMPLE="$SAMPLES_DIR/employee_salaries.csv"

function download_montgomery_county_samples() {
    if [ ! -f "$MONTGOMERY_SALARIES_SAMPLE" ]; then
        wget -O "$MONTGOMERY_SALARIES_SAMPLE" https://data.montgomerycountymd.gov/api/views/54rh-89p8/rows.csv
    fi
}

function download_github_samples() {
    if [ ! -f "$GITHUB_SAMPLE" ]; then
        mkdir -p $SAMPLES_DIR
        cd $SAMPLES_DIR
        wget http://data.githubarchive.org/2015-01-01-15.json.gz && gunzip -f 2015-01-01-15.json.gz
        cd $CWD
    fi
}

function recreate_db() {
  psql -U ${DB_USER} -c "drop database if exists ${DB_NAME};"
  psql -U ${DB_USER} -c "create database ${DB_NAME};"
}

function import_csv() {
    local table=$1
    local filename=$2
    pgfutter --table "$table" --schema $DB_SCHEMA --db $DB_NAME --user $DB_USER csv "$filename"
    if [ $? -ne 0 ]; then
        echo "pgfutter could not import $filename"
        exit 300
    else
        echo "Imported $filename into $table"
    fi
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

function import_montgomery_county_samples() {
    download_montgomery_county_samples
    import_csv "employee_salaries" "$MONTGOMERY_SALARIES_SAMPLE"
}

function import_github_samples() {
    download_github_samples
    import_json "github_events" "$GITHUB_SAMPLE"
}

function test_json_lines_export() {
    local query="SELECT e.data->'repo'->>'name' as name, json_agg(c->>'sha') as commmits FROM github_events AS e, json_array_elements(e.data->'payload'->'commits') AS c WHERE e.data->>'type' = 'PushEvent' GROUP BY e.data->'repo'->>'name'"
    local filename="push_events.json"
    pgclimb -d $DB_NAME -U $DB_USER -c "$query" -o "$filename" jsonlines
    echo "Exported JSON lines to $filename"
}

function test_json_doc_export {
    local query="SELECT e.data FROM github_events e WHERE e.data->>'type' = 'PushEvent'"
    local filename="push_event_docs.json"
    pgclimb --dbname $DB_NAME --username $DB_USER --command "$query" -o "$filename" json
    echo "Exported JSON to $filename"

}

function test_csv_export() {
local query="SELECT position_title, COUNT(*) AS employees, round(AVG(replace(current_annual_salary, '$', '')::numeric)) AS avg_salary FROM employee_salaries GROUP BY position_title ORDER BY 3 DESC"
    local filename="montgomery_average_salaries.csv"
    echo "$query" | pgclimb -d $DB_NAME -U $DB_USER -o "$filename" csv
    echo "Exported CSV to $filename"
}

function main() {
    recreate_db
    import_github_samples
    import_montgomery_county_samples
    test_csv_export
    test_json_lines_export
    test_json_doc_export
}

main
