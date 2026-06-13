From ClickhouseProof Require Import Stmt.

Lemma normalize_idempotent : forall cs,
  normalize (normalize cs) = normalize cs.
Proof.
  admit.
Admitted.
