From ClickhouseProof Require Import init.

Section proofs.
Context `{!heapGS Σ} {sem : go.Semantics}.
Local Set Default Proof Using "All".

Lemma wp_first_error_stops :
  test_fun_ok txexec.testFirstErrorStops.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_first_error_all_success :
  test_fun_ok txexec.testFirstErrorAllSuccess.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_first_error_preserved :
  test_fun_ok txexec.testFirstErrorPreserved.
Proof. clickhouse_semantics_auto. Qed.

End proofs.
