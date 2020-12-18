DROP TABLE IF EXISTS play_history;

CREATE TABLE play_history (
    played_at datetime,
    track_name text,
    track_id varchar(22),
    artist_name text,
    artist_id text
);
