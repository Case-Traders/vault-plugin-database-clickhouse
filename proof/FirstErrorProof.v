From Coq Require Import List Arith Lia.
From ClickhouseProof Require Import FirstError.
Import ListNotations.

Lemma first_error_stops_qed : forall (A : Type) (xs : list A) (i : nat),
  i < length xs ->
  exists prefix : list A, length prefix = i.
Proof.
  intros A xs i H.
  exists (@firstn A i xs).
  rewrite (@firstn_length A i xs).
  apply Nat.min_l.
  lia.
Qed.
