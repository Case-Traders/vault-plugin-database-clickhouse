From ClickhouseProof Require Import init.

Section proofs.
Context `{!heapGS Σ} {sem : go.Semantics}.
Local Set Default Proof Using "All".

Lemma wp_choose_configured :
  test_fun_ok choose.testChooseConfigured.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_choose_single_discovery :
  test_fun_ok choose.testChooseSingleDiscovery.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_choose_empty_error :
  test_fun_ok choose.testChooseEmptyError.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_choose_ambiguous :
  test_fun_ok choose.testChooseAmbiguous.
Proof. clickhouse_semantics_auto. Qed.

End proofs.
