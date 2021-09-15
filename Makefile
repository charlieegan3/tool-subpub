test-watch:
	find . | grep go | entr bash -c 'clear; go test ./...'
