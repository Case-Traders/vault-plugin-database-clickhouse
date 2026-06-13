From ClickhouseProof Require Import init.
From New.generatedproof.vault_plugin_database_clickhouse.internal Require Import vars.

Section proofs.
Context `{hG: heapGS Σ, !ffi_semantics _ _}.
Context {sem : go.Semantics} {package_sem : vars.vars.Assumptions}.
#[global] Instance : IsPkgInit (iProp Σ) vars.pkg_id.vars := define_is_pkg_init True%I.
#[global] Instance : GetIsPkgInitWf (iProp Σ) vars.pkg_id.vars := build_get_is_pkg_init_wf.
Local Set Default Proof Using "All".

Lemma wp_new_user_has_required_keys :
  test_fun_ok vars.vars.testNewUserHasRequiredKeys.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_all_ops_have_required_keys :
  test_fun_ok vars.vars.testAllOpsHaveRequiredKeys.
Proof. clickhouse_semantics_auto. Qed.

End proofs.
