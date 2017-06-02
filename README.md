# Foldup

<p>
    <a href="https://travis-ci.org/SeerUK/foldup">
        <img src="https://api.travis-ci.org/SeerUK/foldup.svg?branch=master" />
    </a>
    <a href="https://goreportcard.com/report/github.com/SeerUK/foldup">
        <img src="https://goreportcard.com/badge/github.com/SeerUK/foldup" />
    </a>
    <a href="https://github.com/SeerUK/foldup/releases">
        <img src="https://img.shields.io/github/release/SeerUK/foldup.svg" />
    </a>
</p>

Backup folders as archives to durable cloud storage buckets.

## Usage

As a user, it's probably easiest to just use the Docker image:

```
docker run --rm \
    -v /path/to/backup:/backup \
    -v /path/to/google/creds.json:/root/creds.json \
    -e GOOGLE_APPLICATION_CREDENTIALS=/root/creds.json \ 
    seeruk/foldup \
    backup /backup \ 
        --bucket=backups-sierra
        --schedule="0 * * * *"
```

You _must_ provide some [credentials for GCS][1]. The above example uses a service account key, 
mounting it into the container, then using the environment variable `GOOGLE_APPLICATION_CREDENTIALS`
to specify the location of that key file.

The folder to backup is specified in the command. You could use it as it is shown above, or you 
could mount individual directories to a parent folder, like this:

```
docker run --rm \
    -v /1st/to/backup:/backup/1st \
    -v /2nd/to/backup:/backup/2nd \
    -v /3rd/to/backup:/backup/3rd \
    ...
    backup /backup \ 
    ...
```

Foldup backs up all folders in a folder, so you can use this approach easily in a Docker Compose 
environment, or any other environment where you may have several volumes to back up. This can be 
especially useful for backing up the built-in Docker volumes.

## Todo

* Encrypted backups (maybe)
* Listing backed up folders
* Restoring backed up folders

## License

MIT

[1]: https://developers.google.com/identity/protocols/application-default-credentials
