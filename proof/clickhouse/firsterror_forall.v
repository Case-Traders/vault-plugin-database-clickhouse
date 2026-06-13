From ClickhouseProof Require Import init.
From New.generatedproof.vault_plugin_database_clickhouse.internal Require Import firsterror.

Section proofs.
Context `{hG: heapGS Σ, !ffi_semantics _ _}.
Context {sem : go.Semantics} {package_sem : firsterror.firsterror.Assumptions}.
#[global] Instance : IsPkgInit (iProp Σ) firsterror.pkg_id.firsterror := define_is_pkg_init True%I.
#[global] Instance : GetIsPkgInitWf (iProp Σ) firsterror.pkg_id.firsterror := build_get_is_pkg_init_wf.
Local Set Default Proof Using "All".

Lemma wp_first_error_stops :
  test_fun_ok firsterror.firsterror.testFirstErrorStops.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_first_error_all_success :
  test_fun_ok firsterror.firsterror.testFirstErrorAllSuccess.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_first_error_preserved :
  test_fun_ok firsterror.firsterror.testFirstErrorPreserved.
Proof. clickhouse_semantics_auto. Qed.

End proofs.
