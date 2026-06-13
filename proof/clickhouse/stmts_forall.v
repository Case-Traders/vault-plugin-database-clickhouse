From ClickhouseProof Require Import init.

Section proofs.
Context `{!heapGS Σ} {sem : go.Semantics}.
Local Set Default Proof Using "All".

Lemma wp_statements_or_default_fallback :
  test_fun_ok stmts.testStatementsOrDefaultFallback.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_statements_or_default_provided :
  test_fun_ok stmts.testStatementsOrDefaultProvided.
Proof. clickhouse_semantics_auto. Qed.

End proofs.
