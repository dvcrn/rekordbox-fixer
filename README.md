# rekordbox-fixer

This is a hacked together script to deal with rekordbox relocation limitations. 

Mainly, that rekordbox can't re-link files that have a different extension, so with iTunes/Apple Music, when a song is matched with the cloud, Apple is giving you a new .m4a file instead of whatever you had if it can upgrade your song for you. 

Because of this, when migrating to a new computer and restoring from the Cloud, it can happen that some songs that used to be .mp3 are now .m4a, and rekordbox refuses to re-link the file. 

This script generates a new `rekordbox.xml` file with the filepathes fixed to their .m4a (or anything else) version which you can import into rekordbox to one-by-one fix your files. 

## How to use? 

Uhm... yeah this is not in a 'ready to use' state. You'll have to manually go into the sourcecode and edit the constants at the top of `main.go`, but if you're willing to try: 

1. Clone this repo, update the constants in `main.go` to match your `rekordbox.xml` location, and directory where to search for music
2. Run `go run main.go` which will generate a new rekordbox.xml for you
3. Go into rekordbox, open the settings->advanced and select the new `rekordbox.xml` at the "rekordbox xml" section
4. Now in the sidebar you should see a "rekordbox xml" section that contains all the playlists and songs that were remapped
5. Drag&Drop them to their respective playlist and delete the original 

### Limitations

- Pioneer in their wisdom decided that MyTags aren't exported in the XML so all hope of just automatically replacing all existing files is out the window
- Since the new files might still be a little bit different than their previous ones, the grid might no longer match. Make sure you fix that while re-importing the files

## Other notes

Come on pioneer, it's ridiculous that I have to write a program to fix cases like this. Give rekordbox some love 🤎