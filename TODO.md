# TODO

## Multiuser

First priority is multiuser.  I think that's as simple as attaching user IDs to item lists like a bad imitation SQL.  Then you pass the user id (name?) on the CLI and it populates the registry with the associated yamls. (Load all the yamls and then see if the name matches).

This bullshit with name-based items is probably going to bite us in the ass though.  Do we need SQL sooner rather than later, with YAML as a way to initialize it?  Eh let's see how it goes. pack list is just names and yes/no so it really should be fine.

the main issue is as the db changes, you can't load old lists.  which is not *really* a problem at the moment. The way to fix this is each packing list eventually stores ALL of the item data, not just the name / pack status.  That explodes the YAML size but who cares.

Then when you load an old pack list, it's gotta be like read-only.  Once the registry gets out of sync you're stuck.

## Web

second priority is web interface.  Once we have a super simple auth setup we'll have a user ID and we can implement the extremely simple flow of creating lists and packing.  No item editing in the web interface to start.

## context / property split

Need to separate contexts and properties -- I still think it's ok for a context to have an automatic associated property for requirement checking though.

They are separate types, which is nice, but maybe they need to be split apart in the YAML struct? We just have a list of Properties, so the context list gets lost. We can still kick this can down the road though.
