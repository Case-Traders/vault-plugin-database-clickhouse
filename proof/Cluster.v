From Coq Require Import List String.
Import ListNotations.

Open Scope string_scope.

Inductive cluster_result :=
  | Inl : string -> cluster_result
  | Inr : string -> cluster_result.

Parameter empty_msg : string.
Parameter ambiguous_msg : string.

Fixpoint mem (x : string) (l : list string) : bool :=
  match l with
  | nil => false
  | cons h t => if String.eqb x h then true else mem x t
  end.

Fixpoint distinct_nonempty (discovered : list string) : list string :=
  match discovered with
  | nil => nil
  | cons h t =>
      let t' := distinct_nonempty t in
      if mem h t' then t' else h :: t'
  end.

Definition choose_cluster (configured : option string) (discovered : list string)
    : cluster_result :=
  match configured with
  | Some name => Inl name
  | None =>
      let names := distinct_nonempty discovered in
      match names with
      | nil => Inr empty_msg
      | cons h nil => Inl h
      | _ => Inr ambiguous_msg
      end
  end.

Theorem choose_configured : forall name disc,
  choose_cluster (Some name) disc = Inl name.
Proof. intros. reflexivity. Qed.

Theorem choose_single_discovery : forall name,
  choose_cluster None [name] = Inl name.
Proof. intros. simpl. reflexivity. Qed.

Theorem choose_empty_error :
  choose_cluster None nil = Inr empty_msg.
Proof. simpl. reflexivity. Qed.
