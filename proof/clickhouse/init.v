From New.proof Require Import slices_proof.slices_init.
From New.proof Require Import strings.
From New.proof Require Import errors.
From New.proof Require Import fmt.
From New.proof Require Import sort_proof.sort_init.
From New.proof Require Export proof_prelude.

Unset Printing Records.

Section wps.
Context `{hG: heapGS Σ, !ffi_semantics _ _}.
Context {sem : go.Semantics}.
Local Set Default Proof Using "All".

Definition test_fun_ok (name : go_string) :=
  ∀ Φ, Φ #true -∗ WP @! name #() {{ Φ }}.

End wps.

Ltac _cleanup :=
  repeat rewrite -> decide_True by (auto; word);
  repeat rewrite -> decide_False by (auto; word).

Ltac wp_call_auto :=
  first [ wp_func_call; wp_call
        | wp_method_call; wp_call
        | wp_call ].

Ltac clickhouse_slice_step :=
  first
    [ wp_apply wp_slice_literal; iSplitR; [done|]; iIntros "% [? _]"; wp_auto
    | wp_apply wp_slice_make2 as "%? [? ?]"; [word|]; wp_auto ].

Ltac clickhouse_steps := repeat (
  wp_call_auto ||
  wp_auto ||
  let x := fresh "x" in wp_alloc x as "?" ||
  clickhouse_slice_step ||
  _cleanup ||
  solve [ iExactEq "HΦ"; done ]
).

Ltac clickhouse_semantics_auto :=
  iIntros (Φ) "HΦ";
  wp_func_call; wp_call;
  clickhouse_steps;
  try solve [ iExactEq "HΦ"; done ].
