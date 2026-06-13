From ClickhouseProof Require Import init.

Section proofs.
Context `{!heapGS Σ} {sem : go.Semantics}.
Local Set Default Proof Using "All".

Lemma wp_use_custom_revocation_false_for_empty :
  test_fun_ok deletepath.testUseCustomRevocationFalseForEmpty.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_use_custom_revocation_true_when_provided :
  test_fun_ok deletepath.testUseCustomRevocationTrueWhenProvided.
Proof. clickhouse_semantics_auto. Qed.

End proofs.
