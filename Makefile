PROJECT := "go-magicinfo"
USER := rptl8sr
EMAIL := $(USER)@gmail.com

.PHONY: git-init
git-init:
	gh repo create $(PROJECT) --private
	git init
	git config user.name "$(USER)"
	git config user.email "$(EMAIL)"
	git add go.mod Makefile
	git commit -m "Init commit"
	git remote add origin git@github.com:$(USER)/$(PROJECT).git
	git remote -v
	git push -u origin master


BN ?= dev
# make git-checkout BN=dev
.PHONY: git-checkout
git-checkout:
	git checkout -b $(BN)


M ?=
# make git-commit M="commit text"
.PHONY: git-commit
git-commit:
	if [ -z "$(M)" ]; then echo "MESSAGE is not set, aborting commit..."; exit 1; fi
	git commit -m "$(M)"


H ?= 1
# make git-reset H=1
.PHONY: git-reset
git-reset:
	git reset --soft HEAD~$(H)


.PHONY: git-push-head
# push to same branch as local
git-push-head:
	git push origin HEAD


T ?= "Pull request"
B ?= "-"
BRANCH ?= main
.PHONY: git-pr
# create pull (github)/merge (gitlab) request
git-pr:
	gh pr create --title "$(T)" --body "$(B)" --base $(BRANCH)


OS ?= darwin
NAME ?= $(PROJECT)
# linux
.PHONY: build
build:
	GOOS=$(OS) GOARCH=amd64 go build -o $(NAME) ./cmd


PORT ?= 5432
PATH ?= "c:\tmp\backup.sql"
USER ?= magicinfo
DBNAME ?= magicinfo
.PHONY: mi-backup
mi-backup:
	pg_dump.exe -p $(PORT) -f $(PATH) -c -U $(USER) -E UNICODE $(DBNAME)


.PHONY: mi-restore
mi-restore:
	# like this: psql.exe -p 5432 -U magicinfo -d magicinfo -f backup.sql
	psql.exe -p $(PORT) -U $(USER) -d $(DBNAME) -f $(PATH)

M ?= ""
.PHONY: change-url
change-url:
	go run ./cmd/main.go -a $(M)

P ?= "mac_address_example.txt"
.PHONY: change-url
change-url:
	go run ./cmd/main.go -p $(P)