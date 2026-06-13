From ClickhouseProof Require Import init.

Section proofs.
Context `{!heapGS Σ} {sem : go.Semantics}.
Local Set Default Proof Using "All".

Lemma wp_update_user_valid :
  test_fun_ok validate.testUpdateUserValid.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_update_user_missing_username :
  test_fun_ok validate.testUpdateUserMissingUsername.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_creation_statements_requires_one :
  test_fun_ok validate.testCreationStatementsRequiresOne.
Proof. clickhouse_semantics_auto. Qed.

End proofs.
