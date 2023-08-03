CREATE TABLE IF NOT EXISTS photos (
    photoId VARCHAR (255) PRIMARY KEY,
    creationTime TIMESTAMP NOT NULL,
    data BYTEA NOT NULL,
    name VARCHAR (255)
);

CREATE TABLE IF NOT EXISTS albums (
    albumId VARCHAR (255) PRIMARY KEY,
    creationTime TIMESTAMP NOT NULL,
    name VARCHAR (255) NOT NULL
);

CREATE TABLE IF NOT EXISTS album_memberships (
    membership_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    albumId VARCHAR (255),
    photoId VARCHAR (255),

    CONSTRAINT fk_album FOREIGN KEY(albumId) REFERENCES albums(albumId),
    CONSTRAINT fk_photo FOREIGN KEY(photoId) REFERENCES photos(photoId)
);