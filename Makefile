#
# Makefile for generating the watch-cli run program
#

all: 
	go build

install:  watch-cli
	cp watch-cli ~/bin

watch-cli: main.go util_bsd.go util_linux.go util_windows.go
	go build 


run-it:
	go build
	./watch-cli -c "echo 'hello silly world'" test/x.x test/x.y >,log 2>&1 &

run-fg:
	go build
	./watch-cli -c "echo 'hello silly world'" test/x.x test/x.y

# From: http://blog.jgc.org/2015/04/the-one-line-you-should-add-to-every.html
print-%: ; @echo $*=$($*)

#// BuildNo: 000

upd_BuildNo:
	./updBuildNo.sh note.1
	( cd . ; make )
	git commit -m "Set Build No on note.1" .
	git push origin master

test1:
	rm ~/bin/watch-cli
	go build
	cp watch-cli ~/bin
	( cd ~/Projects/pongo2/lexie/cli/ ; ~/bin/watch-cli -c "make www_lexie_com" site/www_lexie_com/index.html & )


