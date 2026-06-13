.PHONY: proof

proof:
	@cd proof && coq_makefile -f _CoqProject -o Makefile.coq && $(MAKE) -f Makefile.coq clean && $(MAKE) -f Makefile.coq
