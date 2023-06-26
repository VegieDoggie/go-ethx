SHELL=/bin/bash
all:
	@make help

.PHONY: help ## Display help commands
help:
	@printf 'Usage:\n'
	@printf '  make <tagert>\n'
	@printf '\n'
	@printf 'Targets:\n'
	@IFS=$$'\n' ; \
    help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//'`); \
    for help_line in $${help_lines[@]}; do \
        IFS=$$'#' ; \
        help_split=($$help_line) ; \
        help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		IFS=$$':' ; \
		phony_command=($$help_split); \
        help_command=`echo $${phony_command[1]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf "  %-50s %s\n" $$help_command $$help_info ; \
    done


.PHONY: verify-tidy ##verify and download  the Dependencies
verify-tidy:
	go mod verify && go mod tidy


golangci_lint_cmd=go run github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: lint-fix  ## fix static syntax detection error
lint-fix:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	$(golangci_lint_cmd) run --fix --out-format=tab --issues-exit-code=0

.PHONY: lint-go  ## Static syntax detection
lint-go:
	go get  github.com/golangci/golangci-lint/cmd/golangci-lint
	echo $(GIT_DIFF)
	$(golangci_lint_cmd) run --out-format=tab $(GIT_DIFF)
