From ClickhouseProof Require Import init.

Section proofs.
Context `{!heapGS Σ} {sem : go.Semantics}.
Local Set Default Proof Using "All".

Lemma wp_new_user_has_required_keys :
  test_fun_ok vars.testNewUserHasRequiredKeys.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_all_ops_have_required_keys :
  test_fun_ok vars.testAllOpsHaveRequiredKeys.
Proof. clickhouse_semantics_auto. Qed.

End proofs.
