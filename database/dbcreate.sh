# #!/bin/bash

# Define the file paths and names
DATASET_1="2021-05.csv"
DATASET_2="2021-06.csv"
DATASET_3="2021-07.csv"
DATASET_4="726277c507ef4914b0aec3cbcfcbfafc_0.csv"
DB_NAME="database/hsk-city-bike-app.db"
ALL_JOURNEYS="all_journeys"

mkdir -p datasets

for JOURNEYS_DATASETS in $DATASET_1 $DATASET_2 $DATASET_3; do
# Delete CSV files of journeys if they exist & download their newest version
if [ -f "./datasets/$JOURNEYS_DATASETS" ]; then
  echo "$JOURNEYS_DATASETS removed"
  rm -vf "./datasets/$JOURNEYS_DATASETS"
fi
echo "Downloading $JOURNEYS_DATASETS ..."
wget https://dev.hsl.fi/citybikes/od-trips-2021/$JOURNEYS_DATASETS -P "./datasets/"
done
# Delete CSV file of stations it it exist & download its newest version
if [ -f "./datasets/$DATASET_4" ]; then
  rm -vf "./datasets/$DATASET_4"
fi
echo "Downloading $DATASET_4 ..."
wget -q https://opendata.arcgis.com/datasets/$DATASET_4 -P "./datasets"

# Loop over each dataset of journeys
for DATASET in $DATASET_1 $DATASET_2 $DATASET_3; do
  # Extract the year and month from the dataset filename
  YEAR=$(echo $DATASET | cut -d "-" -f 1)
  MONTH=$(echo $DATASET | cut -d "-" -f 2 | cut -d "." -f 1)

  # Create the table name
  TABLE_NAME="journeys_$YEAR$MONTH"

  # Import the dataset into a raw table
sqlite3 $DB_NAME <<EOF
CREATE TABLE raw_journeys$YEAR$MONTH (
  "Departure",
  "Return",
  "Departure station id" INTEGER,
  "Departure station name",
  "Return station id" INTEGER,
  "Return station name",
  "Covered distance (m)" INTEGER,
  "Duration (sec.)" INTEGER
);

.mode csv
.import ./datasets/$DATASET raw_journeys$YEAR$MONTH

CREATE TABLE $TABLE_NAME AS
SELECT
  row_number() OVER (ORDER BY "Departure") AS id,
  "Departure",
  "Return",
  "Departure station id",
  "Departure station name",
  "Return station id",
  "Return station name",
  "Covered distance (m)",
  "Duration (sec.)"
FROM raw_journeys$YEAR$MONTH
WHERE "Duration (sec.)" >= 10 AND "Covered distance (m)" >= 10
GROUP BY
  "Departure",
  "Return",
  "Departure station id",
  "Departure station name",
  "Return station id",
  "Return station name",
  "Covered distance (m)",
  "Duration (sec.)";

DELETE FROM $TABLE_NAME WHERE rowid = (SELECT max(rowid) FROM $TABLE_NAME);
EOF
echo "data from $DATASET imported to $TABLE_NAME table in ./$DB_NAME"
done

sqlite3 $DB_NAME <<EOF
.mode csv
CREATE TABLE $ALL_JOURNEYS (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "Departure",
  "Return",
  "Departure station id",
  "Departure station name",
  "Return station id",
  "Return station name",
  "Covered distance (m)",
  "Duration (sec.)"
);

INSERT INTO $ALL_JOURNEYS
SELECT
  row_number() OVER (ORDER BY "Departure") AS id,
  "Departure",
  "Return",
  "Departure station id",
  "Departure station name",
  "Return station id",
  "Return station name",
  "Covered distance (m)",
  "Duration (sec.)"
FROM (
  SELECT *
  FROM journeys_202105
  UNION ALL
  SELECT *
  FROM journeys_202106
  UNION ALL
  SELECT *
  FROM journeys_202107
) t;

SELECT "data from all journeys imported to $ALL_JOURNEYS in ./$DB_NAME";

CREATE TABLE stations (
  "FID" INTEGER,
  "ID" INTEGER,
  "Nimi",
  "Namn",
  "Name",
  "Osoite",
  "Adress",
  "Kaupunki",
  "Stad",
  "Operaattor",
  "Kapasiteet" INTEGER,
  "x",
  "y"
);
.import ./datasets/$DATASET_4 stations
DELETE FROM stations WHERE rowid = 1;

SELECT "data from $DATASET_4 imported to stations table in ./$DB_NAME";

ALTER TABLE stations ADD COLUMN JourneysFrom INTEGER;
ALTER TABLE stations ADD COLUMN JourneysTo INTEGER;
UPDATE stations
SET JourneysFrom = (
    SELECT COUNT(*)
    FROM $ALL_JOURNEYS
    WHERE $ALL_JOURNEYS."Departure station id" = stations.ID
),
JourneysTo = (
    SELECT COUNT(*)
    FROM $ALL_JOURNEYS
    WHERE $ALL_JOURNEYS."Return station id" = stations.ID
);

SELECT "data created JourneysFrom and JourneysTo at stations table in ./$DB_NAME";
.save $DB_NAME
EOF