From Coq Require Import List String.
Import ListNotations.

Definition required_new_user : list string :=
  ["name"%string; "username"%string; "password"%string; "expiration"%string; "cluster"%string].

Definition required_update_password : list string :=
  ["name"%string; "username"%string; "password"%string; "cluster"%string].

Definition required_update_expiration : list string :=
  ["name"%string; "username"%string; "expiration"%string; "cluster"%string].

Definition required_delete_user : list string :=
  ["name"%string; "username"%string].

Fixpoint keys_present (required provided : list string) : bool :=
  match required with
  | nil => true
  | cons k rest => existsb (String.eqb k) provided && keys_present rest provided
  end.

Lemma new_user_keys_complete :
  keys_present required_new_user required_new_user = true.
Proof. reflexivity. Qed.
