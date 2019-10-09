
# Add foo as a parent node to bar and baz
foo.done: bar.todo baz.todo
	grep '( )' $@ && return 1 || touch foo.done

clean:
	rm -f *.done
