From ClickhouseProof Require Import init.
From New.generatedproof.vault_plugin_database_clickhouse.internal Require Import deletepath.

Section proofs.
Context `{hG: heapGS Σ, !ffi_semantics _ _}.
Context {sem : go.Semantics} {package_sem : deletepath.deletepath.Assumptions}.
#[global] Instance : IsPkgInit (iProp Σ) deletepath.pkg_id.deletepath := define_is_pkg_init True%I.
#[global] Instance : GetIsPkgInitWf (iProp Σ) deletepath.pkg_id.deletepath := build_get_is_pkg_init_wf.
Local Set Default Proof Using "All".

Lemma wp_use_custom_revocation_false_for_empty :
  test_fun_ok deletepath.deletepath.testUseCustomRevocationFalseForEmpty.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_use_custom_revocation_true_when_provided :
  test_fun_ok deletepath.deletepath.testUseCustomRevocationTrueWhenProvided.
Proof. clickhouse_semantics_auto. Qed.

End proofs.
