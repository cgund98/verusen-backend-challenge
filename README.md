# Photolib

For this challenge, I created a simple photo organization application. Users can submit photos and save them to various albums.

## Design

### Entities

This application has only two primary entities, `photo` and `album`. An photo can belong to multiple albums.

**photo**

- `photoId`: _string_, unique identifier (Primary Key).
- `creationTime`: _string_, ISO 8601 encoded time at which the entity was created.
- `data`: _binary_, the image file itself.
- `name`: _string || null_, optional description of the photo itself.

**album**

- `albumId`: _string_, unique identifier (Primary Key).
- `creationTime`: _string_, ISO 8601 encoded time at which the entity was created.
- `name`: _string_, name of the photo album.

**album_membership**

The album membership entity is used to reflect which photos belong to which album (many-to-many relationship).

- `membershipId`: _string_, unique identifier (Primary Key).
- `albumId`: _string_, album to which this membership refers to (Foreign Key).
- `photoId`: _string_, photo to which this membership refers to (Foreign Key).

### Endpoints

Each entity has its corresponding CRUD API endpoints.

**Photos**

```
GET     /photos             - List photos
POST    /photos             - Upload a new photo
GET     /photos/<photoId>   - Get an individual photo
PUT     /photos/<photoId>   - Update an individual photo
DELETE  /photos/<photoId>   - Delete an individual photo
```

**Albums**

```
GET     /albums             - List albums
POST    /albums             - Create a new album
GET     /albums/<albumId>   - Get an individual album
PUT     /albums/<albumId>   - Update an individual album, including photo memberships
DELETE  /albums/<albumId>   - Delete an individual album
```

### Tech stack

**Go**: Go is my language of choice for API development.

**PostgreSQL**: PostgreSQL is the storage backend for this application. It's a powerful, general-purpose database that can handle pretty much anything you throw at it.

**Docker**: The application is built and deployed via containers. Containers are preferred for their ease of deployment and platform-independence.

## Running the application

Building and deploying the application is quite simple. The only dependencies are `docker` and `docker compose`.

```bash
# Build the application
docker compose build

# Start the application
docker compose up
```

## Extras
