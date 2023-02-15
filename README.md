# Mazda CMU Updater

All features are written as library first and then the CLI is just convenience.
The CLI will lag behind and potentially never receive all features, since it's
not the primary focus of this work and I'd rather focus my time on analyzing
firmware than polishing the utility.

# Features

Library
* Parse .ini files and make some sense
* Extract files according to files.ini

# TODO

*Library*

* [ ] Unpack/copy binary.ini files
* [ ] Functions for special steps (rootfs unpack)
    * Simulate shellscript based on research?
    * Execute in isolated environment (Docker?) and then move files?
    * Alias cp/dd etc. and log those steps/simulate.
* [ ] Check hashes where provided
* [ ] Repacking

*CLI*

* [ ] Extract files
