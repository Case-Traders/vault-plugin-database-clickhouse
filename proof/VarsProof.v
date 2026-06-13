From ClickhouseProof Require Import Vars.

Theorem update_password_keys_complete :
  keys_present required_update_password required_update_password = true.
Proof. reflexivity. Qed.

Theorem update_expiration_keys_complete :
  keys_present required_update_expiration required_update_expiration = true.
Proof. reflexivity. Qed.

Theorem delete_user_keys_complete :
  keys_present required_delete_user required_delete_user = true.
Proof. reflexivity. Qed.
