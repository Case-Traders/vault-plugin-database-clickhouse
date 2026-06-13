From ClickhouseProof Require Import init.

Section proofs.
Context `{!heapGS Σ} {sem : go.Semantics}.
Local Set Default Proof Using "All".

Lemma wp_normalize_idempotent :
  test_fun_ok stmt.testNormalizeIdempotent.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_normalize_idempotent_empty :
  test_fun_ok stmt.testNormalizeIdempotentEmpty.
Proof. clickhouse_semantics_auto. Qed.

End proofs.
