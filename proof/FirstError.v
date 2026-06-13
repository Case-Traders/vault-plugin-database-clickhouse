From Coq Require Import List Arith Lia.
Import ListNotations.

(** Model of firstError: optional first failure in a list. *)
Fixpoint first_error {A} (xs : list A) (failed : bool) (acc : option nat)
    : option nat :=
  match xs, failed with
  | nil, _ => acc
  | cons _ nil, true => Some 0
  | cons _ (_ :: _), true => Some 0
  | cons _ rest, false => first_error rest false acc
  end.

Definition exec_count {A} (n : nat) (xs : list A) : nat :=
  match n with
  | O => 0
  | S _ => if leb n (length xs) then n else length xs
  end.

Lemma first_error_stops : forall (A : Type) (xs : list A) (i : nat),
  i < length xs ->
  exists prefix : list A, length prefix = i.
Proof.
  intros A xs i H.
  exists (@firstn A i xs).
  rewrite (@firstn_length A i xs).
  apply Nat.min_l.
  lia.
Qed.
