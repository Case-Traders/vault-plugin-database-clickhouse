From ClickhouseProof Require Import init.
From New.generatedproof.vault_plugin_database_clickhouse.internal Require Import stmt.

Section proofs.
Context `{hG: heapGS Σ, !ffi_semantics _ _}.
Context {sem : go.Semantics} {package_sem : stmt.stmt.Assumptions}.
#[global] Instance : IsPkgInit (iProp Σ) stmt.pkg_id.stmt := define_is_pkg_init True%I.
#[global] Instance : GetIsPkgInitWf (iProp Σ) stmt.pkg_id.stmt := build_get_is_pkg_init_wf.
Local Set Default Proof Using "All".

Lemma wp_normalize_idempotent :
  test_fun_ok stmt.stmt.testNormalizeIdempotent.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_normalize_idempotent_empty :
  test_fun_ok stmt.stmt.testNormalizeIdempotentEmpty.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_statements_or_default_fallback :
  test_fun_ok stmt.stmt.testStatementsOrDefaultFallback.
Proof. clickhouse_semantics_auto. Qed.

Lemma wp_statements_or_default_provided :
  test_fun_ok stmt.stmt.testStatementsOrDefaultProvided.
Proof. clickhouse_semantics_auto. Qed.

End proofs.
