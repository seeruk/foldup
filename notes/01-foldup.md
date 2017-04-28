# Foldup

So, what is Foldup? It's a backup tool for containers that uses cloud storage buckets. Docker + 
Bucket = Foldup. It's not exclusively for Docker though I guess?

The aim is for Foldup to be a flexible tool, able to run in multiple environments (locally / k8s, 
un-opinionated), that should make backing up volumes from containers to some kind of cloud storage 
bucket really easy.

Here's the general idea of how it would work:

1. Run inside of a container (of some kind) as a long-running process.
2. Volumes from other containers will be mounted into a specific directory (read only).
3. This directory will be scanned periodically (like on a cron-type schedule), and each folder 
(volume) will be compressed (likely to a `.tar.gz` file or something).
4. Each file will be uploaded to some cloud storage bucket (perhaps concurrently?)
5. After it's been uploaded, it'll be removed locally to free up space.

The cron-like syntax (I'm sure there's a cron-syntax parser for Go somewhere...) will be passed in 
via an environment variable (hooray for support in `eidolon/console`).
 
The configuration for the cloud service will need to be passed in too, that will more than likely 
also be in the form of environment variables.

All in all, it should mean that it's as simple as mounting some volumes (read only) onto a container
and leaving it to do it's thing.

## Other Thoughts

Object eviction could be handled with Foldup, but it might make more sense for that to just have 
been set up beforehand anyway. It would complicate the interaction with cloud providers as they 
would all handle it differently. Let's leave that be for now...
