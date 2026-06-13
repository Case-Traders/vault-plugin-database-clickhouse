From ClickhouseProof Require Import init.
From New.generatedproof.vault_plugin_database_clickhouse.internal.cluster Require Import choose.

Section proofs.
Context `{hG: heapGS Σ, !ffi_semantics _ _}.
Context {sem : go.Semantics} {package_sem : choose.choose.Assumptions}.
#[global] Instance : IsPkgInit (iProp Σ) choose.pkg_id.choose := define_is_pkg_init True%I.
#[global] Instance : GetIsPkgInitWf (iProp Σ) choose.pkg_id.choose := build_get_is_pkg_init_wf.
Local Set Default Proof Using "All".

Lemma wp_choose_configured :
  test_fun_ok choose.choose.testChooseConfigured.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_choose_single_discovery :
  test_fun_ok choose.choose.testChooseSingleDiscovery.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_choose_empty_error :
  test_fun_ok choose.choose.testChooseEmptyError.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_choose_ambiguous :
  test_fun_ok choose.choose.testChooseAmbiguous.
Proof. clickhouse_semantics_auto. Qed.

End proofs.
