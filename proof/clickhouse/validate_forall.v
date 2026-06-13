From ClickhouseProof Require Import init.
From New.generatedproof.vault_plugin_database_clickhouse.internal Require Import validate.

Section proofs.
Context `{hG: heapGS Σ, !ffi_semantics _ _}.
Context {sem : go.Semantics} {package_sem : validate.validate.Assumptions}.
#[global] Instance : IsPkgInit (iProp Σ) validate.pkg_id.validate := define_is_pkg_init True%I.
#[global] Instance : GetIsPkgInitWf (iProp Σ) validate.pkg_id.validate := build_get_is_pkg_init_wf.
Local Set Default Proof Using "All".

Lemma wp_update_user_valid :
  test_fun_ok validate.validate.testUpdateUserValid.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_update_user_missing_username :
  test_fun_ok validate.validate.testUpdateUserMissingUsername.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_creation_statements_requires_one :
  test_fun_ok validate.validate.testCreationStatementsRequiresOne.
Proof. clickhouse_semantics_auto. Qed.

End proofs.
