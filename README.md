# todo2

This simple script makes us of directed acyclic graph solving, to express TODO lists in a way that allows multiple hierarchies of parent/child relationships between different atomic tasks.

Consider the following traditional TODO list.

```
# Birthday Party TODO
[ ] eat cake
[ ] blow out candles
[ ] give presents
```

Normally, the other steps involved in making the above happen would be listed; perhaps even heirarchally:

```
# Birthday Party TODO
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

## .todo files

A .todo file contains a checklist of items to complete an atomic task. For example, `myprogram/.todo` may contain the following:

```
#92930:
BUG - Trying to seek returns incomprehensible error when calling SomeCall
[ ] polygon.go: MyType should implement ReaderAt
[ ] polygon.go: Check return value of SomeCall and try to recover
[ ] polygon.go: Hand back better error info to the calling function

#012099:
BUG - Blob considered too small by users
[ ] blob.go: Make blob 2x larger when it feels threatened
[ ] blob.go: Implement room bounds checking
```

To add an entry into a particular task, such as the last line of BUG #92930 is trivial:
`todo task add 'BUG: [#92930]' "polygon.go: Hand back better error info to the calling function"`

You can also use stdin:

`some function | todo task add someentry`

If for some reason you wish to remove an entry, you can use the complement to `add`, `rm`.
`todo task rm 'BUG: [#92930]' <index|string>`

Since .todo files are simply text files, you can also achieve any of the above by opening and editing the file by hand.

The schema for an entry is non-complicated. Given a line in your code, such as: 
`BUG(halfwit) - There's an issue with my code and I don't understand why`
, an entry would be populated as follows.

```
general:
BUG(halfwit) - There's an issue with my code and I don't know why
```

Additionally, tags are used as a header. Multiple tags are fine as well!

```
// BUG(halfwit): [enhancement] [#354] There's an issue with my code and I don't know why
// requires [#353]
```

Would result with the following:

```
[enhancement] [#354]: requires [#353]
BUG(halfwit) - There's an issue with my code and I don't know why
```

A more interesting example generated like this would have related nodes

```
[#123]: requires [#353] [#124]
BUG(halfwit) - There's an issue with my code and I seem to have figured it out

[#124]:
TODO - write the first bit of code

[#354]: requires [#352]
TODO - Write the code

[#353]: requires [#352]
TODO - write other code

[#352]:
TODO - write the other code
```

## Completing tasks in .todo files

The following command modifies the state of a given entry's task:

 - `todo task toggle <entry> <index|string>`

For example, to toggle the completion of MyType implementing ReaderAt:

`todo task toggle '#92930' 'MyType should implement ReaderAt'`

You'll notice the file name is elided, for brevity a partial string match is all that is required. `todo task toggle '#92930' MyType` would suffice in our example.
By index, it's similar: `todo task toggle '#92930' 1` would update the .todo file to read `[x] polygon.go: MyType should implement ReaderAt`.

## Automatically generate .todo files

Now, no one wants to write all that out all the time. Shouldn't this be able to pull a generic `// TODO:` tag, or `# TODO:`, etc? Yeah, I thought so too.

`todo generate` will walk the source tree, looking `TODO`, `BUG`, `RELEASE`, and other related tags. If they have a tag, such as `BUG: [#92930] Trying to seek returns incomprehensible error when calling SomeCall`   an entry will be made in an appropriate .todo file in that particular source directory/subdirectory, which will be created as necessary.
 - files with extensions will have the extension replaced with .todo, and files without extensions will gain the extension .todo. Name collisions should be benign here.
 - more extensive coverage of tagging conventions can be Pull Requested as needed by users

## Managing dependancy hierarchy

`todo add` is the base command we'll be using, to add hierarchies. (A complementary `todo rm` can be used to undo anything `todo add` can do!)

To add something as a child dependency, `todo add child <parent name> <child name>`. From the above examples, you would issue `todo add child '#01209' '#92930'`, or equivilently `todo add parent '#92930' '#01209'`. 

Arbitrarily many dependencies can be added between arbitrary pieces within the current working source tree.

You end up with a 

## Commands of intrigue
 - `todo init` will create, if none exists, the backing Makefile-like file used internally
 - `todo list` will output a list of all current leaves in the graph, (nodes which have no pending dependencies) 
 - `todo listall` will list every piece of the graph in a flat structure. This is very useful in defining hierarchies on complicated TODO list
 - `todo dot` can be piped to `graphviz` and friends to view your current task hierarchies, with the completed nodes elided
 - `todo dotall` will output a complete graph of all task hierarchies

## Sharing
The best method is to include your existing todo2 main file into your current revision system, and symbolically link it into the `$XDG_DATA_HOME/todo2/myproject` directory. 
