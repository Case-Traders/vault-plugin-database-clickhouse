From Coq Require Import List String.
From ClickhouseProof Require Import Cluster.
Import ListNotations.

Open Scope string_scope.

Lemma distinct_nonempty_two : forall c1 c2,
  c1 <> c2 ->
  distinct_nonempty (c1 :: c2 :: nil) = [c1; c2].
Proof.
  intros c1 c2 H. cbn [distinct_nonempty mem].
  destruct (String.eqb c1 c2) eqn:E.
  - apply String.eqb_eq in E. congruence.
  - reflexivity.
Qed.

Theorem choose_ambiguous_error : forall c1 c2,
  c1 <> c2 ->
  choose_cluster None (c1 :: c2 :: nil) = Inr ambiguous_msg.
Proof.
  intros c1 c2 H.
  unfold choose_cluster.
  rewrite (distinct_nonempty_two c1 c2 H).
  reflexivity.
Qed.
