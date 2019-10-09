# todo2

This simple script makes us of Make's directed acyclic graph solving, to express TODO lists in a way that allows multiple hierarchies of parent/child relationships between different atomic tasks.

For example, here's a simple list.

```
# Birthday Party TODO
[ ] eat cake
[ ] blow out candles
[ ] give presents
```

Now traditionally, the other steps involved in making the above happen would be listed; perhaps even heirarchally:

```
# Birthday Party Prep
[ ] Get cake mix
  [ ] Bake cake
    [ ] Eat cake
[ ] Buy candles
  [ ] Light candles
    [ ] Blow out candles 
[ ] Buy presents
  [ ] Wrap presents
    [ ] Give presents
```

Generally, that's enough for many tasks. But in the programming world, we develop hierarchies of things where multiple pieces, may depend on multiple, disperate pieces; a flat hierarchy doesn't suffice.

Now, many other methods use tagging to achieve this:

```
# Big project 
[ ] Fix bug #92930
  [ ] Ship update #3
[ ] Fix bug #01209
  [ ] Ship update #3
```

This works wonderfully well, but it is not robust enough for multiple parent, multiple child relationships that may arise in real life. If Fix #92930 depends on 3 other things, one of which depends on #01209 we wouldn't be able to express that #01209 blocks the parent dependencies of #92930 cleanly. 

Make solves these sort of dependency hierarchies cleanly, and since I didn't want to write a solver itself at this time, why not just use it! 

The scripts wrap the editing of a proper Makefile, though it is named config.todo to avoid name clobbering; which lives in the root of your project's source tree. Subsequent foo.todo files can be placed where you like them so long as they exist in the tree. 

## .todo files

A .todo file contains a checklist of items to complete an atomic task. For example, `polygon.todo` may contain the following:

```
BUG #92930 - Trying to seek returns incomprehensible error when calling SomeCall
[ ] polygon.go: MyType should implement ReaderAt
[ ] polygon.go: Check return value of SomeCall and try to recover
[ ] polygon.go: Hand back better error info to the calling function

BUG #012099 - Blob considered too small by users
[ ] blob.go: Make blob 2x larger when it feels threatened
[ ] blob.go: Implement room bounds checking
```

## Automatically generate .todo files

Now, no one wants to write all that out all the time. Shouldn't this be able to pull a generic `// TODO:` tag, or `# TODO:`, etc? Yeah, I thought so too.

`todo generate` will walk the source tree, looking `TODO`, `BUG`, `RELEASE`, and other related tags. If they have a tag, such as `BUG #92930: Trying to seek returns incomprehensible error when calling SomeCall`   an entry will be made in an appropriate .todo file in that particular source directory/subdirectory, which will be created as necessary.
 - files with extensions will have the extension replaced with .todo, and files without extensions will gain the extension .todo. Name collisions should be benign here.
 - more extensive coverage of tagging conventions can be Pull Requested as needed by users

## Managing dependancy hierarchy

To say a particular .todo depends on another, we can manually edit the Makefile, or we can use the wrapper scripts to simplify the process. 

`todo add` is the base command we'll be using, to add hierarchies. (A complementary `todo rm` can be used to undo anything `todo add` can do!)

To add something as a child dependency, `todo add child <parent name> <child name>`. From the above examples, you would issue `todo add child 'BUG #01209' 'BUG #92930'`, or equivilently `todo add parent 'BUG #92930' 'BUG #01209'`. 

Arbitrarily many dependencies can be added between arbitrary pieces within the current working source tree. 

## Commands of intrigue
 - `todo list` will output a list of all current leaves in the graph, (nodes which have no pending dependencies) 
 - `todo listall` will list every piece of the graph in a flat structure. This is very useful in defining hierarchies on complicated TODO list
 - `todo dot` can be piped to `graphviz` and friends to view your current task hierarchies, with the completed nodes elided
