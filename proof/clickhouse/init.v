From New.generatedproof.vault_plugin_database_clickhouse.internal Require Import stmt.
From New.generatedproof.vault_plugin_database_clickhouse.internal.cluster Require Import choose.
From New.generatedproof.vault_plugin_database_clickhouse.internal Require Import txexec.
From New.generatedproof.vault_plugin_database_clickhouse.internal Require Import vars.
From New.generatedproof.vault_plugin_database_clickhouse.internal Require Import stmts.
From New.generatedproof.vault_plugin_database_clickhouse.internal Require Import validate.
From New.generatedproof.vault_plugin_database_clickhouse.internal Require Import deletepath.
From New.proof Require Import proof_prelude strings.

Unset Printing Records.

Section wps.
Context `{!heapGS Σ}.
Context {sem : go.Semantics}.
Local Set Default Proof Using "All".

Definition test_fun_ok (name : go_string) :=
  ∀ Φ, Φ #true -∗ WP @! name #() {{ Φ }}.

#[global] Instance stmt_pkg_init : IsPkgInit stmt := define_is_pkg_init True%I.
#[global] Instance choose_pkg_init : IsPkgInit choose := define_is_pkg_init True%I.
#[global] Instance txexec_pkg_init : IsPkgInit txexec := define_is_pkg_init True%I.
#[global] Instance vars_pkg_init : IsPkgInit vars := define_is_pkg_init True%I.
#[global] Instance stmts_pkg_init : IsPkgInit stmts := define_is_pkg_init True%I.
#[global] Instance validate_pkg_init : IsPkgInit validate := define_is_pkg_init True%I.
#[global] Instance deletepath_pkg_init : IsPkgInit deletepath := define_is_pkg_init True%I.

End wps.

Ltac clickhouse_semantics_auto :=
  iIntros (Φ) "HΦ";
  wp_func_call; wp_call;
  repeat (wp_call || wp_auto || solve [ iExactEq "HΦ"; done ]).
